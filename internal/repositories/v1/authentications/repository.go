package authentications

import (
	"context"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/api"
)

type AuthenticationRepository struct {
	database *bun.DB
}

func NewAuthenticationRepository(ds *api.DataSource) *AuthenticationRepository {
	return &AuthenticationRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *AuthenticationRepository) SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: tenantName}

	err := repo.database.NewSelect().Model(tenant).ColumnExpr("name, ed25519_private_key").WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *AuthenticationRepository) SelectAccountByPK(ctx context.Context, tenantName, username string) (*entities.Account, error) {
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

func (repo *AuthenticationRepository) SelectMasterByPK(ctx context.Context, username string) (*entities.Master, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	master := &entities.Master{Username: username}
	err := repo.database.NewSelect().Model(master).WherePK().Scan(ctx)
	if err != nil {
		return master, err
	}
	return master, nil
}
