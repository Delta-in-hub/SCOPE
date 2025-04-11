#include <argp.h>
#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <errno.h>
#include <limits.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <time.h>
#include <unistd.h>

#include "ggml_cpu.h"
#include "ggml_cpu.skel.h" // Needs to be regenerated after BPF code changes

#include "../epoch.h" // 用于 UnixNanoNow()
#include "../ipc_models.h"
#include "../zmqsender.h" // 用于 ZMQ 和 MessagePack 发布

const char *ENDPOINT = "ipc:///tmp/zmq_ipc_pubsub.sock";

// Environment struct, argp setup, libbpf_print_fn, sig_handler remain the
// same...
static struct env {
    pid_t pid;
    char filter_comm[TASK_COMM_LEN];
    char target_path[PATH_MAX];
    bool verbose;
} env = {
    .pid = 0,
    .filter_comm = "",
    .target_path = DEFAULT_TARGET_LIB,
    .verbose = false,
};

const char *argp_program_version = "ggml_cpu 0.1";
const char *argp_program_bug_address = "DeltaMail@qq.com";
const char argp_program_doc[] =
    "ggml_cpu: Monitor ggml_graph_compute calls.\n\n" // Doc update
    "USAGE: ./ggml_cpu [-p PID] [-c COMM] [-f TARGET_LIB_PATH] [-v]\n"
    "       Default target library: " DEFAULT_TARGET_LIB "\n";

static const struct argp_option opts[] = {
    {"pid", 'p', "PID", 0, "Filter by process PID"},
    {"comm", 'c', "COMMAND", 0, "Filter by process command name"},
    {"file", 'f', "TARGET_LIB_PATH", 0, "Path to the target library to probe"},
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
            fprintf(stderr, "Command name too long (max %d): %s\n",
                    TASK_COMM_LEN - 1, arg);
            argp_usage(state);
        }
        strncpy(env.filter_comm, arg, TASK_COMM_LEN - 1);
        env.filter_comm[TASK_COMM_LEN - 1] = '\0';
        break;
    case 'f':
        if (strlen(arg) >= PATH_MAX) {
            fprintf(stderr, "Target file path too long (max %d): %s\n",
                    PATH_MAX - 1, arg);
            argp_usage(state);
        }
        strncpy(env.target_path, arg, PATH_MAX - 1);
        env.target_path[PATH_MAX - 1] = '\0';
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
    .options = opts, .parser = parse_arg, .doc = argp_program_doc};
static int libbpf_print_fn(enum libbpf_print_level level, const char *format,
                           va_list args) {
    if (level == LIBBPF_DEBUG && !env.verbose)
        return 0;
    return vfprintf(stderr, format, args);
}
static volatile bool exiting = false;
static void sig_handler(int sig) { exiting = true; }

// Helper to get order string (unchanged)
const char *get_order_str(enum ggml_cgraph_eval_order order) {
    switch (order) {
    case GGML_CGRAPH_EVAL_ORDER_LEFT_TO_RIGHT:
        return "L->R"; // Shorter
    case GGML_CGRAPH_EVAL_ORDER_RIGHT_TO_LEFT:
        return "R->L"; // Shorter
    case GGML_CGRAPH_EVAL_ORDER_COUNT:
        return "COUNT(?)";
    default:
        return "UNK"; // Shorter
    }
}

// MODIFIED: Ring buffer event handling callback - Processes the combined event
static int handle_event(void *ctx, void *data, size_t data_sz) {
    if (exiting)
        return -1;

    const struct event *e = data; // Event now contains all info
    zmq_pub_handle_t *handle = ctx;

    if (env.verbose) {
        char ts[32];
        time_t t = time(NULL);
        struct tm *tm_info = localtime(&t);

        if (tm_info == NULL) {
            perror("localtime failed");
            strcpy(ts, "ERR_TS");
        } else {
            // Format time slightly differently maybe
            strftime(ts, sizeof(ts), "%H:%M:%S", tm_info);
        }

        // Print the combined information received at function exit
        printf("%-8s %-7d %-16s | Sz:%-5d Nodes:%-5d Leafs:%-5d Ord:%-4s | "
               "Cost:%llu ns\n",
               ts, e->pid, e->comm, e->graph_size, e->graph_n_nodes,
               e->graph_n_leafs, get_order_str(e->graph_order), e->cost_ns);
    }

    struct ggml_graph_compute_event event = {.timestamp_ns = UnixNanoNow(),
                                             .pid = e->pid,
                                             .comm = "",
                                             .graph_size = e->graph_size,
                                             .graph_n_nodes = e->graph_n_nodes,
                                             .graph_n_leafs = e->graph_n_leafs,
                                             .graph_order = e->graph_order,
                                             .cost_ns = e->cost_ns};
    strncpy(event.comm, e->comm, sizeof(event.comm));

    zmq_pub_send(handle, "ggml_graph_compute", &event,
                 ggml_graph_compute_event_pack);

    return 0; // Continue processing
}

