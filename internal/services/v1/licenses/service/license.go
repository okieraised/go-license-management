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
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/licenses/models"
	"go-license-management/internal/services/v1/licenses/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
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

// Create handles new license logics.
func (svc *LicenseService) Create(ctx *gin.Context, input *models.LicenseRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-product-by-id")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", utils.DerefPointer(input.ProductID)))
	product, err := svc.repo.SelectProductByPK(ctx, uuid.MustParse(utils.DerefPointer(input.ProductID)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy-by-id")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
	if err != nil {
		cSpan.End()
		svc.logger.GetLogger().Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrPolicyIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrPolicyIDIsInvalid]
			return resp, cerrors.ErrPolicyIDIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "generate-new-license")
	svc.logger.GetLogger().Info("generating new license")
	license, err := svc.generateLicense(ctx, input, tenant, product, policy)
	if err != nil {
		cSpan.End()
		svc.logger.GetLogger().Error(err.Error())
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-license")
	svc.logger.GetLogger().Info("inserting new license to database")
	err = svc.repo.InsertNewLicense(ctx, license)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
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
			PolicyScheme:       policy.Scheme,
			PolicyPublicKey:    policy.PublicKey,
			ExpirationStrategy: policy.ExpirationStrategy,
			CheckInInterval:    policy.CheckInInterval,
			OverageStrategy:    policy.OverageStrategy,
			HeartbeatBasis:     policy.HeartbeatBasis,
			RenewalBasis:       policy.RenewalBasis,
			RequireCheckIn:     policy.RequireCheckIn,
			RequireHeartbeat:   policy.RequireHeartbeat,
			Strict:             policy.Strict,
			Floating:           policy.Floating,
			UsePool:            policy.UsePool,
			RateLimited:        policy.RateLimited,
			Encrypted:          policy.Encrypted,
			Protected:          policy.Protected,
			Duration:           policy.Duration,
			MaxMachines:        policy.MaxMachines,
			MaxUses:            policy.MaxUses,
			MaxUsers:           policy.MaxUsers,
			HeartbeatDuration:  policy.HeartbeatDuration,
		},
		LicenseProduct: models.LicenseProductOutput{
			Name:                 product.Name,
			DistributionStrategy: product.DistributionStrategy,
			Code:                 product.Code,
			URL:                  product.URL,
			Platforms:            product.Platforms,
			Metadata:             product.Metadata,
			CreatedAt:            product.CreatedAt,
			UpdatedAt:            product.UpdatedAt,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

// Update handles the license update logics.
func (svc *LicenseService) Update(ctx *gin.Context, input *models.LicenseUpdateInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-license")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying license [%s]", utils.DerefPointer(input.LicenseID)))
	license, err := svc.repo.SelectLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIDIsInvalid]
			return resp, cerrors.ErrLicenseIDIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	policy := license.Policy
	// validate policyID exists and replace existing license values with new policy values
	if input.PolicyID != nil {
		_, cSpan = input.Tracer.Start(rootCtx, "check-policy-exist")
		policyID := uuid.MustParse(utils.DerefPointer(input.PolicyID))
		if policyID != license.PolicyID {
			policy, err = svc.repo.SelectPolicyByPK(ctx, policyID)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				if errors.Is(err, sql.ErrNoRows) {
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrPolicyIDIsInvalid]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrPolicyIDIsInvalid]
					return resp, cerrors.ErrPolicyIDIsInvalid
				} else {
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
			license.MaxUsers = policy.MaxUsers
			license.MaxMachines = policy.MaxMachines
			license.MaxUses = policy.MaxUses
			license.PolicyID = policy.ID
		}
		cSpan.End()
	}

	// validate productID exists
	if input.ProductID != nil {
		_, cSpan = input.Tracer.Start(rootCtx, "check-policy-exist")
		productID := uuid.MustParse(utils.DerefPointer(input.ProductID))

		exist, err := svc.repo.CheckProductExist(ctx, productID)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}

		if !exist {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
			return resp, cerrors.ErrProductIDIsInvalid
		} else {
			license.ProductID = productID
		}
		cSpan.End()
	}

	// Update expiration if specified
	var expiry time.Time
	if input.Expiry != nil {
		expiry, _ = time.Parse(time.RFC3339, utils.DerefPointer(input.Expiry))
	}

	if policy.Duration > 0 {
		if expiry.IsZero() {
			license.Expiry = time.Now().Add(time.Duration(policy.Duration) * time.Second)
		} else {
			license.Expiry = expiry.Add(time.Duration(policy.Duration) * time.Second)
		}
	} else {
		// Only update is expiry is not zero
		if !expiry.IsZero() {
			license.Expiry = expiry
		}
	}

	// Update max_users if specified
	if input.MaxUsers != nil {
		license.MaxUsers = utils.DerefPointer(input.MaxUsers)
	}

	// Update max_uses if specified
	if input.MaxUses != nil {
		license.MaxUsers = utils.DerefPointer(input.MaxUses)
	}

	// Update max_machines if specified
	if input.MaxMachines != nil {
		license.MaxMachines = utils.DerefPointer(input.MaxMachines)
	}

	// Update license name is specified
	if input.Name != nil {
		license.Name = utils.DerefPointer(input.Name)
	}

	// Update metadata is specified
	if input.Metadata != nil {
		license.Metadata = input.Metadata
	}

	_, cSpan = input.Tracer.Start(rootCtx, "update-license")
	license, err = svc.repo.UpdateLicenseByPK(ctx, license)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
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
			PolicyScheme:       license.Policy.Scheme,
			PolicyPublicKey:    license.Policy.PublicKey,
			ExpirationStrategy: license.Policy.ExpirationStrategy,
			CheckInInterval:    license.Policy.CheckInInterval,
			OverageStrategy:    license.Policy.OverageStrategy,
			HeartbeatBasis:     license.Policy.HeartbeatBasis,
			RenewalBasis:       license.Policy.RenewalBasis,
			RequireCheckIn:     license.Policy.RequireCheckIn,
			RequireHeartbeat:   license.Policy.RequireHeartbeat,
			Strict:             license.Policy.Strict,
			Floating:           license.Policy.Floating,
			UsePool:            license.Policy.UsePool,
			RateLimited:        license.Policy.RateLimited,
			Encrypted:          license.Policy.Encrypted,
			Protected:          license.Policy.Protected,
			Duration:           license.Policy.Duration,
			MaxMachines:        license.Policy.MaxMachines,
			MaxUses:            license.Policy.MaxUses,
			MaxUsers:           license.Policy.MaxUsers,
			HeartbeatDuration:  license.Policy.HeartbeatDuration,
		},
		LicenseProduct: models.LicenseProductOutput{
			Name:                 license.Product.Name,
			DistributionStrategy: license.Product.DistributionStrategy,
			Code:                 license.Product.Code,
			URL:                  license.Product.URL,
			Platforms:            license.Product.Platforms,
			Metadata:             license.Product.Metadata,
			CreatedAt:            license.Product.CreatedAt,
			UpdatedAt:            license.Product.UpdatedAt,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

// Retrieve handles the license retrieval logics.
func (svc *LicenseService) Retrieve(ctx *gin.Context, input *models.LicenseRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-license")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying license [%s]", utils.DerefPointer(input.LicenseID)))
	license, err := svc.repo.SelectLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIDIsInvalid]
			return resp, cerrors.ErrLicenseIDIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
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
			PolicyScheme:       license.Policy.Scheme,
			PolicyPublicKey:    license.Policy.PublicKey,
			ExpirationStrategy: license.Policy.ExpirationStrategy,
			CheckInInterval:    license.Policy.CheckInInterval,
			OverageStrategy:    license.Policy.OverageStrategy,
			HeartbeatBasis:     license.Policy.HeartbeatBasis,
			RenewalBasis:       license.Policy.RenewalBasis,
			RequireCheckIn:     license.Policy.RequireCheckIn,
			RequireHeartbeat:   license.Policy.RequireHeartbeat,
			Strict:             license.Policy.Strict,
			Floating:           license.Policy.Floating,
			UsePool:            license.Policy.UsePool,
			RateLimited:        license.Policy.RateLimited,
			Encrypted:          license.Policy.Encrypted,
			Protected:          license.Policy.Protected,
			Duration:           license.Policy.Duration,
			MaxMachines:        license.Policy.MaxMachines,
			MaxUses:            license.Policy.MaxUses,
			MaxUsers:           license.Policy.MaxUsers,
			HeartbeatDuration:  license.Policy.HeartbeatDuration,
		},
		LicenseProduct: models.LicenseProductOutput{
			Name:                 license.Product.Name,
			DistributionStrategy: license.Product.DistributionStrategy,
			Code:                 license.Product.Code,
			URL:                  license.Product.URL,
			Platforms:            license.Product.Platforms,
			Metadata:             license.Product.Metadata,
			CreatedAt:            license.Product.CreatedAt,
			UpdatedAt:            license.Product.UpdatedAt,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *LicenseService) Delete(ctx *gin.Context, input *models.LicenseDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "delete-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-license")
	svc.logger.GetLogger().Info(fmt.Sprintf("deleting license [%s]", utils.DerefPointer(input.LicenseID)))
	err = svc.repo.DeleteLicenseByPK(ctx, uuid.MustParse(utils.DerefPointer(input.LicenseID)))
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

