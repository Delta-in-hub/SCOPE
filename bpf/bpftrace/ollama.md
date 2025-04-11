# Probes

1.  **`uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_compute`** (CPU 后端)
    *   **原因**: 这是 GGML 在 CPU 上执行计算图的核心函数。监控此函数的进入和退出，可以了解 CPU 计算任务的开始和结束，以及大致的执行时间。其调用频率（40次）虽然不高，但每次调用都代表一个重要的计算阶段。


2.  **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_vec_q`** 和 **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_q`** (CUDA 后端)
    *   **原因**: 这两个函数（分别是 quantized matrix-vector 和 matrix-matrix 乘法）是 LLM 推理中最核心、计算量最大的操作之一，尤其是在 GPU 上。它们的高调用频率（分别为 9975 和 165 次，加起来非常显著）表明了它们的重要性。监控这些函数可以深入了解 GPU 的计算负载和性能瓶颈，例如可以统计调用次数、计算每次调用的耗时。

**内存管理与数据移动:**

5.  **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::alloc`** 和 **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::free`** (CUDA 后端)
    *   **原因**: 这两个函数直接关联到 CUDA 虚拟内存管理（VMM）池的内存分配和释放。监控它们可以精确地追踪 GPU 显存的使用情况，了解显存的动态分配/释放模式，对于诊断显存不足或碎片化问题非常有价值。高调用频率（10473 次）说明 GPU 内存操作非常频繁。
6.  **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cpy`** (CUDA 后端)
    *   **原因**: 这个函数负责在 GPU 内部或 Host（CPU内存）与 Device（GPU显存）之间复制数据。数据传输是常见的性能瓶颈之一。监控此函数可以了解数据传输的方向、大小和频率，有助于识别不必要的传输或优化传输策略。调用频率（4320 次）也相当高。

7.  **`uprobe:/usr/lib/ollama/libggml-base.so:ggml_aligned_malloc`** 和 **`uprobe:/usr/lib/ollama/libggml-base.so:ggml_aligned_free`** (基础库)
    *   **原因**: 这是底层的 CPU 内存分配和释放。虽然 GPU 内存通常更关键，但监控 CPU 内存分配有助于了解整体资源使用情况，特别是在模型加载或 CPU / GPU 混合执行的场景下。

**初始化与配置:**

8.  **`uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_set_device`** (CUDA 后端)
    *   **原因**: 如果系统中有多个 GPU，这个函数用于选择当前活动的 GPU 设备。监控它可以确认 Ollama 是否正确地使用了预期的 GPU，或者是否存在不必要的设备切换。
