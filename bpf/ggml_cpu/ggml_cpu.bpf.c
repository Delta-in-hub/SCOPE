#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "ggml_cpu.h" // Include our header last




char LICENSE[] SEC("license") = "Dual BSD/GPL";

// Optional filtering parameters (runtime configurable)
const volatile pid_t filter_pid = 0;
const volatile char filter_comm[TASK_COMM_LEN] = {}; // Initialized to empty

// Ring buffer map to send events to user space
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1024 * 1024); // Adjust size as needed
} rb SEC(".maps");


struct entry_data {
    uint64_t entry_ts; // Entry timestamp
    int graph_size;
    int graph_n_nodes;
    int graph_n_leafs;
    enum ggml_cgraph_eval_order graph_order;
    // Store comm here as well, in case it changes between entry/exit? Optional.
    // char comm[TASK_COMM_LEN];
};

// Hash map to store entry data, keyed by PID
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, MAX_ENTRIES);
	__type(key, pid_t);
	__type(value, struct entry_data); // Store the new struct
} entry_data_map SEC(".maps");



// Helper function for command filtering
static __always_inline int comm_allowed(const char *comm) {
    #pragma unroll
    for (int i = 0; i < TASK_COMM_LEN && filter_comm[i] != '\0'; i++) {
        if (comm[i] != filter_comm[i])
            return 0;
    }
    return 1;
}

// Helper function for process filtering
static __always_inline int process_allowed(pid_t pid, const char *comm) {
    if (filter_pid != 0 && pid != filter_pid)
        return 0;
    if (filter_comm[0] != '\0' && !comm_allowed(comm))
        return 0;
    return 1;
}


SEC("uprobe")
int BPF_KPROBE(uprobe_ggml_graph_compute, struct ggml_cgraph * cgraph, void* cplan)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    // Apply filters
    if (!process_allowed(pid, comm))
        return 0;

    // Get entry timestamp
	uint64_t ts = bpf_ktime_get_ns();

    // Prepare struct to store entry data
    struct entry_data entry = {}; // Initialize to zero
    entry.entry_ts = ts;

    // Read ggml_cgraph fields from user space
    long err;
    int size, n_nodes, n_leafs;
    enum ggml_cgraph_eval_order order;

    if (cgraph == NULL) {
        bpf_printk("uprobe ggml_graph_compute: cgraph is NULL for pid %d\n", pid);
        return 0; // Don't store anything if pointer is bad
    }

    err = bpf_probe_read_user(&size, sizeof(size), &cgraph->size);
    if (err != 0) {
        bpf_printk("uprobe ggml_graph_compute: Failed read size: %ld\n", err);
        return 0; // Don't store incomplete data
    }
    err = bpf_probe_read_user(&n_nodes, sizeof(n_nodes), &cgraph->n_nodes);
     if (err != 0) {
        bpf_printk("uprobe ggml_graph_compute: Failed read n_nodes: %ld\n", err);
        return 0;
    }
    err = bpf_probe_read_user(&n_leafs, sizeof(n_leafs), &cgraph->n_leafs);
     if (err != 0) {
        bpf_printk("uprobe ggml_graph_compute: Failed read n_leafs: %ld\n", err);
        return 0;
    }
    err = bpf_probe_read_user(&order, sizeof(order), &cgraph->order);
     if (err != 0) {
        bpf_printk("uprobe ggml_graph_compute: Failed read order: %ld\n", err);
        return 0;
    }

    // Populate entry_data struct
    entry.graph_size = size;
    entry.graph_n_nodes = n_nodes;
    entry.graph_n_leafs = n_leafs;
    entry.graph_order = order;
    // Optionally store comm: bpf_core_read_str(&entry.comm, sizeof(entry.comm), comm);

    // Store the entry data in the map
    bpf_map_update_elem(&entry_data_map, &pid, &entry, BPF_ANY);

	return 0;
}


// Uretprobe handler - retrieves data, calculates cost, sends combined event
SEC("uretprobe")
int BPF_KRETPROBE(uretprobe_ggml_graph_compute, int ret) // ret is the return value
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;

    // Look up entry data stored by the uprobe
    struct entry_data *entry_ptr = bpf_map_lookup_elem(&entry_data_map, &pid);
    if (!entry_ptr) {
        // Entry data not found - maybe uprobe failed, filtered, or map error
        return 0; // Cannot proceed
    }

    // Got entry data, calculate duration
    uint64_t exit_ts = bpf_ktime_get_ns();
    uint64_t cost_ns = exit_ts - entry_ptr->entry_ts;

    // Data retrieved, delete the entry from the map *now*
    bpf_map_delete_elem(&entry_data_map, &pid);

    // Prepare the combined event for the ring buffer
    struct event *e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
    if (!e) {
        bpf_printk("uretprobe ggml_graph_compute: Failed reserve ringbuf for pid %d\n", pid);
        return 0; // Skip event if buffer is full
    }

    // Populate the event struct
    e->pid = pid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    e->cost_ns = cost_ns;

    // Copy details from the retrieved entry data
    e->graph_size = entry_ptr->graph_size;
    e->graph_n_nodes = entry_ptr->graph_n_nodes;
    e->graph_n_leafs = entry_ptr->graph_n_leafs;
    e->graph_order = entry_ptr->graph_order;

    // Submit the combined event
    bpf_ringbuf_submit(e, 0);

	return 0;
}