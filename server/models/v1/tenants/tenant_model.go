package tenants

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/server/v1/tenants/models"
	"go.opentelemetry.io/otel/trace"
)

type TenantRegistrationRequest struct {
	Name *string `form:"name" validate:"required" example:"test"`
}

func (req *TenantRegistrationRequest) Validate() error {
	if req.Name == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	return nil
}

func (req *TenantRegistrationRequest) ToTenantRegistrationInput(ctx context.Context, tracer trace.Tracer) *models.TenantRegistrationInput {
	return &models.TenantRegistrationInput{
		TracerCtx: ctx,
		Tracer:    tracer,
		Name:      req.Name,
	}
}

type TenantRetrievalRequest struct {
	TenantName *string `uri:"tenant_name" binding:"required"`
}

func (req *TenantRetrievalRequest) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

func (req *TenantRetrievalRequest) ToTenantRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.TenantRetrievalInput {
	return &models.TenantRetrievalInput{
		TracerCtx: ctx,
		Tracer:    tracer,
		Name:      req.TenantName,
	}
}

type TenantDeletionRequest struct {
	TenantName *string `uri:"tenant_name" binding:"required"`
}

func (req *TenantDeletionRequest) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

func (req *TenantDeletionRequest) ToTenantDeletionInput(ctx context.Context, tracer trace.Tracer) *models.TenantDeletionInput {
	return &models.TenantDeletionInput{
		TracerCtx: ctx,
		Tracer:    tracer,
		Name:      req.TenantName,
	}
}
