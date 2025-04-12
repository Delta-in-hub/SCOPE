#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "Ollamabin.h"

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024);
} rb SEC(".maps");

// struct {
// 	__uint(type, BPF_MAP_TYPE_HASH);
// 	__uint(max_entries, MAX_ENTRIES);
// 	__type(key, int);  // pid
// 	__type(value, struct logident);
// } llamaLogmap SEC(".maps");


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

/*
// extern void llamaLog(int level, char* text, void* user_data);
uprobe:/usr/bin/ollama:llamaLog
{
    $level = arg0;
    $text = arg1;
    printf("%d %s llamaLog %d %s\n", pid, comm, $level, str($text));
}
*/

SEC("uprobe")
int BPF_KPROBE(uprobe_llamaLog, int level, char* text, void* user_data)
{
    int pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));
    if (!process_allowed(pid, comm))
        return 0;
	uint64_t ts = bpf_ktime_get_ns();  // Return the time elapsed since system boot, in nanoseconds.

    // struct logident ident = {
    //     .pid = pid,
    //     .ts = ts,
    //     .textp = text
    // };
    // bpf_map_update_elem(&llamaLogmap, &pid, &ident, BPF_ANY);
    struct event *ev = bpf_ringbuf_reserve(&rb, sizeof(struct event), 0);
    if (ev == NULL) {
        return 0;
    }
    ev->pid = pid;
    bpf_get_current_comm(&ev->comm, sizeof(ev->comm));
    bpf_probe_read_user_str(ev->text, sizeof(ev->text), text);
    bpf_ringbuf_submit(ev, 0);
	return 0;
}


// SEC("uretprobe")
// int BPF_KRETPROBE(uretprobe_llamaLog, int ret)
// {
//     int pid = bpf_get_current_pid_tgid() >> 32;
//     uint64_t ts = bpf_ktime_get_ns();
//     struct logident *ident = bpf_map_lookup_elem(&llamaLogmap, &pid);
//     if (ident == NULL) {
//         return 0;
//     }
//     struct event *ev = bpf_ringbuf_reserve(&rb, sizeof(struct event), 0);
//     if (ev == NULL) {
//         return 0;
//     }
//     ev->pid = pid;
//     bpf_get_current_comm(&ev->comm, sizeof(ev->comm));
//     ev->cost_ns = ts - ident->ts;
//     bpf_probe_read_user_str(ev->text, sizeof(ev->text), ident->textp);
//     bpf_ringbuf_submit(ev, 0);
//     bpf_map_delete_elem(&llamaLogmap, &pid);
// 	return 0;
// }