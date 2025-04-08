#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "ggml_base.h"

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// --- Maps ---
// 1. Ring Buffer Map: 发送事件到用户空间
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // 1MB Ring Buffer
} rb SEC(".maps");

// 2. Hash Map: 在 malloc uprobe 和 uretprobe 之间传递 size
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ENTRIES); // 使用头文件中定义的常量
    __type(key, pid_t);               // Key: process ID
    __type(value, size_t);            // Value: allocation size
} malloc_size_map SEC(".maps");


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



// --- BPF 程序 ---

// void * ggml_aligned_malloc(size_t size)
SEC("uprobe")
int BPF_KPROBE(uprobe_ggml_aligned_malloc, size_t size) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    // 应用过滤规则
    if (!process_allowed(pid, comm)) {
        return 0;
    }

    // 将 size 存储到 map 中，以 pid 为 key
    int ret = bpf_map_update_elem(&malloc_size_map, &pid, &size, BPF_ANY);
    if (ret != 0) {
        bpf_printk("WARN: malloc: failed to update map, pid %d, ret %d\n", pid, ret);
    }
    // bpf_printk("DEBUG: malloc entry: pid %d, comm %s, size %lld\n", pid, comm, size);
    return 0;
}

SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_ggml_aligned_malloc, void *ret) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    // 查找对应的 size
    size_t *size_ptr = bpf_map_lookup_elem(&malloc_size_map, &pid);
    if (!size_ptr) {
        // bpf_printk("WARN: malloc ret: size not found for pid %d\n", pid);
        // 可能是 uprobe 失败或过滤掉了，这里不需要删除 map，因为它不存在
        return 0; // 没有找到对应的 entry probe 数据
    }

    size_t size = *size_ptr;

    // 清理 map 中的条目
    bpf_map_delete_elem(&malloc_size_map, &pid);


    // 检查 malloc 是否成功 (ret == NULL 表示失败)
    if (ret == NULL) {
        // bpf_printk("DEBUG: malloc ret: pid %d, comm %s, size %lld, FAILED (ret=NULL)\n", pid, comm, size);
        return 0; // 分配失败，不记录事件
    }

    // 记录成功分配的事件
    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        bpf_printk("WARN: malloc ret: ringbuf reserve failed, pid %d\n", pid);
        return 0;
    }

    e->type = EVENT_MALLOC;
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    e->size = size;
    e->ptr = (unsigned long long)ret; // 将指针转换为 u64
    bpf_ringbuf_submit(e, 0);
    return 0;
}


// void ggml_aligned_free(void * ptr, size_t size)
// 注意：根据 bpftrace 脚本，我们假设 free 函数有两个参数
SEC("uprobe")
int BPF_KPROBE(uprobe_ggml_aligned_free, void *ptr, size_t size) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    // 应用过滤规则
    if (!process_allowed(pid, comm)) {
        return 0;
    }

    // bpf_printk("DEBUG: free entry: pid %d, comm %s, size %lld, ptr %llx\n", pid, comm, size, (unsigned long long)ptr);

    // 记录 free 事件
    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        bpf_printk("WARN: free: ringbuf reserve failed, pid %d\n", pid);
        return 0;
    }

    e->type = EVENT_FREE;
    e->pid = pid;
    // 直接使用上面获取的 comm
    __builtin_memcpy(e->comm, comm, sizeof(e->comm));
    e->size = size;
    e->ptr = (unsigned long long)ptr; // 将指针转换为 u64

    bpf_ringbuf_submit(e, 0);
    return 0;
}