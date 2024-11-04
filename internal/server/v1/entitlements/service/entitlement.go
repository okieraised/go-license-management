package service

import (
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/server/v1/entitlements/repository"
)

type EntitlementService struct {
	repo   repository.IEntitlement
	logger *logging.Logger
}

func NewEntitlementService(options ...func(*EntitlementService)) *EntitlementService {
	svc := &EntitlementService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IEntitlement) func(*EntitlementService) {
	return func(c *EntitlementService) {
		c.repo = repo
	}
}
