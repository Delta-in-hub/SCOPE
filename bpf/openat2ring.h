
#include <stdbool.h> // 用于 bool 类型
#define MAX_FILENAME_LEN 128
#define TASK_COMM_LEN 16 // 内核任务 comm 长度

// 定义要通过 ringbuf 发送到用户空间的数据结构
struct event {
    int pid;
    long ret;                        // fexit 的返回值
    bool is_exit;                    // 标记是 fentry 还是 fexit 事件
    char comm[TASK_COMM_LEN];        // 进程名
    char filename[MAX_FILENAME_LEN]; // 文件名
};