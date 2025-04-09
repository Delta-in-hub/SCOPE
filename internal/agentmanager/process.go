package agentmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"scope/internal/models"
	"scope/internal/platform"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

var machineIDOnce sync.Once
var cachedMachineID string

func getMachineID() string {
	machineIDOnce.Do(func() {
		id, err := machineid.ID()
		if err != nil {
			log.Printf("Error getting machine ID: %v", err)
			cachedMachineID = uuid.New().String()
		} else {
			cachedMachineID = id
		}
	})
	return cachedMachineID
}

// --- Processor Goroutine (Handles Different Message Types, including Array) ---
// Reads raw messages, unmarshals topic and payload, and processes.
func Processor(msgChan <-chan RawMessage, wg *sync.WaitGroup, config Config, redisClient *goredis.Client) {
	defer wg.Done()

	ctx := context.Background()

	for rawMsg := range msgChan {
		payloadBytes := rawMsg.Payload

		// --- Unmarshal the Topic first ---
		// (Topic is still expected to be a msgpack encoded string)
		var topic string
		err := msgpack.Unmarshal(rawMsg.Topic, &topic)
		if err != nil {
			log.Printf("Processor: Error unmarshaling topic: %v (Raw Topic Bytes: %x)", err, rawMsg.Topic)
			continue // Skip message with unparseable topic
		}
		// --- Topic successfully unmarshaled ---

		// Prepare data for Redis Stream
		var eventData map[string]interface{}

		// Use a switch to handle different topics
		switch topic {
		case models.VfsOpenTopic:
			// Use the untagged struct expecting an array format based on field order.
			var event models.VfsOpenEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				// If this error occurs often, double-check C packing order vs Go struct order
				log.Printf("Processor: Error unmarshaling VfsOpenEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			// Format timestamp for better readability
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Print received event data - access fields by name as usual
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"filename":  event.Filename,
			}

			// Only print if verbose mode is enabled
			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', File='%s'\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Filename)

				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}

		case models.SyscallsTopic:
			var event models.SyscallsEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling SyscallsEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"syscall":   event.SyscallName,
			}

			// Only print if verbose mode is enabled
			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Syscall='%s'\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.SyscallName)
			}

			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}

		case models.SchedTopic:
			var event models.SchedEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling SchedEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"cpu":       event.Cpu,
				"type":      "switch_in",
			}
			if event.Type == 0 {
				// Create event data for Redis
				eventData["type"] = "switch_in"
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Into cpu(%d)\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Cpu)
				}
			} else if event.Type == 1 {
				// Create event data for Redis
				eventData["type"] = "switch_out"

				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Out cpu(%d)\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Cpu)
				}
			} else {
				eventData["type"] = "unknown"
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Unknown type(%d)\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Type)
				}
			}
			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}

		case models.OllamabinTopic:
			var event models.LlamaLogEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling LlamaLogEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"text":      event.Text,
			}
			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Text='%s'\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Text)
			}
			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}
		case models.GGMLCudaTopic:
			var event models.GGMLCudaEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling GGMLCudaEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":       topic,
				"timestamp":   event.TimestampNs,
				"pid":         event.PID,
				"comm":        event.Comm,
				"cmdline":     cmdline,
				"operation":   event.FuncName,
				"func_name":   event.FuncName,
				"duration_ns": event.DurationNs,
			}
			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Func='%s', Duration=%d ns\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.FuncName, event.DurationNs)
			}
			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}

		case models.GGMLCpuTopic:
			var event models.GGMLCpuEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling GGMLCpuEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			/*
				enum ggml_cgraph_eval_order {
					GGML_CGRAPH_EVAL_ORDER_LEFT_TO_RIGHT = 0,
					GGML_CGRAPH_EVAL_ORDER_RIGHT_TO_LEFT,
					GGML_CGRAPH_EVAL_ORDER_COUNT // Should be 2
				};
			*/
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":       topic,
				"timestamp":   event.TimestampNs,
				"pid":         event.PID,
				"comm":        event.Comm,
				"cmdline":     cmdline,
				"operation":   "ggml_graph_compute",
				"graph_size":  event.GraphSize,
				"graph_nodes": event.GraphNodes,
				"graph_leafs": event.GraphLeafs,
				"cost_ns":     event.CostNs,
			}

			switch event.GraphOrder {
			case 0:
				eventData["graph_order"] = "LEFT_TO_RIGHT"
			case 1:
				eventData["graph_order"] = "RIGHT_TO_LEFT"
			default:
				eventData["graph_order"] = "COUNT"
			}

			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', GraphSize=%d, GraphNodes=%d, GraphLeafs=%d, GraphOrder=%d, Cost=%d ns\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.GraphSize, event.GraphNodes, event.GraphLeafs, event.GraphOrder, event.CostNs)
			}
			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}
		case models.GGMLBaseTopic:
			var event models.GGMLBaseEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling GGMLBaseEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis with common fields
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"size":      event.Size,
				"ptr":       event.Ptr,
			}

			if event.Type == 0 {
				// Add operation type for malloc
				eventData["operation"] = "ggml_aligned_malloc"
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Aligned Malloc Size=%d, Pointer=0x%x\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Size, event.Ptr)
				}
			} else if event.Type == 1 {
				// Add operation type for free
				eventData["operation"] = "ggml_aligned_free"
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Aligned Free Size=%d, Pointer=0x%x\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Size, event.Ptr)
				}
			}
			if config.Verbose {
				tsnow := time.Now().UnixNano()
				fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
			}
		case models.ExecvTopic:
			var event models.ExecvEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling ExecvEvent (Array Format): %v", err)
				continue
			}
			ppidcomm, _ := platform.GetComm(int(event.Ppid))
			ppidcmdline, _ := platform.GetCmdline(int(event.Ppid))
			pidcomm, _ := platform.GetComm(int(event.PID))
			pidcmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":        topic,
				"timestamp":    event.TimestampNs,
				"ppid":         event.Ppid,
				"ppid_comm":    ppidcomm,
				"ppid_cmdline": ppidcmdline,
				"pid":          event.PID,
				"pid_comm":     pidcomm,
				"pid_cmdline":  pidcmdline,
				"filename":     event.Filename,
				"args":         event.Args,
			}

			if config.Verbose {
				fmt.Printf("%s Process[%d comm:%s cmdline:%s] created subprocess[%d comm:%s cmdline:%s filename:%s args:%s]\n",
					tsFormatted, event.Ppid, ppidcomm, ppidcmdline, event.PID, pidcomm, pidcmdline, event.Filename, event.Args)
			}
		case models.CudaMallocTopic:
			var event models.CudaMallocEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaMallocEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"operation": "cudaMalloc",
				"cmdline":   cmdline,
				"ptr":       event.AllocatedPtr,
				"size":      event.Size,
				"retval":    event.Retval,
			}

			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMalloc AllocatedPtr=0x%x, Size=%d, Retval=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.AllocatedPtr, event.Size, event.Retval)
			}
		case models.CudaFreeTopic:
			var event models.CudaFreeEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaFreeEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"cmdline":   cmdline,
				"operation": "cudaFree",
				"ptr":       event.DevPtr,
			}

			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaFree DevPtr=0x%x\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.DevPtr)
			}
		case models.CudaLaunchKernelTopic:
			var event models.CudaLaunchKernelEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaLaunchKernelEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)

			// 通过读取 /proc/PID/maps 可以获取到 funcptr 属于哪个库, 使用 addr2line 可以获取到函数名
			symbol, err := platform.FindSymbolFromPidPtr(int(event.PID), uintptr(event.FuncPtr))
			if err != nil {
				fmt.Printf("Symbol: Error finding symbol: %v\n", err)
			}

			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":         topic,
				"timestamp":     event.TimestampNs,
				"pid":           event.PID,
				"comm":          event.Comm,
				"operation":     "cudaLaunchKernel",
				"cmdline":       cmdline,
				"func_ptr":      event.FuncPtr,
				"symbol_name":   symbol.SymbolName,
				"symbol_file":   symbol.FilePath,
				"symbol_offset": symbol.Offset,
			}

			// Add source file information if available
			if symbol.SourceLine != 0 {
				eventData["symbol_sourcefile"] = fmt.Sprintf("%s:%d", symbol.SourceFile, symbol.SourceLine)
			}

			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaLaunchKernel FuncPtr=0x%x\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.FuncPtr)
				fmt.Printf("Symbol: Name='%s', File='%s', Offset=0x%x", symbol.SymbolName, symbol.FilePath, symbol.Offset)
				if symbol.SourceLine != 0 {
					fmt.Printf(" SourceFile=%s:%d", symbol.SourceFile, symbol.SourceLine)
				}
				fmt.Println()
			}

		case models.CudaMemcpyTopic:
			var event models.CudaMemcpyEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaMemcpyEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)

			// Create event data for Redis with common fields
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": event.TimestampNs,
				"pid":       event.PID,
				"comm":      event.Comm,
				"operation": "cudaMemcpy",
				"cmdline":   cmdline,
				"src":       event.Src,
				"dst":       event.Dst,
				"size":      event.Size,
				"kind":      event.Kind,
			}

			// Add a human-readable transfer type based on kind
			switch event.Kind {
			case 0:
				eventData["type"] = "host_to_host"
			case 1:
				eventData["type"] = "host_to_device"
			case 2:
				eventData["type"] = "device_to_host"
			case 3:
				eventData["type"] = "device_to_device"
			case 4:
				eventData["type"] = "default"
			default:
				eventData["type"] = "unknown"
			}

			switch event.Kind {
			/*
								enum cuda_memcpy_kind {
					CUDA_MEMCPY_HOST_TO_HOST = 0,
					CUDA_MEMCPY_HOST_TO_DEVICE = 1,
					CUDA_MEMCPY_DEVICE_TO_HOST = 2,
					CUDA_MEMCPY_DEVICE_TO_DEVICE = 3,
					CUDA_MEMCPY_DEFAULT = 4,
				};
			*/
			case 0:
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Host(0x%x) To Host(0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
			case 1:
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Host(0x%x) To Device(0x%x), Size=%d\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
				}
			case 2:
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Device(0x%x) To Host(0x%x), Size=%d\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
				}
			case 3:
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Device(0x%x) To Device(0x%x), Size=%d\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
				}
			case 4:
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy DEFAULT From (0x%x) To (0x%x), Size=%d\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
				}
			default:
				if config.Verbose {
					fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy Unknown Kind=%d, From (0x%x) To (0x%x), Size=%d\n",
						topic, tsFormatted, event.PID, event.Comm, cmdline, event.Kind, event.Src, event.Dst, event.Size)
				}
			}
		case models.CudaSyncTopic:
			var event models.CudaSyncEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaSyncEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Create event data for Redis
			eventData = map[string]interface{}{
				"topic":       topic,
				"timestamp":   event.TimestampNs,
				"pid":         event.PID,
				"comm":        event.Comm,
				"operation":   "cudaDeviceSynchronize",
				"cmdline":     cmdline,
				"duration_ns": event.DurationNs,
			}

			// Only print if verbose mode is enabled
			if config.Verbose {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaSync Duration=%d ns\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.DurationNs)
			}
		default:
			// Handle unexpected topics

			log.Printf("Processor: Warning: Received message with unhandled topic '%s'", topic)

			// Create generic event data for Redis
			eventData = map[string]interface{}{
				"topic":     topic,
				"timestamp": time.Now().UnixNano(),
				"payload":   string(payloadBytes),
			}
		}

		// Send data to Redis Stream if we have event data

		eventData["machineid"] = getMachineID()

		eventJson, err := json.Marshal(eventData)
		if err != nil {
			log.Printf("Error json marshaling event data: %v", err)
			continue
		}

		// Add to Redis Stream with optimized performance
		_, err = redisClient.XAdd(ctx, &goredis.XAddArgs{
			Stream: config.StreamKey,
			Values: string(eventJson),
		}).Result()

		if err != nil {
			log.Printf("Error adding event to Redis Stream: %v", err)
		}

	} // End for range msgChan

	log.Println("Processor goroutine finished (channel closed).")
}
