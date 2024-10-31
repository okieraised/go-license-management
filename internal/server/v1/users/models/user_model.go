package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type UserRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Username  *string `json:"username,omitempty" validate:"required" example:"test"`
	Password  *string `json:"password,omitempty" validate:"required" example:"test"`
	Email     *string `json:"email,omitempty" validate:"required" example:"test"`
	Role      *string `json:"role,omitempty" validate:"required" example:"test"`
}