int main(int argc, char **argv) {
    struct ring_buffer *rb = NULL;
    struct ggml_cpu_bpf *skel = NULL;
    int err;
    LIBBPF_OPTS(bpf_uprobe_opts, uprobe_opts);

    // Parse args, set limits, setup logging/signals (unchanged)...
    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err)
        return err;

    libbpf_set_print(libbpf_print_fn);
    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    zmq_pub_handle_t *zmq_handle = zmq_pub_init(ENDPOINT);
    if (!zmq_handle) {
        fprintf(stderr, "Failed to initialize ZMQ publisher\n");
        return 1;
    }

    // Open, load, set rodata (unchanged)...
    skel = ggml_cpu_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }

    skel->rodata->filter_pid = env.pid;
    strncpy((char *)skel->rodata->filter_comm, env.filter_comm, TASK_COMM_LEN);
    skel->rodata->filter_comm[TASK_COMM_LEN - 1] = '\0';

    err = ggml_cpu_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load BPF skeleton: %s\n", strerror(-err));
        goto cleanup;
    }

    // Attach uprobe (unchanged)...
    uprobe_opts.func_name = TARGET_FUNC_NAME;
    uprobe_opts.retprobe = false;
    skel->links.uprobe_ggml_graph_compute =
        bpf_program__attach_uprobe_opts(skel->progs.uprobe_ggml_graph_compute,
                                        -1, env.target_path, 0, &uprobe_opts);
    if (!skel->links.uprobe_ggml_graph_compute) {
        err = -errno;
        fprintf(stderr, "Failed to attach uprobe to %s:%s: %s\n",
                env.target_path, TARGET_FUNC_NAME, strerror(-err));
        goto cleanup;
    }
    printf("Attached uprobe to %s:%s\n", env.target_path, TARGET_FUNC_NAME);

    // Attach uretprobe (unchanged)...
    uprobe_opts.retprobe = true;
    skel->links.uretprobe_ggml_graph_compute = bpf_program__attach_uprobe_opts(
        skel->progs.uretprobe_ggml_graph_compute, -1, env.target_path, 0,
        &uprobe_opts);
    if (!skel->links.uretprobe_ggml_graph_compute) {
        err = -errno;
        fprintf(stderr, "Failed to attach uretprobe to %s:%s: %s\n",
                env.target_path, TARGET_FUNC_NAME, strerror(-err));
        goto cleanup;
    }
    printf("Attached uretprobe to %s:%s\n", env.target_path, TARGET_FUNC_NAME);

    // Setup Ring Buffer (unchanged)...
    rb = ring_buffer__new(bpf_map__fd(skel->maps.rb), handle_event, zmq_handle,
                          NULL);
    if (!rb) {
        err = -errno;
        fprintf(stderr, "Failed to create ring buffer: %s\n", strerror(-err));
        goto cleanup;
    }

    // Main Event Loop (updated header print)...
    printf("Monitoring %s calls (data sent on exit). Press Ctrl+C to exit...\n",
           TARGET_FUNC_NAME);
    // Adjusted header to match the new output format
    printf("%-8s %-7s %-16s | %-5s %-5s %-5s %-4s | %s\n", "TIME", "PID",
           "COMM", "Sz", "Nodes", "Leafs", "Ord", "Cost (ns)");
    while (!exiting) {
        err = ring_buffer__poll(rb, 100 /* ms */);
        if (err == -EINTR) {
            err = 0;
            break;
        }
        if (err < 0) {
            fprintf(stderr, "Error polling ring buffer: %s\n", strerror(-err));
            // Decide if you want to break on poll errors
            // break;
        }
    }

cleanup:
    // Cleanup (unchanged)...
    printf("\nDetaching probes and cleaning up...\n");
    ring_buffer__free(rb);
    ggml_cpu_bpf__destroy(skel);
    zmq_pub_cleanup(&zmq_handle);
    printf("Exited.\n");
    return err < 0 ? -err : 0;
}