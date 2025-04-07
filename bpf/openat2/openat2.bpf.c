// SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause
/* 基于 fentry.bpf.c 修改 */
#include "vmlinux.h" // 包含内核类型定义，包括 do_sys_openat2 的签名
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h> // 虽然此示例中未使用CORE，但保留以供将来扩展

// 定义一个足够大的缓冲区来读取文件名
#define MAX_FILENAME_LEN 256

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// fentry 钩子，在 do_sys_openat2 函数开始执行时触发
SEC("fentry/do_sys_openat2")
// BPF_PROG 宏简化了 fentry/fexit 程序的定义
// 参数需要与被追踪的内核函数 do_sys_openat2 的签名匹配
// do_sys_openat2(int dfd, const char __user *filename, struct open_how *how)
int BPF_PROG(openat2_entry, int dfd, const char *filename, struct open_how *how)
{
	pid_t pid;
	char fname_buf[MAX_FILENAME_LEN];

	// 获取当前进程的 PID
	pid = bpf_get_current_pid_tgid() >> 32;

	// 从用户空间安全地读取文件名
	// bpf_probe_read_user_str 对于读取用户空间字符串是必要的
    // 注意：直接在bpf_printk中使用 %s 格式化用户空间指针是不安全且通常无效的
	long res = bpf_probe_read_user_str(fname_buf, sizeof(fname_buf), filename);
    if (res < 0) {
        // 读取失败，可能是无效指针或权限问题
        bpf_printk("fentry: openat2 called by PID %d, failed to read filename, err %ld\n", pid, res);
    } else {
        // 使用 bpf_printk 将信息打印到内核跟踪管道
	    bpf_printk("fentry: openat2 called by PID %d, filename: %s\n", pid, fname_buf);
    }

	return 0; // 必须返回 0
}

// fexit 钩子，在 do_sys_openat2 函数执行完毕返回时触发
SEC("fexit/do_sys_openat2")
// 参数与 fentry 相同，但末尾增加了内核函数的返回值 'ret'
int BPF_PROG(openat2_exit, int dfd, const char *filename, struct open_how *how, long ret)
{
	pid_t pid;
    char fname_buf[MAX_FILENAME_LEN];

	// 获取当前进程的 PID
	pid = bpf_get_current_pid_tgid() >> 32;

    // 再次读取文件名（虽然可能在入口处已读取，但为保持一致性或处理特殊情况）
	long res = bpf_probe_read_user_str(fname_buf, sizeof(fname_buf), filename);
    if (res < 0) {
        bpf_printk("fexit: openat2 called by PID %d, failed to read filename (err %ld), ret = %ld\n", pid, res, ret);
    } else {
        // 打印退出信息，包括返回值 (ret 通常是文件描述符或负的错误码)
	    bpf_printk("fexit: openat2 called by PID %d, filename: %s, ret = %ld\n", pid, fname_buf, ret);
    }
	return 0; // 必须返回 0
}