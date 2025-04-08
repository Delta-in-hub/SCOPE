#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "ggml_cuda.h" // 自己的头文件必须在最下面, 顺序很重要

// --- 可配置参数 (由用户空间设置) ---
const volatile pid_t filter_pid = 0;                  // 过滤 PID (0 表示不过滤)
const volatile char filter_comm[TASK_COMM_LEN] = ""; // 过滤进程名 (空字符串表示不过滤)

// --- BPF 程序许可证 ---
char LICENSE[] SEC("license") = "Dual BSD/GPL";

// --- Maps ---

// Ring buffer map: 用于向用户空间发送事件数据
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // 1 MB Ring buffer
} rb SEC(".maps");


// 用于在 uprobe 和 uretprobe 之间传递时间戳的 map 值
struct entry_data {
    uint64_t ts; // 进入函数时的时间戳
};


// Hash map: 用于存储函数进入时间戳 (pid -> entry_data)
// 为每个需要计时的函数创建一个 map
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240); // 根据预期的并发 PID 数量调整
    __type(key, pid_t);
    __type(value, struct entry_data);
} mul_mat_vec_q_entry SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, pid_t);
    __type(value, struct entry_data);
} mul_mat_q_entry SEC(".maps");

// --- 内联辅助函数 ---

// 检查进程名是否匹配过滤规则 (来自示例, 无需修改)
static __always_inline int comm_allowed(const char *comm) {
#pragma unroll
    for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
        if (comm[i] != filter_comm[i])
            return 0;
    }
    return 1;
}

// 检查进程 PID 和名称是否允许被追踪 (来自示例, 无需修改)
static __always_inline int process_allowed(pid_t pid, const char *comm) {
    if (filter_pid != 0 && pid != filter_pid)
        return 0;
    if (filter_comm[0] != '\0' && !comm_allowed(comm))
        return 0;
    return 1;
}

// --- uprobe/uretprobe 函数 ---

// --- ggml_cuda_op_mul_mat_vec_q ---
// 注意：函数签名中的参数在这里并不直接使用，但写出来有助于理解。
// _Z26ggml_cuda_op_mul_mat_vec_qR25ggml_backend_cuda_contextPK11ggml_tensorS3_PS1_PKcPKfS6_PfllllP11CUstream_st
SEC("uprobe")
int BPF_KPROBE(uprobe_ggml_cuda_op_mul_mat_vec_q) // 简化函数名以符合 C 标识符规则
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    if (!process_allowed(pid, comm))
        return 0;

    struct entry_data data = {};
    data.ts = bpf_ktime_get_ns(); // 获取当前时间戳

    // 存储时间戳到 map
    bpf_map_update_elem(&mul_mat_vec_q_entry, &pid, &data, BPF_ANY);
    return 0;
}

SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_ggml_cuda_op_mul_mat_vec_q)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;

    // 查找进入时的时间戳
    struct entry_data *entry_p = bpf_map_lookup_elem(&mul_mat_vec_q_entry, &pid);
    if (!entry_p) {
        return 0; // 没有找到对应的进入记录，可能是在 BPF 程序启动前进入的
    }

    uint64_t end_ts = bpf_ktime_get_ns();
    uint64_t start_ts = entry_p->ts;
    uint64_t duration_ns = end_ts - start_ts;

    // 从 map 中删除记录
    bpf_map_delete_elem(&mul_mat_vec_q_entry, &pid);

    // 准备并发送事件到用户空间
    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        return 0; // Ring buffer 空间不足
    }

    e->type = EVENT_TYPE_FUNC_DURATION;
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    // 注意：这里直接硬编码函数名，因为 BPF 程序知道它在哪个探针里
    // 更复杂的场景可能需要传递函数标识符
    __builtin_memcpy(e->func_duration.func_name, "ggml_cuda_op_mul_mat_vec_q", sizeof("ggml_cuda_op_mul_mat_vec_q"));
    e->func_duration.duration_ns = duration_ns;

    bpf_ringbuf_submit(e, 0);
    return 0;
}

// --- ggml_cuda_op_mul_mat_q ---
// _Z22ggml_cuda_op_mul_mat_qR25ggml_backend_cuda_contextPK11ggml_tensorS3_PS1_PKcPKfS6_PfllllP11CUstream_st
SEC("uprobe")
int BPF_KPROBE(uprobe_ggml_cuda_op_mul_mat_q)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    if (!process_allowed(pid, comm))
        return 0;

    struct entry_data data = {};
    data.ts = bpf_ktime_get_ns();

    bpf_map_update_elem(&mul_mat_q_entry, &pid, &data, BPF_ANY);
    return 0;
}

SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_ggml_cuda_op_mul_mat_q)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;

    struct entry_data *entry_p = bpf_map_lookup_elem(&mul_mat_q_entry, &pid);
    if (!entry_p) {
        return 0;
    }

    uint64_t end_ts = bpf_ktime_get_ns();
    uint64_t start_ts = entry_p->ts;
    uint64_t duration_ns = end_ts - start_ts;

    bpf_map_delete_elem(&mul_mat_q_entry, &pid);

    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    e->type = EVENT_TYPE_FUNC_DURATION;
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    __builtin_memcpy(e->func_duration.func_name, "ggml_cuda_op_mul_mat_q", sizeof("ggml_cuda_op_mul_mat_q"));
    e->func_duration.duration_ns = duration_ns;

    bpf_ringbuf_submit(e, 0);
    return 0;
}

// --- ggml_cuda_set_device ---
// _Z20ggml_cuda_set_devicei
// SEC("uprobe")
// int BPF_KPROBE(uprobe_ggml_cuda_set_device, int device_id) // 获取第一个参数
// {
//     pid_t pid = bpf_get_current_pid_tgid() >> 32;
//     char comm[TASK_COMM_LEN];
//     bpf_get_current_comm(&comm, sizeof(comm));

//     if (!process_allowed(pid, comm))
//         return 0;

//     struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
//     if (!e) {
//         return 0;
//     }

//     e->type = EVENT_TYPE_SET_DEVICE;
//     e->pid = pid;
//     __builtin_memcpy(e->comm, comm, sizeof(e->comm));
//     e->set_device.device_id = device_id;

//     bpf_ringbuf_submit(e, 0);
//     return 0;
// }