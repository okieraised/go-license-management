package service

import (
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/server/v1/machines/repository"
)

type MachineService struct {
	repo   repository.IMachine
	logger *logging.Logger
}

func NewMachineService(options ...func(*MachineService)) *MachineService {
	svc := &MachineService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IMachine) func(*MachineService) {
	return func(c *MachineService) {
		c.repo = repo
	}
}
