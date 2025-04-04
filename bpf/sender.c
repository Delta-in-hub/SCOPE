#define _GNU_SOURCE
#include <assert.h>
#include <msgpack.h> // MessagePack C library
#include <stdint.h>  // For standard integer types (int32_t)
#include <stdio.h>   // For standard I/O (printf, perror)
#include <stdlib.h>  // For general utilities (rand, srand, exit)
#include <string.h>  // For string operations (strnlen, memcpy, snprintf)
#include <time.h>    // For seeding random number generator (time)
#include <unistd.h>  // For POSIX functions (sleep, usleep)
#include <zmq.h>     // ZeroMQ library

// --- Struct definitions ---
// Defines the structure for command messages
typedef struct {
    int32_t command_id;
    char target_device[32]; // Fixed-size buffer for device name
    double parameter;
} CommandPayload;

// Defines the structure for status update messages
typedef struct {
    int32_t source_id;
    char status_code[16]; // Fixed-size buffer for status code
    char details[128];    // Fixed-size buffer for details
} StatusUpdatePayload;

// --- Constants for Message Types (Topics) ---
#define MSG_TYPE_COMMAND "CMD"
#define MSG_TYPE_STATUS "STAT"

// --- Helper function to pack CommandPayload into msgpack format ---
// Takes a msgpack packer and a pointer to the payload struct
void pack_command(msgpack_packer *pk, const CommandPayload *payload) {
    // Pack as a msgpack array with 3 elements
    msgpack_pack_array(pk, 3);
    msgpack_pack_int32(pk, payload->command_id);

    // Pack the string: length first, then the body
    // Ensure we only pack the actual string length, not the buffer size
    size_t device_len =
        strnlen(payload->target_device, sizeof(payload->target_device));
    msgpack_pack_str(pk, device_len);
    msgpack_pack_str_body(pk, payload->target_device, device_len);

    msgpack_pack_double(pk, payload->parameter);
}

// --- Helper function to pack StatusUpdatePayload into msgpack format ---
// Takes a msgpack packer and a pointer to the payload struct
void pack_status(msgpack_packer *pk, const StatusUpdatePayload *payload) {
    // Pack as a msgpack array with 3 elements
    msgpack_pack_array(pk, 3);
    msgpack_pack_int32(pk, payload->source_id);

    // Pack status_code string
    size_t code_len =
        strnlen(payload->status_code, sizeof(payload->status_code));
    msgpack_pack_str(pk, code_len);
    msgpack_pack_str_body(pk, payload->status_code, code_len);

    // Pack details string
    size_t details_len = strnlen(payload->details, sizeof(payload->details));
    msgpack_pack_str(pk, details_len);
    msgpack_pack_str_body(pk, payload->details, details_len);
}

// --- Helper function to send a multipart message (Topic + Payload) via ZMQ ---
// Takes the ZMQ socket, the topic string, and the msgpack buffer containing the
// payload Returns 0 on success, -1 on error
int send_multipart_topic_msg(void *socket, const char *topic_str,
                             const msgpack_sbuffer *payload_buf) {
    // 1. Send the topic frame
    zmq_msg_t topic_msg;
    size_t topic_len = strnlen(topic_str, 128);
    // Initialize ZMQ message with specific size
    if (zmq_msg_init_size(&topic_msg, topic_len) != 0) {
        perror("zmq_msg_init_size (topic) failed");
        return -1;
    }
    // Copy topic string into the ZMQ message
    memcpy(zmq_msg_data(&topic_msg), topic_str, topic_len);
    // Send the message part, indicating more parts will follow (ZMQ_SNDMORE)
    int rc = zmq_msg_send(&topic_msg, socket, ZMQ_SNDMORE);
    zmq_msg_close(
        &topic_msg); // Always close the message after sending or error
    if (rc == -1) {
        perror("zmq_msg_send (topic) failed");
        return -1;
    }

    // 2. Send the payload frame
    zmq_msg_t payload_msg;
    // Initialize ZMQ message with the size of the packed data in sbuffer
    if (zmq_msg_init_size(&payload_msg, payload_buf->size) != 0) {
        perror("zmq_msg_init_size (payload) failed");
        // Note: Topic frame was already sent!
        return -1;
    }
    // Copy the packed data from sbuffer into the ZMQ message
    memcpy(zmq_msg_data(&payload_msg), payload_buf->data, payload_buf->size);
    // Send the final message part (flag = 0)
    rc = zmq_msg_send(&payload_msg, socket, 0);
    zmq_msg_close(&payload_msg); // Always close the message
    if (rc == -1) {
        perror("zmq_msg_send (payload) failed");
        // Note: Topic frame was already sent!
        return -1;
    }

    return 0; // Success
}

