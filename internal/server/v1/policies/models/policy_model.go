package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/policy"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type PolicyRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	ProductID  *string `json:"product_id" validate:"required" example:"test"`
	policy.PolicyAttributeModel
}

type PolicyListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	Limit      *int    `json:"limit" validate:"required" example:"test"`
	Offset     *int    `json:"offset" validate:"required" example:"test"`
}

type PolicyRetrievalInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string   `json:"tenant_name" validate:"required" example:"test"`
	PolicyID   uuid.UUID `json:"Policy_id" validate:"required" example:"test"`
}

type PolicyRetrievalOutput struct {
	ID        uuid.UUID              `json:"id"`
	TenantID  uuid.UUID              `json:"tenant_id"`
	Name      string                 `json:"name,"`
	Code      string                 `json:"code"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type PolicyUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type PolicyDeletionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string   `json:"tenant_name" validate:"required" example:"test"`
	PolicyID   uuid.UUID `json:"Policy_id" validate:"required" example:"test"`
}
