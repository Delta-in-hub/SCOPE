#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>
// 包含共享头文件和骨架头文件
#include "ggml_base.h"
#include "ggml_base.skel.h" // 由 libbpf 生成

#include "../epoch.h" // 用于 UnixNanoNow()
#include "../ipc_models.h"
#include "../zmqsender.h" // 用于 ZMQ 和 MessagePack 发布
const char *ENDPOINT = "ipc:///tmp/zmq_ipc_pubsub.sock";

#define DEFAULT_TARGET_LIB "/usr/lib/ollama/libggml-base.so"
#define MALLOC_FUNC "ggml_aligned_malloc"
#define FREE_FUNC "ggml_aligned_free"

// 环境配置结构体
static struct env {
    pid_t pid;                       // 过滤 PID
    char filter_comm[TASK_COMM_LEN]; // 过滤进程名
    char target_lib[PATH_MAX];       // 目标库路径
    bool verbose;                    // 详细日志
} env = {
    .pid = 0,
    .filter_comm = "",
    .target_lib = DEFAULT_TARGET_LIB,
    .verbose = false,
};

const char *argp_program_version = "ggml_base 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "ggml_base: Monitor ggml_aligned_malloc/free calls in a shared library "
    "using eBPF.\n\n"
    "USAGE: ./ggml_base [-p PID] [-c COMM] [-l LIBRARY_PATH] [-v]\n";

