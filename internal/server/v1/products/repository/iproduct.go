package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IProduct interface {
	InsertNewProduct(ctx context.Context, product *entities.Product) error
	InsertNewProductToken(ctx context.Context, productToken *entities.ProductToken) error
	UpdateProductByPK(ctx context.Context, product *entities.Product) error
	SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckProductExistByCode(ctx context.Context, code string) (bool, error)
	SelectProductByPK(ctx context.Context, productID uuid.UUID) (*entities.Product, error)
	SelectProducts(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Product, int, error)
	DeleteProductByPK(ctx context.Context, productID uuid.UUID) error
}
