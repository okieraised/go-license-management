package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type ProductRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type ProductListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type ProductRetrievalInput struct {
	TracerCtx   context.Context
	Tracer      trace.Tracer
	TenantName  *string `json:"tenant_name" validate:"required" example:"test"`
	ProductName *string `json:"product_name" validate:"required" example:"test"`
}

type ProductDeletionInput struct {
	TracerCtx   context.Context
	Tracer      trace.Tracer
	TenantName  *string `json:"tenant_name" validate:"required" example:"test"`
	ProductName *string `json:"product_name" validate:"required" example:"test"`
}