// 命令行选项定义
static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process PID"},
    {"comm", 'c', "COMMAND", 0,
     "Filter by process command name (max 15 chars)"},
    {"lib", 'l', "LIBRARY_PATH", 0,
     "Path to the target shared library (default: " DEFAULT_TARGET_LIB ")"},
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
        strncpy(env.filter_comm, arg, TASK_COMM_LEN - 1); // 留一位给 \0
        env.filter_comm[TASK_COMM_LEN - 1] = '\0';
        break;
    case 'l':
        if (strlen(arg) >= PATH_MAX) {
            fprintf(stderr, "Target library path too long (max %d): %s\n",
                    PATH_MAX - 1, arg);
            argp_usage(state);
        }
        strncpy(env.target_lib, arg, PATH_MAX - 1);
        env.target_lib[PATH_MAX - 1] = '\0';
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
        return -1; // 收到退出信号后停止处理

    const struct event *e = data;
    zmq_pub_handle_t *zmq_handle = ctx;

    if (env.verbose) {
        char ts[32];
        time_t t = time(NULL);
        struct tm *tm_info = localtime(&t);

        if (tm_info == NULL) {
            perror("localtime failed");
            strcpy(ts, "error");
        } else {
            strftime(ts, sizeof(ts), "%H:%M:%S", tm_info);
        }

        const char *type_str = (e->type == EVENT_MALLOC) ? "MALLOC" : "FREE";

        // 打印事件信息
        // TIME     PID     COMM             TYPE      SIZE       POINTER
        printf("%-8s %-7u %-16s %-8s %-10llu 0x%llx\n", ts, e->pid, e->comm,
               type_str, e->size, e->ptr);
    }

    struct ggml_base_event event = {
        .timestamp_ns = UnixNanoNow(),
        .pid = e->pid,
        .comm = "",
        .type = (e->type == EVENT_MALLOC) ? EVENT_MALLOC : EVENT_FREE,
        .size = e->size,
        .ptr = e->ptr,
    };
    strncpy(event.comm, e->comm, sizeof(event.comm));

    zmq_pub_send(zmq_handle, "ggml_base", &event, ggml_base_event_pack);

    return 0; // 返回 0 表示成功处理事件
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct ggml_base_bpf *skel = NULL;
    int err;
    LIBBPF_OPTS(bpf_uprobe_opts, uprobe_opts); // 初始化 uprobe 选项

    // --- 解析命令行参数 ---
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    // --- 设置 libbpf 日志和信号处理 ---
    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    zmq_pub_handle_t *zmq_handle = zmq_pub_init(ENDPOINT);
    if (!zmq_handle) {
        fprintf(stderr, "Failed to initialize ZMQ publisher\n");
        return 1;
    }

    // --- 打开并加载 BPF 骨架 ---
    skel = ggml_base_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }

    // --- 设置 BPF 程序中的过滤参数 (rodata) ---
    skel->rodata->filter_pid = env.pid;
    strncpy((char *)skel->rodata->filter_comm, env.filter_comm, TASK_COMM_LEN);
    // 确保 null 终止 (虽然 rodata 默认是 0，显式设置更安全)
    ((char *)skel->rodata->filter_comm)[TASK_COMM_LEN - 1] = '\0';

    err = ggml_base_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 附加 BPF 程序 ---
    // 1. ggml_aligned_malloc (uprobe)
    uprobe_opts.func_name = MALLOC_FUNC;
    uprobe_opts.retprobe = false; // 这是入口探针
    skel->links.uprobe_ggml_aligned_malloc =
        bpf_program__attach_uprobe_opts(skel->progs.uprobe_ggml_aligned_malloc,
                                        -1, // pid: -1 表示附加到所有进程
                                        env.target_lib, // 目标库路径
                                        0, // offset: 0 表示函数入口
                                        &uprobe_opts);
    if (!skel->links.uprobe_ggml_aligned_malloc) {
        err = -errno; // libbpf 通常在失败时设置 errno
        fprintf(stderr, "Failed to attach uprobe to %s:%s: %s\n",
                env.target_lib, MALLOC_FUNC, strerror(-err));
        goto cleanup;
    }
    printf("Attached uprobe to %s:%s\n", env.target_lib, MALLOC_FUNC);

    // 2. ggml_aligned_malloc (uretprobe)
    // uprobe_opts.func_name 仍然是 MALLOC_FUNC
    uprobe_opts.retprobe = true; // 这是返回探针
    skel->links.uretprobe_ggml_aligned_malloc = bpf_program__attach_uprobe_opts(
        skel->progs.uretprobe_ggml_aligned_malloc, -1, env.target_lib, 0,
        &uprobe_opts);
    if (!skel->links.uretprobe_ggml_aligned_malloc) {
        err = -errno;
        fprintf(stderr, "Failed to attach uretprobe to %s:%s: %s\n",
                env.target_lib, MALLOC_FUNC, strerror(-err));
        goto cleanup;
    }
    printf("Attached uretprobe to %s:%s\n", env.target_lib, MALLOC_FUNC);

    // 3. ggml_aligned_free (uprobe)
    uprobe_opts.func_name = FREE_FUNC;
    uprobe_opts.retprobe = false; // 这是入口探针
    skel->links.uprobe_ggml_aligned_free =
        bpf_program__attach_uprobe_opts(skel->progs.uprobe_ggml_aligned_free,
                                        -1, env.target_lib, 0, &uprobe_opts);
    if (!skel->links.uprobe_ggml_aligned_free) {
        err = -errno;
        fprintf(stderr, "Failed to attach uprobe to %s:%s: %s\n",
                env.target_lib, FREE_FUNC, strerror(-err));
        goto cleanup;
    }
    printf("Attached uprobe to %s:%s\n", env.target_lib, FREE_FUNC);

    // --- 设置 Ring Buffer ---
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, zmq_handle,
                          NULL);
    if (!rb) {
        err = -errno;
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 事件轮询循环 ---
    printf(
        "Monitoring ggml memory operations in %s (Press Ctrl+C to exit)...\n",
        env.target_lib);
    printf("%-8s %-7s %-16s %-8s %-10s %s\n", "TIME", "PID", "COMM", "TYPE",
           "SIZE", "POINTER");
    while (!exiting) {
        // 设置较短的超时时间（例如 100ms），以便能及时响应 Ctrl+C
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);
        // EINTR 表示被信号中断（例如 Ctrl+C），是正常退出路径
        if (err == -EINTR) {
            err = 0;
            break;
        }
        // 其他负数错误码表示轮询出错
        if (err < 0) {
            fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
            break;
        }
        // err == 0 表示超时，没有事件，继续轮询
    }

cleanup:
    // --- 清理资源 ---
    printf("\nDetaching probes and cleaning up...\n");
    // 销毁 ring buffer
    ring_buffer__free(rb); // 安全起见，即使 rb 为 NULL 也无害

    // 销毁 BPF 骨架 (会自动分离链接并卸载程序/映射)
    ggml_base_bpf__destroy(skel); // 安全起见，即使 skel 为 NULL 也无害

    zmq_pub_cleanup(&zmq_handle);

    printf("Exited.\n");
    // 返回 0 表示成功，非 0 表示有错误发生
    return err < 0 ? -err : 0;
}