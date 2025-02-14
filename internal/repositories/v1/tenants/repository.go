package tenants

import (
	"context"
	"github.com/uptrace/bun"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/api"
	"time"
)

type TenantRepository struct {
	database *bun.DB
}

func NewTenantRepository(ds *api.DataSource) *TenantRepository {
	return &TenantRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *TenantRepository) SelectTenants(ctx context.Context, queryParam constants.QueryCommonParam) ([]entities.Tenant, int, error) {
	var count = 0

	if repo.database == nil {
		return nil, count, cerrors.ErrInvalidDatabaseClient
	}

	tenant := make([]entities.Tenant, 0)
	count, err := repo.database.NewSelect().Model(new(entities.Tenant)).
		Order("created_at DESC").
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		ScanAndCount(ctx, &tenant)
	if err != nil {
		return tenant, count, err
	}
	return tenant, count, nil
}

func (repo *TenantRepository) SelectTenantByPK(ctx context.Context, name string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{
		Name: name,
	}
	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}

func (repo *TenantRepository) CheckTenantExistByPK(ctx context.Context, name string) (bool, error) {
	if repo.database == nil {
		return false, cerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: name}

	exist, err := repo.database.NewSelect().Model(tenant).WherePK().Exists(ctx)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (repo *TenantRepository) InsertNewTenant(ctx context.Context, tenant *entities.Tenant) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(tenant).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TenantRepository) DeleteTenantByPK(ctx context.Context, name string) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: name}

	_, err := repo.database.NewDelete().Model(tenant).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TenantRepository) UpdateTenantByPK(ctx context.Context, tenant *entities.Tenant) (*entities.Tenant, error) {
	if repo.database == nil {
		return tenant, cerrors.ErrInvalidDatabaseClient
	}

	tenant.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(tenant).WherePK().Exec(ctx)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}
