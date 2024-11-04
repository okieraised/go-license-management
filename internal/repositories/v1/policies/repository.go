package policies

import (
	"github.com/uptrace/bun"
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
