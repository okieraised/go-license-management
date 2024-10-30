package machine

import "time"

type MachineAttribute struct {
	Fingerprint       *string     `json:"fingerprint"`
	Cores             *int        `json:"cores"`
	IP                *string     `json:"ip"`
	Hostname          *string     `json:"hostname"`
	Platform          *string     `json:"platform"`
	Name              *string     `json:"name"`
	Metadata          interface{} `json:"metadata"`
	MaxProcesses      *string     `json:"max_processes"`
	RequireHeartbeat  *bool       `json:"require_heartbeat"`
	HeartbeatStatus   *string     `json:"heartbeat_status"`
	HeartbeatDuration *int        `json:"heartbeat_duration"`
	LastHeartbeat     time.Time   `json:"last_heartbeat"`
	NextHeartbeat     time.Time   `json:"next_heartbeat"`
	LastCheckOut      time.Time   `json:"last_checkout"`
	Created           time.Time   `json:"created"`
	Updated           time.Time   `json:"updated"`
}
