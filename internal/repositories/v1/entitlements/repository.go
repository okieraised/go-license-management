package entitlements

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/models"
)

type EntitlementRepository struct {
	database *bun.DB
}

func NewEntitlementRepository(ds *models.DataSource) *EntitlementRepository {
	return &EntitlementRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *EntitlementRepository) SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error) {
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

func (repo *EntitlementRepository) CheckEntitlementExistByCode(ctx context.Context, code string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	entitlement := &entities.Entitlement{
		Code: code,
	}
	exist, err := repo.database.NewSelect().Model(entitlement).Where("code = ?", code).Exists(ctx)
	if err != nil {
		return exist, err
	}

	return exist, nil
}

func (repo *EntitlementRepository) InsertNewEntitlement(ctx context.Context, entitlement *entities.Entitlement) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(entitlement).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EntitlementRepository) SelectEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) (*entities.Entitlement, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	entitlement := &entities.Entitlement{ID: entitlementID}
	err := repo.database.NewSelect().Model(entitlement).WherePK().Scan(ctx)
	if err != nil {
		return entitlement, err
	}
	return entitlement, nil
}

func (repo *EntitlementRepository) DeleteEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	entitlement := &entities.Entitlement{ID: entitlementID}
	_, err := repo.database.NewDelete().Model(entitlement).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EntitlementRepository) SelectEntitlementsByTenant(ctx context.Context, tenantName string, param constants.QueryCommonParam) ([]entities.Entitlement, int, error) {
	var total = 0
	if repo.database == nil {
		return nil, total, comerrors.ErrInvalidDatabaseClient
	}

	entitlements := make([]entities.Entitlement, 0)
	total, err := repo.database.NewSelect().Model(new(entities.Entitlement)).
		Where("tenant_name = ?", tenantName).
		Limit(utils.DerefPointer(param.Limit)).
		Offset(utils.DerefPointer(param.Offset)).
		Order("created_at DESC").
		ScanAndCount(ctx, &entitlements)
	if err != nil {
		return entitlements, total, comerrors.ErrInvalidDatabaseClient
	}
	return entitlements, total, nil
}
