// SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause)
#include <stdio.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>
#include <errno.h>
#include <sys/resource.h>
#include <bpf/libbpf.h>
#include "openat2ring.skel.h"
#include <stdbool.h>

// 定义与 BPF 程序中匹配的数据结构
#define MAX_FILENAME_LEN 512
#define TASK_COMM_LEN 16
struct event {
	pid_t pid;
	long ret;
	bool is_exit;
	char comm[TASK_COMM_LEN];
	char filename[MAX_FILENAME_LEN];
};

// libbpf 打印回调函数
static int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args)
{
	// 保持不变，但可以根据需要添加详细程度
	if (level == LIBBPF_WARN) // 可以过滤掉一些 INFO 级别的信息
		return 0;
	return vfprintf(stderr, format, args);
}

// 信号处理，用于优雅退出
static volatile sig_atomic_t stop;

void sig_int(int signo)
{
	stop = 1;
}

// Ring buffer 回调处理函数
// 当 ring buffer 中有数据时，libbpf 会调用此函数
int handle_event(void *ctx, void *data, size_t data_sz)
{
	const struct event *e = data; // 将 void* 转换为我们的事件结构指针

	// 检查收到的数据大小是否符合预期
	if (data_sz != sizeof(*e)) {
		fprintf(stderr, "Error: Malformed event received (size %zu != %zu)\n",
				data_sz, sizeof(*e));
		return 1; // 返回非零表示处理出错
	}

	// 根据是入口还是出口事件，打印不同格式的信息
	if (e->is_exit) {
		printf("EXIT:  PID: %-6d COMM: %-15s FILENAME: %s RET: %ld\n",
			   e->pid, e->comm, e->filename, e->ret);
	} else {
		printf("ENTRY: PID: %-6d COMM: %-15s FILENAME: %s\n",
			   e->pid, e->comm, e->filename);
	}

	return 0; // 返回 0 表示成功处理
}

int main(int argc, char **argv)
{
	struct openat2ring_bpf *skel; // *** 更改结构体名称 ***
	struct ring_buffer *rb = NULL;  // Ring buffer 管理器
	int err;

	/* 设置 libbpf 的错误和调试信息回调 */
	libbpf_set_print(libbpf_print_fn);

	/* 打开、加载并验证 BPF 应用程序 */
	skel = openat2ring_bpf__open_and_load(); // *** 更改函数名 ***
	if (!skel) {
		fprintf(stderr, "Failed to open and load BPF skeleton\n");
		return 1;
	}

	/* 附加 fentry/fexit 处理程序 */
	err = openat2ring_bpf__attach(skel); // *** 更改函数名 ***
	if (err) {
		fprintf(stderr, "Failed to attach BPF skeleton: %s\n", strerror(-err));
		goto cleanup;
	}

	/* 设置 SIGINT 信号处理程序 */
	if (signal(SIGINT, sig_int) == SIG_ERR) {
		fprintf(stderr, "can't set signal handler: %s\n", strerror(errno));
		err = 1;
		goto cleanup;
	}

	/* 设置 Ring Buffer */
	// - bpf_map__fd(skel->maps.rb): 获取 BPF 程序中名为 'rb' 的 map 的文件描述符
	// - handle_event: 指定处理事件的回调函数
	// - NULL: context 指针，这里不需要
	// - NULL: ring_buffer_opts，使用默认选项
	rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
	if (!rb) {
		err = -errno; // ring_buffer__new 在失败时设置 errno
		fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
		goto cleanup;
	}

	printf("Successfully started! Tracing openat2 calls...\n");
	printf("Press Ctrl+C to exit.\n");
	printf("%-6s %-15s %-6s %s\n", "EVENT", "PID", "COMM", "FILENAME/RET");


	/* 主循环，轮询 ring buffer 等待事件 */
	while (!stop) {
		// ring_buffer__poll() 会检查 ring buffer 中是否有新数据
		// 如果有，它会调用 handle_event 回调函数来处理每个事件
		// 参数是超时时间（毫秒），100ms 意味着它最多阻塞 100ms 等待数据
		// 如果为 0，则立即返回；如果为负，则无限期阻塞直到有数据或出错
		err = ring_buffer__poll(rb, 100 /* timeout, ms */);
		/* Ctrl-C 会中断 poll() 调用，返回 -EINTR */
		if (err == -EINTR) {
			err = 0; // 不是真正的错误，是预期的中断
			break;   // 跳出循环准备退出
		}
		if (err < 0) {
			// 发生了 poll 错误
			fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
			break; // 退出循环
		}
		// err == 0 表示超时，没有新数据，继续循环
	}

	printf("\nExiting...\n");

cleanup:
	/* 清理资源 */
	ring_buffer__free(rb);             // 释放 ring buffer 管理器
	openat2ring_bpf__destroy(skel);    // 销毁 BPF 骨架
	return err < 0 ? -err : err;
}