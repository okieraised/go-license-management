package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/products/models"
	"go-license-management/internal/server/v1/products/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
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

	fmt.Println("tenant", tenant)

	return resp, nil
}

func (svc *ProductService) Retrieve(ctx *gin.Context, input *models.ProductRetrievalInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
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
