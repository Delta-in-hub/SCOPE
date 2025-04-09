#ifndef ZMQ_MSGPACK_PUB_H
#define ZMQ_MSGPACK_PUB_H

#define _GNU_SOURCE // Needed for strnlen if not implicitly defined
#include <errno.h>
#include <msgpack.h> // MessagePack C library
#include <stdint.h>  // For standard integer types
#include <stdio.h>   // For standard I/O (printf, perror, fprintf, stderr)
#include <stdlib.h>  // For general utilities (malloc, free)
#include <string.h>  // For string operations (strnlen, memcpy)
#include <sys/stat.h>
#include <unistd.h>
#include <zmq.h> // ZeroMQ library
// --- 类型定义 ---

// 用于打包用户特定数据的函数指针类型
typedef void (*zmq_packer_func_t)(msgpack_packer *pk, const void *user_data);

// ZMQ 发布者句柄结构体
typedef struct {
    void *context; // ZMQ 上下文
    void *socket;  // ZMQ PUB 套接字

    // --- 可重用缓冲区和打包器 ---
    msgpack_sbuffer payload_sbuf; // 用于 *负载* 数据的缓冲区
    msgpack_packer payload_pk;    // 用于 *负载* 数据的打包器

    msgpack_sbuffer topic_sbuf; // 用于 *主题* 数据的缓冲区 (优化点)
    msgpack_packer topic_pk;    // 用于 *主题* 数据的打包器 (优化点)
    // ---------------------------

    char endpoint[256]; // 存储绑定的端点地址
} zmq_pub_handle_t;

// --- 辅助函数 (内部使用，设为 static inline) ---

/**
 * @brief 发送多部分消息 (MsgPack编码的Topic + MsgPack编码的Payload) - 优化版
 *
 * 使用句柄中预先分配和初始化的缓冲区来编码主题，避免重复分配。
 *
 * @warning 仍然会破坏标准的 ZMQ 基于原始字节的主题过滤。
 *
 * @param handle ZMQ 发布者句柄，包含所有必要的缓冲区和打包器
 * @param topic_str 要编码并作为主题发送的原始字符串
 * @return 0 表示成功, -1 表示错误
 */
static inline int
zmq_internal_send_multipart_packed_topic(zmq_pub_handle_t *handle,
                                         const char *topic_str) {
    int overall_rc = -1;

    // --- 1. 清除并编码主题到句柄的主题缓冲区 ---
    msgpack_sbuffer_clear(&handle->topic_sbuf); // <--- 重用并清除主题缓冲区

    // 使用句柄中预置的主题打包器和缓冲区
    size_t topic_len_raw =
        strnlen(topic_str, 128); // 获取原始长度, 限制最大长度
    msgpack_pack_str(&handle->topic_pk, topic_len_raw); // 打包字符串头
    msgpack_pack_str_body(&handle->topic_pk, topic_str,
                          topic_len_raw); // 打包字符串体

    // 发送编码后的主题帧 (数据在 handle->topic_sbuf 中)
    zmq_msg_t topic_msg;
    if (zmq_msg_init_size(&topic_msg, handle->topic_sbuf.size) != 0) {
        fprintf(stderr, "ERROR: zmq_msg_init_size (packed topic) failed: %s\n",
                zmq_strerror(errno));
        return -1; // 无需 goto，因为没有局部资源需要清理
    }
    memcpy(zmq_msg_data(&topic_msg), handle->topic_sbuf.data,
           handle->topic_sbuf.size);

    int rc = zmq_msg_send(&topic_msg, handle->socket, ZMQ_SNDMORE);
    zmq_msg_close(&topic_msg);

    if (rc == -1) {
        fprintf(stderr, "ERROR: zmq_msg_send (packed topic) failed: %s\n",
                zmq_strerror(errno));
        return -1;
    }

    // --- 2. 发送负载帧 (数据已在 handle->payload_sbuf 中准备好) ---
    zmq_msg_t payload_msg;
    if (zmq_msg_init_size(&payload_msg, handle->payload_sbuf.size) != 0) {
        fprintf(stderr, "ERROR: zmq_msg_init_size (payload) failed: %s\n",
                zmq_strerror(errno));
        // 注意：编码后的主题帧可能已发送!
        return -1;
    }
    memcpy(zmq_msg_data(&payload_msg), handle->payload_sbuf.data,
           handle->payload_sbuf.size);

    rc = zmq_msg_send(&payload_msg, handle->socket, 0); // 最后一个部分
    zmq_msg_close(&payload_msg);

    if (rc == -1) {
        fprintf(stderr, "ERROR: zmq_msg_send (payload) failed: %s\n",
                zmq_strerror(errno));
        // 注意：编码后的主题帧可能已发送!
        return -1;
    }

    overall_rc = 0; // 所有步骤成功
    return overall_rc;
}

