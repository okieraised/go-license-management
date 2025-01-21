package entitlements

import (
	"context"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/entitlement_attribute"
	"go-license-management/internal/services/v1/entitlements/models"
	"go.opentelemetry.io/otel/trace"
)

type EntitlementRegistrationRequest struct {
	Name     *string                `json:"name" validate:"required" example:"test"`
	Code     *string                `json:"code" validate:"required" example:"test"`
	Metadata map[string]interface{} `json:"metadata" validate:"optional"`
}

func (req *EntitlementRegistrationRequest) Validate() error {
	if req.Name == nil {
		return cerrors.ErrEntitlementNameIsEmpty
	}

	if req.Code == nil {
		return cerrors.ErrEntitlementCodeIsEmpty
	}

	return nil
}

func (req *EntitlementRegistrationRequest) ToEntitlementRegistrationInput(ctx context.Context, tracer trace.Tracer, entitlementURI entitlement_attribute.EntitlementCommonURI) *models.EntitlementRegistrationInput {
	return &models.EntitlementRegistrationInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		Name:                 req.Name,
		Code:                 req.Code,
		Metadata:             req.Metadata,
		EntitlementCommonURI: entitlementURI,
	}
}

type EntitlementRetrievalRequest struct {
	entitlement_attribute.EntitlementCommonURI
}

func (req *EntitlementRetrievalRequest) Validate() error {
	if req.EntitlementID == nil {
		return cerrors.ErrEntitlementIDIsEmpty
	}
	return req.EntitlementCommonURI.Validate()
}

func (req *EntitlementRetrievalRequest) ToEntitlementRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.EntitlementRetrievalInput {
	return &models.EntitlementRetrievalInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		EntitlementCommonURI: req.EntitlementCommonURI,
	}
}

type EntitlementDeletionRequest struct {
	entitlement_attribute.EntitlementCommonURI
}

func (req *EntitlementDeletionRequest) Validate() error {
	if req.EntitlementID == nil {
		return cerrors.ErrEntitlementIDIsEmpty
	}
	return req.EntitlementCommonURI.Validate()
}

func (req *EntitlementDeletionRequest) ToEntitlementDeletionInput(ctx context.Context, tracer trace.Tracer) *models.EntitlementDeletionInput {
	return &models.EntitlementDeletionInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		EntitlementCommonURI: req.EntitlementCommonURI,
	}
}

type EntitlementListRequest struct {
	constants.QueryCommonParam
}

func (req *EntitlementListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *EntitlementListRequest) ToEntitlementListInput(ctx context.Context, tracer trace.Tracer, uriParam entitlement_attribute.EntitlementCommonURI) *models.EntitlementListInput {
	return &models.EntitlementListInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		EntitlementCommonURI: uriParam,
		QueryCommonParam:     req.QueryCommonParam,
	}
}
