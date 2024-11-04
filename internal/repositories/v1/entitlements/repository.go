package entitlements

import (
	"github.com/uptrace/bun"
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
