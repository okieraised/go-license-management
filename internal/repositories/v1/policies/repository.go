package policies

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/models"
)

type PolicyRepository struct {
	database *bun.DB
}

func NewPolicyRepository(ds *models.DataSource) *PolicyRepository {
	return &PolicyRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *PolicyRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{}

	err := repo.database.NewSelect().Model(tenant).ColumnExpr("id, name").Where("name = ?", tenantName).Scan(ctx)
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

func (repo *PolicyRepository) CheckProductExistByID(ctx context.Context, tenantID, productID uuid.UUID) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	product := &entities.Product{ID: productID, TenantName: tenantID}

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

func (repo *PolicyRepository) SelectPolicies(ctx context.Context, tenantID uuid.UUID, queryParam constants.QueryCommonParam) ([]entities.Policy, int, error) {
	var total = 0
	if repo.database == nil {
		return nil, total, comerrors.ErrInvalidDatabaseClient
	}

	policies := make([]entities.Policy, 0)
	total, err := repo.database.NewSelect().Model(new(entities.Policy)).
		Where("tenant_id = ?", tenantID.String()).
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		Order("created_at DESC").
		ScanAndCount(ctx, &policies)
	if err != nil {
		return policies, total, err
	}
	return policies, total, nil
}
