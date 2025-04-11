// vfsopen/vfsopen.h
#pragma once

#define TASK_COMM_LEN 16
// 定义一个合理的最大路径长度，bpf_d_path 依赖于此
// PATH_MAX (4096 in linux/limits.h) 可能是安全的，但 BPF 栈空间有限
// 可以根据需要调整，256 或 512 通常足够用于演示
#define MAX_PATH_LEN 256

struct event {
    int pid;
    char comm[TASK_COMM_LEN];
    char filename[MAX_PATH_LEN];
};