// --- 公共 API 函数 ---

/**
 * @brief 初始化 ZMQ 发布者和 MessagePack 环境。
 *
 * 创建 ZMQ 上下文、PUB 套接字，绑定端点，并初始化可重用的
 * MessagePack 缓冲区和打包器（用于 *负载* 和 *主题*）。
 *
 * @param endpoint ZMQ 绑定的端点地址。
 * @return 成功时返回句柄指针，失败时返回 NULL。
 */
static inline zmq_pub_handle_t *zmq_pub_init(const char *endpoint) {
    if (!endpoint) {
        fprintf(stderr, "ERROR: zmq_pub_init: Endpoint cannot be NULL\n");
        return NULL;
    }

    zmq_pub_handle_t *handle =
        (zmq_pub_handle_t *)malloc(sizeof(zmq_pub_handle_t));
    if (!handle) {
        perror("ERROR: zmq_pub_init: Failed to allocate handle");
        return NULL;
    }
    // 安全初始化指针和缓冲区状态
    handle->context = NULL;
    handle->socket = NULL;
    handle->payload_sbuf.data = NULL;
    handle->payload_sbuf.size = 0;
    handle->payload_sbuf.alloc = 0;
    handle->topic_sbuf.data = NULL; // 初始化新增的 topic buffer 状态
    handle->topic_sbuf.size = 0;
    handle->topic_sbuf.alloc = 0;
    memset(handle->endpoint, 0, sizeof(handle->endpoint));

    // 初始化 ZMQ 上下文
    handle->context = zmq_ctx_new();
    if (!handle->context) {
        perror("ERROR: zmq_pub_init: zmq_ctx_new failed");
        free(handle);
        return NULL;
    }

    // 创建 PUB 套接字
    handle->socket = zmq_socket(handle->context, ZMQ_PUB);
    if (!handle->socket) {
        perror("ERROR: zmq_pub_init: zmq_socket (PUB) failed");
        zmq_ctx_destroy(handle->context);
        free(handle);
        return NULL;
    }

    // 绑定套接字
    int rc = zmq_bind(handle->socket, endpoint);
    if (rc != 0) {
        fprintf(stderr, "ERROR: zmq_pub_init: zmq_bind to '%s' failed: %s\n",
                endpoint, zmq_strerror(errno));
        zmq_close(handle->socket);
        zmq_ctx_destroy(handle->context);
        free(handle);
        return NULL;
    }
    strncpy(handle->endpoint, endpoint, sizeof(handle->endpoint) - 1);
    handle->endpoint[sizeof(handle->endpoint) - 1] = '\0';

    // 初始化 MessagePack *负载* 缓冲区和打包器
    msgpack_sbuffer_init(&handle->payload_sbuf);
    msgpack_packer_init(&handle->payload_pk, &handle->payload_sbuf,
                        msgpack_sbuffer_write);

    // 初始化 MessagePack *主题* 缓冲区和打包器 (优化点)
    msgpack_sbuffer_init(&handle->topic_sbuf);
    msgpack_packer_init(&handle->topic_pk, &handle->topic_sbuf,
                        msgpack_sbuffer_write);

    printf("INFO: ZMQ Publisher initialized and bound to %s\n",
           handle->endpoint);

    // --- 新增：修改 IPC Socket 权限 ---
    if (strncmp(endpoint, "ipc://", 6) == 0) {
        const char *socket_path = endpoint + 6; // 跳过 "ipc://"
        // 修改文件权限为 0666 (所有者读写，组读写，其他人读写)
        if (chmod(socket_path, 0666) == -1) {
            perror(
                "WARN: Failed to change permissions of the IPC socket to 0666");
            // 在这种情况下，Go 程序很可能无法连接
        } else {
            printf(
                "INFO: Set IPC socket permissions to world-writable (0666)\n");
        }
    }
    return handle;
}

