package policies

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go-license-management/internal/server/v1/policies/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type PolicyRegistrationRequest struct {
	ProductID *string `json:"product_id" validate:"required" example:"test"`
	policy_attribute.PolicyAttributeModel
	policy_attribute.PolicyCommonURI
}

func (req *PolicyRegistrationRequest) Validate() error {
	// Must have
	if req.Name == nil {
		return comerrors.ErrPolicyNameIsEmpty
	}

	if req.ProductID == nil {
		return comerrors.ErrProductIDIsEmpty
	}

	_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
	if err != nil {
		return comerrors.ErrProductIDIsInvalid
	}
	if req.Strict == nil {
		req.Strict = utils.RefPointer(false)
	}
	if req.Floating == nil {
		req.Floating = utils.RefPointer(false)
	}
	if req.Scheme == nil {
		req.Scheme = utils.RefPointer(constants.PolicySchemeED25519)
	} else {
		if _, ok := constants.ValidPolicySchemeMapper[utils.DerefPointer(req.Scheme)]; !ok {
			return comerrors.ErrPolicySchemeIsInvalid
		}
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
	if req.ExpirationBasis == nil {
		req.ExpirationBasis = utils.RefPointer(constants.PolicyExpirationBasisFromCreation)
	}
	if req.OverageStrategy == nil {
		req.OverageStrategy = utils.RefPointer(constants.PolicyOverageStrategyNoOverage)
	}
	if req.RenewalBasis == nil {
		req.RenewalBasis = utils.RefPointer(constants.PolicyRenewalBasisFromExpiry)
	}
	if req.HeartbeatBasis == nil {
		req.HeartbeatBasis = utils.RefPointer(constants.PolicyHeartbeatBasisFromCreation)
	}
	if req.CheckInInterval == nil {
		req.CheckInInterval = utils.RefPointer(constants.PolicyCheckinIntervalDaily)
	}
	if req.HeartbeatCullStrategy == nil {
		req.HeartbeatCullStrategy = utils.RefPointer(constants.PolicyHeartbeatCullPolicyDeactivateDead)
	}
	if req.HeartbeatResurrectionStrategy == nil {
		req.HeartbeatResurrectionStrategy = utils.RefPointer(constants.PolicyHeartbeatResurrectionPolicyNoRevive)
	}

	// Optional
	if req.RequireCheckIn == nil {
		req.RequireCheckIn = utils.RefPointer(false)
	}
	if req.RequireHeartbeat == nil {
		req.RequireHeartbeat = utils.RefPointer(false)
	}
	if req.UsePool == nil {
		req.UsePool = utils.RefPointer(false)
	}
	if req.Protected == nil {
		req.Protected = utils.RefPointer(false)
	}
	if req.RateLimited == nil {
		req.RateLimited = utils.RefPointer(false)
	}
	if req.Encrypted == nil {
		req.Encrypted = utils.RefPointer(false)
	}
	if req == nil {
		req.Protected = utils.RefPointer(false)
	}
	if req.Duration == nil {
		req.Duration = utils.RefPointer(0)
	}
	if req.MaxMachines == nil {
		req.MaxMachines = utils.RefPointer(0)
	}
	if req.MaxUses == nil {
		req.MaxUses = utils.RefPointer(0)
	}
	if req.HeartbeatDuration == nil {
		req.HeartbeatDuration = utils.RefPointer(0)
	}
	if req.MaxUsers == nil {
		req.MaxUsers = utils.RefPointer(0)
	}
	if req.Concurrent == nil {
		req.Concurrent = utils.RefPointer(true)
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
}

func (req *PolicyUpdateRequest) Validate() error {
	return nil
}

type PolicyDeleteRequest struct{}

func (req *PolicyDeleteRequest) Validate() error {
	return nil
}

type PolicyRetrievalRequest struct{}

func (req *PolicyRetrievalRequest) Validate() error {
	return nil
}

type PolicyAttachmentRequest struct {
}

func (req *PolicyAttachmentRequest) Validate() error {
	return nil
}

type PolicyDetachmentRequest struct{}

func (req *PolicyDetachmentRequest) Validate() error {
	return nil
}

type PolicyEntitlementListRequest struct {
}

func (req *PolicyEntitlementListRequest) Validate() error {
	return nil
}
