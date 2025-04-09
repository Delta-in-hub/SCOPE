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

// 包含我们自己的事件结构定义
#include "ggml_cuda.h"
// 包含 libbpf 生成的骨架头文件
#include "ggml_cuda.skel.h"

#include "../epoch.h" // 用于 UnixNanoNow()
#include "../ipc_models.h"
#include "../zmqsender.h" // 用于 ZMQ 和 MessagePack 发布

const char *ENDPOINT = "ipc:///tmp/zmq_ipc_pubsub.sock";

// 定义默认的探测库路径
#define DEFAULT_TARGET_LIB "/usr/lib/ollama/cuda_v12/libggml-cuda.so"

// 定义要探测的 C++ mangled 符号名称
#define TARGET_FUNC_MUL_MAT_VEC_Q                                              \
    "_Z26ggml_cuda_op_mul_mat_vec_qR25ggml_backend_cuda_contextPK11ggml_"      \
    "tensorS3_PS1_PKcPKfS6_PfllllP11CUstream_st"
#define TARGET_FUNC_MUL_MAT_Q                                                  \
    "_Z22ggml_cuda_op_mul_mat_qR25ggml_backend_cuda_contextPK11ggml_tensorS3_" \
    "PS1_PKcPKfS6_PfllllP11CUstream_st"
#define TARGET_FUNC_SET_DEVICE "_Z20ggml_cuda_set_devicei"

// 环境配置结构体
static struct env {
    pid_t pid;                       // 过滤 PID
    char filter_comm[TASK_COMM_LEN]; // 过滤进程名
    char target_path[PATH_MAX];      // 要探测的目标共享库路径
    bool verbose;                    // 是否启用详细日志
} env = {
    .pid = 0,
    .filter_comm = "",
    .target_path = DEFAULT_TARGET_LIB, // 初始化为默认路径
    .verbose = false,
};

const char *argp_program_version = "ggml_cuda 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com"; // 替换成你的邮箱
const char argp_program_doc[] =
    "ggml_cuda: Monitor specific ggml-cuda functions using eBPF.\n\n"
    "Monitors function call durations and device settings.\n\n"
    "USAGE: ./ggml_cuda [-p PID] [-c COMM] [-f FILE_PATH] [-v]\n";

