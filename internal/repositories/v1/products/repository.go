package products

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/api"
	"time"
)

type ProductRepository struct {
	database *bun.DB
}

func NewProductRepository(ds *api.DataSource) *ProductRepository {
	return &ProductRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *ProductRepository) SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: tenantName}

	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *ProductRepository) InsertNewProduct(ctx context.Context, product *entities.Product) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) SelectProductByPK(ctx context.Context, productID uuid.UUID) (*entities.Product, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID}
	err := repo.database.NewSelect().Model(product).WherePK().Scan(ctx)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (repo *ProductRepository) SelectProducts(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Product, int, error) {
	var total = 0

	if repo.database == nil {
		return nil, total, cerrors.ErrInvalidDatabaseClient
	}

	products := make([]entities.Product, 0)
	total, err := repo.database.NewSelect().Model(new(entities.Product)).
		Where("tenant_name = ?", tenantName).
		Order("created_at DESC").
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		ScanAndCount(ctx, &products)
	if err != nil {
		return products, total, err
	}
	return products, total, nil
}

func (repo *ProductRepository) CheckProductExistByCode(ctx context.Context, code string) (bool, error) {
	if repo.database == nil {
		return false, cerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{Code: code}
	exist, err := repo.database.NewSelect().Model(product).Where("code = ?", code).Exists(ctx)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (repo *ProductRepository) DeleteProductByPK(ctx context.Context, productID uuid.UUID) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID}
	_, err := repo.database.NewDelete().Model(product).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) UpdateProductByPK(ctx context.Context, product *entities.Product) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	product.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(product).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) InsertNewProductToken(ctx context.Context, productToken *entities.ProductToken) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(productToken).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
