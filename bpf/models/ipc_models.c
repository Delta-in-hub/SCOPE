#include "ipc_models.h"

void pack_ipc_model(msgpack_packer *pk, const void *payload) {
    const struct IPC_Model *model = (const struct IPC_Model *)payload;

    // Pack as a msgpack array with 5 elements
    msgpack_pack_array(pk, 5);

    // Pack timestamp
    msgpack_pack_int64(pk, model->nano_since_epoch);

    // Pack process ID
    msgpack_pack_int32(pk, model->pid);

    // Pack process name (comm)
    size_t comm_len = strnlen(model->comm, sizeof(model->comm));
    msgpack_pack_str(pk, comm_len);
    msgpack_pack_str_body(pk, model->comm, comm_len);

    // Pack command line
    if (model->cmdline != NULL) {
        size_t cmdline_len = strlen(model->cmdline);
        msgpack_pack_str(pk, cmdline_len);
        msgpack_pack_str_body(pk, model->cmdline, cmdline_len);
    } else {
        // Empty string if cmdline is NULL
        msgpack_pack_str(pk, 0);
        msgpack_pack_str_body(pk, "", 0);
    }

    // Pack message
    if (model->msg != NULL) {
        size_t msg_len = strlen(model->msg);
        msgpack_pack_str(pk, msg_len);
        msgpack_pack_str_body(pk, model->msg, msg_len);
    } else {
        // Empty string if msg is NULL
        msgpack_pack_str(pk, 0);
        msgpack_pack_str_body(pk, "", 0);
    }
}

struct IPC_Model ipc_model = {.pack = pack_ipc_model};