package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IProduct interface {
	InsertNewProduct(ctx context.Context, product *entities.Product) error
	UpdateProductByPK(ctx context.Context, product *entities.Product) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckProductExistByCode(ctx context.Context, code string) (bool, error)
	SelectProductByPK(ctx context.Context, tenantID, productID uuid.UUID) (*entities.Product, error)
	DeleteProductByPK(ctx context.Context, tenantID, productID uuid.UUID) error
}
