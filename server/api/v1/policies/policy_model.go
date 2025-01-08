package policies

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go-license-management/internal/services/v1/policies/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type PolicyRegistrationRequest struct {
	ProductID *string `json:"product_id" validate:"required" example:"test"`
	policy_attribute.PolicyAttributeModel
}

func (req *PolicyRegistrationRequest) Validate() error {
	// Validate must-have field

	// Policy name
	if req.Name == nil {
		return comerrors.ErrPolicyNameIsEmpty
	}

	// Product ID
	if req.ProductID == nil {
		return comerrors.ErrProductIDIsEmpty
	}

	_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
	if err != nil {
		return comerrors.ErrProductIDIsInvalid
	}

	// Strict policy - all categories must valid in order for the license to be considered valid. Default: false
	if req.Strict == nil {
		req.Strict = utils.RefPointer(false)
	}

	// When true, license that implements the policy will be valid across multiple machines. Default: false
	if req.Floating == nil {
		req.Floating = utils.RefPointer(false)
	}

	// The encryption/signature scheme used on license keys. Default to ED25519
	if req.Scheme == nil {
		req.Scheme = utils.RefPointer(constants.PolicySchemeED25519)
	} else {
		if _, ok := constants.ValidPolicySchemeMapper[utils.DerefPointer(req.Scheme)]; !ok {
			return comerrors.ErrPolicySchemeIsInvalid
		}
	}

	// The strategy for expired licenses during a license validation.
	if req.ExpirationStrategy == nil {
		req.ExpirationStrategy = utils.RefPointer(constants.PolicyExpirationStrategyRevokeAccess)
	} else {
		if _, ok := constants.ValidPolicyExpirationStrategyMapper[utils.DerefPointer(req.ExpirationStrategy)]; !ok {
			return comerrors.ErrPolicyInvalidExpirationStrategy
		}
	}

	// The strategy used for allowing machine overages. Default to not allowed any overage
	if req.OverageStrategy == nil {
		req.OverageStrategy = utils.RefPointer(constants.PolicyOverageStrategyNoOverage)
	} else {
		if _, ok := constants.ValidPolicyOverageStrategyMapper[utils.DerefPointer(req.OverageStrategy)]; !ok {
			return comerrors.ErrPolicyInvalidOverageStrategy
		}
	}

	// Control how a license's expiry is extended during renewal.
	if req.RenewalBasis == nil {
		req.RenewalBasis = utils.RefPointer(constants.PolicyRenewalBasisFromExpiry)
	} else {
		if _, ok := constants.ValidPolicyRenewalBasisMapper[utils.DerefPointer(req.RenewalBasis)]; !ok {
			return comerrors.ErrPolicyInvalidRenewalBasis
		}
	}

	// Control when a machine's initial heartbeat is started.
	if req.HeartbeatBasis == nil {
		req.HeartbeatBasis = utils.RefPointer(constants.PolicyHeartbeatBasisFromCreation)
	} else {
		if _, ok := constants.ValidPolicyHeartbeatBasisMapper[utils.DerefPointer(req.HeartbeatBasis)]; !ok {
			return comerrors.ErrPolicyInvalidHeartbeatBasis
		}
	}

	// Control the time duration between each checkin
	if req.CheckInInterval == nil {
		req.CheckInInterval = utils.RefPointer(constants.PolicyCheckinIntervalDaily)
	} else {
		if _, ok := constants.ValidPolicyCheckinIntervalMapper[utils.DerefPointer(req.CheckInInterval)]; !ok {
			return comerrors.ErrPolicyInvalidCheckinInterval
		}
	}

	// Optional parameters

	// When true, require check-in at a predefined interval to continue to pass validation. Default: false
	if req.RequireCheckIn == nil {
		req.RequireCheckIn = utils.RefPointer(false)
	}

	// Whether the policy requires its machines to maintain a heartbeat.
	if req.RequireHeartbeat == nil {
		req.RequireHeartbeat = utils.RefPointer(false)
	}

	//  Whether to pull license keys from a finite pool of pre-determined keys
	if req.UsePool == nil {
		req.UsePool = utils.RefPointer(false)
	}

	//  Whether the policy is protected.
	if req.Protected == nil {
		req.Protected = utils.RefPointer(false)
	}

	// Whether the policy is for rate limiting feature. Default: false
	if req.RateLimited == nil {
		req.RateLimited = utils.RefPointer(false)
	}

	// Whether to encrypt the license file
	if req.Encrypted == nil {
		req.Encrypted = utils.RefPointer(false)
	}

	// The length of time that a policy is valid. Zero value means there is no expiration
	if req.Duration == nil {
		req.Duration = utils.RefPointer(int64(0))
	} else {
		if utils.DerefPointer(req.Duration) < 0 {
			return comerrors.ErrPolicyDurationIsLessThanZero
		}
	}

	// The maximum number of machines a license implementing the policy can have associated with it.
	if req.MaxMachines == nil {
		req.MaxMachines = utils.RefPointer(0)
	} else {
		if utils.DerefPointer(req.MaxMachines) < 0 {
			return comerrors.ErrPolicyMaxMachinesIsLessThanZero
		}
	}

	// The maximum number of uses a license implementing the policy can have.
	if req.MaxUses == nil {
		req.MaxUses = utils.RefPointer(0)
	} else {
		if utils.DerefPointer(req.MaxUses) < 0 {
			return comerrors.ErrPolicyMaxUsesIsLessThanZero
		}
	}

	// The heartbeat duration for the policy, in seconds.
	if req.HeartbeatDuration == nil {
		req.HeartbeatDuration = utils.RefPointer(0)
	} else {
		if utils.DerefPointer(req.HeartbeatDuration) < 0 {
			return comerrors.ErrPolicyHeartbeatDurationIsLessThanZero
		}
	}

	// The maximum number of users a license implementing the policy can have associated with it
	if req.MaxUsers == nil {
		req.MaxUsers = utils.RefPointer(0)
	} else {
		if utils.DerefPointer(req.MaxUsers) < 0 {
			return comerrors.ErrPolicyMaxUsersIsLessThanZero
		}
	}

	return nil
}

