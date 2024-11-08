package service

import (
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
	licenseKey, err := svc.generateLicenseKey(ctx, licenseID.String(), tenant, product, policy)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	license := &entities.License{
		ID:        licenseID,
		TenantID:  tenant.ID,
		PolicyID:  policy.ID,
		ProductID: product.ID,
		Key:       licenseKey,
		Name:      utils.DerefPointer(input.Name),
		Suspended: utils.DerefPointer(input.Suspended),
		Protected: utils.DerefPointer(input.Protected),
		Metadata:  input.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if input.Expiry != nil {
		expiry, _ := time.Parse(constants.DateFormatISO8601Hyphen, utils.DerefPointer(input.Expiry))
		license.Expiry = expiry
	}

	return license, nil
}

func (svc *LicenseService) generateLicenseKey(ctx *gin.Context, licenseID string, tenant *entities.Tenant, product *entities.Product, policy *entities.Policy) (string, error) {
	licenseKeyContent := &license_key.LicenseKeyContent{
		TenantID:  utils.RefPointer(tenant.ID.String()),
		ProductID: utils.RefPointer(product.ID.String()),
		PolicyID:  utils.RefPointer(policy.ID.String()),
		LicenseID: utils.RefPointer(licenseID),
		CreatedAt: time.Now(),
	}

	var licenseKey = ""
	var err error
	switch policy.Scheme {
	case constants.PolicySchemeED25519:
		licenseKey, err = utils.NewLicenseKeyWithEd25519(policy.PrivateKey, licenseKeyContent)
	case constants.PolicySchemeRSA2048PKCS1:
		licenseKey, err = utils.NewLicenseKeyWithRSA2048PKCS1(policy.PrivateKey, licenseKeyContent)
	case constants.PolicySchemeRSA2048JWTRS256:
		licenseKey, err = utils.NewLicenseKeyWithJWTRS256(policy.PrivateKey, licenseKeyContent)
	default:
		err = comerrors.ErrPolicySchemeIsInvalid
	}
	if err != nil {
		return licenseKey, err
	}

	return licenseKey, nil
}
