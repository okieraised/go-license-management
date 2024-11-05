package licenses

import (
	"github.com/uptrace/bun"
	"go-license-management/server/models"
)

type LicenseRepository struct {
	database *bun.DB
}

func NewLicenseRepository(ds *models.DataSource) *LicenseRepository {
	return &LicenseRepository{
		database: ds.GetDatabase(),
	}
}

//func (repo *LicenseRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
//	if repo.database == nil {
//		return nil, comerrors.ErrInvalidDatabaseClient
//	}
//
//	tenant := &entities.Tenant{}
//
//	err := repo.database.NewSelect().Model(tenant).ColumnExpr("id, name").Where("name = ?", tenantName).Scan(ctx)
//	if err != nil {
//		return tenant, err
//	}
//
//	return tenant, nil
//}
