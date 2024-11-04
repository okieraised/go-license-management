package products

import (
	"github.com/uptrace/bun"
	"go-license-management/server/models"
)

type ProductRepository struct {
	database *bun.DB
}

func NewProductRepository(ds *models.DataSource) *ProductRepository {
	return &ProductRepository{
		database: ds.GetDatabase(),
	}
}
