package repository

import (
	"context"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IAccount interface {
	InsertNewAccount(ctx context.Context, account *entities.Account) error
	UpdateAccountByPK(ctx context.Context, account *entities.Account) error
	SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectAccountsByTenant(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Account, int, error)
	SelectAccountByPK(ctx context.Context, tenantName, username string) (*entities.Account, error)
	CheckAccountExistByPK(ctx context.Context, tenantName, username string) (bool, error)
	DeleteAccountByPK(ctx context.Context, tenantName, username string) error
}
