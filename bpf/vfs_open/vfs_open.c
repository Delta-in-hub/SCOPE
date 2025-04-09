#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>
// 包含 BPF 骨架头文件
#include "vfs_open.h" // 假设这里定义了 struct event, TASK_COMM_LEN, MAX_FILENAME_LEN
#include "vfs_open.skel.h"

// 包含新增的头文件
#include "../epoch.h" // 用于 UnixNanoNow()
#include "../ipc_models.h"
#include "../zmqsender.h" // 用于 ZMQ 和 MessagePack 发布

// 如果 vfs_open.h 没有定义 MAX_FILENAME_LEN，请在此处添加一个合理的定义
#ifndef MAX_FILENAME_LEN
#define MAX_FILENAME_LEN 256 // 或者根据你的 BPF 程序设置
#endif

// --- 全局环境配置 ---
static struct env {
    pid_t pid;
    char parent_comm[TASK_COMM_LEN];
    bool verbose;
    char zmq_endpoint[256]; // 新增 ZMQ 端点配置
} env = {
    .pid = 0,
    .parent_comm = "",
    .verbose = false,
    .zmq_endpoint = "ipc:///tmp/zmq_ipc_pubsub.sock", // 默认 ZMQ 端点
};

const char *argp_program_version = "vfs_open 0.2 (ZMQ enabled)";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "\n"
    "Trace vfs_open calls and publish events via ZeroMQ.\n"
    "\n"
    "USAGE: ./vfs_open [-p PID] [-c PARENT_COMM] [-e ENDPOINT] [-v]\n";

// --- 命令行参数选项 ---
static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by PID calling execve"},
    {"parent-comm", 'c', "PARENT_COMMAND", 0,
     "Filter by parent process command name"},
    {"endpoint", 'e', "ENDPOINT", 0,
     "ZeroMQ PUB socket endpoint (default: tcp://*:5556)"},
    {"verbose", 'v', NULL, 0, "Verbose debug output (prints to console)"},
    {},
};

// --- 命令行参数解析 ---
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
            fprintf(stderr, "Parent command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.parent_comm, arg, TASK_COMM_LEN);
        env.parent_comm[TASK_COMM_LEN - 1] = '\0';
        break;
    case 'e':
        if (strlen(arg) >= sizeof(env.zmq_endpoint)) {
            fprintf(stderr, "ZMQ Endpoint too long (max %zu): %s\n",
                    sizeof(env.zmq_endpoint) - 1, arg);
            argp_usage(state);
        }
        strncpy(env.zmq_endpoint, arg, sizeof(env.zmq_endpoint));
        env.zmq_endpoint[sizeof(env.zmq_endpoint) - 1] = '\0';
        break;
    case 'v':
        env.verbose = true;
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static const struct argp argp = {
    .options = opts,
    .parser = parse_arg,
    .doc = argp_program_doc,
};

// --- libbpf 打印回调 ---
static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    // 仅在 verbose 模式下打印 libbpf 的 DEBUG 信息
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    // 其他级别的 libbpf 信息（如 WARNING, INFO）总是打印到 stderr
    return vfprintf(stderr, format, args);
}

// --- 信号处理 ---
static volatile bool exiting = false;

static void sig_handler(int sig) { exiting = true; }

