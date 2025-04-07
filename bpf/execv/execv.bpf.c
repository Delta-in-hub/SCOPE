
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "execv.h" // 自己的头文件放到最下面

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // 1MB
} rb SEC(".maps");

static __always_inline int comm_allowed(const char *comm) {
#pragma unroll // Hint to unroll the loop if possible
    for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
        if (comm[i] != filter_comm[i])
            return 0;
    }
    return 1;
}

struct syscall_trace_enter {
    struct trace_entry ent;
    int nr;
    long unsigned int args[0];
};

SEC("tracepoint/syscalls/sys_enter_execve")
int tracepoint__syscalls__sys_enter_execve(struct syscall_trace_enter *ctx) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    if (filter_pid && pid != filter_pid)
        return 0;

    struct task_struct *task = (struct task_struct *)bpf_get_current_task();
    struct task_struct *parent_task =
        (struct task_struct *)BPF_CORE_READ(task, real_parent);
    pid_t ppid = (pid_t)BPF_CORE_READ(parent_task, tgid);

    char comm[TASK_COMM_LEN];
    // get ppid's comm
    BPF_CORE_READ_STR_INTO(&comm, parent_task, comm);
    if (filter_comm[0]) {
        if (!comm_allowed(comm))
            return 0;
    }

    struct event *event = bpf_ringbuf_reserve(&rb, sizeof(*event), 0);
    if (!event)
        return 0;

    // --- Populate Basic Event Data ---
    event->pid = pid;
    event->ppid = ppid;
    __builtin_memset(event->filename, 0,
                     sizeof(event->filename)); // Clear buffers
    __builtin_memset(event->args, 0, sizeof(event->args));

    // --- Read Filename ---
    // Syscall arguments are in ctx->args array. args[0] is filename user
    // pointer.
    const char *filename_ptr = (const char *)ctx->args[0];
    bpf_probe_read_user_str(&event->filename, sizeof(event->filename),
                            filename_ptr);

    // --- Read Arguments ---

    const char *const *argv = (const char *const *)ctx->args[1];
    for (int i = 0; i < MAX_ARGS_TO_READ; i++) {
        const char *arg_ptr = NULL;
        if (bpf_probe_read_user(&arg_ptr, sizeof(arg_ptr), &argv[i]) ||
            !arg_ptr)
            break;

        int len = bpf_probe_read_user_str(event->args + i * 16, 16, arg_ptr);
        if (len <= 0)
            continue;
    }

    // --- Submit Event ---
    bpf_ringbuf_submit(event, 0);
    return 0;
}