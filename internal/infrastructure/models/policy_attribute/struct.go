package policy_attribute

type PolicyAttributeModel struct {
	Name                          *string                `json:"name" validate:"required"`                            // Name: name of the policy
	Scheme                        *string                `json:"scheme" validate:"optional"`                          // Scheme: The encryption/signature scheme used on license keys.
	Strict                        *bool                  `json:"strict" validate:"optional"`                          // Strict: All categories must valid in order for the license to be considered valid. Default: false
	Floating                      *bool                  `json:"floating" validate:"optional"`                        // Floating: When true, license that implements the policy will be valid across multiple machines. Default: false
	RequireProductScope           *bool                  `json:"require_product_scope" validate:"optional"`           // RequireProductScope: Currently unused. Default false
	RequirePolicyScope            *bool                  `json:"require_policy_scope" validate:"optional"`            // RequirePolicyScope: Currently unused. Default false
	RequireMachineScope           *bool                  `json:"require_machine_scope" validate:"optional"`           // RequireMachineScope: Currently unused. Default false
	RequireFingerprintScope       *bool                  `json:"require_fingerprint_scope" validate:"optional"`       // RequireFingerprintScope: Currently unused. Default false
	RequireComponentsScope        *bool                  `json:"require_components_scope" validate:"optional"`        // RequireComponentsScope: Currently unused. Default false
	RequireUserScope              *bool                  `json:"require_user_scope" validate:"optional"`              // RequireUserScope: Currently unused. Default false
	RequireChecksumScope          *bool                  `json:"require_checksum_scope" validate:"optional"`          // RequireChecksumScope: Currently unused. Default false
	RequireVersionScope           *bool                  `json:"require_version_scope" validate:"optional"`           // RequireVersionScope: Currently unused. Default false
	RequireCheckIn                *bool                  `json:"require_check_in" validate:"optional"`                // RequireCheckIn: When true, require check-in at a predefined interval to continue to pass validation. Default: false
	UsePool                       *bool                  `json:"use_pool" validate:"optional"`                        // UsePool: Whether to pull license keys from a finite pool of pre-determined keys
	Encrypted                     *bool                  `json:"encrypted" validate:"optional"`                       // Encrypted
	Protected                     *bool                  `json:"protected" validate:"optional"`                       // Protected: Whether the policy is protected.
	RequireHeartbeat              *bool                  `json:"require_heartbeat" validate:"optional"`               // RequireHeartbeat: Whether the policy requires its machines to maintain a heartbeat.
	CheckInInterval               *int                   `json:"check_in_interval" validate:"optional"`               // CheckInInterval
	CheckInIntervalCount          *int                   `json:"check_in_interval_count" validate:"optional"`         // CheckInIntervalCount
	MaxMachines                   *int                   `json:"max_machines" validate:"optional"`                    // MaxMachines: The maximum number of machines a license implementing the policy can have associated with it
	MaxProcesses                  *int                   `json:"max_processes" validate:"optional"`                   // MaxProcesses: The maximum number of machine processes a license implementing the policy can have associated with it
	MaxUsers                      *int                   `json:"max_users" validate:"optional"`                       // MaxUsers: The maximum number of users a license implementing the policy can have associated with it
	MaxCores                      *int                   `json:"max_cores" validate:"optional"`                       // MaxCores: The maximum number of machine CPU cores a license implementing the policy can have associated with it.
	MaxUses                       *int                   `json:"max_uses" validate:"optional"`                        // MaxUses: The maximum number of uses a license implementing the policy can have.
	HeartbeatDuration             *int                   `json:"heartbeat_duration" validate:"optional"`              // HeartbeatDuration: The heartbeat duration for the policy, in seconds. Currently unused.
	Duration                      *int                   `json:"duration" validate:"optional"`                        // Duration
	HeartbeatCullStrategy         *string                `json:"heartbeat_cull_strategy" validate:"optional"`         // HeartbeatCullStrategy: The strategy used for culling dead machines and processes.
	HeartbeatResurrectionStrategy *string                `json:"heartbeat_resurrection_strategy" validate:"optional"` // HeartbeatResurrectionStrategy: The strategy used for controlling whether or not dead machines and processes can be resurrected shortly after death.
	HeartbeatBasis                *string                `json:"heartbeat_basis" validate:"optional"`                 // HeartbeatBasis: Control when a machine's initial heartbeat is started.
	MachineUniquenessStrategy     *string                `json:"machine_uniqueness_strategy" validate:"optional"`     // MachineUniquenessStrategy: The uniqueness validation strategy for machine fingerprints.
	MachineMatchingStrategy       *string                `json:"machine_matching_strategy" validate:"optional"`       // MachineMatchingStrategy: The matching strategy for machine fingerprints supplied during a license validation.
	ComponentUniquenessStrategy   *string                `json:"component_uniqueness_strategy" validate:"optional"`   // ComponentUniquenessStrategy: The uniqueness validation strategy for component fingerprints.
	ComponentMatchingStrategy     *string                `json:"component_matching_strategy" validate:"optional"`     // ComponentMatchingStrategy: The matching strategy for component fingerprints supplied during a license validation.
	ExpirationStrategy            *string                `json:"expiration_strategy" validate:"optional"`             // ExpirationStrategy: The strategy for expired licenses during a license validation and when accessing releases.
	ExpirationBasis               *string                `json:"expiration_basis" validate:"optional"`                // ExpirationBasis: Control when a license's initial expiry is set.
	RenewalBasis                  *string                `json:"renewal_basis" validate:"optional"`                   // RenewalBasis: Control how a license's expiry is extended during renewal.
	TransferStrategy              *string                `json:"transfer_strategy" validate:"optional"`               // TransferStrategy: Control whether a license's expiration is reset when transferred to this policy.
	AuthenticationStrategy        *string                `json:"authentication_strategy" validate:"optional"`         // AuthenticationStrategy: The strategy used for authenticating as a license, for client-side integrations.
	MachineLeasingStrategy        *string                `json:"machine_leasing_strategy" validate:"optional"`        // MachineLeasingStrategy: The strategy used for leasing machines.
	ProcessLeasingStrategy        *string                `json:"process_leasing_strategy" validate:"optional"`        // ProcessLeasingStrategy: The strategy used for leasing machine processes.
	OverageStrategy               *string                `json:"overage_strategy" validate:"optional"`                // OverageStrategy: The strategy used for allowing machine, core and process overages.
	Metadata                      map[string]interface{} `json:"metadata" validate:"optional"`                        // Metadata: Policy metadata.
}
