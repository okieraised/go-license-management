package service

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/server/v1/licenses/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
)

type LicenseService struct {
	repo   repository.ILicense
	logger *logging.Logger
}

func NewLicenseService(options ...func(*LicenseService)) *LicenseService {
	svc := &LicenseService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.ILicense) func(*LicenseService) {
	return func(c *LicenseService) {
		c.repo = repo
	}
}

func (svc *LicenseService) Create(ctx *gin.Context, input *models.LicenseRegistrationInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-id")
	product, err := svc.repo.SelectProductByPK(ctx, tenant.ID, input.ProductID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrProductIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrProductIDIsInvalid]
			return resp, comerrors.ErrProductIDIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy-by-id")
	policy, err := svc.repo.SelectPolicyByPK(ctx, input.PolicyID)
	if err != nil {
		cSpan.End()
		svc.logger.GetLogger().Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrPolicyIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrPolicyIDIsInvalid]
			return resp, comerrors.ErrPolicyIDIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "generate-new-license")
	license, err := svc.generateLicense(ctx, input, tenant, product, policy)
	if err != nil {
		cSpan.End()
		svc.logger.GetLogger().Error(err.Error())
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-license")
	err = svc.repo.InsertNewLicense(ctx, license)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.LicenseRegistrationOutput{
		ID:                            license.ID.String(),
		TenantID:                      tenant.ID.String(),
		ProductID:                     product.ID.String(),
		PolicyID:                      policy.ID.String(),
		Name:                          license.Name,
		Key:                           license.Key,
		MD5:                           fmt.Sprintf("%x", md5.Sum([]byte(license.Key))),
		Sha1:                          fmt.Sprintf("%x", sha1.Sum([]byte(license.Key))),
		Sha256:                        fmt.Sprintf("%x", sha256.Sum256([]byte(license.Key))),
		Scheme:                        policy.Scheme,
		PublicKey:                     policy.PublicKey,
		ExpirationStrategy:            policy.ExpirationStrategy,
		ExpirationBasis:               policy.ExpirationBasis,
		AuthenticationStrategy:        policy.AuthenticationStrategy,
		HeartbeatCullStrategy:         policy.HeartbeatCullStrategy,
		HeartbeatResurrectionStrategy: policy.HeartbeatResurrectionStrategy,
		CheckInInterval:               policy.CheckInInterval,
		TransferStrategy:              policy.TransferStrategy,
		OverageStrategy:               policy.OverageStrategy,
		HeartbeatBasis:                policy.HeartbeatBasis,
		RenewalBasis:                  policy.RenewalBasis,
		RequireCheckIn:                policy.RequireCheckIn,
		Concurrent:                    policy.Concurrent,
		RequireHeartbeat:              policy.RequireHeartbeat,
		Strict:                        policy.Strict,
		Floating:                      policy.Floating,
		UsePool:                       policy.UsePool,
		RateLimited:                   policy.RateLimited,
		Encrypted:                     policy.Encrypted,
		Protected:                     policy.Protected,
		Duration:                      policy.Duration,
		MaxMachines:                   policy.MaxMachines,
		MaxUses:                       policy.MaxUses,
		MaxUsers:                      policy.MaxUsers,
		HeartbeatDuration:             policy.HeartbeatDuration,
		Metadata:                      license.Metadata,
		CreatedAt:                     license.CreatedAt,
		UpdatedAt:                     license.UpdatedAt,
	}
	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *LicenseService) Retrieve(ctx *gin.Context, input *models.LicenseRetrievalInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *LicenseService) Delete(ctx *gin.Context, input *models.LicenseDeletionInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *LicenseService) List(ctx *gin.Context, input *models.LicenseListInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *LicenseService) actions(ctx *gin.Context, input *models.LicenseActionInput) (*response.BaseOutput, error) {

	return nil, nil
}
