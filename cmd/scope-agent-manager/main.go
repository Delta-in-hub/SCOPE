package main

import (
	"fmt"
	"log"
	"os"
	"scope/internal/models"
	"scope/internal/platform"
	"strings"
	"sync" // Import sync package for WaitGroup
	"syscall"
	"time"

	"github.com/joho/godotenv"
	zmq "github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

// --- Struct to pass raw messages between goroutines ---
type RawMessage struct {
	Topic   []byte // Received raw topic bytes (might be msgpack encoded)
	Payload []byte // Received raw payload bytes
}

// Reads from ZMQ socket and sends raw messages to the channel.
func receiver(subscriber *zmq.Socket, msgChan chan<- RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(msgChan)

	fmt.Println("Receiver goroutine started.")

	poller := zmq.NewPoller()
	poller.Add(subscriber, zmq.POLLIN)

	running := true
	for running {
		polledSockets, err := poller.Poll(250 * time.Millisecond)
		if err != nil {
			errno := zmq.AsErrno(err)
			if errno == zmq.ETERM {
				fmt.Println("Receiver: Context terminated during poll, exiting.")
				running = false
				continue
			}
			if errno == zmq.Errno(syscall.EINTR) {
				fmt.Println("Receiver: Poll interrupted, continuing...")
				continue
			}
			log.Printf("Receiver: Error polling socket: %v (errno %d)", err, errno)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if len(polledSockets) > 0 {
			msgParts, err := subscriber.RecvMessageBytes(0)
			if err != nil {
				if zmq.AsErrno(err) == zmq.ETERM {
					fmt.Println("Receiver: Context terminated during receive, exiting.")
					running = false
					continue
				}
				log.Printf("Receiver: Error receiving message: %v", err)
				continue
			}

			if len(msgParts) != 2 {
				log.Printf("Receiver: Error: Received message with %d parts, expected 2 (EncodedTopic, Payload)", len(msgParts))
				continue
			}

			msg := RawMessage{
				Topic:   msgParts[0],
				Payload: msgParts[1],
			}

			msgChan <- msg
		}
	}

	fmt.Println("Receiver goroutine finished.")
}

// --- Processor Goroutine (Handles Different Message Types, including Array) ---
// Reads raw messages, unmarshals topic and payload, and processes.
func processor(msgChan <-chan RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Processor goroutine started.")

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
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', File='%s'\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.Filename)

			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)

		case models.SyscallsTopic:
			var event models.SyscallsEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling SyscallsEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Syscall='%s'\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.SyscallName)

			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)

		case models.SchedTopic:
			var event models.SchedEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling SchedEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			if event.Type == 0 {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Into cpu(%d)\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Cpu)
			} else if event.Type == 1 {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Out cpu(%d)\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Cpu)
			} else {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Sched Unknown type(%d)\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Type)
			}
			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)

		case models.OllamabinTopic:
			var event models.LlamaLogEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling LlamaLogEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Text='%s'\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.Text)
			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
		case models.GGMLCudaTopic:
			var event models.GGMLCudaEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling GGMLCudaEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Func='%s', Duration=%d ns\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.FuncName, event.DurationNs)
			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)

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
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', GraphSize=%d, GraphNodes=%d, GraphLeafs=%d, GraphOrder=%d, Cost=%d ns\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.GraphSize, event.GraphNodes, event.GraphLeafs, event.GraphOrder, event.CostNs)
			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
		case models.GGMLBaseTopic:
			var event models.GGMLBaseEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling GGMLBaseEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			if event.Type == 0 {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Aligned Malloc Size=%d, Pointer=0x%x\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Size, event.Ptr)
			} else if event.Type == 1 {
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', Aligned Free Size=%d, Pointer=0x%x\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Size, event.Ptr)
			}
			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)
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
			fmt.Printf("%s Process[%d comm:%s cmdline:%s] created subprocess[%d comm:%s cmdline:%s filename:%s args:%s]\n",
				tsFormatted, event.Ppid, ppidcomm, ppidcmdline, event.PID, pidcomm, pidcmdline, event.Filename, event.Args)
		case models.CudaMallocTopic:
			var event models.CudaMallocEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaMallocEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMalloc AllocatedPtr=0x%x, Size=%d, Retval=%d\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.AllocatedPtr, event.Size, event.Retval)
		case models.CudaFreeTopic:
			var event models.CudaFreeEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaFreeEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaFree DevPtr=0x%x\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.DevPtr)
		case models.CudaLaunchKernelTopic:
			var event models.CudaLaunchKernelEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaLaunchKernelEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaLaunchKernel FuncPtr=0x%x\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.FuncPtr)
			// 通过读取 /proc/PID/maps 可以获取到 funcptr 属于哪个库, 使用 addr2line 可以获取到函数名
			symbol, err := platform.FindSymbolFromPidPtr(int(event.PID), uintptr(event.FuncPtr))
			if err != nil {
				fmt.Printf("Symbol: Error finding symbol: %v\n", err)
				continue
			}
			fmt.Printf("Symbol: Name='%s', File='%s', Offset=0x%x", symbol.SymbolName, symbol.FilePath, symbol.Offset)
			if symbol.SourceLine != 0 {
				fmt.Printf(" SourceFile=%s:%d", symbol.SourceFile, symbol.SourceLine)
			}
			fmt.Println()

		case models.CudaMemcpyTopic:
			var event models.CudaMemcpyEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CudaMemcpyEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(int(event.PID))
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
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
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Host(0x%x) To Device(0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
			case 2:
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Device(0x%x) To Host(0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
			case 3:
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy From Device(0x%x) To Device(0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
			case 4:
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy DEFAULT From (0x%x) To (0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Src, event.Dst, event.Size)
			default:
				fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaMemcpy Unknown Kind=%d, From (0x%x) To (0x%x), Size=%d\n",
					topic, tsFormatted, event.PID, event.Comm, cmdline, event.Kind, event.Src, event.Dst, event.Size)
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
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', CudaSync Duration=%d ns\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.DurationNs)
		default:
			// TODO
			// 可以处理 , 动态的 eBPF
			// Topic , Payload(是一个字符串)
			// Payload 每个 k:v 一行

			// Handle unexpected topics
			log.Printf("Processor: Warning: Received message with unhandled topic '%s'", topic)
			// log.Printf("Raw Payload for unhandled topic: %x", payloadBytes)
		}
	} // End for range msgChan

	fmt.Println("Processor goroutine finished (channel closed).")
}

