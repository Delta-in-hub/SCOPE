#pragma once
#include <stdio.h>
#include <string.h>

// get process pid's comm(16 bytes max)
static inline void get_comm(char *comm, int pid) {
    char path[64];
    snprintf(path, sizeof(path), "/proc/%d/comm", pid);
    FILE *fp = fopen(path, "r");
    if (fp) {
        if (fgets(comm, 16, fp)) {
            // Remove trailing newline if present
            size_t len = strlen(comm);
            if (len > 0 && comm[len - 1] == '\n') {
                comm[len - 1] = '\0';
            }
        } else {
            comm[0] = '\0';
        }
        fclose(fp);
    } else {
        comm[0] = '\0';
    }
}

// get process pid's cmdline(4096 bytes max)
static inline void get_cmdline(char *cmdline, int pid) {
    char path[64];
    snprintf(path, sizeof(path), "/proc/%d/cmdline", pid);
    FILE *fp = fopen(path, "r");
    if (fp) {
        size_t read_bytes = fread(cmdline, 1, 4095, fp);
        if (read_bytes > 0) {
            // cmdline is NUL-separated, replace NULs with spaces
            for (size_t i = 0; i < read_bytes - 1; i++) {
                if (cmdline[i] == '\0') {
                    cmdline[i] = ' ';
                }
            }
            cmdline[read_bytes] = '\0';
        } else {
            cmdline[0] = '\0';
        }
        fclose(fp);
    } else {
        cmdline[0] = '\0';
    }
}
