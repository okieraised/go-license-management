package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
)

type IMachine interface {
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckLicenseExistByPK(ctx context.Context, licenseID uuid.UUID) (bool, error)
	CheckMachineExistByFingerprintAndLicense(ctx context.Context, licenseKey, fingerprint string) (bool, error)
	SelectLicenseByPK(ctx context.Context, licenseID uuid.UUID) (*entities.License, error)
	SelectLicenseByLicenseKey(ctx context.Context, licenseKey string) (*entities.License, error)
	SelectMachines(ctx context.Context, tenantName string, queryParam constants.QueryCommonParam) ([]entities.Machine, int, error)
	SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error)
	SelectMachineByPK(ctx context.Context, machineID uuid.UUID) (*entities.Machine, error)
	InsertNewMachine(ctx context.Context, machine *entities.Machine) error
	UpdateMachineByPK(ctx context.Context, machine *entities.Machine) (*entities.Machine, error)
	UpdateMachineByPKAndLicense(ctx context.Context, machine *entities.Machine, currentLicense, newLicense *entities.License) (*entities.Machine, error)
	InsertNewMachineAndUpdateLicense(ctx context.Context, machine *entities.Machine) error
	DeleteMachineByPK(ctx context.Context, machineID uuid.UUID) error
	DeleteMachineByPKAndUpdateLicense(ctx context.Context, machineID uuid.UUID) error
}
