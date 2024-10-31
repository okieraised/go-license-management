package users

import (
	"context"
	"errors"
	"go-license-management/internal/server/v1/users/models"
	"go.opentelemetry.io/otel/trace"
)

type UserRegistrationRequest struct {
	Username *string `json:"username,omitempty" validate:"required" example:"test"`
	Password *string `json:"password,omitempty" validate:"required" example:"test"`
	Email    *string `json:"email,omitempty" validate:"required" example:"test"`
	Role     *string `json:"role,omitempty" validate:"required" example:"test"`
}

func (req *UserRegistrationRequest) Validate() error {
	if req.Username == nil {
		return errors.New("username is empty")
	}

	if req.Password == nil {
		return errors.New("password is empty")
	}

	if req.Email == nil {
		return errors.New("email is empty")
	}

	if req.Role == nil {
		return errors.New("user role is empty")
	}

	return nil
}

func (req *UserRegistrationRequest) ToUserRegistrationInput(ctx context.Context, tracer trace.Tracer) *models.UserRegistrationInput {
	return &models.UserRegistrationInput{
		TracerCtx: ctx,
		Tracer:    tracer,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Role:      req.Role,
	}
}