// --- Main Function (Unchanged from previous version) ---
func main() {
	godotenv.Load(".env")

	context, err := zmq.NewContext()
	if err != nil {
		log.Fatalf("Error creating ZMQ context: %v", err)
	}
	defer context.Term()

	subscriber, err := context.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatalf("Error creating ZMQ subscriber socket: %v", err)
	}
	defer subscriber.Close()

	ipcEndpoint := "ipc:///tmp/zmq_ipc_pubsub.sock"

	fmt.Printf("Go Subscriber binding to %s\n", ipcEndpoint)
	err = subscriber.Bind(ipcEndpoint)
	if err != nil {
		log.Fatalf("Failed to bind subscriber to '%s': %v", ipcEndpoint, err)
	}

	// --- 新增：在 Bind 成功后修改 IPC Socket 权限 ---
	// Bind 会创建 socket 文件，需要确保其他用户（运行 C 程序的）有权连接
	if strings.HasPrefix(ipcEndpoint, "ipc://") {
		socketPath := ipcEndpoint[len("ipc://"):] // 获取文件路径
		// 设置权限为 0666 (所有者读写, 组读写, 其他人读写)
		log.Printf("Attempting to set permissions on %s to 0666", socketPath)
		err = os.Chmod(socketPath, 0666)
		if err != nil {
			// 这可能是一个严重问题，因为 C 程序可能无法连接
			log.Printf("WARN: Failed to change permissions of the IPC socket '%s' to 0666: %v. C clients might fail to connect.", socketPath, err)
		} else {
			fmt.Printf("INFO: Set IPC socket permissions for %s to world-writable (0666)\n", socketPath)
		}
	}

	err = subscriber.SetSubscribe("") // Subscribe to all topics
	if err != nil {
		log.Fatalf("Error subscribing to topics: %v", err)
	}

	msgChan := make(chan RawMessage, 20480)
	var wg sync.WaitGroup

	wg.Add(2)
	go receiver(subscriber, msgChan, &wg)
	go processor(msgChan, &wg)

	fmt.Println("Main: Waiting for goroutines to finish...")
	wg.Wait()

	fmt.Println("Main: All goroutines finished. Exiting.")
}
