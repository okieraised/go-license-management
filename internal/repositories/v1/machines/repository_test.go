package machines

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/server/api"
)

func setupTestDB(t *testing.T) (*bun.DB, func()) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	// Create tables
	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*entities.Tenant)(nil)).IfNotExists().Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewCreateTable().Model((*entities.Product)(nil)).IfNotExists().Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewCreateTable().Model((*entities.Policy)(nil)).IfNotExists().Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewCreateTable().Model((*entities.License)(nil)).IfNotExists().Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewCreateTable().Model((*entities.Machine)(nil)).IfNotExists().Exec(ctx)
	require.NoError(t, err)

	cleanup := func() {
		_ = db.Close()
	}

	return db, cleanup
}

func createTestTenant(t *testing.T, db *bun.DB, name string) *entities.Tenant {
	tenant := &entities.Tenant{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := db.NewInsert().Model(tenant).Exec(context.Background())
	require.NoError(t, err)
	return tenant
}

func createTestProduct(t *testing.T, db *bun.DB, tenantName string) *entities.Product {
	product := &entities.Product{
		ID:         uuid.New(),
		TenantName: tenantName,
		Name:       "Test Product",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := db.NewInsert().Model(product).Exec(context.Background())
	require.NoError(t, err)
	return product
}

func createTestPolicy(t *testing.T, db *bun.DB, tenantName string, productID uuid.UUID) *entities.Policy {
	policy := &entities.Policy{
		ID:         uuid.New(),
		TenantName: tenantName,
		ProductID:  productID,
		Name:       "Test Policy",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := db.NewInsert().Model(policy).Exec(context.Background())
	require.NoError(t, err)
	return policy
}

func createTestLicense(t *testing.T, db *bun.DB, tenantName string, policyID, productID uuid.UUID, status string, machinesCount int) *entities.License {
	license := &entities.License{
		ID:            uuid.New(),
		TenantName:    tenantName,
		PolicyID:      policyID,
		ProductID:     productID,
		Key:           uuid.New().String(),
		Status:        status,
		MachinesCount: machinesCount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err := db.NewInsert().Model(license).Exec(context.Background())
	require.NoError(t, err)
	return license
}

func createTestMachine(t *testing.T, db *bun.DB, tenantName string, licenseID uuid.UUID) *entities.Machine {
	machine := &entities.Machine{
		ID:          uuid.New(),
		TenantName:  tenantName,
		LicenseID:   licenseID,
		LicenseKey:  "test-key",
		Fingerprint: "test-fingerprint",
		Name:        "Test Machine",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := db.NewInsert().Model(machine).Exec(context.Background())
	require.NoError(t, err)
	return machine
}

func TestDeleteMachineByPKAndUpdateLicense_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusActive, 2)
	machine := createTestMachine(t, db, tenant.Name, license.ID)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	err := repo.DeleteMachineByPKAndUpdateLicense(context.Background(), machine.ID)
	assert.NoError(t, err)

	// Verify machine is deleted
	var machineExists bool
	machineExists, err = db.NewSelect().Model((*entities.Machine)(nil)).Where("id = ?", machine.ID).Exists(context.Background())
	assert.NoError(t, err)
	assert.False(t, machineExists)

	// Verify license is updated
	updatedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(updatedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusActive, updatedLicense.Status)
}

func TestDeleteMachineByPKAndUpdateLicense_LastMachine(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusActive, 1)
	machine := createTestMachine(t, db, tenant.Name, license.ID)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	err := repo.DeleteMachineByPKAndUpdateLicense(context.Background(), machine.ID)
	assert.NoError(t, err)

	// Verify license status is set to inactive when last machine is removed
	updatedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(updatedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, updatedLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusInactive, updatedLicense.Status)
}

func TestDeleteMachineByPKAndUpdateLicense_RollbackOnError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusActive, 1)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	// Try to delete non-existent machine
	nonExistentID := uuid.New()
	err := repo.DeleteMachineByPKAndUpdateLicense(context.Background(), nonExistentID)
	assert.Error(t, err)

	// Verify license is unchanged (transaction rolled back)
	unchangedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(unchangedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, unchangedLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusActive, unchangedLicense.Status)
}

