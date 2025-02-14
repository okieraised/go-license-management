package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/machines/models"
	"go-license-management/internal/services/v1/machines/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type MachineService struct {
	repo   repository.IMachine
	logger *logging.Logger
}

func NewMachineService(options ...func(*MachineService)) *MachineService {
	svc := &MachineService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IMachine) func(*MachineService) {
	return func(c *MachineService) {
		c.repo = repo
	}
}

func (svc *MachineService) Create(ctx *gin.Context, input *models.MachineRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-license-id")
	license, err := svc.repo.SelectLicenseByLicenseKey(ctx, utils.DerefPointer(input.LicenseKey))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineLicenseIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineLicenseIsInvalid]
			return resp, cerrors.ErrMachineLicenseIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	// If license status is either suspended, banned, or expired, return error
	if license.Suspended {
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIsSuspended]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIsSuspended]
		return resp, cerrors.ErrLicenseIsSuspended
	}

	if license.Status == constants.LicenseStatusBanned {
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIsBanned]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIsBanned]
		return resp, cerrors.ErrLicenseIsBanned
	}

	if license.Status == constants.LicenseStatusExpired {
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseHasExpired]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseHasExpired]
		return resp, cerrors.ErrLicenseHasExpired
	}

	_, cSpan = input.Tracer.Start(rootCtx, "query-machine-by-fingerprint")
	mExists, err := svc.repo.CheckMachineExistByFingerprintAndLicense(ctx, utils.DerefPointer(input.LicenseKey), utils.DerefPointer(input.Fingerprint))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if mExists {
		svc.logger.GetLogger().Info(fmt.Sprintf("machine fingerprint [%s] is already associated with license [%s]", utils.DerefPointer(input.Fingerprint), utils.DerefPointer(input.LicenseKey)))
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineFingerprintAssociatedWithLicense]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineFingerprintAssociatedWithLicense]
		return resp, cerrors.ErrMachineFingerprintAssociatedWithLicense
	}

	// Check max machine of the license
	if license.MachinesCount != 0 {
		if license.MachinesCount+1 > license.MaxMachines && license.Policy.OverageStrategy == constants.PolicyOverageStrategyNoOverage {
			svc.logger.GetLogger().Error("license max machine exceeded")
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseMaxMachineExceeded]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseMaxMachineExceeded]
			return resp, cerrors.ErrLicenseMaxMachineExceeded
		}
	}

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-machine")
	machineID := uuid.New()
	now := time.Now()
	machine := &entities.Machine{
		ID:          machineID,
		TenantName:  tenant.Name,
		LicenseID:   license.ID,
		LicenseKey:  utils.DerefPointer(input.LicenseKey),
		Fingerprint: utils.DerefPointer(input.Fingerprint),
		Metadata:    input.Metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if input.IP != nil {
		machine.IP = utils.DerefPointer(input.IP)
	}

	if input.Hostname != nil {
		machine.Hostname = utils.DerefPointer(input.Hostname)
	}

	if input.Platform != nil {
		machine.Platform = utils.DerefPointer(input.Platform)
	}

	if input.Name != nil {
		machine.Name = utils.DerefPointer(input.Name)
	}

	if input.Cores != nil {
		machine.Cores = utils.DerefPointer(input.Cores)
	}

	err = svc.repo.InsertNewMachineAndUpdateLicense(ctx, machine)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.MachineInfoOutput{
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
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *MachineService) Update(ctx *gin.Context, input *models.MachineUpdateInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-machine")
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineIDIsInvalid]
			return resp, cerrors.ErrMachineIDIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	if input.IP != nil {
		machine.IP = utils.DerefPointer(input.IP)
	}

	if input.Hostname != nil {
		machine.Hostname = utils.DerefPointer(input.Hostname)
	}

	if input.Platform != nil {
		machine.Platform = utils.DerefPointer(input.Platform)
	}

	if input.Name != nil {
		machine.Name = utils.DerefPointer(input.Name)
	}

	if input.Cores != nil {
		machine.Cores = utils.DerefPointer(input.Cores)
	}

	if input.Metadata != nil {
		machine.Metadata = input.Metadata
	}

	if input.Fingerprint != nil {
		machine.Fingerprint = utils.DerefPointer(input.Fingerprint)
	}

	// Update license key if specified
	currentLicense := machine.License
	var newLicense *entities.License
	if input.LicenseKey != nil {
		// Only update if license key is different from the current key
		if machine.LicenseKey != utils.DerefPointer(input.LicenseKey) {
			_, cSpan = input.Tracer.Start(rootCtx, "query-license-id")
			license, err := svc.repo.SelectLicenseByLicenseKey(ctx, utils.DerefPointer(input.LicenseKey))
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				switch {
				case errors.Is(err, sql.ErrNoRows):
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseKeyIsInvalid]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseKeyIsInvalid]
					return resp, cerrors.ErrLicenseKeyIsInvalid
				default:
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
			cSpan.End()
			if license.Suspended {
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIsSuspended]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIsSuspended]
				return resp, cerrors.ErrLicenseIsSuspended
			}

			if license.Status == constants.LicenseStatusBanned {
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIsBanned]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIsBanned]
				return resp, cerrors.ErrLicenseIsBanned
			}

			if license.Status == constants.LicenseStatusExpired {
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseHasExpired]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseHasExpired]
				return resp, cerrors.ErrLicenseHasExpired
			}
			if license.MachinesCount != 0 {
				if license.MachinesCount+1 > license.MaxMachines && license.Policy.OverageStrategy == constants.PolicyOverageStrategyNoOverage {
					svc.logger.GetLogger().Error("license max machine exceeded")
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseMaxMachineExceeded]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseMaxMachineExceeded]
					return resp, cerrors.ErrLicenseMaxMachineExceeded
				}
			}

			machine.LicenseKey = utils.DerefPointer(input.LicenseKey)
			newLicense = license
		}
	}

	_, cSpan = input.Tracer.Start(rootCtx, "update-machine")
	machine, err = svc.repo.UpdateMachineByPKAndLicense(ctx, machine, currentLicense, newLicense)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.MachineInfoOutput{
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
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *MachineService) Retrieve(ctx *gin.Context, input *models.MachineRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "select-machine")
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineIDIsInvalid]
			return resp, cerrors.ErrMachineIDIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.MachineInfoOutput{
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
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *MachineService) Delete(ctx *gin.Context, input *models.MachineDeleteInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "delete-product")
	err = svc.repo.DeleteMachineByPKAndUpdateLicense(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	return resp, nil
}

