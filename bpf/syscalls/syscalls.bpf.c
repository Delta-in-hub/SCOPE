
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "syscalls.h"  // 自己的头文件放到最下面

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // 1MB
} rb SEC(".maps");

static __always_inline int comm_allowed(const char *comm) {
    int i;

    for (i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
        if (comm[i] != filter_comm[i])
            return 0;
    }
    return 1;
}

SEC("tracepoint/raw_syscalls/sys_enter")
int sys_enter(struct trace_event_raw_sys_enter *args) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    if (filter_pid && pid != filter_pid)
        return 0;

    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));
    if (filter_comm[0]) {
        if (!comm_allowed(comm))
            return 0;
    }

    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    e->pid = pid;
    __builtin_memcpy(e->comm, comm, sizeof(e->comm));
    e->syscallid = args->id;
    bpf_ringbuf_submit(e, 0);
    return 0;
}