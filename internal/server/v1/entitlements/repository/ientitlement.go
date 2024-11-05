package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IEntitlement interface {
	InsertNewEntitlement(ctx context.Context, entitlement *entities.Entitlement) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectEntitlementsByTenant(ctx context.Context, tenantID uuid.UUID) ([]entities.Entitlement, int, error)
	SelectEntitlementByPK(ctx context.Context, tenantID, entitlementID uuid.UUID) (*entities.Entitlement, error)
	CheckEntitlementExistByCode(ctx context.Context, code string) (bool, error)
	DeleteEntitlementByPK(ctx context.Context, tenantID, entitlementID uuid.UUID) error
}
