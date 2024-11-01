package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type AccountRegistrationInput struct {
	TracerCtx   context.Context
	Tracer      trace.Tracer
	Name        *string `json:"name" validate:"required" example:"test"`
	Description *string `json:"description"  validate:"optional" example:"test"`
}
