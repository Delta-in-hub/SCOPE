- System Wide (Prometheus)
  - Cpu Usage
  - Memory Usage
  - I/O Usage
  - Network Traffic
  - Gpu Usage
  - Gpu Memory Usage
  - Health Check golang backend ("/health")

- GPU(CUDA)
  - cudaMalloc
  - cudaFree
  - cudaMallocHost
  - cudaFreeHost
  - cudaLaunchKernel
  - cudaMemcpy
  - cudaDeviceSynchronize

> https://docs.python.org/3/howto/instrumentation.html
- Python
  - function__entry
  - function__return


- Ollama
  - uprobe:/usr/bin/ollama:llamaLog
    - extern void llamaLog(int level, char* text, void* user_data);
  - uprobe:/usr/lib/ollama/libggml-base.so:gguf_init_from_file
  - uprobe:/usr/lib/ollama/libggml-base.so:gguf_init_from_file, gguf_init_from_file_impl
    - (高价值 - 模型加载) 指示开始从文件加载 GGUF 模型。可以计时这个函数的执行时间来大致了解模型加载耗时。
  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor*, ggml_set_param
    - (中等价值) 在内存中构建计算图（创建张量、设置参数）。
  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_compute
    - (高价值 - 计算图执行) 非常重要。这通常是触发整个计算图执行的入口点（在特定后端执行前）。监控它的调用标志着推理计算的开始。
  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_graph_compute
    - (高价值 - 后端调度) 如果 Ollama 使用调度器在多个后端（CPU/GPU）之间分配工作，这个函数是调度执行的入口。
  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_alloc, ggml_backend_tensor_copy
    - (高价值 - 跨后端操作) 指示在某个后端（可能是 CPU 或 GPU）分配张量，或在不同后端之间拷贝张量数据。这是监控数据流动的关键点。
  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul_mat
    - (较高价值) 通用的矩阵乘法操作入口。监控它可以了解 MatMul 的总体频率，无论是在 CPU 还是 GPU 上执行。

  - uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_compute, ggml_graph_compute_with_ctx
    - (高价值 - CPU 计算) 这些函数表明计算图正在CPU上执行。监控它们可以确认 CPU 是否被用于推理。
  - uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_threadpool_*
    - (中等价值) CPU 线程池管理，了解 CPU 并行计算的启动和停止。

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_*, ggml_cuda_*_cuda
    - (高价值 - GPU 计算) 这些函数表明特定的 GGML 操作（如矩阵乘法 mul_mat, add, rope, norm 等）正在被卸载到 GPU执行。监控它们可以确认 GPU 是否在工作以及哪些操作在 GPU 上运行。
  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:launch_fattn<*>: (极高价值 - GPU 计算) 非常重要。直接表明正在启动 FlashAttention CUDA kernel。这是现代 LLM 中计算量最大的部分之一。监控它可以确认是否在使用 FlashAttention 以及频率。
  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_init, ggml_cuda_set_device: (中等价值) GPU 初始化和设备选择。
  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cpy, ggml_cuda_cpy_fn: (较高价值 - GPU 数据传输) GGML 层面的 GPU 数据拷贝操作。结合 CUDA API 的 cudaMemcpy* 可以更全面地了解数据移动。

  - 和底层 CUDA 库 (libcudart.so, libcuda.so) 的交互, PID or comm


> https://github.com/iovisor/bcc/tree/master/libbpf-tools
- Process (via pid or comm)
  - all system calls
  - 子进程创建
  - 进程调度, sched_wakeup , sched_switch
  - 文件系统操作, vfs_open 打开文件的路径
  - 网络活动
    -  TCP/UDP 连接建立的延迟和频率、TCP 连接的生命周期、数据收发、DNS 解析延迟。



- reuseport
  - 服务热更新?


- 动态新增 libbpf 程序, 编译/挂载
  - LLM Agent