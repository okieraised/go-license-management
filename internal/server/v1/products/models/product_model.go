package models

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"time"
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
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	ProductID  *string `json:"product_id" validate:"required" example:"test"`
}

type ProductRetrievalOutput struct {
	ID                   uuid.UUID              `json:"id"`
	TenantID             uuid.UUID              `json:"tenant_id"`
	Name                 string                 `json:"name"`
	DistributionStrategy string                 `json:"distribution_strategy"`
	Code                 string                 `json:"code"`
	Platforms            []string               `json:"platform"`
	Metadata             map[string]interface{} `json:"metadata"`
	URL                  string                 `json:"url,type"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type ProductDeletionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	ProductID  *string `json:"product_id" validate:"required" example:"test"`
}

type ProductTokensInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
