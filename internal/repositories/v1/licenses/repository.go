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