// --- Main Program ---
int main() {
    // --- ZMQ Context and Socket Setup (PUB) ---
    printf("Initializing ZeroMQ context...\n");
    void *context = zmq_ctx_new();
    if (!context) {
        perror("zmq_ctx_new failed");
        return 1;
    }

    printf("Creating ZeroMQ PUB socket...\n");
    void *publisher = zmq_socket(context, ZMQ_PUB); // Use PUB socket type
    if (!publisher) {
        perror("zmq_socket (PUB) failed");
        zmq_ctx_destroy(context);
        return 1;
    }

    // Define the IPC endpoint address
    const char *ipc_endpoint = "ipc:///tmp/zmq_ipc_pubsub.sock";
    printf("C Publisher (Multi-Type) binding to %s\n", ipc_endpoint);
    // Bind the PUB socket to the endpoint
    int rc = zmq_bind(publisher, ipc_endpoint);
    if (rc != 0) {
        perror("zmq_bind failed");
        zmq_close(publisher);
        zmq_ctx_destroy(context);
        return 1;
    }

    // Allow time for potential subscribers to connect and set up subscriptions.
    // This is important in PUB/SUB to avoid losing the first few messages (slow
    // joiner syndrome).
    printf("Publisher bound. Waiting 1 second for subscribers...\n");
    sleep(1);

    // --- MessagePack Buffer and Packer Initialization (Outside the loop) ---
    // Initialize the reusable buffer for storing serialized data
    msgpack_sbuffer sbuf;
    msgpack_sbuffer_init(&sbuf);

    // Initialize the reusable packer and associate it with the buffer
    msgpack_packer pk;
    msgpack_packer_init(&pk, &sbuf,
                        msgpack_sbuffer_write); // Uses sbuf internally

    // Seed the random number generator
    srand(time(NULL));

    printf("Starting to publish different types of data (reusing buffer)...\n");

    // --- Main Publishing Loop ---
    for (int i = 0; i < 10; ++i) {
        // NOTE: sbuf and pk are NOT re-initialized inside the loop

        const char *topic = NULL; // Variable to hold the message topic

        // Alternate between sending Command and Status messages
        if (i % 2 == 0) {
            // --- Prepare and Pack Command ---
            topic = MSG_TYPE_COMMAND;
            CommandPayload cmd;
            cmd.command_id = 1000 + i;
            // Safely format the device name string
            snprintf(cmd.target_device, sizeof(cmd.target_device), "Sensor_%d",
                     i / 2);
            cmd.target_device[sizeof(cmd.target_device) - 1] =
                '\0'; // Ensure null termination
            cmd.parameter =
                (double)rand() / RAND_MAX * 10.0; // Example parameter value

            // Pack the command payload into the reusable buffer (sbuf) via the
            // packer (pk)
            pack_command(&pk, &cmd);

            printf("Publishing [%s]: ID=%d, Target='%s', Param=%.2f ", topic,
                   cmd.command_id, cmd.target_device, cmd.parameter);

        } else {
            // --- Prepare and Pack Status ---
            topic = MSG_TYPE_STATUS;
            StatusUpdatePayload stat;
            stat.source_id = 2000 + i;
            // Example status codes
            snprintf(stat.status_code, sizeof(stat.status_code),
                     (i % 4 == 1) ? "OK" : "PENDING");
            stat.status_code[sizeof(stat.status_code) - 1] =
                '\0'; // Ensure null termination
            // Example details string
            snprintf(stat.details, sizeof(stat.details),
                     "Status details update for sequence %d", i);
            stat.details[sizeof(stat.details) - 1] =
                '\0'; // Ensure null termination

            // Pack the status payload into the reusable buffer (sbuf) via the
            // packer (pk)
            pack_status(&pk, &stat);

            // Print limited details to keep output concise
            printf("Publishing [%s]: SrcID=%d, Code='%s', Details='%.30s...' ",
                   topic, stat.source_id, stat.status_code, stat.details);
        }

        // At this point, 'sbuf' contains the packed data for the current
        // message
        printf("(Packed size: %zu bytes)\n", sbuf.size);

        // --- Send the Multipart Message (Topic + Payload) ---
        rc = send_multipart_topic_msg(publisher, topic, &sbuf);

        // --- CRITICAL OPTIMIZATION: Clear the buffer for the next iteration
        // --- This resets sbuf.size to 0 but keeps the allocated memory
        // (sbuf.data) ready for reuse, avoiding repeated malloc/free.
        msgpack_sbuffer_clear(&sbuf);

        // Check if sending failed
        if (rc != 0) {
            fprintf(stderr, "Failed to send message %d. Exiting loop.\n", i);
            break; // Exit the loop on send error
        }

        // Pause briefly between messages
        usleep(500000); // 500 milliseconds
    }

    printf("Finished publishing data loop.\n");

    // --- Cleanup ---
    // Destroy the msgpack buffer (freeing its internally allocated memory)
    // This is done ONCE after the loop finishes.
    printf("Destroying msgpack sbuffer...\n");
    msgpack_sbuffer_destroy(&sbuf);

    // Close the ZMQ socket
    printf("Closing publisher socket...\n");
    zmq_close(publisher);

    // Destroy the ZMQ context
    printf("Destroying ZeroMQ context...\n");
    zmq_ctx_destroy(context);

    printf("Publisher finished cleanly.\n");
    return 0; // Indicate successful execution
}