func (req *PolicyRegistrationRequest) ToPolicyRegistrationInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyRegistrationInput {
	return &models.PolicyRegistrationInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		PolicyCommonURI:      policyURI,
		ProductID:            req.ProductID,
		PolicyAttributeModel: req.PolicyAttributeModel,
	}
}

type PolicyUpdateRequest struct {
	policy_attribute.PolicyAttributeModel
}

func (req *PolicyUpdateRequest) Validate() error {
	if req.Scheme != nil {
		if _, ok := constants.ValidPolicySchemeMapper[utils.DerefPointer(req.Scheme)]; !ok {
			return comerrors.ErrPolicySchemeIsInvalid
		}
	}

	if req.Duration != nil {
		if utils.DerefPointer(req.Duration) < 0 {
			return comerrors.ErrPolicyDurationIsLessThanZero
		}
	}

	if req.MaxMachines != nil {
		if utils.DerefPointer(req.MaxMachines) < 0 {
			return comerrors.ErrPolicyMaxMachinesIsLessThanZero
		}
	}

	if req.MaxUses != nil {
		if utils.DerefPointer(req.MaxUses) < 0 {
			return comerrors.ErrPolicyMaxUsesIsLessThanZero
		}
	}

	if req.HeartbeatDuration != nil {
		if utils.DerefPointer(req.HeartbeatDuration) < 0 {
			return comerrors.ErrPolicyHeartbeatDurationIsLessThanZero
		}
	}

	if req.MaxUsers != nil {
		if utils.DerefPointer(req.MaxUsers) < 0 {
			return comerrors.ErrPolicyMaxUsersIsLessThanZero
		}
	}

	if req.Scheme != nil {
		if _, ok := constants.ValidPolicySchemeMapper[utils.DerefPointer(req.Scheme)]; !ok {
			return comerrors.ErrPolicySchemeIsInvalid
		}
	}

	if req.ExpirationStrategy != nil {
		if _, ok := constants.ValidPolicyExpirationStrategyMapper[utils.DerefPointer(req.ExpirationStrategy)]; !ok {
			return comerrors.ErrPolicyInvalidExpirationStrategy
		}
	}

	if req.OverageStrategy != nil {
		if _, ok := constants.ValidPolicyOverageStrategyMapper[utils.DerefPointer(req.OverageStrategy)]; !ok {
			return comerrors.ErrPolicyInvalidOverageStrategy
		}
	}

	if req.RenewalBasis != nil {
		if _, ok := constants.ValidPolicyRenewalBasisMapper[utils.DerefPointer(req.RenewalBasis)]; !ok {
			return comerrors.ErrPolicyInvalidRenewalBasis
		}
	}

	if req.HeartbeatBasis != nil {
		if _, ok := constants.ValidPolicyHeartbeatBasisMapper[utils.DerefPointer(req.HeartbeatBasis)]; !ok {
			return comerrors.ErrPolicyInvalidHeartbeatBasis
		}
	}

	if req.CheckInInterval != nil {
		if _, ok := constants.ValidPolicyCheckinIntervalMapper[utils.DerefPointer(req.CheckInInterval)]; !ok {
			return comerrors.ErrPolicyInvalidCheckinInterval
		}
	}

	return nil
}

