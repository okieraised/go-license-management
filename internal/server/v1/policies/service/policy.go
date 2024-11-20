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
	"go-license-management/internal/server/v1/policies/models"
	"go-license-management/internal/server/v1/policies/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type PolicyService struct {
	repo   repository.IPolicy
	logger *logging.Logger
}

func NewPolicyService(options ...func(*PolicyService)) *PolicyService {
	svc := &PolicyService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IPolicy) func(*PolicyService) {
	return func(c *PolicyService) {
		c.repo = repo
	}
}

func (svc *PolicyService) Create(ctx *gin.Context, input *models.PolicyRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	// Check if tenant exists
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

	productID := uuid.MustParse(utils.DerefPointer(input.ProductID))

	// Check if productID exists
	_, cSpan = input.Tracer.Start(rootCtx, "check-product-id")
	exists, err := svc.repo.CheckProductExistByID(ctx, productID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if !exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("product id [%s] does not exist in tenant [%s]", productID.String(), tenant.Name))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrProductIDIsInvalid]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrProductIDIsInvalid]
		return resp, comerrors.ErrProductIDIsInvalid
	}

	// Generate new private/public key pair
	var privateKey = ""
	var publicKey = ""
	scheme := utils.DerefPointer(input.Scheme)
	switch scheme {
	case constants.PolicySchemeED25519:
		privateKey, publicKey, err = utils.NewEd25519KeyPair()
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	case constants.PolicySchemeRSA2048PKCS1:
		privateKey, publicKey, err = utils.NewRSA2048PKCS1KeyPair()
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	default:
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrPolicySchemeIsInvalid]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrPolicySchemeIsInvalid]
		return resp, comerrors.ErrPolicySchemeIsInvalid
	}

	// Insert new policy
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-policy")
	policyID := uuid.New()
	now := time.Now()
	policy := &entities.Policy{
		ID:                            policyID,
		TenantName:                    tenant.Name,
		ProductID:                     productID,
		Duration:                      int64(utils.DerefPointer(input.Duration)),
		MaxMachines:                   utils.DerefPointer(input.MaxMachines),
		MaxUses:                       utils.DerefPointer(input.MaxUses),
		HeartbeatDuration:             utils.DerefPointer(input.HeartbeatDuration),
		MaxUsers:                      utils.DerefPointer(input.MaxUsers),
		Strict:                        utils.DerefPointer(input.Strict),
		Floating:                      utils.DerefPointer(input.Floating),
		UsePool:                       utils.DerefPointer(input.UsePool),
		RateLimited:                   utils.DerefPointer(input.RateLimited),
		Encrypted:                     utils.DerefPointer(input.Encrypted),
		Protected:                     utils.DerefPointer(input.Protected),
		RequireCheckIn:                utils.DerefPointer(input.RequireCheckIn),
		Concurrent:                    utils.DerefPointer(input.Concurrent),
		RequireHeartbeat:              utils.DerefPointer(input.RequireHeartbeat),
		PublicKey:                     publicKey,
		PrivateKey:                    privateKey,
		Name:                          utils.DerefPointer(input.Name),
		Scheme:                        scheme,
		ExpirationStrategy:            utils.DerefPointer(input.ExpirationStrategy),
		ExpirationBasis:               utils.DerefPointer(input.ExpirationBasis),
		AuthenticationStrategy:        utils.DerefPointer(input.AuthenticationStrategy),
		HeartbeatCullStrategy:         utils.DerefPointer(input.HeartbeatCullStrategy),
		HeartbeatResurrectionStrategy: utils.DerefPointer(input.HeartbeatResurrectionStrategy),
		CheckInInterval:               utils.DerefPointer(input.CheckInInterval),
		OverageStrategy:               utils.DerefPointer(input.OverageStrategy),
		HeartbeatBasis:                utils.DerefPointer(input.HeartbeatBasis),
		RenewalBasis:                  utils.DerefPointer(input.RenewalBasis),
		Metadata:                      input.Metadata,
		CreatedAt:                     now,
		UpdatedAt:                     now,
	}
	err = svc.repo.InsertNewPolicy(ctx, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.PolicyRegistrationOutput{
		ID:                            policyID.String(),
		TenantID:                      tenant.Name,
		ProductID:                     productID.String(),
		Name:                          policy.Name,
		Scheme:                        policy.Scheme,
		Duration:                      policy.Duration,
		MaxMachines:                   policy.MaxMachines,
		MaxUses:                       policy.MaxUses,
		MaxUsers:                      policy.MaxUsers,
		CheckInIntervalCount:          policy.CheckInIntervalCount,
		HeartbeatDuration:             policy.HeartbeatDuration,
		Strict:                        policy.Strict,
		Floating:                      policy.Floating,
		UsePool:                       policy.UsePool,
		RateLimited:                   policy.RateLimited,
		Encrypted:                     policy.Encrypted,
		Protected:                     policy.Protected,
		RequireCheckIn:                policy.RequireCheckIn,
		Concurrent:                    policy.Concurrent,
		RequireHeartbeat:              policy.RequireHeartbeat,
		PublicKey:                     policy.PublicKey,
		ExpirationStrategy:            policy.ExpirationStrategy,
		ExpirationBasis:               policy.ExpirationBasis,
		AuthenticationStrategy:        policy.AuthenticationStrategy,
		HeartbeatCullStrategy:         policy.HeartbeatCullStrategy,
		HeartbeatResurrectionStrategy: policy.HeartbeatResurrectionStrategy,
		CheckInInterval:               policy.CheckInInterval,
		OverageStrategy:               policy.OverageStrategy,
		HeartbeatBasis:                policy.HeartbeatBasis,
		RenewalBasis:                  policy.RenewalBasis,
		Metadata:                      policy.Metadata,
		CreatedAt:                     policy.CreatedAt,
		UpdatedAt:                     policy.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *PolicyService) List(ctx *gin.Context, input *models.PolicyListInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policies")
	products, total, err := svc.repo.SelectPolicies(ctx, tenant.Name, input.QueryCommonParam)
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

	policiesOutput := make([]models.PolicyRetrievalOutput, 0)
	for _, policy := range products {
		policiesOutput = append(policiesOutput, models.PolicyRetrievalOutput{
			ID:                            policy.ID.String(),
			TenantName:                    policy.TenantName,
			ProductID:                     policy.ProductID.String(),
			PublicKey:                     policy.PublicKey,
			Name:                          policy.Name,
			Scheme:                        policy.Scheme,
			ExpirationStrategy:            policy.ExpirationStrategy,
			ExpirationBasis:               policy.ExpirationBasis,
			AuthenticationStrategy:        policy.AuthenticationStrategy,
			HeartbeatCullStrategy:         policy.HeartbeatCullStrategy,
			HeartbeatResurrectionStrategy: policy.HeartbeatResurrectionStrategy,
			CheckInInterval:               policy.CheckInInterval,
			OverageStrategy:               policy.OverageStrategy,
			HeartbeatBasis:                policy.HeartbeatBasis,
			RenewalBasis:                  policy.RenewalBasis,
			Metadata:                      policy.Metadata,
			Duration:                      policy.Duration,
			MaxMachines:                   policy.MaxMachines,
			MaxUses:                       policy.MaxUses,
			MaxUsers:                      policy.MaxUsers,
			CheckInIntervalCount:          policy.CheckInIntervalCount,
			HeartbeatDuration:             policy.HeartbeatDuration,
			Strict:                        policy.Strict,
			Floating:                      policy.Floating,
			UsePool:                       policy.UsePool,
			RateLimited:                   policy.RateLimited,
			Encrypted:                     policy.Encrypted,
			Protected:                     policy.Protected,
			RequireCheckIn:                policy.RequireCheckIn,
			Concurrent:                    policy.Concurrent,
			RequireHeartbeat:              policy.RequireHeartbeat,
			CreatedAt:                     policy.CreatedAt,
			UpdatedAt:                     policy.UpdatedAt,
		})
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = policiesOutput
	return resp, nil
}

func (svc *PolicyService) Retrieve(ctx *gin.Context, input *models.PolicyRetrievalInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-product")
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrPolicyIDIsInvalid]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrPolicyIDIsInvalid]
			return resp, comerrors.ErrPolicyIDIsInvalid
		default:
			resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
			resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
			return resp, comerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.PolicyRetrievalOutput{
		ID:                            policy.ID.String(),
		TenantName:                    policy.TenantName,
		ProductID:                     policy.ProductID.String(),
		PublicKey:                     policy.PublicKey,
		Name:                          policy.Name,
		Scheme:                        policy.Scheme,
		ExpirationStrategy:            policy.ExpirationStrategy,
		ExpirationBasis:               policy.ExpirationBasis,
		AuthenticationStrategy:        policy.AuthenticationStrategy,
		HeartbeatCullStrategy:         policy.HeartbeatCullStrategy,
		HeartbeatResurrectionStrategy: policy.HeartbeatResurrectionStrategy,
		CheckInInterval:               policy.CheckInInterval,
		OverageStrategy:               policy.OverageStrategy,
		HeartbeatBasis:                policy.HeartbeatBasis,
		RenewalBasis:                  policy.RenewalBasis,
		Metadata:                      policy.Metadata,
		Duration:                      policy.Duration,
		MaxMachines:                   policy.MaxMachines,
		MaxUses:                       policy.MaxUses,
		MaxUsers:                      policy.MaxUsers,
		CheckInIntervalCount:          policy.CheckInIntervalCount,
		HeartbeatDuration:             policy.HeartbeatDuration,
		Strict:                        policy.Strict,
		Floating:                      policy.Floating,
		UsePool:                       policy.UsePool,
		RateLimited:                   policy.RateLimited,
		Encrypted:                     policy.Encrypted,
		Protected:                     policy.Protected,
		RequireCheckIn:                policy.RequireCheckIn,
		Concurrent:                    policy.Concurrent,
		RequireHeartbeat:              policy.RequireHeartbeat,
		CreatedAt:                     policy.CreatedAt,
		UpdatedAt:                     policy.UpdatedAt,
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *PolicyService) Delete(ctx *gin.Context, input *models.PolicyDeletionInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-policy")
	err = svc.repo.DeletePolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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

func (svc *PolicyService) Update(ctx *gin.Context, input *models.PolicyUpdateInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Attach(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Detach(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) ListEntitlements(ctx *gin.Context) (*response.BaseOutput, error) {
	return nil, nil
}
