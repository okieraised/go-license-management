package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IPolicy interface {
	InsertNewPolicy(ctx context.Context, policy *entities.Policy) error
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckProductExistByID(ctx context.Context, tenantID, productID uuid.UUID) (bool, error)
}
