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

	// Query policy
	license := machine.License
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

	machineContent := models.MachineInfoOutput{
		ID:              machine.ID,
		LicenseKey:      machine.LicenseKey,
		TenantName:      machine.TenantName,
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
	}

	svc.logger.GetLogger().Info(fmt.Sprintf("generate new machine file using [%s] scheme", alg))
	var machineLicense string
	switch alg {
	case constants.PolicySchemeED25519:
		machineLicense, err = utils.NewLicenseKeyWithEd25519(policy.PrivateKey, machineContent)
	case constants.PolicySchemeRSA2048PKCS1:
		machineLicense, err = utils.NewLicenseKeyWithRSA2048PKCS1(policy.PrivateKey, machineContent)
	}
	if err != nil {
		return nil, err
	}
	parts := strings.Split(machineLicense, ".")
	signature := parts[0]
	encoded := parts[1]

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
	machineCert := base64.StdEncoding.EncodeToString(bMachineCert)

	if strings.ToLower(ctx.Query("encrypt")) == "true" || policy.Encrypted {
		svc.logger.GetLogger().Info(fmt.Sprintf("encrypting machine certificate file for machine [%s]", machine.ID.String()))
		sha256Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(machineCert)))
		encryptedCert, err := utils.Encrypt([]byte(machineCert), []byte(sha256Hash))
		if err != nil {
			return nil, err
		}
		machineCert = base64.StdEncoding.EncodeToString(encryptedCert)
		ctx.Writer.Header().Add(constants.XMachineChecksumHeader, fmt.Sprintf("sha256=%s", sha256Hash))
	}

	machineCert = fmt.Sprintf(constants.MachineFileFormat, machineCert)
	output := &models.MachineActionCheckoutOutput{
		Certificate: machineCert,
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
