package agentmanager

// --- Struct to pass raw messages between goroutines ---
type RawMessage struct {
	Topic   []byte // Received raw topic bytes (might be msgpack encoded)
	Payload []byte // Received raw payload bytes
}

// --- Configuration struct for the application ---
type Config struct {
	Verbose       bool   // Whether to print verbose output
	RedisAddr     string // Redis server address
	RedisDB       int    // Redis database number
	RedisPassword string // Redis password
	StreamKey     string // Redis stream key
	IPCEndpoint   string // ZMQ IPC endpoint
}
