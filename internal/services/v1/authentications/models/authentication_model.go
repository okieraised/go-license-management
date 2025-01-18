package models

import (
	"context"
	"go-license-management/internal/infrastructure/models/authentication_attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthenticationLoginInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	authentication_attribute.AuthenticationCommonURI
	Username *string `json:"username" validate:"required" example:"test"`
	Password *string `json:"password" validate:"required" example:"test"`
}

type AuthenticationLoginOutput struct {
	Access   string `json:"access"`
	ExpireAt int64  `json:"expire_at"`
}
