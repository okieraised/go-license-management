package authentications

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/server/v1/authentications/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type AuthenticationLoginRequest struct {
	Username *string `form:"username" validate:"required" example:"test"`
	Password *string `form:"password" validate:"required" example:"test"`
}

func (req *AuthenticationLoginRequest) Validate() error {
	if req.Username == nil {
		return comerrors.ErrAccountUsernameIsEmpty
	}

	if req.Password == nil {
		return comerrors.ErrAccountPasswordIsEmpty
	}

	return nil
}

func (req *AuthenticationLoginRequest) ToAuthenticationLoginInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.AuthenticationLoginInput {
	return &models.AuthenticationLoginInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
		Username:   req.Username,
		Password:   req.Password,
	}
}
