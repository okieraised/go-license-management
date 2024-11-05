package policies

import (
	"context"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/models"
)

type PolicyRepository struct {
	database *bun.DB
}

func NewPolicyRepository(ds *models.DataSource) *PolicyRepository {
	return &PolicyRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *PolicyRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
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

func (repo *PolicyRepository) InsertNewPolicy(ctx context.Context, policy *entities.Policy) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(policy).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
