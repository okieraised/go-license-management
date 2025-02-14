package policy_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/utils"
)

type PolicyCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	PolicyID   *string `uri:"policy_id"`
}

func (req *PolicyCommonURI) Validate() error {
	if req.TenantName == nil {
		return cerrors.ErrTenantNameIsEmpty
	}

	if req.PolicyID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.PolicyID)); err != nil {
			return cerrors.ErrPolicyIDIsInvalid
		}
	}

	return nil
}

type PolicyAttributeModel struct {
	Name               *string                `json:"name" validate:"required"`                // Name: name of the policy
	Scheme             *string                `json:"scheme" validate:"optional"`              // Scheme: The encryption/signature scheme used on license keys.
	Strict             *bool                  `json:"strict" validate:"optional"`              // Strict: All categories must valid in order for the license to be considered valid. Default: false
	RateLimited        *bool                  `json:"rate_limited" validate:"optional"`        // RateLimited: Whether the policy is for rate limiting feature. Default: false
	Floating           *bool                  `json:"floating" validate:"optional"`            // Floating: When true, license that implements the policy will be valid across multiple machines. Default: false
	UsePool            *bool                  `json:"use_pool" validate:"optional"`            // UsePool: Whether to pull license keys from a finite pool of pre-determined keys
	Encrypted          *bool                  `json:"encrypted" validate:"optional"`           // Encrypted: Whether to encrypt the license file
	Protected          *bool                  `json:"protected" validate:"optional"`           // Protected: Whether the policy is protected.
	RequireCheckIn     *bool                  `json:"require_check_in" validate:"optional"`    // RequireCheckIn: When true, require check-in at a predefined interval to continue to pass validation. Default: false
	RequireHeartbeat   *bool                  `json:"require_heartbeat" validate:"optional"`   // RequireHeartbeat: Whether the policy requires its machines to maintain a heartbeat.
	MaxMachines        *int                   `json:"max_machines" validate:"optional"`        // MaxMachines: The maximum number of machines a license implementing the policy can have associated with it
	MaxUsers           *int                   `json:"max_users" validate:"optional"`           // MaxUsers: The maximum number of users a license implementing the policy can have associated with it
	MaxUses            *int                   `json:"max_uses" validate:"optional"`            // MaxUses: The maximum number of uses a license implementing the policy can have.
	HeartbeatDuration  *int                   `json:"heartbeat_duration" validate:"optional"`  // HeartbeatDuration: The heartbeat duration for the policy, in seconds.
	Duration           *int64                 `json:"duration" validate:"optional"`            // Duration: The length of time that a policy is valid
	CheckInInterval    *string                `json:"check_in_interval" validate:"optional"`   // CheckInInterval: The time duration between each checkin
	HeartbeatBasis     *string                `json:"heartbeat_basis" validate:"optional"`     // HeartbeatBasis: Control when a machine's initial heartbeat is started.
	ExpirationStrategy *string                `json:"expiration_strategy" validate:"optional"` // ExpirationStrategy: The strategy for expired licenses during a license validation.
	RenewalBasis       *string                `json:"renewal_basis" validate:"optional"`       // RenewalBasis: Control how a license's expiry is extended during renewal.
	OverageStrategy    *string                `json:"overage_strategy" validate:"optional"`    // OverageStrategy: The strategy used for allowing machine overages.
	Metadata           map[string]interface{} `json:"metadata" validate:"optional"`            // Metadata: Policy metadata.
}
