package models

import (
	"context"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/account_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type AccountRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	account_attribute.AccountCommonURI
	Username  *string                `json:"username" validate:"required" example:"test"`
	Password  *string                `json:"password" validate:"required" example:"test"`
	FirstName *string                `json:"first_name" validate:"required" example:"test"`
	LastName  *string                `json:"lastName" validate:"required" example:"test"`
	Email     *string                `json:"email" validate:"required" example:"test"`
	Role      *string                `json:"role" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"required" example:"test"`
}

type AccountRetrievalOutput struct {
	Username  string                 `json:"username"`
	RoleName  string                 `json:"role_name"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
	BannedAt  time.Time              `json:"banned_at"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AccountRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	account_attribute.AccountCommonURI
}

type AccountDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	account_attribute.AccountCommonURI
}

type AccountUpdateInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	account_attribute.AccountCommonURI
	Password  *string                `json:"password" validate:"required" example:"test"`
	FirstName *string                `json:"first_name" validate:"optional" example:"test"`
	LastName  *string                `json:"lastName" validate:"optional" example:"test"`
	Email     *string                `json:"email" validate:"required" example:"test"`
	Role      *string                `json:"role" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type AccountListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	account_attribute.AccountCommonURI
	constants.QueryCommonParam
}
