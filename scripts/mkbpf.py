import argparse
import os

# 模板内容定义
HEADER_TEMPLATE = """\
#pragma once

#define TASK_COMM_LEN 16

struct event {{
    int pid;

}};
"""

BPF_C_TEMPLATE = """\
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "{app_name}.h"

const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN];

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {{
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024);
}} rb SEC(".maps");

static __always_inline int comm_allowed(const char *comm) {{
    #pragma unroll
    for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\\0'; i++) {{
        if (comm[i] != filter_comm[i])
            return 0;
    }}
    return 1;
}}

static __always_inline int is_kernel_thread(pid_t pid, const char *comm) {{
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
}}

static __always_inline int process_allowed(pid_t pid, const char *comm) {{
    // 忽略内核线程和进程
    if (is_kernel_thread(pid, comm))
        return 0;

    if (filter_pid != 0 && pid != filter_pid)
        return 0;
    if (filter_comm[0] != '\\0' && !comm_allowed(comm))
        return 0;
    return 1;
}}

"""

C_TEMPLATE = """\
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

#include "{app_name}.h"
#include "{app_name}.skel.h"

static struct env {{
    pid_t pid;
    char parent_comm[TASK_COMM_LEN];
    bool verbose;
}} env = {{
    .pid = 0,
    .parent_comm = "",
    .verbose = false,
}};

const char *argp_program_version = "{app_name} 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "\\n"
    "USAGE: ./{app_name} [-p PID] [-c PARENT_COMM] [-v]\\n";

static const struct argp_option opts[] = {{
    {{"pid", 'p', "PID", 0, "Filter by PID calling execve"}},
    {{"parent-comm", 'c', "PARENT_COMMAND", 0, "Filter by parent process command name"}},
    {{"verbose", 'v', NULL, 0, "Verbose debug output"}},
    {{}},
}};

static error_t parse_arg(int key, char *arg, struct argp_state *state) {{
    long pid_in;
    switch (key) {{
    case 'p':
        errno = 0;
        pid_in = strtol(arg, NULL, 10);
        if (errno || pid_in <= 0) {{
            fprintf(stderr, "Invalid PID: %s\\n", arg);
            argp_usage(state);
        }}
        env.pid = (pid_t)pid_in;
        break;
    case 'c':
        if (strlen(arg) >= TASK_COMM_LEN) {{
            fprintf(stderr, "Parent command name too long (max %d): %s\\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }}
        strncpy(env.parent_comm, arg, TASK_COMM_LEN);
        env.parent_comm[TASK_COMM_LEN - 1] = '\\0';
        break;
    case 'v':
        env.verbose = true;
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }}
    return 0;
}}

static const struct argp argp = {{
    .options = opts,
    .parser = parse_arg,
    .doc = argp_program_doc,
}};

static int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args) {{
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}}

static volatile bool exiting = false;

static void sig_handler(int sig) {{ exiting = true; }}

static int handle_event(void *ctx, void *data, size_t data_sz) {{
    if (exiting) return -1;
    
    const struct event *e = data;
    char ts[32];
    time_t t = time(NULL);
    strftime(ts, sizeof(ts), "%H:%M:%S", localtime(&t));
    

    return 0;
}}

int main(int argc, char **argv) {{
    struct ring_buffer *rb = NULL;
    struct {app_name}_bpf *skel = NULL;
    int err;

    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err) return err;

    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    skel = {app_name}_bpf__open();
    if (!skel) {{
        fprintf(stderr, "Failed to open BPF skeleton\\n");
        return 1;
    }}

    skel->rodata->filter_pid = env.pid;
    memcpy((char *)skel->rodata->filter_comm, env.parent_comm, TASK_COMM_LEN);

    err = {app_name}_bpf__load(skel);
    if (err) {{
        fprintf(stderr, "Failed to load skeleton: %s\\n", strerror(-err));
        goto cleanup;
    }}

    err = {app_name}_bpf__attach(skel);
    if (err) {{
        fprintf(stderr, "Failed to attach BPF: %s\\n", strerror(-err));
        goto cleanup;
    }}

    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
    if (!rb) {{
        err = -errno;
        fprintf(stderr, "Failed to create ring buffer: %s\\n", strerror(-err));
        goto cleanup;
    }}

    printf("%-8s %-7s %-7s %-20s %s\\n", "TIME", "PID", "PPID", "FILENAME", "ARGS");

    while (!exiting) {{
        err = ring_buffer__poll(rb, 100);
        if (err == -EINTR) {{
            err = 0;
            break;
        }}
        if (err < 0) break;
    }}

cleanup:
    ring_buffer__free(rb);
    {app_name}_bpf__destroy(skel);
    return err ? -err : 0;
}}
"""


def create_app_structure(app_name):
    # 创建应用目录
    os.makedirs(app_name, exist_ok=False)

    # 创建三个文件并使用模板填充内容
    templates = {
        f"{app_name}.h": HEADER_TEMPLATE,
        f"{app_name}.bpf.c": BPF_C_TEMPLATE,
        f"{app_name}.c": C_TEMPLATE,
    }

    for filename, template in templates.items():
        file_path = os.path.join(app_name, filename)
        content = template.format(app_name=app_name)
        with open(file_path, "w") as f:
            f.write(content)

    print(f"成功创建应用 '{app_name}' 目录结构，包含以下文件：")
    print("\n".join(f" - {f}" for f in templates.keys()))


def main():
    parser = argparse.ArgumentParser(description="创建BPF应用目录结构")
    parser.add_argument("app_name", help="应用程序名称（将作为目录和文件名前缀）")
    args = parser.parse_args()

    create_app_structure(args.app_name)


if __name__ == "__main__":
    main()
