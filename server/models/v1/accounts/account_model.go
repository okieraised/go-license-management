package accounts

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/server/v1/accounts/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type AccountRegistrationRequest struct {
	Username  *string                `json:"username" validate:"required" example:"test"`
	Password  *string                `json:"password" validate:"required" example:"test"`
	FirstName *string                `json:"first_name" validate:"optional" example:"test"`
	LastName  *string                `json:"lastName" validate:"optional" example:"test"`
	Email     *string                `json:"email" validate:"required" example:"test"`
	Role      *string                `json:"role" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *AccountRegistrationRequest) Validate() error {
	if req.Username == nil {
		return comerrors.ErrAccountUsernameIsEmpty
	}
	if req.Password == nil {
		return comerrors.ErrAccountPasswordIsEmpty
	}

	if req.Email == nil {
		return comerrors.ErrAccountEmailIsEmpty
	}
	if req.Role == nil {
		req.Role = utils.RefPointer(constants.RoleUser)
	} else {
		if _, ok := constants.ValidRoleMapper[utils.DerefPointer(req.Role)]; !ok {
			return comerrors.ErrAccountRoleIsInvalid
		}
	}
	if req.FirstName == nil {
		req.FirstName = utils.RefPointer("")
	}
	if req.LastName == nil {
		req.LastName = utils.RefPointer("")
	}

	return nil
}

func (req *AccountRegistrationRequest) ToAccountRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.AccountRegistrationInput {
	return &models.AccountRegistrationInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
		Username:   req.Username,
		Password:   req.Password,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Role:       req.Role,
		Metadata:   req.Metadata,
	}
}
