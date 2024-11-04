package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/license_key"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/tenants/models"
	"go-license-management/internal/server/v1/tenants/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type TenantService struct {
	repo   repository.ITenant
	logger *logging.Logger
}

func NewTenantService(options ...func(*TenantService)) *TenantService {
	svc := &TenantService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.ITenant) func(*TenantService) {
	return func(c *TenantService) {
		c.repo = repo
	}
}

func (svc *TenantService) Create(ctx *gin.Context, input *models.TenantRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	exists, err := svc.repo.CheckTenantExistByName(ctx, utils.DerefPointer(input.Name))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// If tenant name already exists, return with error
	if exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("tenant [%s] already exists", utils.DerefPointer(input.Name)))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrTenantNameAlreadyExist]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrTenantNameAlreadyExist]
		return resp, comerrors.ErrTenantNameAlreadyExist
	}

	// If not, generate additional required info
	_, cSpan = input.Tracer.Start(rootCtx, "generate-tenant-key")
	privateKey, publicKey, err := license_key.NewEd25519KeyPair()
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// Insert new tenant
	_, cSpan = input.Tracer.Start(rootCtx, "insert-tenant")
	tenandID := uuid.New()
	tenant := &entities.Tenant{
		ID:                tenandID,
		Name:              utils.DerefPointer(input.Name),
		Protected:         utils.DerefPointer(input.Protected),
		Ed25519PublicKey:  publicKey,
		Ed25519PrivateKey: privateKey,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = svc.repo.InsertNewTenant(ctx, tenant)
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
		"tenant_name": utils.DerefPointer(input.Name),
	}

	return resp, nil
}

func (svc *TenantService) List(ctx *gin.Context, input *models.TenantListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant")
	tenants, count, err := svc.repo.SelectTenants(ctx)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "convert-tenants-to-output")
	respData := make([]models.TenantRetrievalOutput, 0)
	for _, tenant := range tenants {
		respData = append(respData, models.TenantRetrievalOutput{
			ID:                             tenant.ID.String(),
			Name:                           tenant.Name,
			Protected:                      tenant.Protected,
			Ed25519PublicKey:               tenant.Ed25519PublicKey,
			LastLowActivityLifelineSentAt:  tenant.LastLowActivityLifelineSentAt,
			LastTrialWillEndSentAt:         tenant.LastTrialWillEndSentAt,
			LastLicenseLimitExceededSentAt: tenant.LastLicenseLimitExceededSentAt,
			LastRequestLimitExceededSentAt: tenant.LastRequestLimitExceededSentAt,
			LastPromptForReviewSentAt:      tenant.LastPromptForReviewSentAt,
			CreatedAt:                      tenant.CreatedAt,
			UpdatedAt:                      tenant.UpdatedAt,
		})
	}
	cSpan.End()

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Count = count
	resp.Data = respData

	return resp, nil
}

func (svc *TenantService) Retrieve(ctx *gin.Context, input *models.TenantRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "retrieval-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	tenant, err := svc.repo.SelectTenantByName(ctx, utils.DerefPointer(input.Name))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "convert-tenant-to-output")
	respData := models.TenantRetrievalOutput{
		ID:                             tenant.ID.String(),
		Name:                           tenant.Name,
		Protected:                      tenant.Protected,
		Ed25519PublicKey:               tenant.Ed25519PublicKey,
		LastLowActivityLifelineSentAt:  tenant.LastLowActivityLifelineSentAt,
		LastTrialWillEndSentAt:         tenant.LastTrialWillEndSentAt,
		LastLicenseLimitExceededSentAt: tenant.LastLicenseLimitExceededSentAt,
		LastRequestLimitExceededSentAt: tenant.LastRequestLimitExceededSentAt,
		LastPromptForReviewSentAt:      tenant.LastPromptForReviewSentAt,
		CreatedAt:                      tenant.CreatedAt,
		UpdatedAt:                      tenant.UpdatedAt,
	}
	cSpan.End()

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *TenantService) Delete(ctx *gin.Context, input *models.TenantDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "retrieval-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "delete-tenant-by-name")
	err := svc.repo.DeleteTenantByName(ctx, utils.DerefPointer(input.Name))
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
