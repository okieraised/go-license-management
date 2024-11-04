package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type AccountRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string                `json:"tenant_name" validate:"required" example:"test"`
	Username   *string                `json:"username" validate:"required" example:"test"`
	Password   *string                `json:"password" validate:"required" example:"test"`
	FirstName  *string                `json:"first_name" validate:"required" example:"test"`
	LastName   *string                `json:"lastName" validate:"required" example:"test"`
	Email      *string                `json:"email" validate:"required" example:"test"`
	Role       *string                `json:"role" validate:"required" example:"test"`
	Metadata   map[string]interface{} `json:"metadata" validate:"required" example:"test"`
}

type AccountListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type AccountRetrievalInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type AccountDeletionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type AccountUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
