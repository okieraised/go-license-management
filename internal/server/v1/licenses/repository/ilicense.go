package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type ILicense interface {
	InsertNewLicense(ctx context.Context, license *entities.License) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectProductByPK(ctx context.Context, tenantID, productID uuid.UUID) (*entities.Product, error)
	SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error)
}
