#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "cuda.h" // 自己的头文件必须在最下面

// --- 可配置参数 ---
const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN] = "";

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// --- Maps ---
struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1024 * 1024); // 1 MB Ring buffer
} rb SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 10240);
	__type(key, pid_t);
	__type(value, struct malloc_entry_data);
} malloc_entries SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 10240);
	__type(key, pid_t);
	__type(value, struct sync_entry_data);
} sync_entries SEC(".maps");

// --- 内联辅助函数 (来自示例) ---
static __always_inline int comm_allowed(const char *comm) {
#pragma unroll
	for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
		if (comm[i] != filter_comm[i])
			return 0;
	}
	return 1;
}

static __always_inline int process_allowed(pid_t pid, const char *comm) {
	if (filter_pid != 0 && pid != filter_pid)
		return 0;
	if (filter_comm[0] != '\0' && !comm_allowed(comm))
		return 0;
	return 1;
}

// --- 探针函数 ---

// cudaMalloc(void** devPtr, size_t size)
SEC("uprobe")
int BPF_KPROBE(uprobe_cudaMalloc, void** devPtr_addr, size_t size) {
	pid_t pid = bpf_get_current_pid_tgid() >> 32;
	char comm[TASK_COMM_LEN];
	bpf_get_current_comm(&comm, sizeof(comm));

	if (!process_allowed(pid, comm))
		return 0;

	// 存储 devPtr 的地址，供 uretprobe 读取实际分配的指针
	struct malloc_entry_data entry = {};
	entry.user_dev_ptr_addr = devPtr_addr;
    entry.size = size;
	bpf_map_update_elem(&malloc_entries, &pid, &entry, BPF_ANY);

	return 0;
}

SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_cudaMalloc, int ret) { // cudaError_t is typically int
	pid_t pid = bpf_get_current_pid_tgid() >> 32;

	// 查找对应的入口数据 (devPtr 的地址)
	struct malloc_entry_data* entry_p = bpf_map_lookup_elem(&malloc_entries, &pid);
	if (!entry_p) {
		return 0; // 没有找到入口记录
	}
	// 获取到地址后立即删除 map 条目
    void** user_addr = entry_p->user_dev_ptr_addr;
    size_t size = entry_p->size;
	bpf_map_delete_elem(&malloc_entries, &pid);

	// 准备 Malloc Exit 事件
	struct event* e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e) {
		return 0;
	}

	e->type = EVENT_TYPE_MALLOC;
	e->pid = pid;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));
    e->malloc.size = size;
	e->malloc.retval = ret;
	e->malloc.allocated_ptr = NULL; // 默认为 NULL

	// 如果 cudaMalloc 成功 (retval == 0)，尝试读取实际分配的指针
	if (ret == 0 && user_addr != NULL) {
		void* actual_dev_ptr = NULL;
        // 从用户空间 devPtr_addr 指向的位置读取指针值
		bpf_probe_read_user(&actual_dev_ptr, sizeof(actual_dev_ptr), user_addr);
        // bpf_probe_read_user 可能失败，但即使失败，我们仍然发送事件
        // (allocated_ptr 将保持 NULL)
		e->malloc.allocated_ptr = actual_dev_ptr;
	}

	bpf_ringbuf_submit(e, 0);
	return 0;
}

// cudaFree(void* devPtr)
SEC("uprobe")
int BPF_KPROBE(uprobe_cudaFree, void* devPtr) {
	pid_t pid = bpf_get_current_pid_tgid() >> 32;
	char comm[TASK_COMM_LEN];
	bpf_get_current_comm(&comm, sizeof(comm));

	if (!process_allowed(pid, comm))
		return 0;

	struct event* e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e) {
		return 0;
	}

	e->type = EVENT_TYPE_FREE;
	e->pid = pid;
	__builtin_memcpy(&e->comm, comm, sizeof(comm));
	e->free.dev_ptr = devPtr;

	bpf_ringbuf_submit(e, 0);
	return 0;
}

// cudaLaunchKernel(const void *func, dim3 gridDim, dim3 blockDim, void **args, size_t sharedMem, cudaStream_t stream)
// 我们只关心 func (arg0)
SEC("uprobe")
int BPF_KPROBE(uprobe_cudaLaunchKernel, const void* func) {
	pid_t pid = bpf_get_current_pid_tgid() >> 32;
	char comm[TASK_COMM_LEN];
	bpf_get_current_comm(&comm, sizeof(comm));

	if (!process_allowed(pid, comm))
		return 0;

	struct event* e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e) {
		return 0;
	}

	e->type = EVENT_TYPE_LAUNCH_KERNEL;
	e->pid = pid;
	__builtin_memcpy(&e->comm, comm, sizeof(comm));
	e->launch_kernel.func_ptr = func;

	bpf_ringbuf_submit(e, 0);
	return 0;
}

// cudaMemcpy(void* dst, const void* src, size_t size, enum cudaMemcpyKind kind)
SEC("uprobe")
int BPF_KPROBE(uprobe_cudaMemcpy, void* dst, const void* src, size_t size, int kind) { // kind is enum, treat as int
	pid_t pid = bpf_get_current_pid_tgid() >> 32;
	char comm[TASK_COMM_LEN];
	bpf_get_current_comm(&comm, sizeof(comm));

	if (!process_allowed(pid, comm))
		return 0;

    // 验证 kind 是否在预期范围内 (可选但推荐)
    if (kind < CUDA_MEMCPY_HOST_TO_HOST || kind > CUDA_MEMCPY_DEFAULT) {
        // 可以选择记录一个警告或忽略无效 kind
        bpf_printk("cudaMemcpy unknown kind: %d", kind);
        kind = -1; // 标记为未知
    }


	struct event* e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e) {
		return 0;
	}

	e->type = EVENT_TYPE_MEMCPY;
	e->pid = pid;
	__builtin_memcpy(&e->comm, comm, sizeof(comm));
	e->memcpy.dst = dst;
	e->memcpy.src = src;
	e->memcpy.size = size;
	e->memcpy.kind = (enum cuda_memcpy_kind)kind;

	bpf_ringbuf_submit(e, 0);
	return 0;
}

// cudaDeviceSynchronize()
SEC("uprobe")
int BPF_KPROBE(uprobe_cudaDeviceSynchronize) {
	pid_t pid = bpf_get_current_pid_tgid() >> 32;
	char comm[TASK_COMM_LEN];
	bpf_get_current_comm(&comm, sizeof(comm));

	if (!process_allowed(pid, comm))
		return 0;

    uint64_t ts = bpf_ktime_get_ns();

	// 存储入口时间戳
	struct sync_entry_data entry = {};
	entry.entry_ts = ts;
	bpf_map_update_elem(&sync_entries, &pid, &entry, BPF_ANY);

	return 0;
}

SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_cudaDeviceSynchronize, int ret) {
	pid_t pid = bpf_get_current_pid_tgid() >> 32;

	// 查找入口时间戳
	struct sync_entry_data* entry_p = bpf_map_lookup_elem(&sync_entries, &pid);
	if (!entry_p) {
		return 0;
	}
    uint64_t entry_ts = entry_p->entry_ts;
	bpf_map_delete_elem(&sync_entries, &pid);

    uint64_t exit_ts = bpf_ktime_get_ns();
    uint64_t duration_ns = exit_ts - entry_ts;

	// 发送 Sync Exit 事件
	struct event* e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e) {
		return 0;
	}

	e->type = EVENT_TYPE_SYNC;
	e->pid = pid;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));
	e->sync.duration_ns = duration_ns;

	bpf_ringbuf_submit(e, 0);
	return 0;
}