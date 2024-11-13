package machines

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/models"
	"time"
)

type MachineRepository struct {
	database *bun.DB
}

func NewMachineRepository(ds *models.DataSource) *MachineRepository {
	return &MachineRepository{
		database: ds.GetDatabase(),
	}
}

func (repo *MachineRepository) SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error) {
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

func (repo *MachineRepository) SelectMachineByPK(ctx context.Context, machineID uuid.UUID) (*entities.Machine, error) {
	if repo.database == nil {
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	machine := &entities.Machine{ID: machineID}

	err := repo.database.NewSelect().Model(machine).WherePK().Scan(ctx)
	if err != nil {
		return machine, err
	}

	return machine, nil
}

func (repo *MachineRepository) InsertNewMachine(ctx context.Context, machine *entities.Machine) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
	}

	_, err := repo.database.NewInsert().Model(machine).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MachineRepository) DeleteMachineByPK(ctx context.Context, machineID uuid.UUID) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
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
		return comerrors.ErrInvalidDatabaseClient
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
		return false, comerrors.ErrInvalidDatabaseClient
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
		return nil, comerrors.ErrInvalidDatabaseClient
	}

	license := &entities.License{ID: licenseID}

	err := repo.database.NewSelect().Model(license).WherePK().Scan(ctx)
	if err != nil {
		return license, err
	}

	return license, nil
}

func (repo *MachineRepository) CheckMachineExistByFingerprintAndLicense(ctx context.Context, licenseID uuid.UUID, fingerprint string) (bool, error) {
	if repo.database == nil {
		return false, comerrors.ErrInvalidDatabaseClient
	}

	exists, err := repo.database.NewSelect().Model(new(entities.Machine)).
		Where("license_id = ?", licenseID).
		Where("fingerprint = ?", fingerprint).
		Exists(ctx)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *MachineRepository) InsertNewMachineAndUpdateLicense(ctx context.Context, machine *entities.Machine) error {
	if repo.database == nil {
		return comerrors.ErrInvalidDatabaseClient
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
		license.UpdatedAt = time.Now()
		license.MachinesCount += 1
	}

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
