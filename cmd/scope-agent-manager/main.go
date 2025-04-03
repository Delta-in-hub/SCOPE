package main

import (
	"fmt"
	"log"
	"sync" // Import sync package for WaitGroup
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

// --- Struct definitions (identical to before) ---
type CommandPayload struct {
	_msgpack     struct{} `msgpack:",array"`
	CommandID    int32
	TargetDevice string
	Parameter    float64
}

type StatusUpdatePayload struct {
	_msgpack   struct{} `msgpack:",array"`
	SourceID   int32
	StatusCode string
	Details    string
}

const (
	MsgTypeCommand = "CMD"
	MsgTypeStatus  = "STAT"
)

// --- End Struct definitions ---

// --- Struct to pass raw messages between goroutines ---
type RawMessage struct {
	Topic   []byte
	Payload []byte
}

// --- Receiver Goroutine ---
// Reads from ZMQ socket and sends raw messages to the channel.
func receiver(subscriber *zmq.Socket, msgChan chan<- RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()      // Signal WaitGroup when this goroutine exits
	defer close(msgChan) // Close the channel when receiver exits to signal processor

	fmt.Println("Receiver goroutine started.")

	// Poller Setup within the goroutine
	poller := zmq.NewPoller()
	poller.Add(subscriber, zmq.POLLIN)

	running := true
	for running {
		// Poll the sockets with a timeout
		// Using a timeout allows checking for context termination implicitly via ETERM error
		polledSockets, err := poller.Poll(250 * time.Millisecond) // Poll longer, less busy-wait
		if err != nil {
			if zmq.AsErrno(err) == zmq.ETERM {
				fmt.Println("Receiver: Context terminated during poll, exiting.")
				running = false // Exit loop cleanly
				continue
			}
			if zmq.AsErrno(err) == zmq.Errno(syscall.EINTR) {
				fmt.Println("Receiver: Poll interrupted, continuing...")
				continue // Interrupted by signal, loop again
			}
			// Log other polling errors but try to continue
			log.Printf("Receiver: Error polling socket: %v (errno %d)", err, zmq.AsErrno(err))
			// Maybe add a small delay here if errors persist
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// If Poll returned any ready sockets
		if len(polledSockets) > 0 {
			// Receive message parts (Topic + Payload)
			msgParts, err := subscriber.RecvMessageBytes(0) // Blocking receive is fine after successful poll
			if err != nil {
				if zmq.AsErrno(err) == zmq.ETERM {
					fmt.Println("Receiver: Context terminated during receive, exiting.")
					running = false // Exit loop cleanly
					continue
				}
				// Handle other potential receive errors
				log.Printf("Receiver: Error receiving message: %v", err)
				continue // Skip this message and try again
			}

			// Validate message structure
			if len(msgParts) != 2 {
				log.Printf("Receiver: Error: Received message with %d parts, expected 2 (Topic, Payload)", len(msgParts))
				continue // Skip malformed message
			}

			// Create RawMessage and send it on the channel
			// This might block if the channel buffer is full
			msg := RawMessage{
				Topic:   msgParts[0],
				Payload: msgParts[1],
			}
			// Use a select with a timeout or default case if you want to prevent
			// the receiver from blocking indefinitely if the processor is stuck.
			// For simplicity now, we assume the processor keeps up or the buffer handles bursts.
			select {
			case msgChan <- msg:
				// Message successfully sent
				// Example of non-blocking send (might drop messages if processor is slow):
				// default:
				//  log.Println("Receiver: Warning - Processor channel full, dropping message.")
			}
		} // End if len(polledSockets) > 0
	} // End for running

	fmt.Println("Receiver goroutine finished.")
}

// --- Processor Goroutine ---
// Reads raw messages from the channel, unmarshals, and prints.
func processor(msgChan <-chan RawMessage, wg *sync.WaitGroup) {
	defer wg.Done() // Signal WaitGroup when this goroutine exits
	fmt.Println("Processor goroutine started.")

	// Loop reading from the channel. The loop automatically exits when msgChan is closed.
	for rawMsg := range msgChan {
		topic := string(rawMsg.Topic)
		payloadBytes := rawMsg.Payload

		// Use a switch to handle different topics
		switch topic {
		case MsgTypeCommand:
			var cmd CommandPayload
			err := msgpack.Unmarshal(payloadBytes, &cmd)
			if err != nil {
				log.Printf("Processor: Error unmarshaling CommandPayload: %v", err)
			} else {
				// Print received command data
				fmt.Printf("Processed [%s]: ID=%d, Target='%s', Param=%.2f\n",
					topic, cmd.CommandID, cmd.TargetDevice, cmd.Parameter)
				// TODO: Add actual command processing logic here
			}

		case MsgTypeStatus:
			var stat StatusUpdatePayload
			err := msgpack.Unmarshal(payloadBytes, &stat)
			if err != nil {
				log.Printf("Processor: Error unmarshaling StatusUpdatePayload: %v", err)
			} else {
				// Print received status data
				fmt.Printf("Processed [%s]: SrcID=%d, Code='%s', Details='%.20s...'\n",
					topic, stat.SourceID, stat.StatusCode, stat.Details)
				// TODO: Add actual status processing logic here
			}

		default:
			log.Printf("Processor: Error: Received message with unexpected topic '%s'", topic)
		}
	} // End for range msgChan

	fmt.Println("Processor goroutine finished (channel closed).")
}

// --- Main Function ---
func main() {
	// ZMQ Context and Socket Setup (SUB)
	context, err := zmq.NewContext()
	if err != nil {
		log.Fatalf("Error creating ZMQ context: %v", err)
	}
	// Use defer context.Term() to ensure it's called even on panic,
	// but we will also call it explicitly during graceful shutdown.
	defer context.Term()

	subscriber, err := context.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatalf("Error creating ZMQ subscriber socket: %v", err)
	}
	defer subscriber.Close()

	ipcEndpoint := "ipc:///tmp/zmq_ipc_pubsub.sock"
	fmt.Printf("Go Subscriber (Multi-Type Concurrent) connecting to %s\n", ipcEndpoint)

	err = subscriber.Connect(ipcEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	fmt.Println("Connection initiated.")

	// Subscribe to all topics ("" prefix)
	err = subscriber.SetSubscribe("")
	if err != nil {
		log.Fatalf("Error subscribing to topics: %v", err)
	}
	fmt.Println("Subscribed to all topics. Waiting for messages...")

	// --- Channel for communication between goroutines ---
	// Use a buffered channel to decouple receiver and processor slightly
	// Adjust buffer size based on expected load and processing time
	const channelBufferSize = 100
	msgChan := make(chan RawMessage, channelBufferSize)

	// --- WaitGroup for synchronizing goroutine shutdown ---
	var wg sync.WaitGroup

	// --- Start Goroutines ---
	wg.Add(2) // Expect two goroutines to finish
	go receiver(subscriber, msgChan, &wg)
	go processor(msgChan, &wg)

	// Wait for all goroutines to finish.
	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()

	fmt.Println("All goroutines finished. Exiting.")
	// Deferred subscriber.Close() and context.Term() will run here if not already called
}
