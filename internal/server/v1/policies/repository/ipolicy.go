package repository

import (
	"context"
	"go-license-management/internal/infrastructure/database/entities"
)

type IPolicy interface {
	InsertNewPolicy(ctx context.Context, policy *entities.Policy) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
}
