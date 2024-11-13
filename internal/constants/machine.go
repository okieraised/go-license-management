package constants

const (
	// MachineActionCheckout - Action to checkout a machine. This will generate a snapshot of the machine at time of checkout,
	// encoded into a machine file certificate that can be decoded and used for licensing offline and
	// air-gapped environments. The algorithm will depend on the license policy's scheme.
	// Machine files can be distributed using email or USB drives to air-gapped devices.
	MachineActionCheckout = "check-out"

	// MachineActionPingHeartbeat - Action to begin or maintain a machine heartbeat monitor. When a machine has not performed a
	// heartbeat ping within the monitor window, it will automatically be deactivated. This can be utilized for machine leasing,
	// where a license has a limited number of machines allowed, and each machine must maintain heartbeat pings in order to remain active.
	MachineActionPingHeartbeat = "ping-heartbeat"

	// MachineActionResetHeartbeat - Action to reset and stop the machine's heartbeat monitor. This will not deactivate the machine.
	MachineActionResetHeartbeat = "reset-heartbeat"
)

var ValidMachineActionsMapper = map[string]bool{
	MachineActionCheckout:       true,
	MachineActionPingHeartbeat:  true,
	MachineActionResetHeartbeat: true,
}
