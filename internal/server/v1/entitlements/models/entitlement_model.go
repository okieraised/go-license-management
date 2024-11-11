package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/entitlement_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type EntitlementRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string                `json:"name" validate:"required" example:"test"`
	Code      *string                `json:"code" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
	entitlement_attribute.EntitlementCommonURI
}

type EntitlementListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	entitlement_attribute.EntitlementCommonURI
	constants.QueryCommonParam
}

type EntitlementRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	entitlement_attribute.EntitlementCommonURI
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

type EntitlementDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	entitlement_attribute.EntitlementCommonURI
}
