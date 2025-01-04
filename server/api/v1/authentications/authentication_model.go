package authentications

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/models/authentication_attribute"
	"go-license-management/internal/services/v1/authentications/models"
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

func (req *AuthenticationLoginRequest) ToAuthenticationLoginInput(ctx context.Context, tracer trace.Tracer, uriReq authentication_attribute.AuthenticationCommonURI) *models.AuthenticationLoginInput {
	return &models.AuthenticationLoginInput{
		TracerCtx:               ctx,
		Tracer:                  tracer,
		AuthenticationCommonURI: uriReq,
		Username:                req.Username,
		Password:                req.Password,
	}
}
