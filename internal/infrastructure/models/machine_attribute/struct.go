package machine_attribute

// MachineAttributeModel contains information about the machine. Machines can be used to track and manage where your users are allowed to use your product.
type MachineAttributeModel struct {
	Fingerprint  *string     `json:"fingerprint"`   // The fingerprint of the machine. This can be an arbitrary string, but must be unique within the scope of the license it belongs to.
	Cores        *int        `json:"cores"`         // The number of CPU cores for the machine.
	Name         *string     `json:"name"`          // The human-readable name of the machine.
	IP           *string     `json:"ip"`            // The IP of the machine.
	Hostname     *string     `json:"hostname"`      // The hostname of the machine.
	Platform     *string     `json:"platform"`      // The platform of the machine.
	Metadata     interface{} `json:"metadata"`      // Object containing machine metadata.
	MaxProcesses *string     `json:"max_processes"` // The maximum number of processes the machine can have associated with it. Inherited from its license.
}
