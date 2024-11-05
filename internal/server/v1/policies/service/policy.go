package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/policies/models"
	"go-license-management/internal/server/v1/policies/repository"
)

type PolicyService struct {
	repo   repository.IPolicy
	logger *logging.Logger
}

func NewPolicyService(options ...func(*PolicyService)) *PolicyService {
	svc := &PolicyService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IPolicy) func(*PolicyService) {
	return func(c *PolicyService) {
		c.repo = repo
	}
}

func (svc *PolicyService) Create(ctx *gin.Context, input *models.PolicyRegistrationInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) List(ctx *gin.Context, input *models.PolicyListInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Retrieve(ctx *gin.Context, input *models.PolicyRetrievalInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Delete(ctx *gin.Context, input *models.PolicyDeletionInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Update(ctx *gin.Context, input *models.PolicyUpdateInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Attach(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Detach(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) ListEntitlements(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}
