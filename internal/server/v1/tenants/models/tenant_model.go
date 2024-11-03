package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type TenantRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}
