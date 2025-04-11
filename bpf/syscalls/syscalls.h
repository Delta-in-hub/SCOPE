#pragma once

#define TASK_COMM_LEN 16

struct event {
	int pid;
    int syscallid;
	char comm[TASK_COMM_LEN];
};