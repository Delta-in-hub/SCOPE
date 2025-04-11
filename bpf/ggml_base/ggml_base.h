#ifndef __GGML_BASE_H
#define __GGML_BASE_H

#define TASK_COMM_LEN 16
#define MAX_ENTRIES   10240 // Map size, adjust if needed

// 定义事件类型
enum event_type {
    EVENT_MALLOC = 0,
    EVENT_FREE,
};

// 定义通过 Ring Buffer 传递的事件结构体
struct event {
    enum event_type type;
    int pid;
    char comm[TASK_COMM_LEN];
    size_t size;            // 内存大小
    unsigned long long ptr; // 内存指针地址 (使用 ull 保证足够大小)
};

#endif /* __GGML_BASE_H */