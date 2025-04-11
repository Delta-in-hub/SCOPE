#pragma once

#define TASK_COMM_LEN 16
#define MAX_ENTRIES	10240
#define DEFAULT_TARGET_LIB "/usr/lib/ollama/libggml-cpu-alderlake.so"
#define TARGET_FUNC_NAME "ggml_graph_compute"

enum ggml_cgraph_eval_order {
    GGML_CGRAPH_EVAL_ORDER_LEFT_TO_RIGHT = 0,
    GGML_CGRAPH_EVAL_ORDER_RIGHT_TO_LEFT,
    GGML_CGRAPH_EVAL_ORDER_COUNT // Should be 2
};

struct ggml_hash_set {
    size_t size;
    void * used;
    void ** keys;
};

struct ggml_cgraph {
    int size;
    int n_nodes;
    int n_leafs;

    void ** nodes;
    void ** grads;
    void ** grad_accs;
    void ** leafs;

    struct ggml_hash_set visited_hash_set; // Size/layout might matter if accessed

    enum ggml_cgraph_eval_order order;
};

// --- End of copied structures ---

struct event {
    int pid;
    char comm[TASK_COMM_LEN];

    // Fields from ggml_cgraph (collected at entry, sent at exit)
    int graph_size;
    int graph_n_nodes;
    int graph_n_leafs;
    enum ggml_cgraph_eval_order graph_order;

    // Execution time (calculated at exit)
    uint64_t cost_ns;
};