package products

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/database/entities"
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

func (repo *ProductRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
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

func (repo *ProductRepository) InsertNewProduct(ctx context.Context, product *entities.Product) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) SelectProductByPK(ctx context.Context, tenantID, productID uuid.UUID) (*entities.Product, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID, TenantID: tenantID}
	err := repo.database.NewSelect().Model(product).Scan(ctx)
	if err != nil {
		return product, err
	}
	return product, nil
}
