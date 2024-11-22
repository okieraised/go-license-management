package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IPolicy interface {
	InsertNewPolicy(ctx context.Context, policy *entities.Policy) error
	InsertNewPolicyEntitlement(ctx context.Context, policyEntitlement *entities.PolicyEntitlement) error
	InsertNewPolicyEntitlements(ctx context.Context, policyEntitlement []entities.PolicyEntitlement) error
	UpdatePolicyByPK(ctx context.Context, policy *entities.Policy) error
	SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error)
	SelectEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) (*entities.Entitlement, error)
	SelectEntitlementsByPK(ctx context.Context, entitlementID []uuid.UUID) ([]entities.Entitlement, error)
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	SelectPolicies(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Policy, int, error)
	CheckProductExistByID(ctx context.Context, productID uuid.UUID) (bool, error)
	CheckPolicyEntitlementExistsByPolicyIDAndEntitlementID(ctx context.Context, policyID, entitlementID uuid.UUID) (bool, error)
	DeletePolicyByPK(ctx context.Context, policyID uuid.UUID) error
	DeletePolicyEntitlementByPK(ctx context.Context, policyEntitlementID uuid.UUID) error
	DeletePolicyEntitlementsByPK(ctx context.Context, policyEntitlementID []uuid.UUID) error
}
