package products

import (
	"context"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/product_attribute"
	"go-license-management/internal/services/v1/products/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type ProductRegistrationRequest struct {
	product_attribute.ProductAttribute
}

func (req *ProductRegistrationRequest) Validate() error {
	if req.Name == nil {
		return cerrors.ErrProductNameIsEmpty
	}

	if req.Code == nil {
		return cerrors.ErrProductCodeIsEmpty
	}

	if req.DistributionStrategy == nil {
		req.DistributionStrategy = utils.RefPointer(constants.ProductDistributionStrategyLicensed)
	} else {
		if _, ok := constants.ValidProductDistributionStrategyMapper[utils.DerefPointer(req.DistributionStrategy)]; !ok {
			return cerrors.ErrProductDistributionStrategyIsInvalid
		}
	}
	return nil
}

func (req *ProductRegistrationRequest) ToProductRegistrationInput(ctx context.Context, tracer trace.Tracer, productURI product_attribute.ProductCommonURI) *models.ProductRegistrationInput {
	return &models.ProductRegistrationInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: productURI,
		ProductAttribute: req.ProductAttribute,
	}
}

type ProductUpdateRequest struct {
	product_attribute.ProductAttribute
}

func (req *ProductUpdateRequest) Validate() error {
	if req.DistributionStrategy != nil {
		if _, ok := constants.ValidProductDistributionStrategyMapper[utils.DerefPointer(req.DistributionStrategy)]; !ok {
			return cerrors.ErrProductDistributionStrategyIsInvalid
		}
	}
	return nil
}

func (req *ProductUpdateRequest) ToProductUpdateInput(ctx context.Context, tracer trace.Tracer, productURI product_attribute.ProductCommonURI) *models.ProductUpdateInput {
	return &models.ProductUpdateInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: productURI,
		ProductAttribute: req.ProductAttribute,
	}
}

type ProductListRequest struct {
	constants.QueryCommonParam
}

func (req *ProductListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *ProductListRequest) ToProductListInput(ctx context.Context, tracer trace.Tracer, productURI product_attribute.ProductCommonURI) *models.ProductListInput {
	return &models.ProductListInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: productURI,
		QueryCommonParam: req.QueryCommonParam,
	}
}

type ProductRetrievalRequest struct {
	product_attribute.ProductCommonURI
}

func (req *ProductRetrievalRequest) Validate() error {
	if req.ProductID == nil {
		return cerrors.ErrProductIDIsEmpty
	}
	return req.ProductCommonURI.Validate()
}

func (req *ProductRetrievalRequest) ToProductRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.ProductRetrievalInput {
	return &models.ProductRetrievalInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: req.ProductCommonURI,
	}
}

type ProductDeletionRequest struct {
	product_attribute.ProductCommonURI
}

func (req *ProductDeletionRequest) Validate() error {
	if req.ProductID == nil {
		return cerrors.ErrProductIDIsEmpty
	}
	return req.ProductCommonURI.Validate()
}

func (req *ProductDeletionRequest) ToProductDeletionInput(ctx context.Context, tracer trace.Tracer) *models.ProductDeletionInput {
	return &models.ProductDeletionInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: req.ProductCommonURI,
	}
}

type ProductTokenRequest struct {
	Name        *string  `json:"name" validate:"optional" example:"test"`
	Expiry      *string  `json:"expiry" validate:"optional" example:"test"`
	Permissions []string `json:"permissions" validate:"required" example:"test"`
}

func (req *ProductTokenRequest) Validate() error {

	if req.Expiry != nil {
		_, err := time.Parse(constants.DateFormatISO8601Hyphen, utils.DerefPointer(req.Expiry))
		if err != nil {
			return cerrors.ErrProductTokenExpirationFormatIsInvalid
		}
	}
	if req.Permissions == nil {
		req.Permissions = []string{"*"}
	}

	return nil
}

func (req *ProductTokenRequest) ToProductTokenInput(ctx context.Context, tracer trace.Tracer, productURI product_attribute.ProductCommonURI) *models.ProductTokensInput {
	return &models.ProductTokensInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		ProductCommonURI: productURI,
		Name:             req.Name,
		Expiry:           req.Expiry,
		Permissions:      req.Permissions,
	}
}
