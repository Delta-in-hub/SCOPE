package models

//! [vfs_open] START

// --- NEW: Struct for vfs_open events (Array Format) ---
// NOTE: Field order MUST match the C packer's array order: [timestamp, pid, comm, filename].
// NOTE: No msgpack tags are used; unmarshaling relies on field order.
type VfsOpenEvent struct {
	TimestampNs int64  // Index 0 in the packed array
	PID         int32  // Index 1 in the packed array
	Comm        string // Index 2 in the packed array
	Filename    string // Index 3 in the packed array
}

//! [vfs_open] END
