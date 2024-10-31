package machines

import "time"

// MachineAttributeModel contains information about the machine. Machines can be used to track and manage where your users are allowed to use your product.
type MachineAttributeModel struct {
	Fingerprint       *string     `json:"fingerprint"`        // The fingerprint of the machine. This can be an arbitrary string, but must be unique within the scope of the license it belongs to.
	Cores             *int        `json:"cores"`              // The number of CPU cores for the machine.
	Name              *string     `json:"name"`               // The human-readable name of the machine.
	IP                *string     `json:"ip"`                 // The IP of the machine.
	Hostname          *string     `json:"hostname"`           // The hostname of the machine.
	Platform          *string     `json:"platform"`           // The platform of the machine.
	Metadata          interface{} `json:"metadata"`           // Object containing machine metadata.
	MaxProcesses      *string     `json:"max_processes"`      // The maximum number of processes the machine can have associated with it. Inherited from its license.
	RequireHeartbeat  *bool       `json:"require_heartbeat"`  // Whether the machine requires heartbeat pings, i.e. the policy requires heartbeats, or the machine has an active heartbeat monitor.
	HeartbeatStatus   *string     `json:"heartbeat_status"`   // The status of the machine's heartbeat.
	HeartbeatDuration *int        `json:"heartbeat_duration"` // The policy's heartbeat duration. When a heartbeat monitor is active, the machine must send a heartbeat ping within this timeframe to remain activated.
	LastHeartbeat     time.Time   `json:"last_heartbeat"`     // When the machine last sent a heartbeat ping. This is null if the machine does not require a heartbeat.
	NextHeartbeat     time.Time   `json:"next_heartbeat"`     //The time at which the machine is required to send a heartbeat ping by. This is null if the machine does not require a heartbeat.
	LastCheckOut      time.Time   `json:"last_checkout"`      // When the machine was last checked-out.
	Created           time.Time   `json:"created"`            // When the machine was created.
	Updated           time.Time   `json:"updated"`            // When the machine was last updated.
}
