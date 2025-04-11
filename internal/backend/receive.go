package backend // Assuming this code resides in a 'backend' package

import (
	"context"
	"database/sql" // Import sql package for Null types
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"                  // Ensure pq driver is registered
	goredis "github.com/redis/go-redis/v9" // Import pq for potential error checking if needed, though not strictly required by the rewrite logic itself
)

const (
	redisStreamKey = "SCOPE_STREAM" // Should match producer config
	// Consumer group name
	consumerGroup = "backend-consumers"
	// Consumer name prefix (will be appended with a unique identifier)
	consumerNamePrefix = "consumer"
	// Number of messages to read in a single batch
	batchSize = 100 // Increased batch size for potentially better throughput
	// Maximum wait time for reading from stream
	readTimeout = 2 * time.Second // Slightly longer block time
)

var (
	ackedMessageIDs     map[string]bool
	ackedMessageIDsLock sync.Mutex
)

// Receive reads messages from Redis Stream and inserts them into TimescaleDB.
func Receive(ctx context.Context, wg *sync.WaitGroup, tsdb *sqlx.DB, redisClient *goredis.Client, verbose bool, consumerID int) {
	defer wg.Done()

	ackedMessageIDsLock.Lock()
	if ackedMessageIDs == nil {
		ackedMessageIDs = make(map[string]bool)
	}
	ackedMessageIDsLock.Unlock()

	// Generate a unique consumer name
	consumerName := fmt.Sprintf("%s-%d", consumerNamePrefix, consumerID)
	if verbose {
		log.Printf("Starting Redis Stream consumer (%s) for group (%s) on stream (%s)\n", consumerName, consumerGroup, redisStreamKey)
	}

	// Create consumer group if it doesn't exist (errors ignored if BUSYGROUP)
	err := redisClient.XGroupCreateMkStream(ctx, redisStreamKey, consumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		// Log the error but continue, maybe another consumer created it.
		// If stream doesn't exist, XReadGroup will fail later.
		log.Printf("Warning: Error creating/checking consumer group '%s' on stream '%s': %v", consumerGroup, redisStreamKey, err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled for consumer %s, stopping.", consumerName)
			return
		default:
			// Read new messages using the consumer group. '>' means only new messages.
			streams, err := redisClient.XReadGroup(ctx, &goredis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerName,
				Streams:  []string{redisStreamKey, ">"}, // Read only new messages for this consumer
				Count:    batchSize,
				Block:    readTimeout,
				// NoAck: false, // We will manually ACK after successful processing
			}).Result()

			if err != nil {
				// redis.Nil means timeout, which is expected when no new messages
				if err == goredis.Nil {
					continue // No new messages, loop again
				}
				// Log other errors
				log.Printf("Error reading from Redis Stream for consumer %s: %v", consumerName, err)
				// Optional: Add a small delay before retrying on persistent errors
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Process messages if any were received
			for _, stream := range streams {
				if len(stream.Messages) > 0 {
					if verbose {
						log.Printf("Consumer %s received %d messages from stream %s", consumerName, len(stream.Messages), stream.Stream)
					}
					// Process the batch and get IDs that were successfully processed (or attempted)
					processMessages(ctx, tsdb, stream.Messages, verbose)

					// Acknowledge successfully processed messages
					ackedMessageIDsLock.Lock()
					for _, msg := range stream.Messages {
						ackedMessageIDs[msg.ID] = true
					}
					ackedMessageIDsLock.Unlock()
				}
			}
		}
	}
}