/**
 * @brief 使用 MessagePack 序列化负载并通过 ZMQ
 * 发布（主题也被编码，使用优化缓冲区）。
 *
 * @param handle 指向句柄。
 * @param topic 原始主题字符串，将被内部编码。
 * @param payload_data 指向要序列化的负载数据。
 * @param packer_func 用户提供的负载打包函数。
 * @return 0 成功, -1 失败。
 */
static inline int zmq_pub_send(zmq_pub_handle_t *handle, const char *topic,
                               const void *payload_data,
                               zmq_packer_func_t packer_func) {
    if (!handle || !topic || !payload_data || !packer_func) {
        fprintf(stderr, "ERROR: zmq_pub_send: Invalid arguments (handle, "
                        "topic, data, or packer function is NULL)\n");
        return -1;
    }
    if (!handle->socket || !handle->payload_pk.data || !handle->topic_pk.data) {
        fprintf(stderr, "ERROR: zmq_pub_send: Invalid handle state (socket or "
                        "buffers not initialized)\n");
        return -1;
    }

    // 1. 清除可重用的 *负载* 缓冲区
    msgpack_sbuffer_clear(&handle->payload_sbuf);

    // 2. 使用用户提供的函数打包 *负载* 数据到句柄的 *负载* 缓冲区中
    packer_func(&handle->payload_pk, payload_data);

    // 检查打包后是否有负载数据 (可选警告)
    if (handle->payload_sbuf.size == 0) {
        fprintf(stderr,
                "WARNING: zmq_pub_send: Packer function produced no payload "
                "data for topic '%s'.\n",
                topic);
    }

    // 3. 发送多部分消息 (内部函数将清除并使用句柄的主题缓冲区)
    if (zmq_internal_send_multipart_packed_topic(handle, topic) != 0) {
        // 错误信息已在内部函数中打印
        return -1; // 发送失败
    }

    return 0; // 发送成功
}

/**
 * @brief 清理并释放 ZMQ 发布者相关的资源。
 *
 * 关闭套接字，销毁上下文，释放 *负载* 和 *主题* MessagePack 缓冲区，
 * 并释放句柄本身。
 *
 * @param handle_ptr 指向句柄指针的指针，函数会将其设为 NULL。
 */
static inline void zmq_pub_cleanup(zmq_pub_handle_t **handle_ptr) {
    if (handle_ptr && *handle_ptr) {
        zmq_pub_handle_t *handle = *handle_ptr;

        // 销毁 msgpack *负载* 缓冲区
        if (handle->payload_sbuf.data) {
            msgpack_sbuffer_destroy(&handle->payload_sbuf);
            handle->payload_sbuf.data = NULL;
        }

        // 销毁 msgpack *主题* 缓冲区 (优化点)
        if (handle->topic_sbuf.data) {
            msgpack_sbuffer_destroy(&handle->topic_sbuf);
            handle->topic_sbuf.data = NULL;
        }

        // 关闭 ZMQ 套接字
        if (handle->socket) {
            zmq_close(handle->socket);
            handle->socket = NULL;
        }

        // 销毁 ZMQ 上下文
        if (handle->context) {
            zmq_ctx_destroy(handle->context);
            handle->context = NULL;
        }

        // 释放句柄结构体内存
        free(handle);
        *handle_ptr = NULL;
    }
}

#endif // ZMQ_MSGPACK_PUB_H