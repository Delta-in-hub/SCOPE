#pragma once

#ifndef _GNU_SOURCE
#define _GNU_SOURCE
#endif

#include <msgpack.h>
#include <stdint.h>
#include <string.h>

#ifndef TASK_COMM_LEN
#define TASK_COMM_LEN 16
#endif

#ifndef MAX_FILENAME_LEN
#define MAX_FILENAME_LEN 256
#endif

//! [vfs_open] START

// --- 定义要通过 ZMQ 发送的数据结构 ---
struct vfs_open_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    char filename[MAX_FILENAME_LEN];
};

// --- MessagePack 打包函数 (Array Format) ---
// 将 struct vfs_open_event 打包成 MessagePack array 以节省空间.
// 数组元素顺序: [timestamp_ns, pid, comm, filename]
static void vfs_open_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct vfs_open_event *event =
        (const struct vfs_open_event *)user_data;

    // 1. 打包数组头，包含 4 个元素
    msgpack_pack_array(pk, 4);

    // 2. 按顺序打包各个字段的值
    //    顺序必须与消费者(Go代码)的预期一致

    //    元素 0: timestamp (int64)
    msgpack_pack_int64(pk, event->timestamp_ns);

    //    元素 1: pid
    msgpack_pack_int32(pk, event->pid);

    //    元素 2: comm (string)
    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    //    元素 3: filename (string)
    //    使用 sizeof(event->filename) 作为 strnlen 的安全界限
    size_t fname_len = strnlen(event->filename, sizeof(event->filename));
    msgpack_pack_str(pk, fname_len);
    msgpack_pack_str_body(pk, event->filename, fname_len);
}

//! [syscalls] START

struct syscalls_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    char syscall_name[32];
};

static void syscalls_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct syscalls_event *event =
        (const struct syscalls_event *)user_data;

    // 1. 打包数组头，包含 4 个元素
    msgpack_pack_array(pk, 4);

    // 2. 按顺序打包各个字段的值
    //    顺序必须与消费者(Go代码)的预期一致

    //    元素 0: timestamp (int64)
    msgpack_pack_int64(pk, event->timestamp_ns);

    //    元素 1: pid (int32)
    msgpack_pack_int32(pk, event->pid);

    //    元素 2: comm (string)
    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    //    元素 3: syscall_name (string)
    size_t syscall_name_len =
        strnlen(event->syscall_name, sizeof(event->syscall_name));
    msgpack_pack_str(pk, syscall_name_len);
    msgpack_pack_str_body(pk, event->syscall_name, syscall_name_len);
}

//! [sched]

struct sched_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    int32_t cpu;
    int32_t type; // enum event_type { SWITCH_IN, SWITCH_OUT };
};

static void sched_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct sched_event *event = (const struct sched_event *)user_data;

    // 1. 打包数组头，包含 5 个元素
    msgpack_pack_array(pk, 5);

    // 2. 按顺序打包各个字段的值
    //    顺序必须与消费者(Go代码)的预期一致

    //    元素 0: timestamp (int64)
    msgpack_pack_int64(pk, event->timestamp_ns);

    //    元素 1: pid (int32)
    msgpack_pack_int32(pk, event->pid);

    //    元素 2: comm (string)
    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    //    元素 3: cpu (int32)
    msgpack_pack_int32(pk, event->cpu);

    //    元素 4: type (int32)
    msgpack_pack_int32(pk, event->type);
}

//! [ollamabin]

#ifndef TEXT_LEN
#define TEXT_LEN 256
#endif

struct llamaLog_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    char text[TEXT_LEN];
};

static void llamaLog_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct llamaLog_event *event =
        (const struct llamaLog_event *)user_data;

    msgpack_pack_array(pk, 4);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    size_t text_len = strnlen(event->text, TEXT_LEN);
    msgpack_pack_str(pk, text_len);
    msgpack_pack_str_body(pk, event->text, text_len);
}

//! [ggml_cuda]

#ifndef MAX_FUNC_NAME_LEN
#define MAX_FUNC_NAME_LEN 32
#endif

struct ggml_cuda_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    char func_name[MAX_FUNC_NAME_LEN]; // 函数名
    int64_t duration_ns;               // 函数执行耗时 (纳秒)
};

static void ggml_cuda_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct ggml_cuda_event *event =
        (const struct ggml_cuda_event *)user_data;

    msgpack_pack_array(pk, 5);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    size_t func_name_len = strnlen(event->func_name, MAX_FUNC_NAME_LEN);
    msgpack_pack_str(pk, func_name_len);
    msgpack_pack_str_body(pk, event->func_name, func_name_len);

    msgpack_pack_int64(pk, event->duration_ns);
}

//! [ggml_cpu]

struct ggml_graph_compute_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int pid;
    char comm[TASK_COMM_LEN];

    // Fields from ggml_cgraph (collected at entry, sent at exit)
    int graph_size;
    int graph_n_nodes;
    int graph_n_leafs;
    int graph_order;
    /*
    enum ggml_cgraph_eval_order {
    GGML_CGRAPH_EVAL_ORDER_LEFT_TO_RIGHT = 0,
    GGML_CGRAPH_EVAL_ORDER_RIGHT_TO_LEFT,
    GGML_CGRAPH_EVAL_ORDER_COUNT // Should be 2
};
     */
    int64_t cost_ns;
};

