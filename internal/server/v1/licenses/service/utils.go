package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/utils"
	"time"
)

// generateLicense creates new license with required information for validation embedded
func (svc *LicenseService) generateLicense(ctx *gin.Context, input *models.LicenseRegistrationInput, tenant *entities.Tenant, product *entities.Product, policy *entities.Policy) (*entities.License, error) {
	licenseID := uuid.New()
	now := time.Now()

	// Init new license
	license := &entities.License{
		ID:         licenseID,
		TenantName: tenant.Name,
		PolicyID:   policy.ID,
		ProductID:  product.ID,
		Name:       utils.DerefPointer(input.Name),
		Status:     constants.LicenseStatusNotActivated,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Check for license expiration
	var expiry time.Time
	if input.Expiry != nil {
		expiry, _ = time.Parse(time.RFC3339, utils.DerefPointer(input.Expiry))
	}

	// Check for policy expiration
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying policy [%s] duration", policy.ID))
	policyDuration := policy.Duration
	if policy.Duration > 0 {
		if expiry.IsZero() {
			license.Expiry = now.Add(time.Duration(policyDuration) * time.Second)
		} else {
			license.Expiry = expiry.Add(time.Duration(policyDuration) * time.Second)
		}
	} else {
		license.Expiry = expiry
	}

	// Check for license max uses
	svc.logger.GetLogger().Info("verifying license max uses")
	license.MaxUses = policy.MaxUses
	if input.MaxUses != nil {
		license.MaxUses = utils.DerefPointer(input.MaxUses)
	}

	// Check for license max machines
	svc.logger.GetLogger().Info("verifying license max machines")
	license.MaxMachines = policy.MaxMachines
	if input.MaxMachines != nil {
		license.MaxMachines = utils.DerefPointer(input.MaxMachines)
	}

	// Check for license max users
	svc.logger.GetLogger().Info("verifying license max users")
	license.MaxUsers = policy.MaxUsers
	if input.MaxUsers != nil {
		license.MaxUsers = utils.DerefPointer(input.MaxUsers)
	}

	// Generating license key
	svc.logger.GetLogger().Info("generating license key")
	licenseKey := utils.GenerateToken()
	license.Key = licenseKey

	return license, nil
}

// validateLicense validates a license. This will check the following: if the license is suspended, if the license is expired,
// if the license is overdue for check-in, and if the license meets its machine requirements (if strict).
func (svc *LicenseService) validateLicense(ctx *gin.Context, license *entities.License) (*models.LicenseValidationOutput, error) {
	resp := &models.LicenseValidationOutput{}

	switch license.Status {
	case constants.LicenseStatusNotActivated:
		resp.Valid = false
		if license.MachinesCount == 0 && license.Policy.MaxMachines == 1 {
			resp.Code = constants.LicenseValidationStatusNoMachine
		} else {
			resp.Code = constants.LicenseValidationStatusNoMachine
		}
	case constants.LicenseStatusActive:
		resp.Valid = true
		resp.Code = constants.LicenseValidationStatusValid
	case constants.LicenseStatusInactive:
		resp.Valid = true
		resp.Code = constants.LicenseValidationStatusValid
	case constants.LicenseStatusBanned:
		resp.Valid = false
		resp.Code = constants.LicenseValidationStatusBanned
	case constants.LicenseStatusExpired:
		resp.Valid = false
		resp.Code = constants.LicenseValidationStatusExpired
	case constants.LicenseStatusSuspended:
		resp.Valid = false
		resp.Code = constants.LicenseValidationStatusSuspended
	}

	return resp, nil
}

// suspendLicense updates the active license status to `suspended`
func (svc *LicenseService) suspendLicense(ctx *gin.Context, license *entities.License) (*entities.License, error) {
	if license.Status == constants.LicenseStatusNotActivated {
		return nil, comerrors.ErrLicenseNotActivated
	}

	license.Status = constants.LicenseStatusSuspended
	license.Suspended = false
	svc.logger.GetLogger().Info(fmt.Sprintf("suspending license [%s]", license.ID))

	err := svc.repo.UpdateLicenseByPK(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

// reinstateLicense updates the license status back to `active`
func (svc *LicenseService) reinstateLicense(ctx *gin.Context, license *entities.License) (*entities.License, error) {
	if license.Status != constants.LicenseStatusSuspended && !license.Suspended {
		svc.logger.GetLogger().Info(fmt.Sprintf("license [%s] has status [%s]", license.ID.String(), license.Status))
		return nil, comerrors.ErrLicenseStatusInvalidToReinstate
	}

	license.Status = constants.LicenseStatusActive
	license.Suspended = false

	svc.logger.GetLogger().Info(fmt.Sprintf("reinstating license [%s]", license.ID))
	err := svc.repo.UpdateLicenseByPK(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

// renewLicense extends license expiry by the policy's duration, according to the policy's renewal basis.
// Renewals take the license's current expiry datetime and add, in seconds, the policy's duration,
func (svc *LicenseService) renewLicense(ctx *gin.Context, license *entities.License) (*entities.License, error) {
	policy := license.Policy

	// If the license does not have an expiration, skip
	if license.Expiry.IsZero() {
		svc.logger.GetLogger().Info(fmt.Sprintf("license [%s] does not expiry, skipping", license.ID.String()))
		return license, nil
	}

	svc.logger.GetLogger().Info(fmt.Sprintf("reinstating license [%s]", license.ID))
	if policy.Duration != 0 {
		switch policy.RenewalBasis {
		case constants.PolicyRenewalBasisFromExpiry:
			license.Expiry = license.Expiry.Add(time.Duration(policy.Duration) * time.Second)
		case constants.PolicyRenewalFromNow:
			license.Expiry = time.Now().Add(time.Duration(policy.Duration) * time.Second)
		case constants.PolicyRenewalFromNowIfExpired:
			if time.Now().After(license.Expiry) {
				license.Expiry = time.Now().Add(time.Duration(policy.Duration) * time.Second)
			} else {
				license.Expiry = license.Expiry.Add(time.Duration(policy.Duration) * time.Second)
			}
		}
	} else {
		license.Expiry = time.Time{}
	}

	svc.logger.GetLogger().Info(fmt.Sprintf("renewing license [%s]", license.ID))
	err := svc.repo.UpdateLicenseByPK(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

// checkoutLicense check-outs a license. This will generate a snapshot of the license at time of checkout,
// encoded into a license file certificate that can be decoded and used for licensing offline and air-gapped
// environments. The algorithm will depend on the policy's scheme.
func (svc *LicenseService) checkoutLicense(ctx *gin.Context, license *entities.License) (*entities.License, error) {
	return license, nil
}

// checkinLicense check-ins a license. Sets the license's LastCheckInAt to the current time,
func (svc *LicenseService) checkinLicense(ctx *gin.Context, license *entities.License) (*entities.License, error) {
	license.LastCheckInAt = time.Now()

	svc.logger.GetLogger().Info(fmt.Sprintf("updating license [%s]'s last checked in time", license.ID))
	err := svc.repo.UpdateLicenseByPK(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}
func (svc *LicenseService) incrementUsageLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}

func (svc *LicenseService) decrementUsageLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}
func (svc *LicenseService) resetUsageLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}
