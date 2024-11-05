package policies

import (
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/utils"
)

type PolicyAttributeModel struct {
	Name                          *string                `json:"name" validate:"required"`                            // Name: name of the policy
	Scheme                        *string                `json:"scheme" validate:"optional"`                          // Scheme: The encryption/signature scheme used on license keys.
	Strict                        *bool                  `json:"strict" validate:"optional"`                          // Strict: All categories must valid in order for the license to be considered valid. Default: false
	Floating                      *bool                  `json:"floating" validate:"optional"`                        // Floating: When true, license that implements the policy will be valid across multiple machines. Default: false
	RequireProductScope           *bool                  `json:"require_product_scope" validate:"optional"`           // RequireProductScope:
	RequirePolicyScope            *bool                  `json:"require_policy_scope" validate:"optional"`            // RequirePolicyScope:
	RequireMachineScope           *bool                  `json:"require_machine_scope" validate:"optional"`           // RequireMachineScope
	RequireFingerprintScope       *bool                  `json:"require_fingerprint_scope" validate:"optional"`       // RequireFingerprintScope
	RequireComponentsScope        *bool                  `json:"require_components_scope" validate:"optional"`        // RequireComponentsScope
	RequireUserScope              *bool                  `json:"require_user_scope" validate:"optional"`              // RequireUserScope
	RequireChecksumScope          *bool                  `json:"require_checksum_scope" validate:"optional"`          // RequireChecksumScope
	RequireVersionScope           *bool                  `json:"require_version_scope" validate:"optional"`           // RequireVersionScope
	RequireCheckIn                *bool                  `json:"require_check_in" validate:"optional"`                // RequireCheckIn
	UsePool                       *bool                  `json:"use_pool" validate:"optional"`                        // UsePool
	Encrypted                     *bool                  `json:"encrypted" validate:"optional"`                       // Encrypted
	Protected                     *bool                  `json:"protected" validate:"optional"`                       // Protected
	RequireHeartbeat              *bool                  `json:"require_heartbeat" validate:"optional"`               // RequireHeartbeat
	CheckInInterval               *int                   `json:"check_in_interval" validate:"optional"`               // CheckInInterval
	CheckInIntervalCount          *int                   `json:"check_in_interval_count" validate:"optional"`         // CheckInIntervalCount
	MaxMachines                   *int                   `json:"max_machines" validate:"optional"`                    // MaxMachines
	MaxProcesses                  *int                   `json:"max_processes" validate:"optional"`                   // MaxProcesses
	MaxUsers                      *int                   `json:"max_users" validate:"optional"`                       // MaxUsers
	MaxCores                      *int                   `json:"max_cores" validate:"optional"`                       // MaxCores
	MaxUses                       *int                   `json:"max_uses" validate:"optional"`                        // MaxUses
	HeartbeatDuration             *int                   `json:"heartbeat_duration" validate:"optional"`              // HeartbeatDuration
	Duration                      *int                   `json:"duration" validate:"optional"`                        // Duration
	HeartbeatCullStrategy         *string                `json:"heartbeat_cull_strategy" validate:"optional"`         // HeartbeatCullStrategy
	HeartbeatResurrectionStrategy *string                `json:"heartbeat_resurrection_strategy" validate:"optional"` // HeartbeatResurrectionStrategy
	HeartbeatBasis                *string                `json:"heartbeat_basis" validate:"optional"`                 // HeartbeatBasis
	MachineUniquenessStrategy     *string                `json:"machine_uniqueness_strategy" validate:"optional"`     // MachineUniquenessStrategy
	MachineMatchingStrategy       *string                `json:"machine_matching_strategy" validate:"optional"`       // MachineMatchingStrategy
	ComponentUniquenessStrategy   *string                `json:"component_uniqueness_strategy" validate:"optional"`   // ComponentUniquenessStrategy
	ComponentMatchingStrategy     *string                `json:"component_matching_strategy" validate:"optional"`     // ComponentMatchingStrategy
	ExpirationStrategy            *string                `json:"expiration_strategy" validate:"optional"`             // ExpirationStrategy
	ExpirationBasis               *string                `json:"expiration_basis" validate:"optional"`                // ExpirationBasis
	RenewalBasis                  *string                `json:"renewal_basis" validate:"optional"`                   // RenewalBasis
	TransferStrategy              *string                `json:"transfer_strategy" validate:"optional"`               // TransferStrategy
	AuthenticationStrategy        *string                `json:"authentication_strategy" validate:"optional"`         // AuthenticationStrategy
	MachineLeasingStrategy        *string                `json:"machine_leasing_strategy" validate:"optional"`        // MachineLeasingStrategy
	ProcessLeasingStrategy        *string                `json:"process_leasing_strategy" validate:"optional"`        // ProcessLeasingStrategy
	OverageStrategy               *string                `json:"overage_strategy" validate:"optional"`                // OverageStrategy
	Metadata                      map[string]interface{} `json:"metadata" validate:"optional"`                        // Metadata
}

type PolicyRegistrationRequest struct {
	PolicyAttributeModel
}

func (req *PolicyRegistrationRequest) Validate() error {
	if req.Name == nil {
		return comerrors.ErrPolicyNameIsEmpty
	}
	if req.Strict == nil {
		req.Strict = utils.RefPointer(false)
	}
	if req.Floating == nil {
		req.Floating = utils.RefPointer(false)
	}
	if req.Scheme == nil {
		req.Scheme = utils.RefPointer(constants.PolicySchemeED25519)
	}
	if req.ExpirationStrategy == nil {
		req.ExpirationStrategy = utils.RefPointer(constants.PolicyExpirationStrategyRevokeAccess)
	}
	if req.TransferStrategy == nil {
		req.TransferStrategy = utils.RefPointer(constants.PolicyTransferStrategyResetExpiry)
	}
	if req.AuthenticationStrategy == nil {
		req.AuthenticationStrategy = utils.RefPointer(constants.PolicyAuthenticationStrategyLicense)
	}

	return nil
}
