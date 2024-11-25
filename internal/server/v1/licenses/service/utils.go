package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/utils"
	"time"
)

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
