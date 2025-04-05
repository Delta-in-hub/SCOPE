#pragma once

#ifndef _GNU_SOURCE
#define _GNU_SOURCE
#endif

#include <msgpack.h>
#include <stdint.h>
#include <string.h>

struct IPC_Model {
    int64_t nano_since_epoch; // timestamp
    int32_t pid;              // process id
    char comm[16];            // process name
    char *cmdline;            // process command line
    char *msg;
    // 对应 struct 的 msgpack 函数, 名为 pack
    void (*pack)(msgpack_packer *pk, const void *payload);
};

extern struct IPC_Model ipc_model;