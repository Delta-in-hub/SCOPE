#pragma once

#define TASK_COMM_LEN 16
#define TEXT_LEN 256
#define MAX_ENTRIES	10240



struct event {
    int pid;
    char comm[TASK_COMM_LEN];
    char text[TEXT_LEN];
};


struct logident {
    int pid;
    uint64_t ts;
    char * textp;
};
