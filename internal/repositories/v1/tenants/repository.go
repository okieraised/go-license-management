package tenants

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/models"
)

type TenantRepository struct {
	database *bun.DB
}

func NewTenantRepository(ds *models.DataSource) *TenantRepository {
	return &TenantRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *TenantRepository) SelectTenants(ctx context.Context) ([]entities.Tenant, int, error) {
	var count = 0

	if repo.database == nil {
		return nil, count, comerrors.ErrInvalidDatabaseClient
	}

	tenant := make([]entities.Tenant, 0)
	count, err := repo.database.NewSelect().Model(new(entities.Tenant)).Order("created_at DESC").ScanAndCount(ctx, &tenant)
	if err != nil {
		return tenant, count, err
	}
	return tenant, count, nil
}

func (repo *TenantRepository) SelectTenantByName(ctx context.Context, name string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{}
	err := repo.database.NewSelect().Model(tenant).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}

func (repo *TenantRepository) SelectTenantByPK(ctx context.Context, id string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{
		ID: uuid.MustParse(id),
	}
	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}

func (repo *TenantRepository) CheckTenantExistByName(ctx context.Context, name string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	exist, err := repo.database.NewSelect().Model(new(entities.Tenant)).Where("name = ?", name).Exists(ctx)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (repo *TenantRepository) InsertNewTenant(ctx context.Context, tenant *entities.Tenant) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(tenant).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TenantRepository) DeleteTenantByName(ctx context.Context, name string) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewDelete().Model(new(entities.Tenant)).Where("name = ?", name).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
