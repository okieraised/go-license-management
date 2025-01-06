package constants

const (
	// MachineActionCheckout - Action to check out a machine. This will generate a snapshot of the machine at time of checkout,
	// encoded into a machine file certificate that can be decoded and used for licensing offline and air-gapped environments.
	MachineActionCheckout = "check-out"

	// MachineActionPingHeartbeat - Action to ping server to announce machine's alive status
	MachineActionPingHeartbeat = "ping-heartbeat"

	// MachineActionResetHeartbeat - Action to reset and stop the machine's heartbeat monitor
	MachineActionResetHeartbeat = "reset-heartbeat"
)

var ValidMachineActionsMapper = map[string]bool{
	MachineActionCheckout:       true,
	MachineActionPingHeartbeat:  true,
	MachineActionResetHeartbeat: true,
}
