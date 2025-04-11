package agentmanager

import (
	"fmt"
	"log"
	"sync"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// Reads from ZMQ socket and sends raw messages to the channel.
func ZMQReceiver(subscriber *zmq.Socket, msgChan chan<- RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(msgChan)

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

	log.Println("Receiver goroutine finished.")
}
