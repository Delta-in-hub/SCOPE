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

#include "sched.h"
#include "sched.skel.h"

static struct env {
    pid_t pid;
    char parent_comm[TASK_COMM_LEN];
    bool verbose;
} env = {
    .pid = 0,
    .parent_comm = "",
    .verbose = false,
};

const char *argp_program_version = "sched 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "\n"
    "USAGE: ./sched [-p PID] [-c PARENT_COMM] [-v]\n";

static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by PID calling execve"},
    {"parent-comm", 'c', "PARENT_COMMAND", 0,
     "Filter by parent process command name"},
    {"verbose", 'v', NULL, 0, "Verbose debug output"},
    {},
};

static error_t parse_arg(int key, char *arg, struct argp_state *state) {
    long pid_in;
    switch (key) {
    case 'p':
        errno = 0;
        pid_in = strtol(arg, NULL, 10);
        if (errno || pid_in <= 0) {
            fprintf(stderr, "Invalid PID: %s\n", arg);
            argp_usage(state);
        }
        env.pid = (pid_t)pid_in;
        break;
    case 'c':
        if (strlen(arg) >= TASK_COMM_LEN) {
            fprintf(stderr, "Parent command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.parent_comm, arg, TASK_COMM_LEN);
        env.parent_comm[TASK_COMM_LEN - 1] = '\0';
        break;
    case 'v':
        env.verbose = true;
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

static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}

static volatile bool exiting = false;

static void sig_handler(int sig) { exiting = true; }

static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)
        return -1;

    const struct event *e = data;
    char ts[32];
    time_t t = time(NULL);
    strftime(ts, sizeof(ts), "%H:%M:%S", localtime(&t));

    if (e->type == SWITCH_IN) {
        printf("%-8s %-7d %-7d %-20s %s CPU(%d)\n", ts, e->cpu, e->pid, e->comm, "Sched IN", e->cpu);
    } else {
        printf("%-8s %-7d %-7d %-20s %s CPU(%d)\n", ts, e->cpu, e->pid, e->comm, "Sched OUT", e->cpu);
    }

    return 0;
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct sched_bpf *skel = NULL;
    int err;

    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    skel = sched_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }

    skel->rodata->filter_pid = env.pid;
    memcpy((char *)skel->rodata->filter_comm, env.parent_comm, TASK_COMM_LEN);

    err = sched_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    err = sched_bpf__attach(skel);
    if (err) {
        fprintf(stderr, "Failed to attach BPF: %s\n", strerror(-err));
        goto cleanup;
    }

    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, NULL, NULL);
    if (!rb) {
        err = -errno;
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    while (!exiting) {
        err = ring_buffer__poll(rb, 100);
        if (err == -EINTR) {
            err = 0;
            break;
        }
        if (err < 0)
            break;
    }

cleanup:
    ring_buffer__free(rb);
    sched_bpf__destroy(skel);
    return err ? -err : 0;
}
