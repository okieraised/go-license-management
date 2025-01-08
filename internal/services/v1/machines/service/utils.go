package service

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go-license-management/internal/services/v1/machines/models"
	"go-license-management/internal/utils"
	"strings"
	"time"
)

// checkout checkouts a license for the machine
func (svc *MachineService) checkout(ctx *gin.Context, input *models.MachineActionsInput) (*models.MachineActionCheckoutOutput, error) {
	// query machine info
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	// Query license info
	license, err := svc.repo.SelectLicenseByPK(ctx, machine.LicenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	// Query policy
	policy, err := svc.repo.SelectPolicyByPK(ctx, license.PolicyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	ttl := utils.DerefPointer(input.TTL)
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Duration(ttl) * time.Second)
	alg := policy.Scheme

	licenseContent := machine_attribute.MachineLicenseField{
		TenantName:         machine.TenantName,
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
	b64MachineCert := base64.StdEncoding.EncodeToString(bMachineCert)

	// generate encryption key from hash of signature and machine fingerprint
	svc.logger.GetLogger().Info(fmt.Sprintf("encrypting machine file for machine [%s]", machine.ID.String()))
	h := sha256.New()
	h.Write([]byte(signature + machine.Fingerprint))
	sha := h.Sum(nil)

	// (FE) public key -> sign request body -> services decrypt

	encryptedMachineCert, err := utils.Encrypt([]byte(b64MachineCert), sha)
	if err != nil {
		return nil, err
	}

	finalCertificate := fmt.Sprintf(constants.MachineFileFormat, base64.StdEncoding.EncodeToString(encryptedMachineCert))
	output := &models.MachineActionCheckoutOutput{
		ID:          machine.ID,
		Type:        "machine",
		Certificate: finalCertificate,
		TTL:         ttl,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiredAt,
	}

	machine.LastCheckOutAt = time.Now()
	machine, err = svc.repo.UpdateMachineByPK(ctx, machine)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (svc *MachineService) pingHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*models.MachineInfoOutput, error) {
	// query machine info
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	machine.LastHeartbeatAt = time.Now()
	machine, err = svc.repo.UpdateMachineByPK(ctx, machine)
	if err != nil {
		return nil, err
	}

	return &models.MachineInfoOutput{
		ID:              machine.ID,
		TenantName:      machine.TenantName,
		LicenseKey:      machine.LicenseKey,
		Fingerprint:     machine.Fingerprint,
		IP:              machine.IP,
		Hostname:        machine.Hostname,
		Platform:        machine.Platform,
		Name:            machine.Name,
		Metadata:        machine.Metadata,
		Cores:           machine.Cores,
		LastHeartbeatAt: machine.LastHeartbeatAt,
		LastCheckOutAt:  machine.LastCheckOutAt,
		CreatedAt:       machine.CreatedAt,
		UpdatedAt:       machine.UpdatedAt,
	}, nil
}

func (svc *MachineService) resetHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*models.MachineInfoOutput, error) {
	// query machine info
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	machine.LastHeartbeatAt = time.Time{}
	machine, err = svc.repo.UpdateMachineByPK(ctx, machine)
	if err != nil {
		return nil, err
	}

	return &models.MachineInfoOutput{
		ID:              machine.ID,
		TenantName:      machine.TenantName,
		LicenseKey:      machine.LicenseKey,
		Fingerprint:     machine.Fingerprint,
		IP:              machine.IP,
		Hostname:        machine.Hostname,
		Platform:        machine.Platform,
		Name:            machine.Name,
		Metadata:        machine.Metadata,
		Cores:           machine.Cores,
		LastHeartbeatAt: machine.LastHeartbeatAt,
		LastCheckOutAt:  machine.LastCheckOutAt,
		CreatedAt:       machine.CreatedAt,
		UpdatedAt:       machine.UpdatedAt,
	}, nil
}
