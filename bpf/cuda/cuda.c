#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <limits.h> // PATH_MAX
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>

#include "cuda.h"       // 包含我们自己的定义
#include "cuda.skel.h" // 包含 libbpf 生成的骨架

// 定义默认的 CUDA Runtime 库路径 (可能需要根据系统调整)
// #define DEFAULT_CUDA_LIB_PATH "/usr/local/cuda/lib64/libcudart.so"
// 如果用 /opt/cuda
#define DEFAULT_CUDA_LIB_PATH "/opt/cuda/targets/x86_64-linux/lib/libcudart.so"

// 定义要探测的函数名 (与 BPF 程序中的 SEC 名称无关，这里是库中的符号名)
#define TARGET_FUNC_MALLOC "cudaMalloc"
#define TARGET_FUNC_FREE "cudaFree"
#define TARGET_FUNC_LAUNCH "cudaLaunchKernel"
#define TARGET_FUNC_MEMCPY "cudaMemcpy"
#define TARGET_FUNC_SYNC "cudaDeviceSynchronize"

// 环境配置
static struct env {
	pid_t pid;                       // 过滤 PID
	char filter_comm[TASK_COMM_LEN]; // 过滤进程名
	char target_path[PATH_MAX];      // 目标 CUDA 库路径
	bool verbose;                    // 详细日志
} env = {
    .pid = 0,
    .filter_comm = "",
    .target_path = DEFAULT_CUDA_LIB_PATH,
    .verbose = false,
};

const char* argp_program_version = "cuda 0.1";
const char* argp_program_bug_address = "DeltaMail@qq.com"; // 替换
const char argp_program_doc[] =
    "cuda: Monitor CUDA Runtime API calls using eBPF.\n\n"
    "Traces cudaMalloc, cudaFree, cudaLaunchKernel, cudaMemcpy, cudaDeviceSynchronize.\n\n"
    "USAGE: ./cuda [-p PID] [-c COMM] [-f FILE_PATH] [-v]\n";

// 命令行选项
static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process PID"},
    {"comm", 'c', "COMMAND", 0, "Filter by process command name"},
    {"file", 'f', "FILE_PATH", 0, "Path to the target libcudart.so (default: " DEFAULT_CUDA_LIB_PATH ")"},
    {"verbose", 'v', NULL, 0, "Verbose debug output"},
    {},
};

// 参数解析
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
			fprintf(stderr, "Command name too long (max %d): %s\n", TASK_COMM_LEN - 1, arg);
			argp_usage(state);
		}
		strncpy(env.filter_comm, arg, TASK_COMM_LEN);
		env.filter_comm[TASK_COMM_LEN - 1] = '\0';
		break;
	case 'f':
		if (strlen(arg) >= PATH_MAX) {
			fprintf(stderr, "Target file path too long (max %d): %s\n", PATH_MAX - 1, arg);
			argp_usage(state);
		}
		strncpy(env.target_path, arg, PATH_MAX);
		env.target_path[PATH_MAX - 1] = '\0';
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

// libbpf 日志回调
static int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args) {
	if (level == LIBBPF_DEBUG && !env.verbose)
		return 0;
	return vfprintf(stderr, format, args);
}

static volatile bool exiting = false;

static void sig_handler(int sig) { exiting = true; }

// Helper to convert cudaMemcpyKind to string
static const char* memcpy_kind_to_str(enum cuda_memcpy_kind kind) {
    switch (kind) {
        case CUDA_MEMCPY_HOST_TO_HOST: return "HostToHost";
        case CUDA_MEMCPY_HOST_TO_DEVICE: return "HostToDevice";
        case CUDA_MEMCPY_DEVICE_TO_HOST: return "DeviceToHost";
        case CUDA_MEMCPY_DEVICE_TO_DEVICE: return "DeviceToDevice";
        case CUDA_MEMCPY_DEFAULT: return "Default";
        default: return "Unknown";
    }
}


