package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type AuthenticationLoginInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	Username   *string `json:"username" validate:"required" example:"test"`
	Password   *string `json:"password" validate:"required" example:"test"`
}
