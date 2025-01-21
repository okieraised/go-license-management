package machines

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"go-license-management/server/api"
	"time"
)

type MachineRepository struct {
	database *bun.DB
}

func NewMachineRepository(ds *api.DataSource) *MachineRepository {
	return &MachineRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *MachineRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	tenant := &entities.Tenant{Name: tenantName}

	err := repo.database.NewSelect().Model(tenant).WherePK().Scan(ctx)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

func (repo *MachineRepository) SelectMachineByPK(ctx context.Context, machineID uuid.UUID) (*entities.Machine, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	machine := &entities.Machine{ID: machineID}

	err := repo.database.NewSelect().Model(machine).Relation("License").WherePK().Scan(ctx)
	if err != nil {
		return machine, err
	}

	return machine, nil
}

func (repo *MachineRepository) InsertNewMachine(ctx context.Context, machine *entities.Machine) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(machine).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MachineRepository) SelectMachines(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Machine, int, error) {
	var total = 0

	if repo.database == nil {
		return nil, total, cerrors.ErrInvalidDatabaseClient
	}

	machines := make([]entities.Machine, 0)
	total, err := repo.database.NewSelect().Model(new(entities.Machine)).
		Where("tenant_name = ?", tenantName).
		Order("created_at DESC").
		Limit(utils.DerefPointer(queryParam.Limit)).
		Offset(utils.DerefPointer(queryParam.Offset)).
		ScanAndCount(ctx, &machines)
	if err != nil {
		return machines, total, err
	}
	return machines, total, nil
}

func (repo *MachineRepository) DeleteMachineByPK(ctx context.Context, machineID uuid.UUID) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	machine := &entities.Machine{ID: machineID}

	_, err := repo.database.NewDelete().Model(machine).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MachineRepository) DeleteMachineByPKAndUpdateLicense(ctx context.Context, machineID uuid.UUID) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	tx, err := repo.database.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		cErr := tx.Commit()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	machine := &entities.Machine{ID: machineID}

	err = tx.NewSelect().Model(machine).WherePK().Scan(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	license := &entities.License{ID: machine.LicenseID}
	err = tx.NewSelect().Model(license).WherePK().Scan(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	license.UpdatedAt = time.Now()
	license.MachinesCount -= 1
	if license.MachinesCount == 0 {
		license.Status = constants.LicenseStatusInactive
	}
	_, err = tx.NewUpdate().Model(license).WherePK().Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NewDelete().Model(machine).WherePK().Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (repo *MachineRepository) CheckLicenseExistByPK(ctx context.Context, licenseID uuid.UUID) (bool, error) {
	if repo.database == nil {
		return false, cerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{ID: licenseID}

	exists, err := repo.database.NewSelect().Model(license).WherePK().Exists(ctx)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *MachineRepository) SelectLicenseByPK(ctx context.Context, licenseID uuid.UUID) (*entities.License, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{ID: licenseID}

	err := repo.database.NewSelect().Model(license).WherePK().Scan(ctx)
	if err != nil {
		return license, err
	}

	return license, nil
}

func (repo *MachineRepository) SelectLicenseByLicenseKey(ctx context.Context, licenseKey string) (*entities.License, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{Key: licenseKey}

	err := repo.database.NewSelect().Model(license).Relation("Policy").Relation("Product").Where("key = ?", licenseKey).Scan(ctx)
	if err != nil {
		return license, err
	}

	return license, nil
}

func (repo *MachineRepository) SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error) {
	if repo.database == nil {
		return nil, cerrors.ErrInvalidDatabaseClient
	}

	policy := &entities.Policy{ID: policyID}

	err := repo.database.NewSelect().Model(policy).WherePK().Scan(ctx)
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func (repo *MachineRepository) CheckMachineExistByFingerprintAndLicense(ctx context.Context, licenseKey, fingerprint string) (bool, error) {
	if repo.database == nil {
		return false, cerrors.ErrInvalidDatabaseClient
	}

	exists, err := repo.database.NewSelect().Model(new(entities.Machine)).
		Where("license_key = ?", licenseKey).
		Where("fingerprint = ?", fingerprint).
		Exists(ctx)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *MachineRepository) InsertNewMachineAndUpdateLicense(ctx context.Context, machine *entities.Machine) error {
	if repo.database == nil {
		return cerrors.ErrInvalidDatabaseClient
	}

	tx, err := repo.database.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		cErr := tx.Commit()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	license := &entities.License{ID: machine.LicenseID}
	err = tx.NewSelect().Model(license).WherePK().Scan(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if license.Status == constants.LicenseStatusNotActivated || license.Status == constants.LicenseStatusInactive {
		license.Status = constants.LicenseStatusActive
	}
	license.UpdatedAt = time.Now()
	license.MachinesCount += 1

	_, err = tx.NewUpdate().Model(license).WherePK().Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NewInsert().Model(machine).Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (repo *MachineRepository) UpdateMachineByPK(ctx context.Context, machine *entities.Machine) (*entities.Machine, error) {
	if repo.database == nil {
		return machine, cerrors.ErrInvalidDatabaseClient
	}

	machine.UpdatedAt = time.Now()
	_, err := repo.database.NewUpdate().Model(machine).WherePK().Exec(ctx)
	if err != nil {
		return machine, err
	}
	return machine, nil
}

func (repo *MachineRepository) UpdateMachineByPKAndLicense(ctx context.Context, machine *entities.Machine, currentLicense, newLicense *entities.License) (*entities.Machine, error) {
	if repo.database == nil {
		return machine, cerrors.ErrInvalidDatabaseClient
	}

	tx, err := repo.database.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		cErr := tx.Commit()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	// Only perform update on license if the new license is not nil
	if newLicense != nil {
		currentLicense.UpdatedAt = time.Now()
		currentLicense.MachinesCount -= 1
		if currentLicense.MachinesCount == 0 {
			currentLicense.Status = constants.LicenseStatusInactive
		}
		_, err = tx.NewUpdate().Model(currentLicense).WherePK().Exec(ctx)
		if err != nil {
			_ = tx.Rollback()
			return machine, err
		}

		if newLicense.Status == constants.LicenseStatusNotActivated || newLicense.Status == constants.LicenseStatusInactive {
			newLicense.Status = constants.LicenseStatusActive
		}
		newLicense.UpdatedAt = time.Now()
		newLicense.MachinesCount += 1
		_, err = tx.NewUpdate().Model(newLicense).WherePK().Exec(ctx)
		if err != nil {
			_ = tx.Rollback()
			return machine, err
		}
	}

	machine.UpdatedAt = time.Now()
	_, err = tx.NewUpdate().Model(machine).WherePK().Exec(ctx)
	if err != nil {
		return machine, err
	}
	return machine, nil
}
