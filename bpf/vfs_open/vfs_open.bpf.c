#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "vfs_open.h"

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
    if (comm[0] == 'k' && comm[1] == 's' && comm[2] == 'o' &&
        comm[3] == 'f') // ksoft
        return 1;
    if (comm[0] == 'k' && comm[1] == 'w' && comm[2] == 'o' &&
        comm[3] == 'r') // kworker
        return 1;
    if (comm[0] == 'k' && comm[1] == 's' && comm[2] == 'w' &&
        comm[3] == 'a') // kswapd
        return 1;
    if (comm[0] == 'w' && comm[1] == 'a' && comm[2] == 't' && comm[3] == 'c' &&
        comm[4] == 'h') // watchdog
        return 1;
    if (comm[0] == 'm' && comm[1] == 'i' && comm[2] == 'g' &&
        comm[3] == 'r') // migration
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

SEC("fentry/do_filp_open")
int BPF_PROG(handle_do_filp_open, int dfd, struct filename *pathname,
         const struct open_flags *op) {
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    if (!process_allowed(pid, comm)) {
        return 0;
    }

    // comm ignore List: Xwayland , kwin_* , 

    if (bpf_strncmp(comm, 7, "Xwayland") == 0) {
        return 0;
    }
    if (bpf_strncmp(comm, 5, "kwin_") == 0) {
        return 0;
    }


    // 3. 尝试从 pathname 参数读取文件名
    // struct filename 包含一个指向实际路径字符串的指针 'name'
    const char *fname_ptr;
    // 使用 bpf_core_read 安全地读取 pathname->name 指针
    // 注意: pathname 可能为 NULL，需要检查
    if (pathname == NULL) {
        bpf_printk("WARN: vfs_open: pathname is NULL\n");
        return 0;
    }
    bpf_core_read(&fname_ptr, sizeof(fname_ptr), &pathname->name);

    if (fname_ptr == NULL) {
        bpf_printk(
            "WARN: vfs_open: filename pointer is NULL in struct filename\n");
        return 0; // 获取文件名指针失败
    }

    // 4. 预订 Ring Buffer 空间
    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        bpf_printk("WARN: vfs_open: ringbuf reserve failed\n");
        return 0; // 无法分配空间
    }

    // 5. 填充事件数据
    e->pid = pid;
    __builtin_memcpy(e->comm, comm, sizeof(e->comm));
    e->comm[TASK_COMM_LEN - 1] = '\0'; // 确保 null 结尾

    // 使用 bpf_probe_read_kernel_str 从内核空间读取文件名字符串
    long ret =
        bpf_probe_read_kernel_str(&e->filename, sizeof(e->filename), fname_ptr);
    if (ret < 0) {
        bpf_printk("WARN: vfs_open: failed to read filename string: %ld\n",
                   ret);
        // 即使读取失败，也可能需要发送事件，只是文件名为空
        e->filename[0] = '\0'; // 清空文件名
        // 或者可以选择丢弃事件：
        // bpf_ringbuf_discard(e, 0);
        // return 0;
    } else {
        // 确保 null 结尾 (bpf_probe_read_kernel_str
        // 在成功时会保证，但多一层保险)
        e->filename[MAX_PATH_LEN - 1] = '\0';
    }

    //  忽略 /proc*
    if (bpf_strncmp(e->filename, 5, "/proc") == 0) {
        bpf_ringbuf_discard(e, 0);
        return 0;
    }

    // 6. 提交事件到 Ring Buffer
    bpf_ringbuf_submit(e, 0);

    return 0; // 成功处理
}