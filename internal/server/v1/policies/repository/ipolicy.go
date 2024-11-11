package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IPolicy interface {
	InsertNewPolicy(ctx context.Context, policy *entities.Policy) error
	SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error)
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectPolicies(ctx context.Context, tenantID uuid.UUID, queryParam constants.QueryCommonParam) ([]entities.Policy, int, error)
	CheckProductExistByID(ctx context.Context, tenantID, productID uuid.UUID) (bool, error)
	DeletePolicyByPK(ctx context.Context, policyID uuid.UUID) error
}
