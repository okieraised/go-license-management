package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/tenants/models"
	"go-license-management/internal/server/v1/tenants/repository"
	"log/slog"
)

const (
	tenantGroup = "tenant_group"
)

type TenantService struct {
	repo   repository.ITenant
	logger *slog.Logger
}

func NewTenantService(options ...func(*TenantService)) *TenantService {
	svc := &TenantService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.GetInstance().With(slog.Group(tenantGroup))
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.ITenant) func(*TenantService) {
	return func(c *TenantService) {
		c.repo = repo
	}
}

func (svc *TenantService) Create(ctx *gin.Context, input *models.TenantRegistrationInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}
