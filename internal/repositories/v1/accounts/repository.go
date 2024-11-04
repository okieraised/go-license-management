package accounts

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/models"
)

type AccountRepository struct {
	database *bun.DB
}

func NewAccountRepository(ds *models.DataSource) *AccountRepository {
	return &AccountRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *AccountRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{}

	err := repo.database.NewSelect().Model(tenant).ColumnExpr("id, name").Where("name = ?", tenantName).Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *AccountRepository) CheckAccountExistByPK(ctx context.Context, tenantID uuid.UUID, username string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Username: username, TenantID: tenantID}
	exist, err := repo.database.NewSelect().Model(account).WherePK().Exists(ctx)
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