func (req *PolicyUpdateRequest) ToPolicyUpdateInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyUpdateInput {
	return &models.PolicyUpdateInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		PolicyCommonURI:      policyURI,
		PolicyAttributeModel: req.PolicyAttributeModel,
	}
}

type PolicyDeletionRequest struct {
	policy_attribute.PolicyCommonURI
}

func (req *PolicyDeletionRequest) Validate() error {
	if req.PolicyID == nil {
		return comerrors.ErrPolicyIDIsEmpty
	}
	return req.PolicyCommonURI.Validate()
}

func (req *PolicyDeletionRequest) ToPolicyDeletionInput(ctx context.Context, tracer trace.Tracer) *models.PolicyDeletionInput {
	return &models.PolicyDeletionInput{
		TracerCtx:       ctx,
		Tracer:          tracer,
		PolicyCommonURI: req.PolicyCommonURI,
	}
}

type PolicyRetrievalRequest struct {
	policy_attribute.PolicyCommonURI
}

func (req *PolicyRetrievalRequest) Validate() error {
	if req.PolicyID == nil {
		return comerrors.ErrPolicyIDIsEmpty
	}
	return req.PolicyCommonURI.Validate()
}

func (req *PolicyRetrievalRequest) ToPolicyRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.PolicyRetrievalInput {
	return &models.PolicyRetrievalInput{
		TracerCtx:       ctx,
		Tracer:          tracer,
		PolicyCommonURI: req.PolicyCommonURI,
	}
}

type PolicyListRequest struct {
	constants.QueryCommonParam
}

func (req *PolicyListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *PolicyListRequest) ToPolicyListInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyListInput {
	return &models.PolicyListInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		PolicyCommonURI:  policyURI,
		QueryCommonParam: req.QueryCommonParam,
	}
}

type PolicyAttachmentRequest struct {
	EntitlementID []string `json:"entitlement_id"`
}

func (req *PolicyAttachmentRequest) Validate() error {
	if req.EntitlementID == nil {
		return comerrors.ErrEntitlementIDIsEmpty
	}
	for _, entitlement := range req.EntitlementID {
		if _, err := uuid.Parse(entitlement); err != nil {
			return comerrors.ErrEntitlementIDIsInvalid
		}
	}
	return nil
}

func (req *PolicyAttachmentRequest) ToPolicyAttachmentInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyAttachmentInput {
	return &models.PolicyAttachmentInput{
		TracerCtx:       ctx,
		Tracer:          tracer,
		PolicyCommonURI: policyURI,
		EntitlementID:   req.EntitlementID,
	}
}

type PolicyDetachmentRequest struct {
	ID []string `json:"id"`
}

func (req *PolicyDetachmentRequest) Validate() error {
	if req.ID == nil {
		return comerrors.ErrEntitlementIDIsEmpty
	}
	for _, entitlement := range req.ID {
		if _, err := uuid.Parse(entitlement); err != nil {
			return comerrors.ErrEntitlementIDIsInvalid
		}
	}
	return nil
}

func (req *PolicyDetachmentRequest) ToPolicyDetachmentInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyDetachmentInput {
	return &models.PolicyDetachmentInput{
		TracerCtx:       ctx,
		Tracer:          tracer,
		PolicyCommonURI: policyURI,
		ID:              req.ID,
	}
}

type PolicyEntitlementListRequest struct {
	constants.QueryCommonParam
}

func (req *PolicyEntitlementListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *PolicyEntitlementListRequest) ToPolicyEntitlementListInput(ctx context.Context, tracer trace.Tracer, policyURI policy_attribute.PolicyCommonURI) *models.PolicyEntitlementListInput {
	return &models.PolicyEntitlementListInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		PolicyCommonURI:  policyURI,
		QueryCommonParam: req.QueryCommonParam,
	}
}
