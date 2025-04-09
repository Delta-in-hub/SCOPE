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
    pid_t pid;
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

    //    元素 1: pid (int -> msgpack integer)
    msgpack_pack_int(pk, event->pid); // msgpack 会选择合适的整数编码

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


//! [vfs_open] END