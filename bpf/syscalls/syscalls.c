// syscalls.c
// SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause)
#include "syscalls.h" // Shared header
#include "syscall_helper.h"
#include "syscalls.skel.h" // Generated BPF skeleton header
#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>

// 环境配置结构体，用于存储命令行参数
static struct env {
    pid_t pid;
    char comm[TASK_COMM_LEN];
    bool verbose;
} env = {
    .pid = 0,   // 默认不过滤 PID
    .comm = "", // 默认不过滤进程名
    .verbose = false,
};

const char *argp_program_version = "syscalls 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "Trace syscall entries using BPF.\n"
    "\n"
    "Filters events based on PID and/or command name and prints details.\n"
    "\n"
    "USAGE: ./syscalls [-p PID] [-c COMM] [-v]\n";

static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process ID (TGID)"},
    {"comm", 'c', "COMMAND", 0, "Filter by command name (exact match)"},
    {"verbose", 'v', NULL, 0, "Verbose debug output"},
    {},
};

static error_t parse_arg(int key, char *arg, struct argp_state *state) {
    long long pid_in;
    switch (key) {
    case 'p':
        errno = 0;
        pid_in = strtoll(arg, NULL, 10);
        if (errno || pid_in <= 0) {
            fprintf(stderr, "Invalid PID: %s\n", arg);
            argp_usage(state);
        }
        env.pid = (pid_t)pid_in;
        break;
    case 'c':
        if (strlen(arg) >= TASK_COMM_LEN) {
            fprintf(stderr, "Command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.comm, arg, TASK_COMM_LEN);
        env.comm[TASK_COMM_LEN - 1] = '\0'; // 确保 null 结尾
        break;
    case 'v':
        env.verbose = true;
        break;
    case ARGP_KEY_ARG:
        argp_usage(state);
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static const struct argp argp = {
    .options = opts,
    .parser = parse_arg,
    .doc = argp_program_doc,
};

// libbpf 打印回调函数
static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}

// 信号处理标志
static volatile bool exiting = false;

// 信号处理函数
static void sig_handler(int sig) { exiting = true; }

// Ring Buffer 事件处理回调函数
static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)
        return -1;

    const struct event *e = data;
    struct tm *tm;
    char ts[32];
    time_t t;

    // 获取当前时间戳
    time(&t);
    tm = localtime(&t);
    strftime(ts, sizeof(ts), "%H:%M:%S", tm);
    char name[32];
    syscall_name(e->syscallid, name, sizeof(name));
    // 打印事件信息
    printf("%-8s %-16s %-7d %-5s\n", ts, e->comm, e->pid, name);

    return 0;
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct syscalls_bpf *skel = NULL; // BPF 骨架指针
    int err;

    // 解析命令行参数
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    // 设置 libbpf 的打印回调函数
    libbpf_set_print(libbpf_print_fn);

    // 设置信号处理函数，用于优雅退出
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    init_syscall_names();

    // 1. 打开 BPF 骨架文件
    skel = syscalls_bpf__open();
    if (!skel) {
        fprintf(stderr, "Error: Failed to open BPF skeleton\n");
        return 1;
    }

    // 2. 设置 BPF 程序参数 (.rodata section)
    skel->rodata->filter_pid = env.pid;
    memcpy((char *)skel->rodata->filter_comm, env.comm, TASK_COMM_LEN);

    // 3. 加载并验证 BPF 程序
    err = syscalls_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Error: Failed to load BPF skeleton: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // 4. 附加 BPF 程序到跟踪点
    err = syscalls_bpf__attach(skel);
    if (err) {
        fprintf(stderr, "Error: Failed to attach BPF skeleton: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // 5. 创建 Ring Buffer
    //    - bpf_map__fd(skel->maps.rb): 获取 'rb' map 的文件描述符
    //    - handle_event: 指定事件处理回调函数
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
    if (!rb) {
        err = -errno; // ring_buffer__new 在失败时通常不设置 errno，但返回 NULL
        if (err == 0)
            err = -EINVAL; // 提供一个默认错误码
        fprintf(stderr, "Error: Failed to create ring buffer: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // 打印表头
    printf("%-8s %-16s %-7s %-5s\n", "TIME", "COMM", "PID", "SYSCALL_ID");

    // 6. 轮询 Ring Buffer 获取事件
    while (!exiting) {
        // ring_buffer__poll() 会调用 handle_event 处理收到的事件
        // timeout 设置为 100 毫秒
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);
        /* Ctrl-C 会导致 poll 返回 -EINTR */
        if (err == -EINTR) {
            err = 0;
            break; // 收到中断信号，跳出循环
        }
        if (err < 0) {
            // 发生其他错误
            fprintf(stderr, "Error polling ring buffer: %d (%s)\n", err,
                    strerror(-err));
            break;
        }
        // err == 0 表示超时，没有新事件，继续循环
    }

cleanup:
    // 7. 清理资源
    fprintf(stderr, "Exiting...\n");
    ring_buffer__free(rb);       // 释放 Ring Buffer
    syscalls_bpf__destroy(skel); // 销毁 BPF 骨架 (会自动分离程序并卸载)
    free_syscall_names();
    return err < 0 ? -err : 0; // 返回错误码
}