package main

import (
	"fmt"
	"log"
	"scope/internal/models"
	"scope/internal/platform"
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
		case "vfs_open":
			// Use the untagged struct expecting an array format based on field order.
			var event models.VfsOpenEvent
			err := msgpack.Unmarshal(payloadBytes, &event)
			if err != nil {
				// If this error occurs often, double-check C packing order vs Go struct order
				log.Printf("Processor: Error unmarshaling VfsOpenEvent (Array Format): %v", err)
				continue
			}
			cmdline, _ := platform.GetCmdline(event.PID)
			// Format timestamp for better readability
			tsFormatted := time.Unix(0, event.TimestampNs).Format(time.RFC1123)
			// Print received event data - access fields by name as usual
			fmt.Printf("Processed [%s]: Time=%s, PID=%d, Comm='%s', Cmdline='%s', File='%s'\n",
				topic, tsFormatted, event.PID, event.Comm, cmdline, event.Filename)

			tsnow := time.Now().UnixNano()
			fmt.Printf("Cost %d ns in zmq\n", tsnow-event.TimestampNs)

		default:
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
	fmt.Printf("Go Subscriber (Array Format Capable) connecting to %s\n", ipcEndpoint)
	err = subscriber.Connect(ipcEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect subscriber to '%s': %v", ipcEndpoint, err)
	}
	fmt.Println("Subscriber connected.")

	err = subscriber.SetSubscribe("") // Subscribe to all topics
	if err != nil {
		log.Fatalf("Error subscribing to topics: %v", err)
	}

	msgChan := make(chan RawMessage, 10240)
	var wg sync.WaitGroup

	wg.Add(2)
	go receiver(subscriber, msgChan, &wg)
	go processor(msgChan, &wg)

	fmt.Println("Main: Waiting for goroutines to finish...")
	wg.Wait()

	fmt.Println("Main: All goroutines finished. Exiting.")
}
