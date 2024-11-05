package policies

import "go-license-management/internal/comerrors"

type PolicyAttributeModel struct {
	Name                          *string                `json:"name" validate:"required"`
	Scheme                        *string                `json:"scheme" validate:"optional"`
	Strict                        *bool                  `json:"strict" validate:"optional"`
	Floating                      *bool                  `json:"floating" validate:"optional"`
	RequireProductScope           *bool                  `json:"require_product_scope" validate:"optional"`
	RequirePolicyScope            *bool                  `json:"require_policy_scope" validate:"optional"`
	RequireMachineScope           *bool                  `json:"require_machine_scope" validate:"optional"`
	RequireFingerprintScope       *bool                  `json:"require_fingerprint_scope" validate:"optional"`
	RequireComponentsScope        *bool                  `json:"require_components_scope" validate:"optional"`
	RequireUserScope              *bool                  `json:"require_user_scope" validate:"optional"`
	RequireChecksumScope          *bool                  `json:"require_checksum_scope" validate:"optional"`
	RequireVersionScope           *bool                  `json:"require_version_scope" validate:"optional"`
	RequireCheckIn                *bool                  `json:"require_check_in" validate:"optional"`
	UsePool                       *bool                  `json:"use_pool" validate:"optional"`
	Encrypted                     *bool                  `json:"encrypted" validate:"optional"`
	Protected                     *bool                  `json:"protected" validate:"optional"`
	RequireHeartbeat              *bool                  `json:"require_heartbeat" validate:"optional"`
	CheckInInterval               *int                   `json:"check_in_interval" validate:"optional"`
	CheckInIntervalCount          *int                   `json:"check_in_interval_count" validate:"optional"`
	MaxMachines                   *int                   `json:"max_machines" validate:"optional"`
	MaxProcesses                  *int                   `json:"max_processes" validate:"optional"`
	MaxUsers                      *int                   `json:"max_users" validate:"optional"`
	MaxCores                      *int                   `json:"max_cores" validate:"optional"`
	MaxUses                       *int                   `json:"max_uses" validate:"optional"`
	HeartbeatDuration             *int                   `json:"heartbeat_duration" validate:"optional"`
	Duration                      *int                   `json:"duration" validate:"optional"`
	HeartbeatCullStrategy         *string                `json:"heartbeat_cull_strategy" validate:"optional"`
	HeartbeatResurrectionStrategy *string                `json:"heartbeat_resurrection_strategy" validate:"optional"`
	HeartbeatBasis                *string                `json:"heartbeat_basis" validate:"optional"`
	MachineUniquenessStrategy     *string                `json:"machine_uniqueness_strategy" validate:"optional"`
	MachineMatchingStrategy       *string                `json:"machine_matching_strategy" validate:"optional"`
	ComponentUniquenessStrategy   *string                `json:"component_uniqueness_strategy" validate:"optional"`
	ComponentMatchingStrategy     *string                `json:"component_matching_strategy" validate:"optional"`
	ExpirationStrategy            *string                `json:"expiration_strategy" validate:"optional"`
	ExpirationBasis               *string                `json:"expiration_basis" validate:"optional"`
	RenewalBasis                  *string                `json:"renewal_basis" validate:"optional"`
	TransferStrategy              *string                `json:"transfer_strategy" validate:"optional"`
	AuthenticationStrategy        *string                `json:"authentication_strategy" validate:"optional"`
	MachineLeasingStrategy        *string                `json:"machine_leasing_strategy" validate:"optional"`
	ProcessLeasingStrategy        *string                `json:"process_leasing_strategy" validate:"optional"`
	OverageStrategy               *string                `json:"overage_strategy" validate:"optional"`
	Metadata                      map[string]interface{} `json:"metadata" validate:"optional"`
}

type PolicyRegistrationRequest struct {
	PolicyAttributeModel
}

func (req *PolicyRegistrationRequest) Validate() error {
	if req.Name == nil {
		return comerrors.ErrPolicyNameIsEmpty
	}
	return nil
}
