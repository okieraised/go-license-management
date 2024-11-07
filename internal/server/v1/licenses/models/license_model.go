package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type LicenseRegistrationInput struct {
	TracerCtx    context.Context
	Tracer       trace.Tracer
	TenantName   *string                `json:"tenant_name" validate:"required" example:"test"`
	Name         *string                `json:"name" validate:"required" example:"test"`
	Key          *string                `json:"key" validate:"optional" example:"test"`
	Expiry       *string                `json:"expiry" validate:"optional" example:"test"`
	MaxMachine   *int                   `json:"max_machine" validate:"optional" example:"test"`
	MaxProcesses *int                   `json:"max_processes" validate:"optional" example:"test"`
	MaxUsers     *int                   `json:"max_users" validate:"optional" example:"test"`
	MaxUses      *int                   `json:"max_uses" validate:"optional" example:"test"`
	MaxCores     *int                   `json:"max_cores" validate:"optional" example:"test"`
	Protected    *bool                  `json:"protected" validate:"optional" example:"test"`
	Suspended    *bool                  `json:"suspended" validate:"optional" example:"test"`
	Permissions  []string               `json:"permissions" validate:"optional" example:"test"`
	Metadata     map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type LicenseListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	LicenseID  *string `json:"license_id" validate:"required" example:"test"`
}

type LicenseRetrievalInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseDeletionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseActionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
