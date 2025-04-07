#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "sched.h"

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024);
} rb SEC(".maps");

static __always_inline int comm_allowed(const char *comm) {
#pragma unroll
    for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
        if (comm[i] != filter_comm[i])
            return 0;
    }
    return 1;
}

static __always_inline int is_kernel_thread(pid_t pid, const char *comm) {
    // 内核线程通常 PID <= 2 或以特定前缀开头
    if (pid <= 2)
        return 1;
    
    // 检查是否以 'k' 开头的内核线程命名模式
    if (comm[0] == 'k' && (comm[1] >= '0' && comm[1] <= '9'))
        return 1;
        
    // 检查其他常见内核线程名称
    if (comm[0] == 'k' && comm[1] == 's' && comm[2] == 'o' && comm[3] == 'f') // ksoft
        return 1;
    if (comm[0] == 'k' && comm[1] == 'w' && comm[2] == 'o' && comm[3] == 'r') // kworker
        return 1;
    if (comm[0] == 'k' && comm[1] == 's' && comm[2] == 'w' && comm[3] == 'a') // kswapd
        return 1;
    if (comm[0] == 'w' && comm[1] == 'a' && comm[2] == 't' && comm[3] == 'c' && comm[4] == 'h') // watchdog
        return 1;
    if (comm[0] == 'm' && comm[1] == 'i' && comm[2] == 'g' && comm[3] == 'r') // migration
        return 1;
    
    return 0;
}

static __always_inline int process_allowed(pid_t pid, const char *comm) {
    // 忽略内核线程和进程
    if (is_kernel_thread(pid, comm))
        return 0;

    if (filter_pid != 0 && pid != filter_pid)
        return 0;
    if (filter_comm[0] != '\0' && !comm_allowed(comm))
        return 0;
    return 1;
}

SEC("tracepoint/sched/sched_switch")
int tracepoint__sched__sched_switch(struct trace_event_raw_sched_switch *ctx) {
    pid_t prev_pid = ctx->prev_pid; // 被换出的进程 PID
    char prev_comm[TASK_COMM_LEN];  // 被换出的进程 COMM
    bpf_probe_read_kernel_str(&prev_comm, sizeof(prev_comm), ctx->prev_comm);

    pid_t next_pid = ctx->next_pid; // 被换入的进程 PID
    char next_comm[TASK_COMM_LEN];  // 被换入的进程 COMM
    bpf_probe_read_kernel_str(&next_comm, sizeof(next_comm), ctx->next_comm);

    if (process_allowed(prev_pid, prev_comm)) { // SWITCH_OUT
        struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
        if (!e) {
            bpf_printk("WARN: sched_switch: ringbuf reserve failed\n");
            return 0; // 无法分配空间，丢弃事件
        }

        // 填充事件数据
        e->type = SWITCH_OUT;
        e->cpu = bpf_get_smp_processor_id();
        e->pid = prev_pid;
        __builtin_memcpy(e->comm, prev_comm, sizeof(e->comm));

        // 提交事件到 Ring Buffer
        bpf_ringbuf_submit(e, 0);
    }

    if (process_allowed(next_pid, next_comm)) { // SWITCH_IN
        struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
        if (!e) {
            bpf_printk("WARN: sched_switch: ringbuf reserve failed\n");
            return 0; // 无法分配空间，丢弃事件
        }

        // 填充事件数据
        e->type = SWITCH_IN;
        e->cpu = bpf_get_smp_processor_id();
        e->pid = next_pid;
        __builtin_memcpy(e->comm, next_comm, sizeof(e->comm));

        // 提交事件到 Ring Buffer
        bpf_ringbuf_submit(e, 0);
    }

    return 0;
}