package entitlements

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/server/v1/entitlements/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type EntitlementRegistrationRequest struct {
	Name     *string                `json:"name" validate:"required" example:"test"`
	Code     *string                `json:"code" validate:"required" example:"test"`
	Metadata map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *EntitlementRegistrationRequest) Validate() error {
	if req.Name == nil {
		return errors.New("entitlement name is empty")
	}

	if req.Code == nil {
		return errors.New("entitlement code is empty")
	}

	return nil
}

func (req *EntitlementRegistrationRequest) ToEntitlementRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.EntitlementRegistrationInput {

	return &models.EntitlementRegistrationInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
		Name:       req.Name,
		Code:       req.Code,
		Metadata:   req.Metadata,
	}
}

type EntitlementRetrievalRequest struct {
	EntitlementID *string `uri:"entitlement_id" validate:"required" example:"test"`
	TenantName    *string `uri:"tenant_name" validate:"required" example:"test"`
}

func (req *EntitlementRetrievalRequest) Validate() error {
	if req.EntitlementID == nil {
		return comerrors.ErrEntitlementIDIsEmpty
	}
	_, err := uuid.Parse(utils.DerefPointer(req.EntitlementID))
	if err != nil {
		return comerrors.ErrEntitlementIDIsInvalid
	}

	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

func (req *EntitlementRetrievalRequest) ToEntitlementRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.EntitlementRetrievalInput {
	return &models.EntitlementRetrievalInput{
		TracerCtx:     ctx,
		Tracer:        tracer,
		TenantName:    req.TenantName,
		EntitlementID: uuid.MustParse(utils.DerefPointer(req.EntitlementID)),
	}
}

type EntitlementDeletionRequest struct {
	EntitlementID *string `uri:"entitlement_id" validate:"required" example:"test"`
	TenantName    *string `uri:"tenant_name" validate:"required" example:"test"`
}

func (req *EntitlementDeletionRequest) Validate() error {
	if req.EntitlementID == nil {
		return comerrors.ErrEntitlementIDIsEmpty
	}
	_, err := uuid.Parse(utils.DerefPointer(req.EntitlementID))
	if err != nil {
		return comerrors.ErrEntitlementIDIsInvalid
	}

	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

func (req *EntitlementDeletionRequest) ToEntitlementDeletionInput(ctx context.Context, tracer trace.Tracer) *models.EntitlementDeletionInput {
	return &models.EntitlementDeletionInput{
		TracerCtx:     ctx,
		Tracer:        tracer,
		TenantName:    req.TenantName,
		EntitlementID: uuid.MustParse(utils.DerefPointer(req.EntitlementID)),
	}
}

type EntitlementListRequest struct {
	Limit  *int `form:"limit" validate:"optional" example:"10"`
	Offset *int `form:"offset" validate:"optional" example:"10"`
}

func (req *EntitlementListRequest) Validate() error {
	if req.Limit == nil {
		req.Limit = utils.RefPointer(10)
	}

	if req.Offset == nil {
		req.Offset = utils.RefPointer(0)
	}

	return nil
}

func (req *EntitlementListRequest) ToEntitlementListInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.EntitlementListInput {
	return &models.EntitlementListInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
		Limit:      req.Limit,
		Offset:     req.Offset,
	}
}
