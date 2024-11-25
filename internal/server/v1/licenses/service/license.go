package service

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
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
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", utils.DerefPointer(input.ProductID)))
	product, err := svc.repo.SelectProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
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
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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
	svc.logger.GetLogger().Info("generating new license")
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
	svc.logger.GetLogger().Info("inserting new license to database")
	err = svc.repo.InsertNewLicense(ctx, license)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.LicenseInfoOutput{
		LicenseID:      license.ID.String(),
		ProductID:      product.ID.String(),
		PolicyID:       policy.ID.String(),
		Name:           license.Name,
		LicenseKey:     license.Key,
		MD5Checksum:    fmt.Sprintf("%x", md5.Sum([]byte(license.Key))),
		Sha1Checksum:   fmt.Sprintf("%x", sha1.Sum([]byte(license.Key))),
		Sha256Checksum: fmt.Sprintf("%x", sha256.Sum256([]byte(license.Key))),
		Status:         license.Status,
		Metadata:       license.Metadata,
		Expiry:         license.Expiry,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
		LicensePolicy: models.LicensePolicyOutput{
			PolicyScheme:           policy.Scheme,
			PolicyPublicKey:        policy.PublicKey,
			ExpirationStrategy:     policy.ExpirationStrategy,
			ExpirationBasis:        policy.ExpirationBasis,
			AuthenticationStrategy: policy.AuthenticationStrategy,
			CheckInInterval:        policy.CheckInInterval,
			OverageStrategy:        policy.OverageStrategy,
			HeartbeatBasis:         policy.HeartbeatBasis,
			RenewalBasis:           policy.RenewalBasis,
			RequireCheckIn:         policy.RequireCheckIn,
			RequireHeartbeat:       policy.RequireHeartbeat,
			Strict:                 policy.Strict,
			Floating:               policy.Floating,
			UsePool:                policy.UsePool,
			RateLimited:            policy.RateLimited,
			Encrypted:              policy.Encrypted,
			Protected:              policy.Protected,
			Duration:               policy.Duration,
			MaxMachines:            policy.MaxMachines,
			MaxUses:                policy.MaxUses,
			MaxUsers:               policy.MaxUsers,
			HeartbeatDuration:      policy.HeartbeatDuration,
		},
	}
	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *LicenseService) Retrieve(ctx *gin.Context, input *models.LicenseRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-product")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying license [%s]", utils.DerefPointer(input.LicenseID)))
	license, err := svc.repo.SelectLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrLicenseIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrLicenseIDIsInvalid]
			return resp, comerrors.ErrLicenseIDIsInvalid
		default:
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.LicenseInfoOutput{
		LicenseID:      license.ID.String(),
		ProductID:      license.ProductID.String(),
		PolicyID:       license.PolicyID.String(),
		Name:           license.Name,
		LicenseKey:     license.Key,
		MD5Checksum:    fmt.Sprintf("%x", md5.Sum([]byte(license.Key))),
		Sha1Checksum:   fmt.Sprintf("%x", sha1.Sum([]byte(license.Key))),
		Sha256Checksum: fmt.Sprintf("%x", sha256.Sum256([]byte(license.Key))),
		Status:         license.Status,
		Metadata:       license.Metadata,
		Expiry:         license.Expiry,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
		LicensePolicy: models.LicensePolicyOutput{
			PolicyScheme:           license.Policy.Scheme,
			PolicyPublicKey:        license.Policy.PublicKey,
			ExpirationStrategy:     license.Policy.ExpirationStrategy,
			ExpirationBasis:        license.Policy.ExpirationBasis,
			AuthenticationStrategy: license.Policy.AuthenticationStrategy,
			CheckInInterval:        license.Policy.CheckInInterval,
			OverageStrategy:        license.Policy.OverageStrategy,
			HeartbeatBasis:         license.Policy.HeartbeatBasis,
			RenewalBasis:           license.Policy.RenewalBasis,
			RequireCheckIn:         license.Policy.RequireCheckIn,
			RequireHeartbeat:       license.Policy.RequireHeartbeat,
			Strict:                 license.Policy.Strict,
			Floating:               license.Policy.Floating,
			UsePool:                license.Policy.UsePool,
			RateLimited:            license.Policy.RateLimited,
			Encrypted:              license.Policy.Encrypted,
			Protected:              license.Policy.Protected,
			Duration:               license.Policy.Duration,
			MaxMachines:            license.Policy.MaxMachines,
			MaxUses:                license.Policy.MaxUses,
			MaxUsers:               license.Policy.MaxUsers,
			HeartbeatDuration:      license.Policy.HeartbeatDuration,
		},
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *LicenseService) Delete(ctx *gin.Context, input *models.LicenseDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-license")
	svc.logger.GetLogger().Info(fmt.Sprintf("deleting license [%s]", utils.DerefPointer(input.LicenseID)))
	err = svc.repo.DeleteLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
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

func (svc *LicenseService) List(ctx *gin.Context, input *models.LicenseListInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *LicenseService) Actions(ctx *gin.Context, input *models.LicenseActionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "action-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))
	svc.logger.GetLogger().Info(fmt.Sprintf("received action [%s]", utils.DerefPointer(input.Action)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-license")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying license [%s]", utils.DerefPointer(input.LicenseID)))
	license, err := svc.repo.SelectLicenseByLicenseKey(ctx, utils.DerefPointer(input.LicenseKey))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrLicenseIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrLicenseIDIsInvalid]
			return resp, comerrors.ErrLicenseIDIsInvalid
		} else {
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	fmt.Println("license", license)

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	return resp, nil
}

// validateLicense validates a license. This will check the following: if the license is suspended, if the license is expired,
// if the license is overdue for check-in, and if the license meets its machine requirements (if strict).
func (svc *LicenseService) validateLicense(ctx *gin.Context, license *entities.License) error {
	if license.Status == constants.LicenseStatusNotActivated {

	}

	return nil
}

func (svc *LicenseService) revokeLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}

func (svc *LicenseService) suspendLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}

func (svc *LicenseService) reinstateLicense(ctx *gin.Context, license *entities.License) error {
	return nil
}
