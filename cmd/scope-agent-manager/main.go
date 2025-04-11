package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"scope/database/redis"
	"scope/internal/agentmanager"
	"scope/internal/utils"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	// Load environment variables
	godotenv.Load(".env")

	// Parse command line arguments
	config := agentmanager.Config{
		IPCEndpoint:   "ipc:///tmp/zmq_ipc_pubsub.sock",
		RedisAddr:     utils.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisDB:       1, // 1 for stream message queue
		RedisPassword: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
		StreamKey:     "SCOPE_STREAM",
	}

	// Define command line flags
	verboseFlag := flag.Bool("verbose", false, "Enable verbose output")
	redisAddrFlag := flag.String("redis-addr", config.RedisAddr, "Redis server address")
	redisDBFlag := flag.Int("redis-db", config.RedisDB, "Redis database number")
	redisPasswordFlag := flag.String("redis-password", config.RedisPassword, "Redis password")
	streamKeyFlag := flag.String("stream-key", config.StreamKey, "Redis stream key")
	ipcEndpointFlag := flag.String("ipc-endpoint", config.IPCEndpoint, "ZMQ IPC endpoint")

	// Parse flags
	flag.Parse()

	// Update config with command line arguments
	config.Verbose = *verboseFlag
	config.RedisAddr = *redisAddrFlag
	config.RedisDB = *redisDBFlag
	config.RedisPassword = *redisPasswordFlag
	config.StreamKey = *streamKeyFlag
	config.IPCEndpoint = *ipcEndpointFlag

	// Initialize Redis client
	redisConfig := redis.Config{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}

	redisClient, err := redis.NewClient(redisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	if config.Verbose {
		log.Printf("Connected to Redis at %s, using database %d", config.RedisAddr, config.RedisDB)
		log.Printf("Using Redis stream key: %s", config.StreamKey)
	}

	// Initialize ZMQ
	zmqContext, err := zmq.NewContext()
	if err != nil {
		log.Fatalf("Error creating ZMQ context: %v", err)
	}
	defer zmqContext.Term()

	subscriber, err := zmqContext.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatalf("Error creating ZMQ subscriber socket: %v", err)
	}
	defer subscriber.Close()

	fmt.Printf("Go Subscriber binding to %s\n", config.IPCEndpoint)

	err = subscriber.Bind(config.IPCEndpoint)
	if err != nil {
		log.Fatalf("Failed to bind subscriber to '%s': %v", config.IPCEndpoint, err)
	}

	// Set IPC Socket permissions
	if strings.HasPrefix(config.IPCEndpoint, "ipc://") {
		socketPath := config.IPCEndpoint[len("ipc://"):] // Get file path
		// Set permissions to 0666 (owner read/write, group read/write, others read/write)
		err = os.Chmod(socketPath, 0666)
		if err != nil {
			log.Printf("WARN: Failed to change permissions of the IPC socket '%s' to 0666: %v. C clients might fail to connect.", socketPath, err)
		} else if config.Verbose {
			fmt.Printf("INFO: Set IPC socket permissions for %s to world-writable (0666)\n", socketPath)
		}
	}

	err = subscriber.SetSubscribe("") // Subscribe to all topics
	if err != nil {
		log.Fatalf("Error subscribing to topics: %v", err)
	}

	// Create a buffered channel for message passing
	// Using a large buffer to handle high message rates
	msgChan := make(chan agentmanager.RawMessage, 20480)
	var wg sync.WaitGroup

	// Calculate optimal number of processor goroutines based on CPU cores
	numProcessors := runtime.NumCPU() / 2
	if numProcessors < 1 {
		numProcessors = 1
	}

	// Start processor goroutines
	wg.Add(numProcessors)
	for range numProcessors {
		go agentmanager.Processor(msgChan, &wg, config, redisClient)
	}

	// Start receiver goroutine
	wg.Add(1)
	go agentmanager.ZMQReceiver(subscriber, msgChan, &wg)

	if config.Verbose {
		fmt.Printf("Starting %d processor goroutines...\n", numProcessors)
	}

	port := utils.GetEnvOrDefault("AGENT_PORT", "18090")
	chi := agentmanager.SetupRouter()
	myips := utils.GetMyIpAddrs()
	for _, ip := range myips {
		log.Printf("Starting agent manager on ip http://%s:%s\n", ip, port)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		centerURL := utils.GetEnvOrDefault("CENTER_URL", "http://localhost:18080")
		for {
			token, err := agentmanager.RegisterNodeToCenter(centerURL)
			if err != nil {
				log.Printf("Failed to register node to center: %v", err)
				time.Sleep(5 * time.Second)
			} else {
				log.Printf("Successfully registered node to center %s , token: %s", centerURL, token)
				break
			}
		}
	}(&wg)

	log.Fatal(http.ListenAndServe(":"+port, chi))

	// Wait for all goroutines to complete
	wg.Wait()
}
