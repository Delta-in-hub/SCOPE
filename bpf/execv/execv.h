#pragma once

#define MAX_ARGS_TO_READ 8
struct event {
	int pid;
	int ppid;
	char filename[64];
	char args[MAX_ARGS_TO_READ*16];
};

#define MAX_TOTAL_ARGS_LEN sizeof(((struct event *)0)->args)
#define TASK_COMM_LEN 16