func (svc *MachineService) List(ctx *gin.Context, input *models.MachineListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-pkc")
	machines, total, err := svc.repo.SelectMachines(ctx, tenant.Name, input.QueryCommonParam)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
			return resp, cerrors.ErrProductIDIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	machineOutput := make([]models.MachineListOutput, 0)
	for _, machine := range machines {
		machineOutput = append(machineOutput, models.MachineListOutput{
			ID:              machine.ID,
			LicenseID:       machine.LicenseID,
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
		})
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = machineOutput

	return resp, nil
}

func (svc *MachineService) Actions(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrTenantNameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrTenantNameIsInvalid]
			return resp, cerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	action := utils.DerefPointer(input.MachineCommonURI.MachineAction)
	switch action {
	case constants.MachineActionCheckout:
		_, cSpan = input.Tracer.Start(rootCtx, "action-checkout")
		respData, err := svc.checkout(ctx, input)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			switch {
			case errors.Is(err, sql.ErrNoRows):
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineIDIsInvalid]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineIDIsInvalid]
				return resp, cerrors.ErrMachineIDIsInvalid
			default:
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}

		}
		resp.Data = respData
		cSpan.End()
	case constants.MachineActionResetHeartbeat:
		respData, err := svc.resetHeartbeat(ctx, input)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			switch {
			case errors.Is(err, sql.ErrNoRows):
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineIDIsInvalid]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineIDIsInvalid]
				return resp, cerrors.ErrMachineIDIsInvalid
			default:
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		}
		resp.Data = respData
		cSpan.End()
	case constants.MachineActionPingHeartbeat:
		respData, err := svc.pingHeartbeat(ctx, input)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			switch {
			case errors.Is(err, sql.ErrNoRows):
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrMachineIDIsInvalid]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrMachineIDIsInvalid]
				return resp, cerrors.ErrMachineIDIsInvalid
			default:
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		}
		resp.Data = respData
		cSpan.End()
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	return resp, nil
}
