package licenses

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/models"
	"time"
)

type LicenseRepository struct {
	database *bun.DB
}

func NewLicenseRepository(ds *models.DataSource) *LicenseRepository {
	return &LicenseRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *LicenseRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
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

func (repo *LicenseRepository) SelectProductByPK(ctx context.Context, productID uuid.UUID) (*entities.Product, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID}

	err := repo.database.NewSelect().Model(product).WherePK().Scan(ctx)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (repo *LicenseRepository) SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	policy := &entities.Policy{ID: policyID}

	err := repo.database.NewSelect().Model(policy).WherePK().Scan(ctx)
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func (repo *LicenseRepository) InsertNewLicense(ctx context.Context, license *entities.License) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(license).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *LicenseRepository) SelectLicenseByPK(ctx context.Context, licenseID uuid.UUID) (*entities.License, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{ID: licenseID}

	err := repo.database.NewSelect().Model(license).Relation("Policy").WherePK().Scan(ctx)
	if err != nil {
		return license, err
	}

	return license, nil
}

func (repo *LicenseRepository) SelectLicenses(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.License, int, error) {
	var total = 0

	if repo.database == nil {
		return nil, total, comerrors.ErrInvalidDatabaseClient
	}

	licenses := make([]entities.License, 0)
	total, err := repo.database.NewSelect().Model(new(entities.License)).
		Where("tenant_name = ?", tenantName).
		Order("created_at DESC").
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		ScanAndCount(ctx, &licenses)
	if err != nil {
		return licenses, total, err
	}
	return licenses, total, nil
}

func (repo *LicenseRepository) SelectLicenseByLicenseKey(ctx context.Context, licenseKey string) (*entities.License, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{Key: licenseKey}

	err := repo.database.NewSelect().Model(license).Relation("Policy").Where("key = ?", licenseKey).Scan(ctx)
	if err != nil {
		return license, err
	}

	return license, nil
}

func (repo *LicenseRepository) DeleteLicenseByPK(ctx context.Context, licenseID uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{ID: licenseID}

	_, err := repo.database.NewDelete().Model(license).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *LicenseRepository) UpdateLicenseByPK(ctx context.Context, license *entities.License) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	license.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(license).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
