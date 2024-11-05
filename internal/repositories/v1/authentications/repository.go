package authentications

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/models"
)

type AuthenticationRepository struct {
	database *bun.DB
}

func NewAuthenticationRepository(ds *models.DataSource) *AuthenticationRepository {
	return &AuthenticationRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *AuthenticationRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{}

	err := repo.database.NewSelect().Model(tenant).ColumnExpr("id, name, ed25519_private_key").Where("name = ?", tenantName).Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *AuthenticationRepository) SelectAccountByPK(ctx context.Context, tenantID uuid.UUID, username string) (*entities.Account, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	account := &entities.Account{Username: username, TenantID: tenantID}
	err := repo.database.NewSelect().Model(account).WherePK().Scan(ctx)
	if err != nil {
		return account, err
	}
	return account, nil
}
