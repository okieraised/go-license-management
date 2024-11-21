package models

import (
	"context"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type PolicyRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
	policy_attribute.PolicyAttributeModel
	ProductID *string `json:"product_id" validate:"required" example:"test"`
}

type PolicyListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
	constants.QueryCommonParam
}

type PolicyRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}

type PolicyRetrievalOutput struct {
	ID         string    `json:"id"`
	TenantName string    `json:"tenant_name"`
	PublicKey  string    `json:"public_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	policy_attribute.PolicyAttributeModel
}

type PolicyUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type PolicyDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}

type PolicyAttachmentInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}

type PolicyDetachmentInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}

type PolicyEntitlementListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}
