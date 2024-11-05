package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type LicenseRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
