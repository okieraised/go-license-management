package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/models/license_key"
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
		expiry, _ = time.Parse(constants.DateFormatISO8601Hyphen, utils.DerefPointer(input.Expiry))
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

	// Generating license key
	svc.logger.GetLogger().Info("generating license key")
	licenseKey, err := svc.generateLicenseKey(ctx, licenseID.String(), tenant, product, policy)
	if err != nil {
		return nil, err
	}
	license.Key = licenseKey

	return license, nil
}

func (svc *LicenseService) generateLicenseKey(ctx *gin.Context, licenseID string, tenant *entities.Tenant, product *entities.Product, policy *entities.Policy) (string, error) {
	licenseKeyContent := &license_key.LicenseKeyContent{
		TenantName: utils.RefPointer(tenant.Name),
		ProductID:  utils.RefPointer(product.ID.String()),
		PolicyID:   utils.RefPointer(policy.ID.String()),
		LicenseID:  utils.RefPointer(licenseID),
		CreatedAt:  time.Now(),
	}

	var licenseKey = ""
	var err error
	switch policy.Scheme {
	case constants.PolicySchemeED25519:
		licenseKey, err = utils.NewLicenseKeyWithEd25519(policy.PrivateKey, licenseKeyContent)
	case constants.PolicySchemeRSA2048PKCS1:
		licenseKey, err = utils.NewLicenseKeyWithRSA2048PKCS1(policy.PrivateKey, licenseKeyContent)
	default:
		err = comerrors.ErrPolicySchemeIsInvalid
	}
	if err != nil {
		return licenseKey, err
	}

	return licenseKey, nil
}