// Ring buffer 事件处理回调
static int handle_event(void *ctx, void *data, size_t data_sz) {
	if (exiting) return -1;

	const struct event *e = data;

	// 获取当前用户空间时间用于打印 (注意 BPF 时间戳是 ktime)
	char ts[32];
	time_t t = time(NULL);
	struct tm *tm_info = localtime(&t);
	if (tm_info == NULL) {
		perror("localtime failed");
		strcpy(ts, "ERROR");
	} else {
		strftime(ts, sizeof(ts), "%H:%M:%S", tm_info);
	}

    // 打印通用前缀
    printf("%-8s %-7d %-16s ", ts, e->pid, e->comm);

	// 根据事件类型打印特定信息
	switch (e->type) {
	case EVENT_TYPE_MALLOC:
		// "%d %s cudaMalloc failed\n" or "%d %s cudaMalloc %p\n"
		if (e->malloc.retval != 0) {
			printf("cudaMalloc failed (ret=%d)\n", e->malloc.retval);
		} else {
			printf("cudaMalloc => ptr=%p %zu bytes\n", e->malloc.allocated_ptr, e->malloc.size);
		}
		break;
	case EVENT_TYPE_FREE:
		// "%d %s cudaFree %p\n"
		printf("cudaFree(ptr=%p)\n", e->free.dev_ptr);
		break;
	case EVENT_TYPE_LAUNCH_KERNEL:
		// "%d %s cudaLaunchKernel %p with shared memory %d\n"
		printf("cudaLaunchKernel(func=%p)\n",
		       e->launch_kernel.func_ptr);
		uintptr_t funcp = (uintptr_t)e->launch_kernel.func_ptr;
		break;
	case EVENT_TYPE_MEMCPY:
        // "%d %s cudaMemcpy<Kind> %p -> %p %d\n"
		printf("cudaMemcpy %s(src=%p, dst=%p, size=%zu)\n",
               memcpy_kind_to_str(e->memcpy.kind),
               e->memcpy.src, e->memcpy.dst, e->memcpy.size);
		break;
	case EVENT_TYPE_SYNC:
		// "%d %s call cudaDeviceSynchronize at %u\n" (bpftrace used nsecs, let's just mark entry)
        printf("cudaDeviceSynchronize cost %lu ns\n", e->sync.duration_ns);
		break;
	default:
		fprintf(stderr, "Warning: Unknown event type received: %d\n", e->type);
		break;
	}

	return 0;
}


