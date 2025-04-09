#pragma once

#define TASK_COMM_LEN 16

// 定义事件类型
enum event_type { SWITCH_IN = 0, SWITCH_OUT };

// 定义通过 Ring Buffer 传递的事件结构体
struct event {
    enum event_type type;
    int cpu;
    int pid;
    char comm[TASK_COMM_LEN];
};