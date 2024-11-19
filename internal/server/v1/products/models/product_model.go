package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/product_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type ProductRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	product_attribute.ProductCommonURI
	product_attribute.ProductAttribute
}

type ProductRegistrationOutput struct {
	ID                   string                 `json:"id"`
	TenantName           string                 `json:"tenant_name"`
	Name                 string                 `json:"name"`
	DistributionStrategy string                 `json:"distribution_strategy"`
	Code                 string                 `json:"code"`
	URL                  string                 `json:"url"`
	Platforms            []string               `json:"platform"`
	Metadata             map[string]interface{} `json:"metadata"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type ProductUpdateInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	product_attribute.ProductCommonURI
	product_attribute.ProductAttribute
}

type ProductUpdateOutput struct {
	ID                   string                 `json:"id"`
	TenantName           string                 `json:"tenant_name"`
	Name                 string                 `json:"name"`
	DistributionStrategy string                 `json:"distribution_strategy"`
	Code                 string                 `json:"code"`
	URL                  string                 `json:"url"`
	Platforms            []string               `json:"platform"`
	Metadata             map[string]interface{} `json:"metadata"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type ProductListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	product_attribute.ProductCommonURI
	constants.QueryCommonParam
}

type ProductRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	product_attribute.ProductCommonURI
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
	TracerCtx context.Context
	Tracer    trace.Tracer
	product_attribute.ProductCommonURI
}

type ProductTokensInput struct {
	TracerCtx   context.Context
	Tracer      trace.Tracer
	TenantName  *string   `json:"tenant_name" validate:"required" example:"test"`
	ProductID   uuid.UUID `json:"product_id" validate:"required" example:"test"`
	Name        *string   `json:"name" validate:"optional" example:"test"`
	Expiry      *string   `json:"expiry" validate:"optional" example:"test"`
	Permissions []string  `json:"permissions" validate:"required" example:"test"`
}
