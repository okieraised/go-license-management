package policies

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/policy"
	"go-license-management/internal/server/v1/policies/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type PolicyRegistrationRequest struct {
	ProductID *string `json:"product_id" validate:"required" example:"test"`
	policy.PolicyAttributeModel
}

func (req *PolicyRegistrationRequest) Validate() error {
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

	return nil
}

func (req *PolicyRegistrationRequest) ToPolicyRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.PolicyRegistrationInput {
	return &models.PolicyRegistrationInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		TenantName:           utils.RefPointer(tenantName),
		ProductID:            req.ProductID,
		PolicyAttributeModel: req.PolicyAttributeModel,
	}
}
