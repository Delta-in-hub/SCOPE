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
  - cudaMallocHost (TODO)
  - cudaFreeHost (TODO)
  - cudaLaunchKernel
  - cudaMemcpy
  - cudaDeviceSynchronize

> Works for ollama


> https://docs.python.org/3/howto/instrumentation.html
- Python
  - function__entry
  - function__return


- Ollama
  - uprobe:/usr/bin/ollama:llamaLog
    - extern void llamaLog(int level, char* text, void* user_data);

  - uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_compute
    - 这是 GGML 在 CPU 上执行计算图的核心函数。监控此函数的进入和退出，可以了解 CPU 计算任务的开始和结束，以及大致的执行时间。其调用频率（40次）虽然不高，但每次调用都代表一个重要的计算阶段。

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_vec_q
    - 这是 GGML 在 GPU 上执行 quantized matrix-vector 乘法的核心函数。监控此函数的进入和退出，可以了解 GPU 计算任务的开始和结束，以及大致的执行时间。其调用频率（9975次）非常高，表明了 GPU 计算任务的密集程度。

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_q
    - 这是 GGML 在 GPU 上执行 quantized matrix-matrix 乘法的核心函数。监控此函数的进入和退出，可以了解 GPU 计算任务的开始和结束，以及大致的执行时间。其调用频率（165次）虽然不高，但每次调用都代表一个重要的计算阶段。

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::alloc (TODO)
    - 这两个函数直接关联到 CUDA 虚拟内存管理（VMM）池的内存分配和释放。监控它们可以精确地追踪 GPU 显存的使用情况，了解显存的动态分配/释放模式，对于诊断显存不足或碎片化问题非常有价值。高调用频率（10473 次）说明 GPU 内存操作非常频繁。
  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::free (TODO)

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cpy (TODO)
    - 这个函数负责在 GPU 内部或 Host（CPU内存）与 Device（GPU显存）之间复制数据。数据传输是常见的性能瓶颈之一。监控此函数可以了解数据传输的方向、大小和频率，有助于识别不必要的传输或优化传输策略。调用频率（4320 次）也相当高。

  - uprobe:/usr/lib/ollama/libggml-base.so:ggml_aligned_malloc
    - 这是底层的 CPU 内存分配和释放。虽然 GPU 内存通常更关键，但监控 CPU 内存分配有助于了解整体资源使用情况，特别是在模型加载或 CPU / GPU 混合执行的场景下。

  - uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_set_device (TODO)
    - 如果系统中有多个 GPU，这个函数用于选择当前活动的 GPU 设备。监控它可以确认 Ollama 是否正确地使用了预期的 GPU，或者是否存在不必要的设备切换。


> https://github.com/iovisor/bcc/tree/master/libbpf-tools
- Process (via pid or comm)
  - all system calls (done -> syscalls)
  - 子进程创建 (done -> execv)
  - 进程调度, sched_wakeup , sched_switch (done -> sched)
  - 文件系统操作, vfs_open 打开文件的路径 (done -> vfs_open)
  - 网络活动 (TODO)
    -  TCP/UDP 连接建立的延迟和频率、TCP 连接的生命周期、数据收发、



- reuseport
  - 服务热更新?


- 动态新增 libbpf 程序, 编译/挂载
  - LLM Agent