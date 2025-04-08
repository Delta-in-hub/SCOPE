#pragma once // 在这个头文件中, 不要 include 任何头文件, 避免冲突. 类型尽量用基本类型, 或者 自己 定义.

#define TASK_COMM_LEN 16
#define CUDA_LIB_PATH_MAX 256 // CUDA 库路径最大长度

// CUDA Memcpy Kind Enum (mirrors cudaMemcpyKind)
enum cuda_memcpy_kind {
	CUDA_MEMCPY_HOST_TO_HOST = 0,
	CUDA_MEMCPY_HOST_TO_DEVICE = 1,
	CUDA_MEMCPY_DEVICE_TO_HOST = 2,
	CUDA_MEMCPY_DEVICE_TO_DEVICE = 3,
	CUDA_MEMCPY_DEFAULT = 4,
};

// 事件类型枚举
enum event_type {
	EVENT_TYPE_MALLOC,
	EVENT_TYPE_FREE,
	EVENT_TYPE_LAUNCH_KERNEL,
	EVENT_TYPE_MEMCPY,
	EVENT_TYPE_SYNC,
};

// 事件结构定义
struct event {
	enum event_type type;      // 事件类型
	int pid;                   // 进程 ID
	char comm[TASK_COMM_LEN];  // 进程名

	union {
		// cudaMalloc 返回
		struct {
            void* allocated_ptr; // 实际分配到的设备指针 (成功时)
            size_t size; // 请求分配的大小
			int retval;          // cudaMalloc 的返回值 (错误码)
		} malloc;

		// cudaFree 入口
		struct {
			void* dev_ptr; // 准备释放的指针
		} free;

		// cudaLaunchKernel 入口
		struct {
			const void* func_ptr; // 内核函数指针 (在设备上的地址)
		} launch_kernel;

		// cudaMemcpy 入口
		struct {
			const void* src;
			void* dst;
			size_t size;
			enum cuda_memcpy_kind kind; // 传输类型
		} memcpy;

		// cudaDeviceSynchronize 返回
		struct {
			uint64_t duration_ns; // 函数执行耗时 (纳秒)
		} sync;
	};
};

// 用于在 cudaMalloc uprobe 和 uretprobe 之间传递 devPtr 的地址
struct malloc_entry_data {
	void** user_dev_ptr_addr; // 用户空间传入的 devPtr 参数的地址 (void**)
    size_t size;
};

// 用于在 cudaDeviceSynchronize uprobe 和 uretprobe 之间传递时间戳
struct sync_entry_data {
	uint64_t entry_ts; // 进入函数时的时间戳
};