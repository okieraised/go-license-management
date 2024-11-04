package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IAccount interface {
	InsertNewAccount(ctx context.Context, account *entities.Account) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckAccountExistByPK(ctx context.Context, tenantID uuid.UUID, username string) (bool, error)
}
