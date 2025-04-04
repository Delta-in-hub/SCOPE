#pragma once

#ifndef _POSIX_C_SOURCE
#define _POSIX_C_SOURCE 199309L // 或者更高版本，以确保 clock_gettime 可用
#endif
#include <stdint.h>
#include <stdio.h>
#include <time.h> // 包含 clock_gettime() 和 struct timespec

static inline int64_t UnixMicroNow() {
    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) == -1) {
        perror("clock_gettime failed");
        return -1;
    }
    return (int64_t)ts.tv_sec * 1000000 + (int64_t)ts.tv_nsec / 1000;
}

static inline int64_t UnixNanoNow() {
    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) == -1) {
        perror("clock_gettime failed");
        return -1;
    }
    return (int64_t)ts.tv_sec * 1000000000 + (int64_t)ts.tv_nsec;
}

static inline int64_t UnixMilliNow() {
    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) == -1) {
        perror("clock_gettime failed");
        return -1;
    }
    return (int64_t)ts.tv_sec * 1000 + (int64_t)ts.tv_nsec / 1000000;
}
