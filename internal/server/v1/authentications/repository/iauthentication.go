package repository

import (
	"context"
	"go-license-management/internal/infrastructure/database/entities"
)

type IAuthentication interface {
	SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectAccountByPK(ctx context.Context, tenantName, username string) (*entities.Account, error)
}
