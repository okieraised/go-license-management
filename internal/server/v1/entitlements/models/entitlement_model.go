package models

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type EntitlementRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string                `json:"tenant_name" validate:"required" example:"test"`
	Name       *string                `json:"name" validate:"required" example:"test"`
	Code       *string                `json:"code" validate:"required" example:"test"`
	Metadata   map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type EntitlementListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type EntitlementRetrievalInput struct {
	TracerCtx     context.Context
	Tracer        trace.Tracer
	TenantName    *string   `json:"tenant_name" validate:"required" example:"test"`
	EntitlementID uuid.UUID `json:"entitlement_id" validate:"required" example:"test"`
}

type EntitlementRetrievalOutput struct {
	ID        uuid.UUID              `json:"id"`
	TenantID  uuid.UUID              `json:"tenant_id"`
	Name      string                 `json:"name,"`
	Code      string                 `json:"code"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type EntitlementUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type EntitlementDeletionInput struct {
	TracerCtx     context.Context
	Tracer        trace.Tracer
	TenantName    *string   `json:"tenant_name" validate:"required" example:"test"`
	EntitlementID uuid.UUID `json:"entitlement_id" validate:"required" example:"test"`
}