// processMessages processes a batch of messages from Redis Stream and inserts them into TimescaleDB.
// It returns a slice of message IDs that were processed (successfully or unsuccessfully attempted within the transaction).
func processMessages(ctx context.Context, tsdb *sqlx.DB, messages []goredis.XMessage, verbose bool) {
	var tx *sqlx.Tx
	var err error // Declare err outside loop for deferred rollback check

	// Helper function to safely get string from interface{}
	getString := func(data map[string]interface{}, key string) (string, bool) {
		val, ok := data[key]
		if !ok {
			return "", false
		}
		strVal, ok := val.(string)
		return strVal, ok
	}

	// Helper function to safely get int64 from interface{} (handles float64 from JSON)
	getInt64 := func(data map[string]interface{}, key string) (int64, bool) {
		val, ok := data[key]
		if !ok {
			return 0, false
		}
		switch v := val.(type) {
		case float64: // JSON numbers are often decoded as float64
			return int64(v), true
		case int:
			return int64(v), true
		case int32:
			return int64(v), true
		case int64:
			return v, true
		default:
			if verbose {
				log.Printf("Invalid type for key %s: %T", key, val)
			}
			return 0, false
		}
	}
	// Helper function to safely get int32 from interface{} (handles float64 from JSON)
	getInt32 := func(data map[string]interface{}, key string) (int32, bool) {
		val, ok := data[key]
		if !ok {
			return 0, false
		}
		switch v := val.(type) {
		case float64: // JSON numbers are often decoded as float64
			return int32(v), true
		case int:
			return int32(v), true
		case int32:
			return v, true
		case int64:
			// Be careful about potential overflow if the value is large
			if v >= -2147483648 && v <= 2147483647 {
				return int32(v), true
			}
			return 0, false // Or handle overflow error
		default:
			if verbose {
				log.Printf("Invalid type for key %s: %T", key, val)
			}
			return 0, false
		}
	}

	// Helper function to create sql.NullString
	getNullString := func(data map[string]interface{}, key string) sql.NullString {
		val, ok := getString(data, key)
		if ok && val != "" { // Treat empty string as NULL too? Adjust if needed.
			return sql.NullString{String: val, Valid: true}
		}
		return sql.NullString{Valid: false}
	}

	// Helper function to create sql.NullInt64
	getNullInt64 := func(data map[string]interface{}, key string) sql.NullInt64 {
		val, ok := getInt64(data, key)
		if ok {
			return sql.NullInt64{Int64: val, Valid: true}
		}
		return sql.NullInt64{Valid: false}
	}

	// Helper function to create sql.NullInt32
	getNullInt32 := func(data map[string]interface{}, key string) sql.NullInt32 {
		val, ok := getInt32(data, key)
		if ok {
			return sql.NullInt32{Int32: val, Valid: true}
		}
		return sql.NullInt32{Valid: false}
	}

	// Begin transaction
	tx, err = tsdb.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		// Cannot process any messages in this batch
		return
	}
	// Use a defer func to handle rollback/commit
	defer func() {
		if p := recover(); p != nil {
			log.Printf("Panic recovered during transaction: %v", p)
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Printf("Error during transaction rollback: %v", rbErr)
			}
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				log.Printf("Error committing transaction: %v", commitErr)
				// Potentially mark messages as not fully processed if commit fails
			} else if verbose {
				log.Printf("Transaction committed successfully for %d msgs.", len(messages))
			}
		}
	}()

	// --- Prepare statements within the transaction ---
	// Note: Preparing statements repeatedly in a loop is less efficient than preparing once outside,
	// but preparing inside the transaction ensures they are valid for that transaction context.
	// For high-throughput, consider preparing once and passing `tx` to `Stmt.ExecContext`.

	// Prepare events_os statement
	osStmt, err := tx.PreparexContext(ctx, `
		INSERT INTO events_os (
			ts, machine_id, event_subtype, pid, comm, cmdline, vfs_filename,
			syscall_name, cpu, sched_type, ppid, ppid_comm, ppid_cmdline,
			exec_filename, exec_args
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`)
	if err != nil {
		log.Printf("Error preparing OS statement: %v", err)
		return // Cannot proceed
	}
	defer osStmt.Close() // Close prepared statement when function exits

	// Prepare events_cuda statement
	cudaStmt, err := tx.PreparexContext(ctx, `
		INSERT INTO events_cuda (
			ts, machine_id, event_subtype, pid, comm, cmdline, operation,
			cuda_ptr, cuda_size, cuda_retval, cuda_func_ptr, cuda_symbol_name,
			cuda_symbol_file, cuda_symbol_offset, cuda_symbol_sourcefile,
			cuda_memcpy_src, cuda_memcpy_dst, cuda_memcpy_kind, cuda_memcpy_type,
			cuda_sync_duration_ns
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`)
	if err != nil {
		log.Printf("Error preparing CUDA statement: %v", err)
		return // Cannot proceed
	}
	defer cudaStmt.Close()

	// Prepare events_ggml statement
	ggmlStmt, err := tx.PreparexContext(ctx, `
		INSERT INTO events_ggml (
			ts, machine_id, event_subtype, pid, comm, cmdline, operation,
			ggml_cuda_func_name, ggml_cuda_duration_ns, ggml_graph_size,
			ggml_graph_nodes, ggml_graph_leafs, ggml_graph_order, ggml_cost_ns,
			ggml_mem_size, ggml_mem_ptr
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`)
	if err != nil {
		log.Printf("Error preparing GGML statement: %v", err)
		return // Cannot proceed
	}
	defer ggmlStmt.Close()

	// Prepare events_app_log statement
	appLogStmt, err := tx.PreparexContext(ctx, `
		INSERT INTO events_app_log (
			ts, machine_id, event_subtype, pid, comm, cmdline, log_text
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		log.Printf("Error preparing AppLog statement: %v", err)
		return // Cannot proceed
	}
	defer appLogStmt.Close()

	// --- Process each message in the batch ---
	for _, msg := range messages {
		// Add message ID to processed list early, even if insertion fails,
		// so we ACK it and don't re-process endlessly on non-transient errors.
		// Adjust this logic if you need retries for certain errors.

		dataStr, ok := msg.Values["data"].(string)
		if !ok {
			log.Printf("Error: message data is not a string, skipping message ID: %s", msg.ID)
			continue // Skip this message, move to the next
		}

		var eventData map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(dataStr), &eventData); jsonErr != nil {
			log.Printf("Error parsing JSON data: %v, skipping message ID: %s", jsonErr, msg.ID)
			continue // Skip this message
		}

		// --- Extract common fields and perform essential conversions ---
		topic, topicOk := getString(eventData, "topic")
		timestampNs, tsOk := getInt64(eventData, "timestamp")
		machineID, machineIDOk := getString(eventData, "machineid")
		pid, pidOk := getInt32(eventData, "pid") // Use getInt32
		comm := getNullString(eventData, "comm")
		cmdline := getNullString(eventData, "cmdline")

		// Basic validation: topic, timestamp, machineID, pid are usually essential
		if !topicOk || !tsOk || !machineIDOk || !pidOk {
			log.Printf("Error: Missing essential common fields (topic, timestamp, machineid, pid) in message ID: %s, skipping", msg.ID)
			continue
		}

		// Convert timestamp to time.Time (UTC)
		ts := time.Unix(0, timestampNs).UTC()

		// --- Insert into appropriate table based on topic ---
		switch topic {
		// OS events
		case "vfs_open", "syscalls", "sched", "execv":
			vfsFilename := getNullString(eventData, "filename")     // Used by vfs_open
			syscallName := getNullString(eventData, "syscall")      // Used by syscalls
			cpu := getNullInt32(eventData, "cpu")                   // Used by sched
			schedType := getNullString(eventData, "type")           // Used by sched (processor already converts to string)
			ppid := getNullInt32(eventData, "ppid")                 // Used by execv
			ppidComm := getNullString(eventData, "ppid_comm")       // Used by execv
			ppidCmdline := getNullString(eventData, "ppid_cmdline") // Used by execv
			// execv also uses "filename" and "args", map them specifically
			execFilename := sql.NullString{Valid: false}
			execArgs := sql.NullString{Valid: false}
			if topic == "execv" {
				execFilename = getNullString(eventData, "filename")
				execArgs = getNullString(eventData, "args")
				// Clear vfsFilename if it's an execv event to avoid confusion
				vfsFilename = sql.NullString{Valid: false}
			}

			_, err = osStmt.ExecContext(ctx,
				ts, machineID, topic, int(pid), comm, cmdline, // Common fields first
				vfsFilename,    // vfs_open specific
				syscallName,    // syscalls specific
				cpu, schedType, // sched specific
				ppid, ppidComm, ppidCmdline, // execv specific (ppid info)
				execFilename, execArgs, // execv specific (exec info)
			)
			if err != nil {
				log.Printf("Error inserting OS event (topic: %s, msgID: %s): %v", topic, msg.ID, err)
				// Continue processing other messages in the batch, transaction will be rolled back later
			}

		// CUDA events
		case "cudaMalloc", "cudaFree", "cudaLaunchKernel", "cudaMemcpy", "cudaDeviceSynchronize":
			operation := getNullString(eventData, "operation") // Should be set by processor

			// Handle different keys for pointer based on topic
			cudaPtr := getNullInt64(eventData, "ptr") // Default key, used by cudaFree

			// Handle different keys for size based on topic
			cudaSize := getNullInt64(eventData, "size") // Used by cudaMalloc, cudaMemcpy

			cudaRetval := getNullInt32(eventData, "retval") // Used by cudaMalloc (Processor maps to int64, but schema is INT)

			cudaFuncPtr := getNullInt64(eventData, "func_ptr") // Used by cudaLaunchKernel

			// Symbol info (likely only present for cudaLaunchKernel)
			cudaSymbolName := getNullString(eventData, "symbol_name")
			cudaSymbolFile := getNullString(eventData, "symbol_file")
			cudaSymbolOffset := getNullInt64(eventData, "symbol_offset")
			cudaSymbolSourcefile := getNullString(eventData, "symbol_sourcefile")

			// Memcpy specific
			cudaMemcpySrc := getNullInt64(eventData, "src")
			cudaMemcpyDst := getNullInt64(eventData, "dst")
			cudaMemcpyKind := getNullInt32(eventData, "kind")  // Processor maps to int32, schema is INT
			cudaMemcpyType := getNullString(eventData, "type") // Processor adds this string type

			// Sync specific
			cudaSyncDurationNs := getNullInt64(eventData, "duration_ns")

			_, err = cudaStmt.ExecContext(ctx,
				ts, machineID, topic, int(pid), comm, cmdline, operation, // Common + operation
				cudaPtr, cudaSize, cudaRetval, cudaFuncPtr, // Malloc, Free, LaunchKernel specifics
				cudaSymbolName, cudaSymbolFile, cudaSymbolOffset, cudaSymbolSourcefile, // LaunchKernel specifics
				cudaMemcpySrc, cudaMemcpyDst, cudaMemcpyKind, cudaMemcpyType, // Memcpy specifics
				cudaSyncDurationNs, // Sync specific
			)
			if err != nil {
				log.Printf("Error inserting CUDA event (topic: %s, msgID: %s): %v", topic, msg.ID, err)
			}

		// GGML events
		case "ggml_cuda", "ggml_graph_compute", "ggml_base":
			operation := getNullString(eventData, "operation") // Set by processor

			// ggml_cuda specific
			ggmlCudaFuncName := getNullString(eventData, "func_name")
			ggmlCudaDurationNs := getNullInt64(eventData, "duration_ns")

			// ggml_graph_compute specific
			ggmlGraphSize := getNullInt32(eventData, "graph_size")    // Processor uses int32, schema INT
			ggmlGraphNodes := getNullInt32(eventData, "graph_nodes")  // Processor uses int32, schema INT
			ggmlGraphLeafs := getNullInt32(eventData, "graph_leafs")  // Processor uses int32, schema INT
			ggmlGraphOrder := getNullString(eventData, "graph_order") // Processor converts int to string
			ggmlCostNs := getNullInt64(eventData, "cost_ns")

			// ggml_base specific (Mapped to ggml_mem_* columns)
			ggmlMemSize := getNullInt64(eventData, "size") // Map eventData["size"] to ggml_mem_size
			ggmlMemPtr := getNullInt64(eventData, "ptr")   // Map eventData["ptr"] to ggml_mem_ptr

			_, err = ggmlStmt.ExecContext(ctx,
				ts, machineID, topic, int(pid), comm, cmdline, operation, // Common + operation
				ggmlCudaFuncName, ggmlCudaDurationNs, // ggml_cuda specific
				ggmlGraphSize, ggmlGraphNodes, ggmlGraphLeafs, ggmlGraphOrder, ggmlCostNs, // ggml_graph_compute specific
				ggmlMemSize, ggmlMemPtr, // ggml_base specific
			)
			if err != nil {
				log.Printf("Error inserting GGML event (topic: %s, msgID: %s): %v", topic, msg.ID, err)
			}

		// App Log events
		case "llamaLog": // Assuming this is the only topic for app logs for now
			logText := getNullString(eventData, "text")

			_, err = appLogStmt.ExecContext(ctx,
				ts, machineID, topic, int(pid), comm, cmdline, // Common fields
				logText, // AppLog specific
			)
			if err != nil {
				log.Printf("Error inserting App Log event (topic: %s, msgID: %s): %v", topic, msg.ID, err)
			}

		default:
			if verbose {
				log.Printf("Unknown event topic '%s' encountered in message ID: %s, skipping insertion.", topic, msg.ID)
			}
			// Optionally, insert into a 'dead-letter' table or log more permanently
		}
	}

}

func XDelMessages(ctx context.Context, redisClient *goredis.Client, verbose bool) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	if verbose {
		log.Println("Starting XDelMessages process...")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msgIDs := []string{}
			ackedMessageIDsLock.Lock()
			if verbose && len(ackedMessageIDs) > 0 {
				log.Printf("AckedMessageIDs... len: %d", len(ackedMessageIDs))
			}
			for id := range ackedMessageIDs {
				msgIDs = append(msgIDs, id)
			}
			clear(ackedMessageIDs)
			ackedMessageIDsLock.Unlock()

			if len(msgIDs) > 0 {
				_, err := redisClient.XDel(ctx, redisStreamKey, msgIDs...).Result()
				if err != nil {
					log.Printf("Error deleting messages from stream: %v", err)
				}
				if verbose {
					log.Printf("Deleted %d messages from stream %s", len(msgIDs), redisStreamKey)
				}
			}
		}
	}

}
