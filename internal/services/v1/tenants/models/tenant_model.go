package models

import (
	"context"
	"go-license-management/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type TenantRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}

type TenantRegistrationOutput struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	constants.QueryCommonParam
}

type TenantRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}

type TenantRetrievalOutput struct {
	Name             string    `json:"name"`
	Ed25519PublicKey string    `json:"ed25519_public_key"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TenantDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}

type TenantRegenerationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}
