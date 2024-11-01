package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/accounts/models"
	"go-license-management/internal/server/v1/accounts/repository"
	"log/slog"
)

const (
	accountGroup = "user_group"
)

type AccountService struct {
	repo   repository.IAccount
	logger *slog.Logger
}

func NewAccountService(options ...func(*AccountService)) *AccountService {
	svc := &AccountService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.GetInstance().With(slog.Group(accountGroup))
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IAccount) func(*AccountService) {
	return func(c *AccountService) {
		c.repo = repo
	}
}

func (svc *AccountService) Create(ctx *gin.Context, input *models.AccountRegistrationInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}
