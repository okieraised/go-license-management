package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/products/models"
	"go-license-management/internal/server/v1/products/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type ProductService struct {
	repo   repository.IProduct
	logger *logging.Logger
}

func NewProductService(options ...func(*ProductService)) *ProductService {
	svc := &ProductService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IProduct) func(*ProductService) {
	return func(c *ProductService) {
		c.repo = repo
	}
}

func (svc *ProductService) Create(ctx *gin.Context, input *models.ProductRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-code")
	exists, err := svc.repo.CheckProductExistByCode(ctx, utils.DerefPointer(input.Code))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}

	// If product code/tenant combo already exists, return with error
	if exists {
		cSpan.End()
		svc.logger.GetLogger().Info(fmt.Sprintf("product code [%s] already exists in tenant [%s]", utils.DerefPointer(input.Code), utils.DerefPointer(input.TenantName)))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrProductCodeAlreadyExist]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrProductCodeAlreadyExist]
		return resp, comerrors.ErrProductCodeAlreadyExist
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-product")
	productID := uuid.New()
	now := time.Now()
	product := &entities.Product{
		ID:                   productID,
		TenantID:             tenant.ID,
		Name:                 utils.DerefPointer(input.Name),
		DistributionStrategy: utils.DerefPointer(input.DistributionStrategy),
		Code:                 utils.DerefPointer(input.Code),
		Platforms:            input.Platforms,
		Metadata:             input.Metadata,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if input.Url != nil {
		product.URL = utils.DerefPointer(input.Url)
	}
	err = svc.repo.InsertNewProduct(ctx, product)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = map[string]interface{}{
		"product_id": productID,
	}

	return resp, nil
}

func (svc *ProductService) Retrieve(ctx *gin.Context, input *models.ProductRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "select-product")
	product, err := svc.repo.SelectProductByPK(ctx, tenant.ID, uuid.MustParse(utils.DerefPointer(input.ProductID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := &models.ProductRetrievalOutput{
		ID:                   product.ID,
		TenantID:             product.ID,
		Name:                 product.Name,
		DistributionStrategy: product.DistributionStrategy,
		Code:                 product.Code,
		Platforms:            product.Platforms,
		Metadata:             product.Metadata,
		URL:                  product.URL,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *ProductService) List(ctx *gin.Context, input *models.ProductListInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *ProductService) Delete(ctx *gin.Context, input *models.ProductDeletionInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *ProductService) Update(ctx *gin.Context, input *models.ProductUpdateInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *ProductService) Tokens(ctx *gin.Context, input *models.ProductTokensInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}
