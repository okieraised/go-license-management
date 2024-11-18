package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/machines/models"
	"go-license-management/internal/server/v1/machines/repository"
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
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-license-id")
	license, err := svc.repo.SelectLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrMachineLicenseIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrMachineLicenseIsInvalid]
			return resp, comerrors.ErrMachineLicenseIsInvalid
		default:
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	if license.Suspended {
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrLicenseIsSuspended]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrLicenseIsSuspended]
		return resp, comerrors.ErrLicenseIsSuspended
	}

	if license.Status == constants.LicenseStatusBanned {
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrLicenseIsBanned]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrLicenseIsBanned]
		return resp, comerrors.ErrLicenseIsBanned
	}

	if license.Status == constants.LicenseStatusExpired {
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrLicenseIsExpired]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrLicenseIsExpired]
		return resp, comerrors.ErrLicenseIsExpired
	}

	_, cSpan = input.Tracer.Start(rootCtx, "query-machine-by-fingerprint")
	mExists, err := svc.repo.CheckMachineExistByFingerprintAndLicense(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)), utils.DerefPointer(input.Fingerprint))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if mExists {
		svc.logger.GetLogger().Info(fmt.Sprintf("machine fingerprint [%s] is already associated with license [%s]", utils.DerefPointer(input.Fingerprint), utils.DerefPointer(input.LicenseID)))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrMachineFingerprintAssociatedWithLicense]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrMachineFingerprintAssociatedWithLicense]
		return resp, comerrors.ErrMachineFingerprintAssociatedWithLicense
	}

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-machine")
	machineID := uuid.New()
	now := time.Now()
	machine := &entities.Machine{
		ID:          machineID,
		TenantID:    tenant.ID,
		LicenseID:   uuid.MustParse(utils.DerefPointer(input.LicenseID)),
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
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.MachineInfoOutput{
		ID:                   machine.ID,
		TenantID:             machine.TenantID,
		LicenseID:            machine.LicenseID,
		Fingerprint:          machine.Fingerprint,
		IP:                   machine.IP,
		Hostname:             machine.Hostname,
		Platform:             machine.Platform,
		Name:                 machine.Name,
		Metadata:             machine.Metadata,
		Cores:                machine.Cores,
		LastHeartbeatAt:      machine.LastHeartbeatAt,
		LastDeathEventSentAt: machine.LastDeathEventSentAt,
		LastCheckOutAt:       machine.LastCheckOutAt,
		CreatedAt:            machine.CreatedAt,
		UpdatedAt:            machine.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *MachineService) Update(ctx *gin.Context, input *models.MachineUpdateInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *MachineService) Retrieve(ctx *gin.Context, input *models.MachineRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
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
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrMachineIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrMachineIDIsInvalid]
			return resp, comerrors.ErrMachineIDIsInvalid
		default:
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.MachineInfoOutput{
		ID:                   machine.ID,
		TenantID:             machine.TenantID,
		LicenseID:            machine.LicenseID,
		Fingerprint:          machine.Fingerprint,
		IP:                   machine.IP,
		Hostname:             machine.Hostname,
		Platform:             machine.Platform,
		Name:                 machine.Name,
		Metadata:             machine.Metadata,
		Cores:                machine.Cores,
		LastHeartbeatAt:      machine.LastHeartbeatAt,
		LastDeathEventSentAt: machine.LastDeathEventSentAt,
		LastCheckOutAt:       machine.LastCheckOutAt,
		CreatedAt:            machine.CreatedAt,
		UpdatedAt:            machine.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *MachineService) Delete(ctx *gin.Context, input *models.MachineDeleteInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "delete-product")
	err = svc.repo.DeleteMachineByPKAndUpdateLicense(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	return resp, nil
}

func (svc *MachineService) List(ctx *gin.Context, input *models.MachineListInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *MachineService) Actions(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	_, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.TenantName))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsInvalid]
			return resp, comerrors.ErrTenantNameIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
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
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
		resp.Data = respData
		cSpan.End()
	case constants.MachineActionResetHeartbeat:
	case constants.MachineActionPingHeartbeat:
	default:
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	return resp, nil
}
