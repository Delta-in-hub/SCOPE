package models

import "time"

// NodeInfo represents information about an agent node
type NodeInfo struct {
	ID       string            `json:"id"`                // Unique identifier for the agent (machine ID)
	IPs      map[string]string `json:"ips"`               // IP addresses of the agent (interface name -> IP)
	LastSeen time.Time         `json:"last_seen"`         // Last time the agent was seen
	Status   string            `json:"status"`            // Status of the agent (online, offline)
	Token    string            `json:"token,omitempty"`   // Authentication token
	Latency  time.Duration     `json:"latency,omitempty"` // Latency of the agent
}