// --- BPF 事件处理回调 ---
// 这个函数现在会发送 ZMQ 消息，并且只有在 verbose 模式下才打印
static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)   // 检查退出标志
        return -1; // 返回 -1 会停止 ring_buffer__consume/poll

    // 从 ctx 获取 ZMQ 句柄
    zmq_pub_handle_t *zmq_handle = (zmq_pub_handle_t *)ctx;
    if (!zmq_handle) {
        fprintf(stderr, "ERROR: ZMQ handle is NULL in callback!\n");
        return 0; // 继续处理，但记录错误
    }

    // 从 data 获取 BPF 事件
    const struct event *e = data;

    // 准备要发布的数据
    struct vfs_open_event pub_event;
    pub_event.timestamp_ns = UnixNanoNow(); // 获取高精度时间戳
    pub_event.pid = e->pid;

    // 复制 comm 和 filename, 确保 null 终止
    strncpy(pub_event.comm, e->comm, TASK_COMM_LEN);
    pub_event.comm[TASK_COMM_LEN - 1] = '\0';
    strncpy(pub_event.filename, e->filename, sizeof(pub_event.filename));
    pub_event.filename[sizeof(pub_event.filename) - 1] = '\0';

    // --- 发送 ZMQ 消息 ---
    // 使用 "vfs_open" 作为主题
    int rc =
        zmq_pub_send(zmq_handle, "vfs_open", &pub_event, vfs_open_event_pack);
    if (rc != 0) {
        // zmq_pub_send 内部应该已经打印了错误信息
        fprintf(stderr, "Warning: Failed to send event via ZMQ for PID %d\n",
                e->pid);
        // 即使发送失败，我们通常也希望继续处理其他事件
    }

    // --- 条件打印到控制台 ---
    if (env.verbose) {
        char ts_str[32];
        // 将纳秒时间戳转换为 time_t 用于 strftime
        time_t t_sec = (time_t)(pub_event.timestamp_ns / 1000000000);
        struct tm *tm_info = localtime(&t_sec);
        // 格式化时间 H:M:S
        strftime(ts_str, sizeof(ts_str), "%H:%M:%S", tm_info);

        // 打印包括 cmdline 的详细信息
        printf("%-8s %-7d %-16s %-40s\n", ts_str, pub_event.pid, pub_event.comm,
               pub_event.filename);
    }

    // int64_t cost = UnixNanoNow() - pub_event.timestamp_ns;
    // printf("Cost %ld ns in zmq pub\n", cost);

    return 0; // 返回 0 表示继续处理事件
}

// --- 主函数 ---
int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct vfs_open_bpf *skel = NULL;
    zmq_pub_handle_t *zmq_handle = NULL; // ZMQ 句柄
    int err = 0;

    // 解析命令行参数
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    // 设置 libbpf 打印回调
    libbpf_set_print(libbpf_print_fn);

    // 设置信号处理
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    // --- 初始化 ZMQ 发布者 ---
    zmq_handle = zmq_pub_init(env.zmq_endpoint);
    if (!zmq_handle) {
        fprintf(stderr, "Failed to initialize ZeroMQ publisher on %s\n",
                env.zmq_endpoint);
        err = -1;     // 使用一个非零错误码
        goto cleanup; // 跳转到清理
    }
    printf("INFO: Publishing vfs_open events to ZMQ endpoint: %s\n",
           env.zmq_endpoint);

    // --- 打开、加载和附加 BPF 程序 ---
    skel = vfs_open_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        err = 1;
        goto cleanup;
    }

    // 设置 BPF 程序中的过滤参数 (如果 BPF 程序支持)
    // 注意：原来的代码有 rodata，这里假设仍然适用
    skel->rodata->filter_pid = env.pid;
    memcpy(skel->rodata->filter_comm, env.parent_comm, TASK_COMM_LEN);

    err = vfs_open_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    err = vfs_open_bpf__attach(skel);
    if (err) {
        fprintf(stderr, "Failed to attach BPF programs: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 设置 Ring Buffer ---
    // 注意: 将 zmq_handle 作为上下文参数传递给 handle_event
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, zmq_handle,
                          NULL);
    if (!rb) {
        err = -errno; // ring_buffer__new 在失败时设置 errno
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    // 仅在 verbose 模式下打印表头
    if (env.verbose) {
        printf("%-8s %-7s %-16s %-40s %s\n", "TIME", "PID", "COMM", "FILENAME",
               "CMDLINE");
    } else {
        printf("INFO: Tracing vfs_open calls. Publishing via ZMQ. Use -v for "
               "console output.\n");
    }

    // --- 事件处理循环 ---
    while (!exiting) {
        // 轮询 Ring Buffer，超时时间 100ms
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);
        // EINTR 表示被信号中断 (例如 Ctrl+C)，是正常退出路径
        if (err == -EINTR) {
            err = 0; // 重置错误码为 0，表示正常退出
            break;   // 跳出循环
        }
        // 其他负数错误表示轮询失败
        if (err < 0) {
            fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
            break; // 跳出循环
        }
        // err == 0 表示超时，没有事件，继续循环
        // err > 0 表示处理了 err 个事件，继续循环
    }

// --- 清理资源 ---
cleanup:
    printf("\nINFO: Exiting...\n");
    // 释放 Ring Buffer
    ring_buffer__free(rb); // 安全地处理 NULL 指针

    // 销毁 BPF 骨架 (分离和卸载程序，销毁 maps)
    vfs_open_bpf__destroy(skel); // 安全地处理 NULL 指针

    // 清理 ZMQ 资源
    zmq_pub_cleanup(&zmq_handle); // 安全地处理 NULL 指针

    return err < 0 ? -err : 0; // 返回正数的 POSIX 错误码
}