int main(int argc, char **argv) {
	struct ring_buffer *rb = NULL;
	struct cuda_bpf *skel = NULL;
	int err;
	LIBBPF_OPTS(bpf_uprobe_opts, uprobe_opts);

	// --- 初始化 ---
	err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
	if (err) return err;

	libbpf_set_print(libbpf_print_fn);
	signal(SIGINT, sig_handler);
	signal(SIGTERM, sig_handler);

	// 提高资源限制
	struct rlimit rlim_new = { .rlim_cur = RLIM_INFINITY, .rlim_max = RLIM_INFINITY };
	if (setrlimit(RLIMIT_MEMLOCK, &rlim_new)) {
		fprintf(stderr, "Warning: Failed to increase RLIMIT_MEMLOCK limit!\n");
	}

	// --- 打开、加载 BPF 程序 ---
	skel = cuda_bpf__open();
	if (!skel) {
		fprintf(stderr, "Failed to open BPF skeleton\n");
		return 1;
	}

	// 设置 BPF 程序过滤参数
	skel->rodata->filter_pid = env.pid;
	strncpy((char *)skel->rodata->filter_comm, env.filter_comm, TASK_COMM_LEN);

	err = cuda_bpf__load(skel);
	if (err) {
		fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
		goto cleanup;
	}

	// --- 附加 BPF 探针 ---
	printf("Attaching probes to %s ...\n", env.target_path);

    // Macro to simplify attaching uprobe/uretprobe pairs
    // prog_uprobe_name should be the member name like "uprobe_cudaMalloc"
    // prog_uretprobe_name should be the member name like "uretprobe_cudaMalloc"
    #define ATTACH_UPROBE_URETPROBE(skel, sym, prog_uprobe_name, prog_uretprobe_name) \
        do { \
            /* Attach Uprobe */ \
            uprobe_opts.func_name = sym; \
            uprobe_opts.retprobe = false; \
            /* Use the provided uprobe name directly */ \
            skel->links.prog_uprobe_name = bpf_program__attach_uprobe_opts( \
                skel->progs.prog_uprobe_name, -1, env.target_path, 0, &uprobe_opts); \
            if (!skel->links.prog_uprobe_name) { \
                err = -errno; \
                fprintf(stderr, "Failed to attach uprobe %s: %s\n", sym, strerror(-err)); \
                goto cleanup; \
            } \
            /* Attach Uretprobe */ \
            /* func_name (sym) remains the same, just set retprobe=true */ \
            uprobe_opts.retprobe = true; \
            /* Use the provided uretprobe name directly */ \
            skel->links.prog_uretprobe_name = bpf_program__attach_uprobe_opts( \
                skel->progs.prog_uretprobe_name, -1, env.target_path, 0, &uprobe_opts); \
            if (!skel->links.prog_uretprobe_name) { \
                err = -errno; \
                fprintf(stderr, "Failed to attach uretprobe %s: %s\n", sym, strerror(-err)); \
                goto cleanup; \
            } \
        } while(0)

    // Macro to simplify attaching only uprobe (保持不变或微调参数名)
    #define ATTACH_UPROBE(skel, sym, prog_uprobe_name) \
        do { \
            uprobe_opts.func_name = sym; \
            uprobe_opts.retprobe = false; \
            skel->links.prog_uprobe_name = bpf_program__attach_uprobe_opts( \
                skel->progs.prog_uprobe_name, -1, env.target_path, 0, &uprobe_opts); \
            if (!skel->links.prog_uprobe_name) { \
                err = -errno; \
                fprintf(stderr, "Failed to attach uprobe %s: %s\n", sym, strerror(-err)); \
                goto cleanup; \
            } \
        } while(0)

    // Attach cudaMalloc (uprobe + uretprobe)
    // Pass both the uprobe and uretprobe member names
    ATTACH_UPROBE_URETPROBE(skel, TARGET_FUNC_MALLOC, uprobe_cudaMalloc, uretprobe_cudaMalloc);

    // Attach cudaFree (uprobe only - 调用不变)
    ATTACH_UPROBE(skel, TARGET_FUNC_FREE, uprobe_cudaFree);

    // Attach cudaLaunchKernel (uprobe only - 调用不变)
    ATTACH_UPROBE(skel, TARGET_FUNC_LAUNCH, uprobe_cudaLaunchKernel);

    // Attach cudaMemcpy (uprobe only - 调用不变)
    ATTACH_UPROBE(skel, TARGET_FUNC_MEMCPY, uprobe_cudaMemcpy);

    // Attach cudaDeviceSynchronize (uprobe + uretprobe)
    // Pass both the uprobe and uretprobe member names
    ATTACH_UPROBE_URETPROBE(skel, TARGET_FUNC_SYNC, uprobe_cudaDeviceSynchronize, uretprobe_cudaDeviceSynchronize);

	printf("Successfully attached probes.\n");

	// --- 设置 Ring Buffer ---
	rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
	if (!rb) {
		err = -errno;
		fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
		goto cleanup;
	}

	// --- 事件轮询循环 ---
	printf("Monitoring CUDA API calls (Press Ctrl+C to exit)...\n");
    printf("%-8s %-7s %-16s %s\n", "TIME", "PID", "COMM", "EVENT DETAILS");
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
	printf("\nDetaching probes and cleaning up...\n");
	ring_buffer__free(rb);
	cuda_bpf__destroy(skel); // Automatically detaches links in skel->links

	printf("Exited.\n");
	return err < 0 ? -err : 0;
}