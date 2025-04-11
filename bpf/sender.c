// main_packed_topic.c (使用 zmq_msgpack_pub.h 的示例，主题被编码)
#include "zmqsender.h" // 包含修改后的头文件库

#include <assert.h>
#include <time.h>   // For srand, time
#include <unistd.h> // For sleep/usleep

// --- 用户定义的数据结构 ---
typedef struct {
    int32_t command_id;
    char target_device[32];
    double parameter;
} CommandPayload;

typedef struct {
    int32_t source_id;
    char status_code[16];
    char details[128];
} StatusUpdatePayload;

// --- 用户定义的 MessagePack *负载* 打包函数 ---
// 这些函数保持不变，只负责打包负载数据

// CommandPayload 的打包函数
void pack_command_callback(msgpack_packer *pk, const void *data) {
    const CommandPayload *payload = (const CommandPayload *)data;
    assert(pk && payload); // 基本检查

    msgpack_pack_array(pk, 3);
    msgpack_pack_int32(pk, payload->command_id);

    size_t device_len =
        strnlen(payload->target_device, sizeof(payload->target_device));
    msgpack_pack_str(pk, device_len);
    msgpack_pack_str_body(pk, payload->target_device, device_len);

    msgpack_pack_double(pk, payload->parameter);
}

// StatusUpdatePayload 的打包函数
void pack_status_callback(msgpack_packer *pk, const void *data) {
    const StatusUpdatePayload *payload = (const StatusUpdatePayload *)data;
    assert(pk && payload); // 基本检查

    msgpack_pack_array(pk, 3);
    msgpack_pack_int32(pk, payload->source_id);

    size_t code_len =
        strnlen(payload->status_code, sizeof(payload->status_code));
    msgpack_pack_str(pk, code_len);
    msgpack_pack_str_body(pk, payload->status_code, code_len);

    size_t details_len = strnlen(payload->details, sizeof(payload->details));
    msgpack_pack_str(pk, details_len);
    msgpack_pack_str_body(pk, payload->details, details_len);
}

// --- 常量 ---
#define MSG_TYPE_COMMAND "CMD" // 原始主题字符串
#define MSG_TYPE_STATUS "STAT" // 原始主题字符串
#define IPC_ENDPOINT                                                           \
    "ipc:///tmp/zmq_ipc_pubsub_lib_packed.sock" // 使用不同端点区分

int main() {
    // 1. 初始化 ZMQ 发布者 (使用修改后的库)
    zmq_pub_handle_t *publisher = zmq_pub_init(IPC_ENDPOINT);
    if (!publisher) {
        fprintf(stderr, "Failed to initialize publisher. Exiting.\n");
        return 1;
    }

    // 可选：给订阅者启动和连接的时间
    printf("Publisher initialized (with packed topic). Waiting 1 second for "
           "subscribers...\n");
    sleep(1);

    // 随机数种子
    srand(time(NULL));

    printf("Starting to publish using the library (topic will be msgpack "
           "encoded)...\n");

    // 2. 循环发送不同类型的消息
    for (int i = 0; i < 10; ++i) {
        int rc = -1; // 发送结果
        const char *current_topic = NULL;

        if (i % 2 == 0) {
            // 准备 CommandPayload 数据
            current_topic = MSG_TYPE_COMMAND;
            CommandPayload cmd;
            cmd.command_id = 1000 + i;
            snprintf(cmd.target_device, sizeof(cmd.target_device), "Sensor_%d",
                     i / 2);
            cmd.target_device[sizeof(cmd.target_device) - 1] = '\0';
            cmd.parameter = (double)rand() / RAND_MAX * 10.0;

            printf("Sending raw topic [%s] (will be packed): ID=%d, "
                   "Target='%s', Param=%.2f\n",
                   current_topic, cmd.command_id, cmd.target_device,
                   cmd.parameter);

            // 使用库函数发送，传入 *原始* 主题和对应的负载打包回调
            rc = zmq_pub_send(publisher, current_topic, &cmd,
                              pack_command_callback);

        } else {
            // 准备 StatusUpdatePayload 数据
            current_topic = MSG_TYPE_STATUS;
            StatusUpdatePayload stat;
            stat.source_id = 2000 + i;
            snprintf(stat.status_code, sizeof(stat.status_code),
                     (i % 4 == 1) ? "OK" : "PENDING");
            stat.status_code[sizeof(stat.status_code) - 1] = '\0';
            snprintf(stat.details, sizeof(stat.details),
                     "Status details update seq %d", i);
            stat.details[sizeof(stat.details) - 1] = '\0';

            printf("Sending raw topic [%s] (will be packed): SrcID=%d, "
                   "Code='%s', Details='%.30s...'\n",
                   current_topic, stat.source_id, stat.status_code,
                   stat.details);

            // 使用库函数发送，传入 *原始* 主题和对应的负载打包回调
            rc = zmq_pub_send(publisher, current_topic, &stat,
                              pack_status_callback);
        }

        // 检查发送结果
        if (rc != 0) {
            fprintf(stderr, "Failed to send message %d. Exiting loop.\n", i);
            break;
        }

        usleep(500000); // 暂停 500 毫秒
    }

    printf("Finished publishing data loop.\n");

    // 3. 清理资源
    //    注意传递句柄指针的地址，以便函数内部能将其设为 NULL
    zmq_pub_cleanup(&publisher);

    // 验证句柄是否已被设为 NULL
    assert(publisher == NULL);

    printf("Publisher finished cleanly.\n");
    return 0;
}