#pragma once // 在这个头文件中, 不要 include 任何头文件, 避免冲突. 类型尽量用基本类型, 或者 自己 定义.

#define TASK_COMM_LEN 16
#define MAX_FUNC_NAME_LEN 32 // 足够存储函数名的长度

// 事件类型枚举
enum event_type {
    EVENT_TYPE_FUNC_DURATION,
    EVENT_TYPE_SET_DEVICE,
};

// 事件结构定义
struct event {
    enum event_type type;      // 事件类型
    int pid;                   // 进程 ID
    char comm[TASK_COMM_LEN];  // 进程名

    union {
        // 用于 EVENT_TYPE_FUNC_DURATION
        struct {
            char func_name[MAX_FUNC_NAME_LEN]; // 函数名
            uint64_t duration_ns;             // 函数执行耗时 (纳秒)
        } func_duration;

        // 用于 EVENT_TYPE_SET_DEVICE
        struct {
            int device_id;                     // 设置的设备 ID
        } set_device;
    };
};

