#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <limits.h> // For PATH_MAX
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>

// 包含你的事件结构定义和 TASK_COMM_LEN
#include "ollamabin.h"
// 包含 libbpf 生成的骨架头文件
#include "ollamabin.skel.h"

#include "../epoch.h" // 用于 UnixNanoNow()
#include "../ipc_models.h"
#include "../zmqsender.h" // 用于 ZMQ 和 MessagePack 发布

const char *ENDPOINT = "ipc:///tmp/zmq_ipc_pubsub.sock";

// 定义默认的探测路径
#define DEFAULT_OLLAMA_PATH "/usr/bin/ollama"
// 定义要探测的函数名
#define TARGET_FUNC_NAME "llamaLog"

// 环境配置结构体
static struct env {
    pid_t pid;                       // 过滤 PID
    char filter_comm[TASK_COMM_LEN]; // 过滤进程名 (修正了原先的 parent_comm)
    char target_path[PATH_MAX];      // 要探测的目标二进制文件路径
    bool verbose;                    // 是否启用详细日志
} env = {
    .pid = 0,
    .filter_comm = "",
    .target_path = DEFAULT_OLLAMA_PATH, // 初始化为默认路径
    .verbose = false,
};

const char *argp_program_version = "ollamabin 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "ollamabin: Monitor Ollama's llamaLog function using eBPF.\n\n"
    "USAGE: ./ollamabin [-p PID] [-c COMM] [-f FILE_PATH] [-v]\n";

// 命令行选项定义
static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process PID"},
    {"comm", 'c', "COMMAND", 0, "Filter by process command name"},
    {"file", 'f', "FILE_PATH", 0,
     "Path to the Ollama binary to probe (default: " DEFAULT_OLLAMA_PATH ")"},
    {"verbose", 'v', NULL, 0, "Verbose debug output"},
    {},
};

// 解析命令行参数
static error_t parse_arg(int key, char *arg, struct argp_state *state) {
    long pid_in;
    switch (key) {
    case 'p':
        errno = 0;
        pid_in = strtol(arg, NULL, 10);
        if (errno || pid_in <= 0) {
            fprintf(stderr, "Invalid PID: %s\n", arg);
            argp_usage(state);
        }
        env.pid = (pid_t)pid_in;
        break;
    case 'c':
        if (strlen(arg) >= TASK_COMM_LEN) {
            fprintf(stderr, "Command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.filter_comm, arg, TASK_COMM_LEN);
        // 不需要手动添加 null 终止符，因为 env.filter_comm 初始化时已清零
        break;
    case 'f':
        if (strlen(arg) >= PATH_MAX) {
            fprintf(stderr, "Target file path too long (max %d): %s\n",
                    PATH_MAX - 1, arg);
            argp_usage(state);
        }
        strncpy(env.target_path, arg, PATH_MAX);
        // 不需要手动添加 null 终止符，因为 env.target_path 初始化时已清零
        break;
    case 'v':
        env.verbose = true;
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

// argp 解析器结构
static const struct argp argp = {
    .options = opts,
    .parser = parse_arg,
    .doc = argp_program_doc,
};

// libbpf 打印回调函数
static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    // 只在 verbose 模式下打印 DEBUG 级别的消息
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}

// 退出标志
static volatile bool exiting = false;

// 信号处理函数
static void sig_handler(int sig) { exiting = true; }

// Ring buffer 事件处理回调函数
static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)
        return -1;

    // 将接收到的原始数据转换为 event 结构体指针
    const struct event *e = data;

    if (e->pid == getpid())
        return 0;

    if (env.verbose) {
        // 获取当前时间并格式化
        char ts[32];
        time_t t = time(NULL);
        struct tm *tm_info = localtime(&t);
        // 检查 localtime 返回值
        if (tm_info == NULL) {
            perror("localtime failed");
            // 可以选择继续，但时间戳会不正确，或者返回错误
            // return -1; // 返回负值会停止 ring_buffer__poll
            strcpy(ts, "error"); // 或提供一个错误指示符
        } else {
            strftime(ts, sizeof(ts), "%H:%M:%S", tm_info);
        }
        // 打印事件信息
        printf("%-8s %-7d %-16s %s\n", ts, e->pid, e->comm, e->text);
    }

    zmq_pub_handle_t *zmq_handle = (zmq_pub_handle_t *)ctx;
    struct llamaLog_event event = {.timestamp_ns = (int64_t)UnixNanoNow(),
                                   .pid = e->pid,
                                   .comm = "",
                                   .text = ""};
    strncpy(event.comm, e->comm, TASK_COMM_LEN);
    strncpy(event.text, e->text, TEXT_LEN);

    zmq_pub_send(zmq_handle, "llamaLog", &event, llamaLog_event_pack);

    return 0; // 返回 0 表示成功处理事件
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct ollamabin_bpf *skel = NULL;

    int err;
    LIBBPF_OPTS(bpf_uprobe_opts, uprobe_opts);

    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    zmq_pub_handle_t *zmq_handle = zmq_pub_init(ENDPOINT);
    if (!zmq_handle) {
        fprintf(stderr, "Failed to initialize ZeroMQ publisher on %s\n",
                ENDPOINT);
        return 1;
    }
    printf("INFO: Publishing ollamabin events to ZMQ endpoint: %s\n", ENDPOINT);

    skel = ollamabin_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }

    // 设置 BPF 程序中的过滤参数 (rodata)
    skel->rodata->filter_pid = env.pid;
    strncpy((char *)skel->rodata->filter_comm, env.filter_comm, TASK_COMM_LEN);
    skel->rodata->filter_comm[TASK_COMM_LEN - 1] = '\0';

    err = ollamabin_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    uprobe_opts.func_name = TARGET_FUNC_NAME;
    uprobe_opts.retprobe = false;

    skel->links.uprobe_llamaLog =
        bpf_program__attach_uprobe_opts(skel->progs.uprobe_llamaLog,
                                        -1, // Attach globally, filter in BPF
                                        env.target_path, // Dynamic target path
                                        0, // Offset 0 for function entry
                                        &uprobe_opts); // Options with func_name

    if (!skel->links.uprobe_llamaLog) {
        err = -errno;
        fprintf(stderr, "Failed to attach uprobe: %d\n", err);
        goto cleanup;
    }

    printf("Successfully attached uprobe to %s:%s\n", env.target_path,
           TARGET_FUNC_NAME);

    // --- 设置 Ring Buffer ---
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, zmq_handle,
                          NULL);
    if (!rb) {
        err = -errno;
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 事件轮询循环 ---
    printf("Monitoring Ollama logs (Press Ctrl+C to exit)...\n");

    if (env.verbose)
        printf("%-8s %-7s %-16s %s\n", "TIME", "PID", "COMM", "LOG_TEXT");

    while (!exiting) {
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);
        if (err == -EINTR) {
            err = 0;
            break;
        }
        if (err < 0) {
            fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
            break;
        }
    }

cleanup:
    // 清理资源
    printf("\nDetaching probes and cleaning up...\n");
    ring_buffer__free(rb);

    ollamabin_bpf__destroy(skel);

    zmq_pub_cleanup(&zmq_handle);

    printf("Exited.\n");
    return err < 0 ? -err : 0;
}