func TestInsertNewMachineAndUpdateLicense_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusNotActivated, 0)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	newMachine := &entities.Machine{
		ID:          uuid.New(),
		TenantName:  tenant.Name,
		LicenseID:   license.ID,
		LicenseKey:  license.Key,
		Fingerprint: "new-fingerprint",
		Name:        "New Machine",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.InsertNewMachineAndUpdateLicense(context.Background(), newMachine)
	assert.NoError(t, err)

	// Verify machine is created
	var machineExists bool
	machineExists, err = db.NewSelect().Model((*entities.Machine)(nil)).Where("id = ?", newMachine.ID).Exists(context.Background())
	assert.NoError(t, err)
	assert.True(t, machineExists)

	// Verify license is updated
	updatedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(updatedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusActive, updatedLicense.Status)
}

func TestInsertNewMachineAndUpdateLicense_ActivatesInactiveLicense(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusInactive, 0)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	newMachine := &entities.Machine{
		ID:          uuid.New(),
		TenantName:  tenant.Name,
		LicenseID:   license.ID,
		LicenseKey:  license.Key,
		Fingerprint: "new-fingerprint",
		Name:        "New Machine",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.InsertNewMachineAndUpdateLicense(context.Background(), newMachine)
	assert.NoError(t, err)

	// Verify license status changed from inactive to active
	updatedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(updatedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, constants.LicenseStatusActive, updatedLicense.Status)
}

func TestInsertNewMachineAndUpdateLicense_RollbackOnError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	// Try to insert machine with non-existent license
	newMachine := &entities.Machine{
		ID:          uuid.New(),
		TenantName:  tenant.Name,
		LicenseID:   uuid.New(), // Non-existent license
		LicenseKey:  "fake-key",
		Fingerprint: "new-fingerprint",
		Name:        "New Machine",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.InsertNewMachineAndUpdateLicense(context.Background(), newMachine)
	assert.Error(t, err)

	// Verify machine was not created (transaction rolled back)
	var machineExists bool
	machineExists, err = db.NewSelect().Model((*entities.Machine)(nil)).Where("id = ?", newMachine.ID).Exists(context.Background())
	assert.NoError(t, err)
	assert.False(t, machineExists)
}

func TestUpdateMachineByPKAndLicense_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	oldLicense := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusActive, 1)
	newLicense := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusNotActivated, 0)
	machine := createTestMachine(t, db, tenant.Name, oldLicense.ID)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	machine.LicenseID = newLicense.ID
	machine.LicenseKey = newLicense.Key

	updatedMachine, err := repo.UpdateMachineByPKAndLicense(context.Background(), machine, oldLicense, newLicense)
	assert.NoError(t, err)
	assert.NotNil(t, updatedMachine)

	// Verify old license is updated
	updatedOldLicense := &entities.License{ID: oldLicense.ID}
	err = db.NewSelect().Model(updatedOldLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, updatedOldLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusInactive, updatedOldLicense.Status)

	// Verify new license is updated
	updatedNewLicense := &entities.License{ID: newLicense.ID}
	err = db.NewSelect().Model(updatedNewLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedNewLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusActive, updatedNewLicense.Status)
}

func TestUpdateMachineByPKAndLicense_NoLicenseChange(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tenant := createTestTenant(t, db, "test-tenant")
	product := createTestProduct(t, db, tenant.Name)
	policy := createTestPolicy(t, db, tenant.Name, product.ID)
	license := createTestLicense(t, db, tenant.Name, policy.ID, product.ID, constants.LicenseStatusActive, 1)
	machine := createTestMachine(t, db, tenant.Name, license.ID)

	ds := &api.DataSource{}
	ds.SetDatabase(db)
	repo := NewMachineRepository(ds)

	machine.Name = "Updated Machine Name"

	updatedMachine, err := repo.UpdateMachineByPKAndLicense(context.Background(), machine, license, nil)
	assert.NoError(t, err)
	assert.NotNil(t, updatedMachine)

	// Verify license is unchanged
	unchangedLicense := &entities.License{ID: license.ID}
	err = db.NewSelect().Model(unchangedLicense).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, unchangedLicense.MachinesCount)
	assert.Equal(t, constants.LicenseStatusActive, unchangedLicense.Status)

	// Verify machine is updated
	verifyMachine := &entities.Machine{ID: machine.ID}
	err = db.NewSelect().Model(verifyMachine).WherePK().Scan(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Machine Name", verifyMachine.Name)
}
