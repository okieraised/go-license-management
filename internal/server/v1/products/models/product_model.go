package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type ProductRegistrationInput struct {
	TracerCtx            context.Context
	Tracer               trace.Tracer
	TenantName           *string                `json:"tenant_name" validate:"required" example:"test"`
	Name                 *string                `json:"name" validate:"required" example:"test"`
	Code                 *string                `json:"code" validate:"required" example:"test"`
	DistributionStrategy *string                `json:"distribution_strategy" validate:"optional" example:"test"`
	Url                  *string                `json:"url" validate:"optional" example:"test"`
	Permissions          []string               `json:"permissions" validate:"optional" example:"test"`
	Platforms            []string               `json:"platforms" validate:"optional" example:"test"`
	Metadata             map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type ProductUpdateInput struct {
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

type ProductTokensInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
