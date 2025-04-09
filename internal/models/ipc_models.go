package models

//! [vfs_open] START

const VfsOpenTopic = "vfs_open"

// --- NEW: Struct for vfs_open events (Array Format) ---
// NOTE: Field order MUST match the C packer's array order: [timestamp, pid, comm, filename].
// NOTE: No msgpack tags are used; unmarshaling relies on field order.
type VfsOpenEvent struct {
	TimestampNs int64  // Index 0 in the packed array
	PID         int32  // Index 1 in the packed array
	Comm        string // Index 2 in the packed array
	Filename    string // Index 3 in the packed array
}

//! [vfs_open] END

//! [syscalls] START

const SyscallsTopic = "syscalls"

type SyscallsEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	SyscallName string
}

//! [syscalls] END

//! [sched]

const SchedTopic = "sched"

type SchedEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	Cpu         int32
	Type        int32 // enum event_type { SWITCH_IN, SWITCH_OUT };
}

//! [ollamabin]

const OllamabinTopic = "llamaLog"

type LlamaLogEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	Text        string
}

//! [ggml_cuda]

const GGMLCudaTopic = "ggml_cuda"

type GGMLCudaEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	FuncName    string
	DurationNs  int64
}

//! [ggml_cpu]

const GGMLCpuTopic = "ggml_graph_compute"

type GGMLCpuEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	GraphSize   int32
	GraphNodes  int32
	GraphLeafs  int32
	GraphOrder  int32
	CostNs      int64
}

//! [ggml_base]

const GGMLBaseTopic = "ggml_base"

type GGMLBaseEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	Type        int32
	Size        uint64
	Ptr         uint64
}

//! [execv]

const ExecvTopic = "execv"

type ExecvEvent struct {
	TimestampNs int64
	PID         int32
	Ppid        int32
	Filename    string // for PID
	Args        string // for PID
}

//! [cuda]

const CudaMallocTopic = "cudaMalloc"

type CudaMallocEvent struct {
	TimestampNs  int64
	PID          int32
	Comm         string
	AllocatedPtr uint64
	Size         uint64
	Retval       int
}

const CudaFreeTopic = "cudaFree"

type CudaFreeEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	DevPtr      uint64
}

const CudaLaunchKernelTopic = "cudaLaunchKernel"

type CudaLaunchKernelEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	FuncPtr     uint64
}

const CudaMemcpyTopic = "cudaMemcpy"

type CudaMemcpyEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	Src         uint64
	Dst         uint64
	Size        uint64
	Kind        int
}

const CudaSyncTopic = "cudaDeviceSynchronize"

type CudaSyncEvent struct {
	TimestampNs int64
	PID         int32
	Comm        string
	DurationNs  uint64
}
