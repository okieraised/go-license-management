package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/products/models"
	"go-license-management/internal/server/v1/products/repository"
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

	return &response.BaseOutput{}, nil
}
