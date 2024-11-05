package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IAuthentication interface {
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectAccountByPK(ctx context.Context, tenantID uuid.UUID, username string) (*entities.Account, error)
}
