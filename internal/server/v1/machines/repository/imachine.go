package repository

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/database/entities"
)

type IMachine interface {
	SelectTenantByName(ctx context.Context, tenantName string) (*entities.Tenant, error)
	CheckLicenseExistByPK(ctx context.Context, licenseID uuid.UUID) (bool, error)
	CheckMachineExistByFingerprintAndLicense(ctx context.Context, licenseID uuid.UUID, fingerprint string) (bool, error)
	SelectLicenseByPK(ctx context.Context, licenseID uuid.UUID) (*entities.License, error)
	SelectPolicyByPK(ctx context.Context, policyID uuid.UUID) (*entities.Policy, error)
	SelectMachineByPK(ctx context.Context, machineID uuid.UUID) (*entities.Machine, error)
	InsertNewMachine(ctx context.Context, machine *entities.Machine) error
	UpdateMachineByPK(ctx context.Context, machine *entities.Machine) error
	InsertNewMachineAndUpdateLicense(ctx context.Context, machine *entities.Machine) error
	DeleteMachineByPK(ctx context.Context, machineID uuid.UUID) error
	DeleteMachineByPKAndUpdateLicense(ctx context.Context, machineID uuid.UUID) error
}
