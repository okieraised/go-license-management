package accounts

import (
	"context"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/api"
	"time"
)

type AccountRepository struct {
	database *bun.DB
}

func NewAccountRepository(ds *api.DataSource) *AccountRepository {
	return &AccountRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *AccountRepository) SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: tenantName}

	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *AccountRepository) SelectAccountsByTenant(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Account, int, error) {
	var count = 0
	if repo.database == nil {
		return nil, count, comerrors.ErrInvalidDatabaseClient
	}

	accounts := make([]entities.Account, 0)
	count, err := repo.database.NewSelect().Model(new(entities.Account)).
		Where("tenant_name = ?", tenantName).
		Order("created_at DESC").
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		ScanAndCount(ctx, &accounts)
	if err != nil {
		return accounts, count, nil
	}

	return accounts, count, nil
}

func (repo *AccountRepository) SelectAccountByPK(ctx context.Context, tenantName, username string) (*entities.Account, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Username: username, TenantName: tenantName}
	err := repo.database.NewSelect().Model(account).WherePK().Scan(ctx)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (repo *AccountRepository) UpdateAccountByPK(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	if repo.database == nil {
		return account, comerrors.ErrInvalidDatabaseClient
	}

	account.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(account).WherePK().Exec(ctx)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (repo *AccountRepository) CheckAccountExistByPK(ctx context.Context, tenantName, username string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Username: username, TenantName: tenantName}
	exist, err := repo.database.NewSelect().Model(account).WherePK().Exists(ctx)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (repo *AccountRepository) CheckAccountEmailExistByPK(ctx context.Context, tenantName, email string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Email: email, TenantName: tenantName}
	exist, err := repo.database.NewSelect().Model(account).Where("tenant_name = ? AND email = ?", tenantName, email).Exists(ctx)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (repo *AccountRepository) InsertNewAccount(ctx context.Context, account *entities.Account) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(account).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AccountRepository) DeleteAccountByPK(ctx context.Context, tenantName, username string) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Username: username, TenantName: tenantName}
	_, err := repo.database.NewDelete().Model(account).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
