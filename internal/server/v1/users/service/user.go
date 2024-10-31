package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/users/models"
	"go-license-management/internal/server/v1/users/repository"
	"log/slog"
)

const (
	userGroup = "user_group"
)

type UserService struct {
	repo   repository.IUser
	logger *slog.Logger
}

func NewUserService(options ...func(*UserService)) *UserService {
	ekycSvc := &UserService{}

	for _, opt := range options {
		opt(ekycSvc)
	}
	logger := logging.GetInstance().With(slog.Group(userGroup))
	ekycSvc.logger = logger

	return ekycSvc
}

func WithRepository(repo repository.IUser) func(*UserService) {
	return func(c *UserService) {
		c.repo = repo
	}
}

func (svc *UserService) Create(ctx *gin.Context, input *models.UserRegistrationInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}
