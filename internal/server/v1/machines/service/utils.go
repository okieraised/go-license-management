package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/machines/models"
	"go-license-management/internal/utils"
	"strings"
	"time"
)

func (svc *MachineService) checkout(ctx *gin.Context, input *models.MachineActionsInput) (*models.MachineActionCheckoutOutput, error) {
	// query machine info
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		return nil, err
	}
	// Query license info
	license, err := svc.repo.SelectLicenseByPK(ctx, machine.LicenseID)
	if err != nil {
		return nil, err
	}
	// Query policy
	policy, err := svc.repo.SelectPolicyByPK(ctx, license.PolicyID)
	if err != nil {
		return nil, err
	}
	ttl := utils.DerefPointer(input.TTL)
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Duration(ttl) * time.Second)
	alg := policy.Scheme

	licenseContent := machine_attribute.MachineLicenseField{
		TenantID:           machine.TenantID.String(),
		ProductID:          policy.ProductID.String(),
		PolicyID:           policy.ID.String(),
		LicenseID:          license.ID.String(),
		MachineFingerprint: machine.Fingerprint,
		Metadata: map[string]interface{}{
			"machine": &machine,
		},
		TTL:       ttl,
		Expiry:    expiredAt,
		CreatedAt: issuedAt,
	}

	svc.logger.GetLogger().Info(fmt.Sprintf("generate new machine file using [%s] scheme", alg))
	var machineLicense string
	switch alg {
	case constants.PolicySchemeED25519:
		machineLicense, err = utils.NewLicenseKeyWithEd25519(policy.PrivateKey, licenseContent)
	case constants.PolicySchemeRSA2048PKCS1:
		machineLicense, err = utils.NewLicenseKeyWithRSA2048PKCS1(policy.PrivateKey, licenseContent)
	}
	if err != nil {
		return nil, err
	}
	parts := strings.Split(machineLicense, ".")
	signature := parts[1]
	encoded := parts[0]

	jsonMachineCert := machine_attribute.MachineLicenseFileContent{
		Enc: encoded,
		Sig: signature,
		Alg: alg,
	}
	bMachineCert, err := json.Marshal(jsonMachineCert)
	if err != nil {
		return nil, err
	}
	// convert the cert to base64
	b64MachineCert := base64.URLEncoding.EncodeToString(bMachineCert)

	// generate encryption key from hash of signature and machine fingerprint
	svc.logger.GetLogger().Info("encrypting machine file")
	h := sha256.New()
	h.Write([]byte(signature + machine.Fingerprint))
	sha := h.Sum(nil)

	encryptedMachineCert, err := utils.Encrypt([]byte(b64MachineCert), sha)
	if err != nil {
		return nil, err
	}

	finalCertificate := fmt.Sprintf(constants.MachineFileFormat, base64.URLEncoding.EncodeToString(encryptedMachineCert))
	output := &models.MachineActionCheckoutOutput{
		ID:          machine.ID,
		Type:        "machine",
		Certificate: finalCertificate,
		TTL:         ttl,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiredAt,
	}

	machine.LastCheckOutAt = time.Now()
	err = svc.repo.UpdateMachineByPK(ctx, machine)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (svc *MachineService) pingHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *MachineService) resetHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}
