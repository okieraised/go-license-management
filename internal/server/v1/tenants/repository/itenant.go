package repository

import (
	"context"
	"go-license-management/internal/infrastructure/database/entities"
)

type ITenant interface {
	InsertNewTenant(ctx context.Context, tenant *entities.Tenant) error
	SelectTenantByPK(ctx context.Context, id string) (*entities.Tenant, error)
	SelectTenantByName(ctx context.Context, name string) (*entities.Tenant, error)
	SelectTenants(ctx context.Context) ([]entities.Tenant, int, error)
	CheckTenantExistByName(ctx context.Context, name string) (bool, error)
	DeleteTenantByName(ctx context.Context, name string) error
}
