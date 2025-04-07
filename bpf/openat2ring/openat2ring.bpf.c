
// include 的 头文件 需要严格控制, 仅限以下

#include "vmlinux.h"
#include <bpf/bpf_core_read.h> // 用于 bpf_probe_read_user_str
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
// #include <linux/bpf.h>     // 不要 include 这个头文件
#include "openat2ring.h"

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// 定义 ringbuf map
// 用户空间程序将从此 map 读取事件
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // 1MB
} rb SEC(".maps");

// --- fentry 钩子 ---
// fentry/fexit 必须要使用 BPF_PROG
// 函数签名必须和内核函数一致, 但是注意类似 __user 不要加到 bpf 的函数签名中
SEC("fentry/do_sys_openat2")
int BPF_PROG(openat2_entry, int dfd, const char *filename,
             struct open_how *how) {
    pid_t pid;
    struct event *e;

    // 预留 ringbuf 空间
    // 第三个参数 flags 通常为 0
    e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        // 无法预留空间，可能缓冲区已满
        return 0;
    }

    // 获取 PID 和进程名
    pid = bpf_get_current_pid_tgid() >> 32;
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    // 标记为入口事件
    e->is_exit = false;
    e->ret = 0; // 入口没有返回值

    // 从用户空间安全地读取文件名
    long res =
        bpf_probe_read_user_str(&e->filename, sizeof(e->filename), filename);
    if (res < 0) {
        bpf_ringbuf_discard(e, 0);
        return 0;
    }
    // 如果 res >= 0, 文件名已成功（或部分成功）读取到 e->filename

    // 提交事件到 ringbuf
    // 第三个参数 flags 通常为 0
    bpf_ringbuf_submit(e, 0);

    return 0;
}

// --- fexit 钩子 ---
SEC("fexit/do_sys_openat2")
int BPF_PROG(openat2_exit, int dfd, const char *filename, struct open_how *how,
             long ret) {
    pid_t pid;
    struct event *e;

    // 预留 ringbuf 空间
    e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    // 获取 PID 和进程名
    pid = bpf_get_current_pid_tgid() >> 32;
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    // 标记为出口事件，并记录返回值
    e->is_exit = true;
    e->ret = ret;

    // 再次从用户空间读取文件名（通常与入口相同，但保持处理逻辑）
    long res =
        bpf_probe_read_user_str(&e->filename, sizeof(e->filename), filename);
    if (res < 0) {
        bpf_ringbuf_discard(e, 0);
        return 0;
    }

    // 提交事件到 ringbuf
    bpf_ringbuf_submit(e, 0);

    return 0;
}