// 命令行选项定义
static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process PID"},
    {"comm", 'c', "COMMAND", 0, "Filter by process command name"},
    {"file", 'f', "FILE_PATH", 0,
     "Path to the target libggml-cuda.so (default: " DEFAULT_TARGET_LIB ")"},
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
        env.filter_comm[TASK_COMM_LEN - 1] = '\0'; // 确保 null 结尾
        break;
    case 'f':
        if (strlen(arg) >= PATH_MAX) {
            fprintf(stderr, "Target file path too long (max %d): %s\n",
                    PATH_MAX - 1, arg);
            argp_usage(state);
        }
        strncpy(env.target_path, arg, PATH_MAX);
        env.target_path[PATH_MAX - 1] = '\0'; // 确保 null 结尾
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
        return -1; // 如果正在退出，停止处理

    const struct event *e = data;
    zmq_pub_handle_t *zmq_handle = ctx;

    if (env.verbose) {
        // 获取当前时间并格式化
        char ts[32];
        time_t t = time(NULL);
        struct tm *tm_info = localtime(&t);
        if (tm_info == NULL) {
            perror("localtime failed");
            strcpy(ts, "ERROR"); // 或提供一个错误指示符
        } else {
            strftime(ts, sizeof(ts), "%H:%M:%S", tm_info);
        }

        // 根据事件类型打印不同的信息
        switch (e->type) {
        case EVENT_TYPE_FUNC_DURATION:
            printf("%-8s %-7d %-16s FUNC %-30s cost %lu ns\n", ts, e->pid,
                   e->comm, e->func_duration.func_name,
                   e->func_duration.duration_ns);
            break;
        case EVENT_TYPE_SET_DEVICE:
            printf("%-8s %-7d %-16s CALL ggml_cuda_set_device(%d)\n", ts,
                   e->pid, e->comm, e->set_device.device_id);
            break;
        default:
            fprintf(stderr, "Warning: Unknown event type received: %d\n",
                    e->type);
            break;
        }
    }

    struct ggml_cuda_event event = {
        .timestamp_ns = UnixNanoNow(),
        .pid = e->pid,
        .comm = "",
        .func_name = "",
        .duration_ns = e->func_duration.duration_ns,
    };
    strncpy(event.comm, e->comm, TASK_COMM_LEN);
    strncpy(event.func_name, e->func_duration.func_name, MAX_FUNC_NAME_LEN);

    zmq_pub_send(zmq_handle, "ggml_cuda", &event, ggml_cuda_event_pack);

    return 0; // 返回 0 表示成功处理事件
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct ggml_cuda_bpf *skel = NULL; // 使用 skeleton 生成的类型
    int err;
    LIBBPF_OPTS(bpf_uprobe_opts, uprobe_opts); // 用于附加 uprobe/uretprobe

    // --- 初始化 ---
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    // 提高资源限制，对于 BPF 程序加载和运行有时是必要的
    struct rlimit rlim_new = {
        .rlim_cur = RLIM_INFINITY,
        .rlim_max = RLIM_INFINITY,
    };
    if (setrlimit(RLIMIT_MEMLOCK, &rlim_new)) {
        fprintf(stderr, "Failed to increase RLIMIT_MEMLOCK limit!\n");
        // return 1; // 可以选择警告而不是退出
    }

    zmq_pub_handle_t *zmq_handle = zmq_pub_init(ENDPOINT);
    if (!zmq_handle) {
        fprintf(stderr, "Failed to initialize ZMQ publisher %s\n", ENDPOINT);
        return 1;
    }

    // --- 打开、加载 BPF 程序 ---
    skel = ggml_cuda_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }

    // 设置 BPF 程序中的过滤参数 (rodata)
    skel->rodata->filter_pid = env.pid;
    strncpy((char *)skel->rodata->filter_comm, env.filter_comm, TASK_COMM_LEN);
    // BPF 端 filter_comm 已经初始化为全零，strncpy 会复制 null 终止符
    // (如果源字符串长度小于 TASK_COMM_LEN)

    err = ggml_cuda_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 附加 BPF 探针 ---
    printf("Attaching probes to %s ...\n", env.target_path);

    // 附加 ggml_cuda_op_mul_mat_vec_q (uprobe + uretprobe)
    uprobe_opts.func_name = TARGET_FUNC_MUL_MAT_VEC_Q;
    uprobe_opts.retprobe = false; // 这是入口探针
    skel->links.uprobe_ggml_cuda_op_mul_mat_vec_q =
        bpf_program__attach_uprobe_opts(
            skel->progs.uprobe_ggml_cuda_op_mul_mat_vec_q, // BPF 程序
            -1,              // pid = -1 表示全局附加 (将在 BPF 程序内过滤)
            env.target_path, // 目标二进制/库文件路径
            0,               // offset = 0 表示附加到函数入口
            &uprobe_opts);   // 包含 func_name 等选项
    if (!skel->links.uprobe_ggml_cuda_op_mul_mat_vec_q) {
        err = -errno;
        fprintf(stderr, "Failed to attach uprobe %s: %s\n",
                TARGET_FUNC_MUL_MAT_VEC_Q, strerror(-err));
        goto cleanup;
    }
    // 附加返回探针 (使用相同的符号名，但指定 retprobe=true)
    uprobe_opts.retprobe = true; // 这是返回探针
    skel->links
        .uretprobe_ggml_cuda_op_mul_mat_vec_q = bpf_program__attach_uprobe_opts(
        skel->progs.uretprobe_ggml_cuda_op_mul_mat_vec_q, // BPF 返回探针程序
        -1, env.target_path,
        0, // offset 仍然是 0，libbpf 会处理
        &uprobe_opts);
    if (!skel->links.uretprobe_ggml_cuda_op_mul_mat_vec_q) {
        err = -errno;
        fprintf(stderr, "Failed to attach uretprobe %s: %s\n",
                TARGET_FUNC_MUL_MAT_VEC_Q, strerror(-err));
        goto cleanup;
    }

    // 附加 ggml_cuda_op_mul_mat_q (uprobe + uretprobe)
    uprobe_opts.func_name = TARGET_FUNC_MUL_MAT_Q;
    uprobe_opts.retprobe = false;
    skel->links.uprobe_ggml_cuda_op_mul_mat_q = bpf_program__attach_uprobe_opts(
        skel->progs.uprobe_ggml_cuda_op_mul_mat_q, -1, env.target_path, 0,
        &uprobe_opts);
    if (!skel->links.uprobe_ggml_cuda_op_mul_mat_q) {
        err = -errno;
        fprintf(stderr, "Failed to attach uprobe %s: %s\n",
                TARGET_FUNC_MUL_MAT_Q, strerror(-err));
        goto cleanup;
    }
    uprobe_opts.retprobe = true;
    skel->links.uretprobe_ggml_cuda_op_mul_mat_q =
        bpf_program__attach_uprobe_opts(
            skel->progs.uretprobe_ggml_cuda_op_mul_mat_q, -1, env.target_path,
            0, &uprobe_opts);
    if (!skel->links.uretprobe_ggml_cuda_op_mul_mat_q) {
        err = -errno;
        fprintf(stderr, "Failed to attach uretprobe %s: %s\n",
                TARGET_FUNC_MUL_MAT_Q, strerror(-err));
        goto cleanup;
    }

    // 附加 ggml_cuda_set_device (uprobe only)
    // uprobe_opts.func_name = TARGET_FUNC_SET_DEVICE;
    // uprobe_opts.retprobe = false;
    // skel->links.uprobe_ggml_cuda_set_device =
    // bpf_program__attach_uprobe_opts(
    //     skel->progs.uprobe_ggml_cuda_set_device, -1, env.target_path, 0,
    //     &uprobe_opts);
    // if (!skel->links.uprobe_ggml_cuda_set_device) {
    //     err = -errno;
    //     fprintf(stderr, "Failed to attach uprobe %s: %s\n",
    //     TARGET_FUNC_SET_DEVICE, strerror(-err)); goto cleanup;
    // }

    printf("Successfully attached probes.\n");

    // --- 设置 Ring Buffer ---
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, zmq_handle,
                          NULL);
    if (!rb) {
        err = -errno; // ring_buffer__new sets errno on failure
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    // --- 事件轮询循环 ---
    printf("Monitoring ggml-cuda functions (Press Ctrl+C to exit)...\n");

    if (env.verbose)
        printf("%-8s %-7s %-16s %s\n", "TIME", "PID", "COMM", "DETAILS");

    while (!exiting) {
        // Poll the ring buffer with a timeout (e.g., 100ms)
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);
        // Ctrl-C will cause -EINTR
        if (err == -EINTR) {
            err = 0;
            break;
        }
        // Handle other errors
        if (err < 0) {
            fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
            break;
        }
        // err >= 0 indicates number of events processed (or 0 if timeout)
    }

cleanup:
    // --- 清理资源 ---
    printf("\nDetaching probes and cleaning up...\n");

    // ring_buffer__free 会处理 ring buffer 相关资源
    ring_buffer__free(rb);

    // ggml_cuda_bpf__destroy 会自动 detach 通过 skel->links
    // 附加的探针并释放所有 BPF 对象
    ggml_cuda_bpf__destroy(skel);
    zmq_pub_cleanup(&zmq_handle);

    printf("Exited.\n");
    return err < 0 ? -err : 0;
}