static void ggml_graph_compute_event_pack(msgpack_packer *pk,
                                          const void *user_data) {
    const struct ggml_graph_compute_event *event =
        (const struct ggml_graph_compute_event *)user_data;

    msgpack_pack_array(pk, 8);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_int32(pk, event->graph_size);
    msgpack_pack_int32(pk, event->graph_n_nodes);
    msgpack_pack_int32(pk, event->graph_n_leafs);
    msgpack_pack_int32(pk, event->graph_order);
    msgpack_pack_int64(pk, event->cost_ns);
}

//! ggml_base

struct ggml_base_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    int32_t type;  // 0 : aligned_malloc, 1 : aligned_free
    uint64_t size; // 内存大小
    uint64_t ptr;  // 内存指针地址 (使用 ull 保证足够大小)
};

static void ggml_base_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct ggml_base_event *event =
        (const struct ggml_base_event *)user_data;

    msgpack_pack_array(pk, 6);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_int32(pk, event->type);
    msgpack_pack_uint64(pk, event->size);
    msgpack_pack_uint64(pk, event->ptr);
}

//! execv

struct execv_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    int32_t ppid;
    char filename[64];
    char args[128];
};

static void execv_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct execv_event *event = (const struct execv_event *)user_data;

    msgpack_pack_array(pk, 5);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);
    msgpack_pack_int32(pk, event->ppid);

    size_t filename_len = strnlen(event->filename, sizeof(event->filename));
    msgpack_pack_str(pk, filename_len);
    msgpack_pack_str_body(pk, event->filename, filename_len);

    size_t args_len = strnlen(event->args, sizeof(event->args));
    msgpack_pack_str(pk, args_len);
    msgpack_pack_str_body(pk, event->args, args_len);
}

//! [cuda]

struct cuda_malloc_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    uint64_t allocated_ptr; // 实际分配到的设备指针 (成功时)
    size_t size;            // 请求分配的大小
    int retval;             // cudaMalloc 的返回值 (错误码)
};
static void cuda_malloc_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct cuda_malloc_event *event =
        (const struct cuda_malloc_event *)user_data;

    msgpack_pack_array(pk, 6);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_uint64(pk, event->allocated_ptr);
    msgpack_pack_uint64(pk, event->size);
    msgpack_pack_int32(pk, event->retval);
}

struct cuda_free_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    uint64_t dev_ptr; // 准备释放的指针
};
static void cuda_free_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct cuda_free_event *event =
        (const struct cuda_free_event *)user_data;

    msgpack_pack_array(pk, 4);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_uint64(pk, event->dev_ptr);
}
struct cuda_launch_kernel_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    uint64_t func_ptr; // 内核函数指针 (在设备上的地址)
};
static void cuda_launch_kernel_event_pack(msgpack_packer *pk,
                                          const void *user_data) {
    const struct cuda_launch_kernel_event *event =
        (const struct cuda_launch_kernel_event *)user_data;

    msgpack_pack_array(pk, 4);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_uint64(pk, event->func_ptr);
}
struct cuda_memcpy_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    uint64_t src;
    uint64_t dst;
    size_t size;
    int kind; // 传输类型
    /*
    enum cuda_memcpy_kind {
    CUDA_MEMCPY_HOST_TO_HOST = 0,
    CUDA_MEMCPY_HOST_TO_DEVICE = 1,
    CUDA_MEMCPY_DEVICE_TO_HOST = 2,
    CUDA_MEMCPY_DEVICE_TO_DEVICE = 3,
    CUDA_MEMCPY_DEFAULT = 4,
};
     */
};
static void cuda_memcpy_event_pack(msgpack_packer *pk, const void *user_data) {

    const struct cuda_memcpy_event *event =
        (const struct cuda_memcpy_event *)user_data;

    msgpack_pack_array(pk, 7);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_uint64(pk, event->src);
    msgpack_pack_uint64(pk, event->dst);
    msgpack_pack_uint64(pk, event->size);
    msgpack_pack_int32(pk, event->kind);
}
struct cuda_sync_event {
    int64_t timestamp_ns; // 纳秒时间戳
    int32_t pid;
    char comm[TASK_COMM_LEN];
    uint64_t duration_ns; // 函数执行耗时 (纳秒)
};

static void cuda_sync_event_pack(msgpack_packer *pk, const void *user_data) {
    const struct cuda_sync_event *event =
        (const struct cuda_sync_event *)user_data;

    msgpack_pack_array(pk, 4);

    msgpack_pack_int64(pk, event->timestamp_ns);
    msgpack_pack_int32(pk, event->pid);

    size_t comm_len = strnlen(event->comm, TASK_COMM_LEN);
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, event->comm, comm_len);

    msgpack_pack_uint64(pk, event->duration_ns);
}