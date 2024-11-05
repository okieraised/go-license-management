package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/server/v1/licenses/repository"
)

type LicenseService struct {
	repo   repository.ILicense
	logger *logging.Logger
}

func NewLicenseService(options ...func(*LicenseService)) *LicenseService {
	svc := &LicenseService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.ILicense) func(*LicenseService) {
	return func(c *LicenseService) {
		c.repo = repo
	}
}

func (svc *LicenseService) Create(ctx *gin.Context, input *models.LicenseRegistrationInput) (*response.BaseOutput, error) {
	return nil, nil
}
