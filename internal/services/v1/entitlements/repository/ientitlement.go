package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IEntitlement interface {
	InsertNewEntitlement(ctx context.Context, entitlement *entities.Entitlement) error
	SelectTenantByPK(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectEntitlementsByTenant(ctx context.Context, tenantName string, param constants.QueryCommonParam) ([]entities.Entitlement, int, error)
	SelectEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) (*entities.Entitlement, error)
	CheckEntitlementExistByCode(ctx context.Context, code string) (bool, error)
	DeleteEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) error
}
