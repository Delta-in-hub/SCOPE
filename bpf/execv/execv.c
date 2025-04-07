// execv.c
// SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause)
#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>

#include "execv.h" // Shared header
#include "execv.skel.h" // Generated BPF skeleton header (assuming bpf file is execv.bpf.c)

// Environment configuration struct to store command-line arguments
static struct env {
    pid_t pid; // PID of the process whose execve calls we trace
    char parent_comm[TASK_COMM_LEN]; // Filter by parent command name
    bool verbose;
} env = {
    .pid = 0,          // Default: don't filter by PID
    .parent_comm = "", // Default: don't filter by parent command name
    .verbose = false,
};

const char *argp_program_version = "execv_tracer 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "Trace execve syscalls using BPF.\n"
    "\n"
    "Filters events based on the PID initiating execve and/or the command "
    "name\n"
    "of the parent process. Prints PID, PPID, filename, and arguments.\n"
    "\n"
    "USAGE: ./execv [-p PID] [-c PARENT_COMM] [-v]\n";

static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by PID calling execve"},
    {"parent-comm", 'c', "PARENT_COMMAND", 0,
     "Filter by parent process command name (exact match)"},
    {"verbose", 'v', NULL, 0, "Verbose debug output"},
    {},
};

static error_t parse_arg(int key, char *arg, struct argp_state *state) {
    long long pid_in;
    switch (key) {
    case 'p':
        errno = 0;
        pid_in = strtoll(arg, NULL, 10);
        if (errno || pid_in <= 0) { // execve is called by existing PIDs
            fprintf(stderr, "Invalid PID: %s\n", arg);
            argp_usage(state);
        }
        env.pid = (pid_t)pid_in;
        break;
    case 'c': // Filter by parent command name
        if (strlen(arg) >= TASK_COMM_LEN) {
            fprintf(stderr, "Parent command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.parent_comm, arg, TASK_COMM_LEN);
        env.parent_comm[TASK_COMM_LEN - 1] = '\0'; // Ensure null termination
        break;
    case 'v':
        env.verbose = true;
        break;
    case ARGP_KEY_ARG:
        argp_usage(state);
        break;
    default:
        return ARGP_ERR_UNKNOWN;
    }
    return 0;
}

static const struct argp argp = {
    .options = opts,
    .parser = parse_arg,
    .doc = argp_program_doc,
};

// libbpf print callback function
static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}

// Signal handling flag
static volatile bool exiting = false;

// Signal handler
static void sig_handler(int sig) { exiting = true; }

// Ring Buffer event handler callback
static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)   // Check exit flag first
        return -1; // Stop processing if exiting

    const struct event *e = data;
    struct tm *tm;
    char ts[32];
    time_t t;

    // Get current timestamp
    time(&t);
    tm = localtime(&t);
    strftime(ts, sizeof(ts), "%H:%M:%S", tm);

    // Print event information
    printf("%-8s %-7d %-7d %-20s", ts, e->pid, e->ppid, e->filename);

    for (int i = 0; i < MAX_ARGS_TO_READ; i++) {
        const char *p = e->args + i * 8;
        if (*p == '\0')
            continue;
        printf(" %s", p);
    }
    printf("\n");
    return 0; // Success
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct execv_bpf *skel =
        NULL; // BPF skeleton pointer (adjust name if needed)
    int err;

    // Parse command line arguments
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    // Set up libbpf's print callback function
    libbpf_set_print(libbpf_print_fn);

    // Set up signal handler for graceful exit
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    // 1. Open BPF skeleton file
    skel =
        execv_bpf__open(); // Make sure execv_bpf matches your .bpf.c filename
    if (!skel) {
        fprintf(stderr, "Error: Failed to open BPF skeleton\n");
        return 1;
    }

    // 2. Set BPF program parameters (.rodata section) before loading
    skel->rodata->filter_pid = env.pid;
    // Copy the parent command filter string into the BPF program's data section
    memcpy((char *)skel->rodata->filter_comm, env.parent_comm, TASK_COMM_LEN);

    // 3. Load and verify BPF programs
    err = execv_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Error: Failed to load BPF skeleton: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // 4. Attach BPF program to tracepoint
    err = execv_bpf__attach(skel);
    if (err) {
        fprintf(stderr, "Error: Failed to attach BPF skeleton: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // 5. Create Ring Buffer
    //    - bpf_map__fd(skel->maps.rb): Get the file descriptor for the 'rb' map
    //    - handle_event: Specify the event handler callback
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
    if (!rb) {
        err = -errno; // ring_buffer__new usually returns NULL on failure, check
                      // errno
        if (err == 0)
            err = -EINVAL; // Provide a default error code if errno isn't set
        fprintf(stderr, "Error: Failed to create ring buffer: %d (%s)\n", err,
                strerror(-err));
        goto cleanup;
    }

    // Print table header
    printf("%-8s %-7s %-7s %-20s %s\n", "TIME", "PID", "PPID", "FILENAME",
           "ARGS");

    // 6. Poll Ring Buffer for events
    while (!exiting) {
        // ring_buffer__poll() will call handle_event for received events
        // timeout set to 100 milliseconds
        err = ring_buffer__poll(rb, 100 /* timeout, ms */);

        if (err == -EINTR) {
            /* Syscall interrupted by signal (e.g., Ctrl+C) */
            fprintf(stderr, "\nInterrupted by signal (EINTR).\n");
            err = 0; // Clear error state
            break;   // Exit loop for cleanup
        }

        if (err < 0) {
            /* Handle other errors or callback requested stop */
            if (exiting && err == -1) {
                /* Expected stop requested by handle_event due to signal */
                fprintf(stderr,
                        "\nPoll stopped by callback due to exit signal.\n");
                err = 0; // Clear error state for normal exit
                break;   // Exit loop for cleanup
            } else {
                /* Other unexpected poll error */
                fprintf(stderr, "\nError polling ring buffer: %d (%s)\n", err,
                        strerror(-err));
                break; // Exit loop
            }
        }
        /* err == 0: Timeout, no new events */
        /* Continue next iteration */
    }

cleanup:
    // 7. Clean up resources
    fprintf(stderr, "\nExiting...\n");
    ring_buffer__free(rb);    // Free Ring Buffer
    execv_bpf__destroy(skel); // Destroy BPF skeleton (detaches and unloads)

    return err < 0 ? -err : 0; // Return error code
}