package policies

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/api"
	"time"
)

type PolicyRepository struct {
	database *bun.DB
}

func NewPolicyRepository(ds *api.DataSource) *PolicyRepository {
	return &PolicyRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *PolicyRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: tenantName}

	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *PolicyRepository) SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	policy := &entities.Policy{ID: policyID}

	err := repo.database.NewSelect().Model(policy).WherePK().Scan(ctx)
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func (repo *PolicyRepository) CheckPolicyEntitlementExistsByPolicyIDAndEntitlementID(ctx context.Context, policyID, entitlementID uuid.UUID) (bool, error) {
	var err error
	exists := false

	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	exists, err = repo.database.NewSelect().
		Model(new(entities.PolicyEntitlement)).
		Where("policy_id = ? AND entitlement_id = ?", policyID, entitlementID).
		Exists(ctx)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *PolicyRepository) SelectEntitlementByPK(ctx context.Context, entitlementID uuid.UUID) (*entities.Entitlement, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	entitlement := &entities.Entitlement{ID: entitlementID}

	err := repo.database.NewSelect().Model(entitlement).WherePK().Scan(ctx)
	if err != nil {
		return entitlement, err
	}

	return entitlement, nil
}

func (repo *PolicyRepository) SelectEntitlementsByPK(ctx context.Context, entitlementID []uuid.UUID) ([]entities.Entitlement, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	entitlements := make([]entities.Entitlement, 0)
	for _, id := range entitlementID {
		entitlements = append(entitlements, entities.Entitlement{ID: id})
	}
	err := repo.database.NewSelect().Model(&entitlements).WherePK().Scan(ctx)
	if err != nil {
		return entitlements, err
	}

	return entitlements, nil
}

func (repo *PolicyRepository) InsertNewPolicyEntitlement(ctx context.Context, policyEntitlement *entities.PolicyEntitlement) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(policyEntitlement).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PolicyRepository) InsertNewPolicyEntitlements(ctx context.Context, policyEntitlement []entities.PolicyEntitlement) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(&policyEntitlement).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PolicyRepository) DeletePolicyByPK(ctx context.Context, policyID uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	policy := &entities.Policy{ID: policyID}

	_, err := repo.database.NewDelete().Model(policy).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PolicyRepository) DeletePolicyEntitlementByPK(ctx context.Context, policyEntitlementID uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	policyEntitlement := &entities.PolicyEntitlement{ID: policyEntitlementID}

	_, err := repo.database.NewDelete().Model(policyEntitlement).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PolicyRepository) SelectPolicyEntitlements(ctx context.Context, policyID uuid.UUID, queryParam constants.QueryCommonParam) ([]entities.PolicyEntitlement, int, error) {
	var total int

	if repo.database == nil {
		return nil, total, comerrors.ErrInvalidDatabaseClient
	}

	policyEntitlements := make([]entities.PolicyEntitlement, 0)
	total, err := repo.database.NewSelect().Model(new(entities.PolicyEntitlement)).
		Where("policy_id = ?", policyID).
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		Order("created_at DESC").
		ScanAndCount(ctx, &policyEntitlements)
	if err != nil {
		return policyEntitlements, total, err
	}
	return policyEntitlements, total, err
}

func (repo *PolicyRepository) DeletePolicyEntitlementsByPK(ctx context.Context, policyEntitlementID []uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	policyEntitlements := make([]entities.PolicyEntitlement, 0)
	for _, id := range policyEntitlementID {
		policyEntitlements = append(policyEntitlements, entities.PolicyEntitlement{ID: id})
	}

	_, err := repo.database.NewDelete().Model(&policyEntitlements).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PolicyRepository) CheckProductExistByID(ctx context.Context, productID uuid.UUID) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID}

	exists, err := repo.database.NewSelect().Model(product).WherePK().Exists(ctx)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *PolicyRepository) InsertNewPolicy(ctx context.Context, policy *entities.Policy) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(policy).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PolicyRepository) UpdatePolicyByPK(ctx context.Context, policy *entities.Policy) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	policy.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(policy).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PolicyRepository) SelectPolicies(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Policy, int, error) {
	var total = 0
	if repo.database == nil {
		return nil, total, comerrors.ErrInvalidDatabaseClient
	}

	policies := make([]entities.Policy, 0)
	total, err := repo.database.NewSelect().Model(new(entities.Policy)).
		Where("tenant_name = ?", tenantName).
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		Order("created_at DESC").
		ScanAndCount(ctx, &policies)
	if err != nil {
		return policies, total, err
	}
	return policies, total, nil
}