func (svc *LicenseService) List(ctx *gin.Context, input *models.LicenseListInput) (*response.BaseOutput, error) {
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
	licenses, total, err := svc.repo.SelectLicenses(ctx, tenant.Name, input.QueryCommonParam)
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

	licenseOutput := make([]models.LicenseInfoOutput, 0)
	for _, license := range licenses {
		licenseOutput = append(licenseOutput, models.LicenseInfoOutput{
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
				PolicyScheme:       license.Policy.Scheme,
				PolicyPublicKey:    license.Policy.PublicKey,
				ExpirationStrategy: license.Policy.ExpirationStrategy,
				CheckInInterval:    license.Policy.CheckInInterval,
				OverageStrategy:    license.Policy.OverageStrategy,
				HeartbeatBasis:     license.Policy.HeartbeatBasis,
				RenewalBasis:       license.Policy.RenewalBasis,
				RequireCheckIn:     license.Policy.RequireCheckIn,
				RequireHeartbeat:   license.Policy.RequireHeartbeat,
				Strict:             license.Policy.Strict,
				Floating:           license.Policy.Floating,
				UsePool:            license.Policy.UsePool,
				RateLimited:        license.Policy.RateLimited,
				Encrypted:          license.Policy.Encrypted,
				Protected:          license.Policy.Protected,
				Duration:           license.Policy.Duration,
				MaxMachines:        license.Policy.MaxMachines,
				MaxUses:            license.Policy.MaxUses,
				MaxUsers:           license.Policy.MaxUsers,
				HeartbeatDuration:  license.Policy.HeartbeatDuration,
			},
			LicenseProduct: models.LicenseProductOutput{
				Name:                 license.Product.Name,
				DistributionStrategy: license.Product.DistributionStrategy,
				Code:                 license.Product.Code,
				URL:                  license.Product.URL,
				Platforms:            license.Product.Platforms,
				Metadata:             license.Product.Metadata,
				CreatedAt:            license.Product.CreatedAt,
				UpdatedAt:            license.Product.UpdatedAt,
			},
		})
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = licenseOutput

	return resp, nil
}

