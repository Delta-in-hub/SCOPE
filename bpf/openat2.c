// SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause)
/* 基于 fentry.c 修改 */
#include <stdio.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>
#include <errno.h>
#include <sys/resource.h>
#include <bpf/libbpf.h>
#include "openat2.skel.h" // 注意：文件名要匹配，这里是 openat2.skel.h

// libbpf 打印回调函数
static int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args)
{
	// 可以根据需要增加详细程度控制 (例如，检查 env.verbose)
	return vfprintf(stderr, format, args);
}

// 信号处理，用于优雅退出
static volatile sig_atomic_t stop;

void sig_int(int signo)
{
	stop = 1;
}

int main(int argc, char **argv)
{
	struct openat2_bpf *skel; // 结构体名称基于 BPF 文件名
	int err;

	/* 设置 libbpf 的错误和调试信息回调 */
	libbpf_set_print(libbpf_print_fn);

	/* 打开、加载并验证 BPF 应用程序 */
	skel = openat2_bpf__open_and_load();
	if (!skel) {
		fprintf(stderr, "Failed to open and load BPF skeleton\n");
		return 1;
	}

	/* 附加 fentry/fexit 处理程序 */
	err = openat2_bpf__attach(skel);
	if (err) {
		fprintf(stderr, "Failed to attach BPF skeleton\n");
		goto cleanup;
	}

	/* 设置 SIGINT 信号处理程序 */
	if (signal(SIGINT, sig_int) == SIG_ERR) {
		fprintf(stderr, "can't set signal handler: %s\n", strerror(errno));
		err = 1; // 将 errno 映射到一个非零错误码
		goto cleanup;
	}

	printf("Successfully started! Please run `sudo cat /sys/kernel/debug/tracing/trace_pipe` "
	       "to see output of the BPF programs.\n");
    printf("Try opening files (e.g., `ls /tmp`, `cat some_file`) to trigger the probes.\n");

	/* 主循环，等待退出信号 */
	while (!stop) {
		// 打印点以表示程序仍在运行
		fprintf(stderr, ".");
		fflush(stderr); // 确保点被立即打印出来
		sleep(1);
	}

	printf("\nExiting...\n");

cleanup:
	/* 清理资源 */
	openat2_bpf__destroy(skel);
	// 返回负的错误码（如果发生错误），否则返回 0
	return err < 0 ? -err : err;
}