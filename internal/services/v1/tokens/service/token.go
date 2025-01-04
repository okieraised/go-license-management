package service

import (
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/services/v1/tokens/repository"
)

type TokenService struct {
	repo   repository.IToken
	logger *logging.Logger
}

func NewTokenService(options ...func(*TokenService)) *TokenService {
	svc := &TokenService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IToken) func(*TokenService) {
	return func(c *TokenService) {
		c.repo = repo
	}
}
