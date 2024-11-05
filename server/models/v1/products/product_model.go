package products

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/server/v1/products/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type ProductRegistrationRequest struct {
	Name                 *string                `json:"name" validate:"required" example:"test"`
	Code                 *string                `json:"code" validate:"required" example:"test"`
	DistributionStrategy *string                `json:"distribution_strategy" validate:"optional" example:"test"`
	Url                  *string                `json:"url" validate:"optional" example:"test"`
	Permissions          []string               `json:"permissions" validate:"optional" example:"test"`
	Platforms            []string               `json:"platforms" validate:"optional" example:"test"`
	Metadata             map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *ProductRegistrationRequest) Validate() error {
	if req.Name == nil {
		return comerrors.ErrProductNameIsEmpty
	}

	if req.Code == nil {
		return comerrors.ErrProductCodeIsEmpty
	}

	if req.DistributionStrategy == nil {
		req.DistributionStrategy = utils.RefPointer(constants.ProductDistributionStrategyLicensed)
	} else {
		if _, ok := constants.ValidProductDistributionStrategyMapper[utils.DerefPointer(req.DistributionStrategy)]; !ok {
			return comerrors.ErrProductDistributionStrategyIsInvalid
		}
	}
	return nil
}

func (req *ProductRegistrationRequest) ToProductRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.ProductRegistrationInput {
	return &models.ProductRegistrationInput{
		TracerCtx:            ctx,
		Tracer:               tracer,
		TenantName:           utils.RefPointer(tenantName),
		Name:                 req.Name,
		Code:                 req.Code,
		DistributionStrategy: req.DistributionStrategy,
		Url:                  req.Url,
		Platforms:            req.Platforms,
		Permissions:          req.Permissions,
		Metadata:             req.Metadata,
	}
}

type ProductListRequest struct {
}

func (req *ProductListRequest) Validate() error {
	return nil
}

type ProductRetrievalRequest struct {
	ProductID  *string `uri:"product_id" validate:"required" example:"test"`
	TenantName *string `uri:"tenant_name" validate:"required" example:"test"`
}

func (req *ProductRetrievalRequest) Validate() error {
	if req.ProductID == nil {
		return comerrors.ErrProductIDIsEmpty
	}
	_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
	if err != nil {
		return comerrors.ErrProductIDIsInvalid
	}

	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

func (req *ProductRetrievalRequest) ToProductRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.ProductRetrievalInput {
	return &models.ProductRetrievalInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: req.TenantName,
		ProductID:  req.ProductID,
	}
}

type ProductDeletionRequest struct {
	ProductID  *string `uri:"product_id" validate:"required" example:"test"`
	TenantName *string `uri:"tenant_name" validate:"required" example:"test"`
}

func (req *ProductDeletionRequest) Validate() error {
	if req.ProductID == nil {
		return comerrors.ErrProductIDIsEmpty
	}
	_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
	if err != nil {
		return comerrors.ErrProductIDIsInvalid
	}

	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}
	return nil
}

type ProductTokenRequest struct {
}

func (req *ProductTokenRequest) Validate() error {
	return nil
}