func (svc *LicenseService) Actions(ctx *gin.Context, input *models.LicenseActionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "action-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)
	svc.logger.GetLogger().Info(fmt.Sprintf("received action [%s]", utils.DerefPointer(input.Action)))

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying tenant [%s]", utils.DerefPointer(input.TenantName)))
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-license")
	license, err := svc.repo.SelectLicenseByLicenseKey(ctx, utils.DerefPointer(input.LicenseKey))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseKeyIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseKeyIsInvalid]
			return resp, cerrors.ErrLicenseKeyIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "perform-license-action")
	licenseAction := utils.DerefPointer(input.Action)
	switch licenseAction {
	case constants.LicenseActionValidate:
		output, err := svc.validateLicense(ctx, license)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		resp.Data = output
	case constants.LicenseActionCheckout:
		output, err := svc.checkoutLicense(ctx, license)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		resp.Data = output
	default:
		var output *entities.License
		switch licenseAction {
		case constants.LicenseActionCheckin:
			output, err = svc.checkinLicense(ctx, license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		case constants.LicenseActionSuspend:
			output, err = svc.suspendLicense(ctx, license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				switch {
				case errors.Is(err, cerrors.ErrLicenseNotActivated):
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseNotActivated]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseNotActivated]
					return resp, cerrors.ErrLicenseNotActivated
				default:
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
		case constants.LicenseActionReinstate:
			output, err = svc.reinstateLicense(ctx, license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				switch {
				case errors.Is(err, cerrors.ErrLicenseStatusInvalidToReinstate):
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseStatusInvalidToReinstate]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseStatusInvalidToReinstate]
					return resp, cerrors.ErrLicenseStatusInvalidToReinstate
				default:
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
		case constants.LicenseActionRenew:
			output, err = svc.renewLicense(ctx, license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		case constants.LicenseActionIncrementUsage:
			output, err = svc.incrementUsageLicense(ctx, utils.DerefPointer(input.Increment), license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				switch {
				case errors.Is(err, cerrors.ErrLicenseMaxUsesExceeded):
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseMaxUsesExceeded]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseMaxUsesExceeded]
					return resp, cerrors.ErrLicenseMaxUsesExceeded
				default:
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
		case constants.LicenseActionDecrementUsage:
			output, err = svc.decrementUsageLicense(ctx, utils.DerefPointer(input.Decrement), license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				switch {
				case errors.Is(err, cerrors.ErrLicenseIncrementIsInvalid):
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrLicenseIncrementIsInvalid]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrLicenseIncrementIsInvalid]
					return resp, cerrors.ErrLicenseIncrementIsInvalid
				default:
					resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
					resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
					return resp, cerrors.ErrGenericInternalServer
				}
			}
		case constants.LicenseActionResetUsage:
			output, err = svc.resetUsageLicense(ctx, license)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		}
		resp.Data = models.LicenseInfoOutput{
			LicenseID:      output.ID.String(),
			ProductID:      output.ProductID.String(),
			PolicyID:       output.PolicyID.String(),
			Name:           output.Name,
			LicenseKey:     output.Key,
			MD5Checksum:    fmt.Sprintf("%x", md5.Sum([]byte(output.Key))),
			Sha1Checksum:   fmt.Sprintf("%x", sha1.Sum([]byte(output.Key))),
			Sha256Checksum: fmt.Sprintf("%x", sha256.Sum256([]byte(output.Key))),
			Status:         output.Status,
			Metadata:       output.Metadata,
			Expiry:         output.Expiry,
			CreatedAt:      output.CreatedAt,
			UpdatedAt:      output.UpdatedAt,
			LicensePolicy: models.LicensePolicyOutput{
				PolicyScheme:       license.Policy.Scheme,
				PolicyPublicKey:    license.Policy.PublicKey,
				ExpirationStrategy: license.Policy.ExpirationStrategy,
				CheckInInterval:    license.Policy.CheckInInterval,
				OverageStrategy:    license.Policy.OverageStrategy,
				HeartbeatBasis:     license.Policy.HeartbeatBasis,
				RenewalBasis:       license.Policy.RenewalBasis,
				RequireCheckIn:     license.Policy.RequireCheckIn,
				RequireHeartbeat:   license.Policy.RequireHeartbeat,
				Strict:             license.Policy.Strict,
				Floating:           license.Policy.Floating,
				UsePool:            license.Policy.UsePool,
				RateLimited:        license.Policy.RateLimited,
				Encrypted:          license.Policy.Encrypted,
				Protected:          license.Policy.Protected,
				Duration:           license.Policy.Duration,
				MaxMachines:        license.Policy.MaxMachines,
				MaxUses:            license.Policy.MaxUses,
				MaxUsers:           license.Policy.MaxUsers,
				HeartbeatDuration:  license.Policy.HeartbeatDuration,
			},
			LicenseProduct: models.LicenseProductOutput{
				Name:                 license.Product.Name,
				DistributionStrategy: license.Product.DistributionStrategy,
				Code:                 license.Product.Code,
				URL:                  license.Product.URL,
				Platforms:            license.Product.Platforms,
				Metadata:             license.Product.Metadata,
				CreatedAt:            license.Product.CreatedAt,
				UpdatedAt:            license.Product.UpdatedAt,
			},
		}
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	return resp, nil
}
