package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/entitlements/models"
	"go-license-management/internal/server/v1/entitlements/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type EntitlementService struct {
	repo   repository.IEntitlement
	logger *logging.Logger
}

func NewEntitlementService(options ...func(*EntitlementService)) *EntitlementService {
	svc := &EntitlementService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IEntitlement) func(*EntitlementService) {
	return func(c *EntitlementService) {
		c.repo = repo
	}
}

func (svc *EntitlementService) Create(ctx *gin.Context, input *models.EntitlementRegistrationInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-entitlement-by-code")
	exists, err := svc.repo.CheckEntitlementExistByCode(ctx, utils.DerefPointer(input.Code))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if exists {
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrEntitlementCodeAlreadyExist]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrEntitlementCodeAlreadyExist]
		return resp, comerrors.ErrEntitlementCodeAlreadyExist
	}

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-entitlement")
	entitlementID := uuid.New()
	now := time.Now()
	entitlement := &entities.Entitlement{
		ID:        entitlementID,
		TenantID:  tenant.ID,
		Name:      utils.DerefPointer(input.Name),
		Code:      utils.DerefPointer(input.Code),
		Metadata:  input.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = svc.repo.InsertNewEntitlement(ctx, entitlement)
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
	resp.Data = map[string]interface{}{
		"entitlement_id": entitlementID.String(),
	}
	return resp, nil
}

func (svc *EntitlementService) Update(ctx *gin.Context, input *models.EntitlementUpdateInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *EntitlementService) List(ctx *gin.Context, input *models.EntitlementListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-entitlements")
	accounts, count, err := svc.repo.SelectEntitlementsByTenant(ctx, tenant.ID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "convert-entitlements-to-output")
	respData := make([]models.EntitlementRetrievalOutput, 0)
	for _, entitlement := range accounts {
		respData = append(respData, models.EntitlementRetrievalOutput{
			ID:        entitlement.ID,
			TenantID:  entitlement.ID,
			Name:      entitlement.Name,
			Code:      entitlement.Code,
			Metadata:  entitlement.Metadata,
			CreatedAt: entitlement.CreatedAt,
			UpdatedAt: entitlement.UpdatedAt,
		})
	}
	cSpan.End()

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Count = count
	resp.Data = respData

	return resp, nil
}

func (svc *EntitlementService) Retrieve(ctx *gin.Context, input *models.EntitlementRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-account")
	entitlement, err := svc.repo.SelectEntitlementByPK(ctx, tenant.ID, input.EntitlementID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := &models.EntitlementRetrievalOutput{
		ID:        entitlement.ID,
		TenantID:  entitlement.ID,
		Name:      entitlement.Name,
		Code:      entitlement.Code,
		Metadata:  entitlement.Metadata,
		CreatedAt: entitlement.CreatedAt,
		UpdatedAt: entitlement.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *EntitlementService) Delete(ctx *gin.Context, input *models.EntitlementDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-entitlement")
	err = svc.repo.DeleteEntitlementByPK(ctx, tenant.ID, input.EntitlementID)
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
