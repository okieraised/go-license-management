package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/products/models"
	"go-license-management/internal/services/v1/products/repository"
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
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-code")
	exists, err := svc.repo.CheckProductExistByCode(ctx, utils.DerefPointer(input.Code))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}

	// If product code/tenant combo already exists, return with error
	if exists {
		cSpan.End()
		svc.logger.GetLogger().Info(fmt.Sprintf("product code [%s] already exists in tenant [%s]", utils.DerefPointer(input.Code), utils.DerefPointer(input.TenantName)))
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductCodeAlreadyExist]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductCodeAlreadyExist]
		return resp, cerrors.ErrProductCodeAlreadyExist
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-product")
	svc.logger.GetLogger().Info("inserting new product to database")
	productID := uuid.New()
	now := time.Now()
	product := &entities.Product{
		ID:                   productID,
		TenantName:           tenant.Name,
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
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.ProductRegistrationOutput{
		ID:                   productID.String(),
		TenantName:           tenant.Name,
		Name:                 product.Name,
		DistributionStrategy: product.DistributionStrategy,
		Code:                 product.Code,
		URL:                  product.URL,
		Platforms:            product.Platforms,
		Metadata:             product.Metadata,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *ProductService) Retrieve(ctx *gin.Context, input *models.ProductRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "select-product")
	product, err := svc.repo.SelectProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
			return resp, cerrors.ErrProductIDIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.ProductRetrievalOutput{
		ID:                   product.ID,
		TenantName:           tenant.Name,
		Name:                 product.Name,
		DistributionStrategy: product.DistributionStrategy,
		Code:                 product.Code,
		Platforms:            product.Platforms,
		Metadata:             product.Metadata,
		URL:                  product.URL,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *ProductService) List(ctx *gin.Context, input *models.ProductListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-pkc")
	products, total, err := svc.repo.SelectProducts(ctx, tenant.Name, input.QueryCommonParam)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
			return resp, cerrors.ErrProductIDIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	productOutput := make([]models.ProductRetrievalOutput, 0)
	for _, product := range products {
		productOutput = append(productOutput, models.ProductRetrievalOutput{
			ID:                   product.ID,
			TenantName:           tenant.Name,
			Name:                 product.Name,
			DistributionStrategy: product.DistributionStrategy,
			Code:                 product.Code,
			Platforms:            product.Platforms,
			Metadata:             product.Metadata,
			URL:                  product.URL,
			CreatedAt:            product.CreatedAt,
			UpdatedAt:            product.UpdatedAt,
		})
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = productOutput

	return resp, nil
}

func (svc *ProductService) Delete(ctx *gin.Context, input *models.ProductDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	_, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "delete-product")
	err = svc.repo.DeleteProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	return resp, nil
}

func (svc *ProductService) Update(ctx *gin.Context, input *models.ProductUpdateInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-pkc")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", utils.DerefPointer(input.ProductID)))
	product, err := svc.repo.SelectProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "update-existing-product")
	svc.logger.GetLogger().Info(fmt.Sprintf("updating product [%s]", utils.DerefPointer(input.ProductID)))
	if input.Code != nil {
		exists, err := svc.repo.CheckProductExistByCode(ctx, utils.DerefPointer(input.Code))
		if err != nil {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		if !exists {
			product.Code = utils.DerefPointer(input.Code)
		} else {
			svc.logger.GetLogger().Info(fmt.Sprintf("product code [%s] already exists", product.Code))
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductCodeAlreadyExist]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductCodeAlreadyExist]
			return resp, cerrors.ErrProductCodeAlreadyExist
		}
	}

	if input.Name != nil {
		product.Name = utils.DerefPointer(input.Name)
	}

	if input.Url != nil {
		product.URL = utils.DerefPointer(input.Url)
	}

	if input.Platforms != nil {
		product.Platforms = input.Platforms
	}

	if input.Metadata != nil {
		product.Metadata = input.Metadata
	}

	if input.DistributionStrategy != nil {
		product.DistributionStrategy = utils.DerefPointer(input.DistributionStrategy)
	}

	err = svc.repo.UpdateProductByPK(ctx, product)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.ProductUpdateOutput{
		ID:                   product.ID.String(),
		TenantName:           tenant.Name,
		Name:                 product.Name,
		DistributionStrategy: product.DistributionStrategy,
		Code:                 product.Code,
		URL:                  product.URL,
		Platforms:            product.Platforms,
		Metadata:             product.Metadata,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *ProductService) Tokens(ctx *gin.Context, input *models.ProductTokensInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	// Check tenant
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	// Check product
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", input.ProductID))
	_, cSpan = input.Tracer.Start(rootCtx, "query-product")
	product, err := svc.repo.SelectProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
			return resp, cerrors.ErrProductIDIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	// Generate token
	_, cSpan = input.Tracer.Start(rootCtx, "generate-product-token")
	token := utils.GenerateToken()
	id := uuid.New()
	productToken := &entities.ProductToken{
		ID:         id,
		TenantName: tenant.Name,
		ProductID:  product.ID,
		Token:      token,
		CreatedAt:  time.Now(),
	}
	err = svc.repo.InsertNewProductToken(ctx, productToken)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = models.ProductTokenOutput{
		ID:    id.String(),
		Token: token,
	}
	return resp, nil
}
