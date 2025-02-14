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
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/policies/models"
	"go-license-management/internal/services/v1/policies/repository"
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
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	// Check if tenant exists
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

	// Check if productID exists
	_, cSpan = input.Tracer.Start(rootCtx, "check-product-id")
	productID := uuid.MustParse(utils.DerefPointer(input.ProductID))
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", productID))
	exists, err := svc.repo.CheckProductExistByID(ctx, productID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}

	if !exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("product id [%s] does not exist in tenant [%s]", productID.String(), tenant.Name))
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrProductIDIsInvalid]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrProductIDIsInvalid]
		return resp, cerrors.ErrProductIDIsInvalid
	}
	cSpan.End()

	// Generate new private/public key pair
	var privateKey = ""
	var publicKey = ""
	scheme := utils.DerefPointer(input.Scheme)
	svc.logger.GetLogger().Info(fmt.Sprintf("generating private/public key pair using [%s] algorithm", scheme))
	switch scheme {
	case constants.PolicySchemeED25519:
		privateKey, publicKey, err = utils.NewEd25519KeyPair()
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	case constants.PolicySchemeRSA2048PKCS1:
		privateKey, publicKey, err = utils.NewRSA2048PKCS1KeyPair()
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	default:
		svc.logger.GetLogger().Error(fmt.Sprintf("invalid supported sheme [%s]", scheme))
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrPolicySchemeIsInvalid]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrPolicySchemeIsInvalid]
		return resp, cerrors.ErrPolicySchemeIsInvalid
	}

	// Insert new policy
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-policy")
	svc.logger.GetLogger().Info("inserting new policy to database")
	policyID := uuid.New()
	now := time.Now()
	policy := &entities.Policy{
		ID:                 policyID,
		ProductID:          productID,
		TenantName:         tenant.Name,
		PublicKey:          publicKey,
		PrivateKey:         privateKey,
		Name:               utils.DerefPointer(input.Name),
		Scheme:             scheme,
		ExpirationStrategy: utils.DerefPointer(input.ExpirationStrategy),
		CheckInInterval:    utils.DerefPointer(input.CheckInInterval),
		OverageStrategy:    utils.DerefPointer(input.OverageStrategy),
		HeartbeatBasis:     utils.DerefPointer(input.HeartbeatBasis),
		RenewalBasis:       utils.DerefPointer(input.RenewalBasis),
		Duration:           utils.DerefPointer(input.Duration),
		MaxMachines:        utils.DerefPointer(input.MaxMachines),
		MaxUses:            utils.DerefPointer(input.MaxUses),
		MaxUsers:           utils.DerefPointer(input.MaxUsers),
		HeartbeatDuration:  utils.DerefPointer(input.HeartbeatDuration),
		Strict:             utils.DerefPointer(input.Strict),
		Floating:           utils.DerefPointer(input.Floating),
		UsePool:            utils.DerefPointer(input.UsePool),
		RateLimited:        utils.DerefPointer(input.RateLimited),
		Encrypted:          utils.DerefPointer(input.Encrypted),
		Protected:          utils.DerefPointer(input.Protected),
		RequireCheckIn:     utils.DerefPointer(input.RequireCheckIn),
		RequireHeartbeat:   utils.DerefPointer(input.RequireHeartbeat),
		Metadata:           input.Metadata,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = svc.repo.InsertNewPolicy(ctx, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.PolicyRetrievalOutput{
		ID:         policyID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:               utils.RefPointer(policy.Name),
			Scheme:             utils.RefPointer(policy.Scheme),
			Strict:             utils.RefPointer(policy.Strict),
			RateLimited:        utils.RefPointer(policy.RateLimited),
			Floating:           utils.RefPointer(policy.Floating),
			UsePool:            utils.RefPointer(policy.UsePool),
			Encrypted:          utils.RefPointer(policy.Encrypted),
			Protected:          utils.RefPointer(policy.Protected),
			RequireCheckIn:     utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:   utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:        utils.RefPointer(policy.MaxMachines),
			MaxUsers:           utils.RefPointer(policy.MaxUsers),
			MaxUses:            utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:  utils.RefPointer(policy.HeartbeatDuration),
			Duration:           utils.RefPointer(policy.Duration),
			CheckInInterval:    utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:     utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy: utils.RefPointer(policy.ExpirationStrategy),
			RenewalBasis:       utils.RefPointer(policy.RenewalBasis),
			OverageStrategy:    utils.RefPointer(policy.OverageStrategy),
			Metadata:           policy.Metadata,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *PolicyService) List(ctx *gin.Context, input *models.PolicyListInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policies")
	products, total, err := svc.repo.SelectPolicies(ctx, tenant.Name, input.QueryCommonParam)
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

	policiesOutput := make([]models.PolicyRetrievalOutput, 0)
	for _, policy := range products {
		policiesOutput = append(policiesOutput, models.PolicyRetrievalOutput{
			ID:         policy.ID.String(),
			TenantName: policy.TenantName,
			PublicKey:  policy.PublicKey,
			CreatedAt:  policy.CreatedAt,
			UpdatedAt:  policy.UpdatedAt,
			PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
				Name:               utils.RefPointer(policy.Name),
				Scheme:             utils.RefPointer(policy.Scheme),
				Strict:             utils.RefPointer(policy.Strict),
				RateLimited:        utils.RefPointer(policy.RateLimited),
				Floating:           utils.RefPointer(policy.Floating),
				UsePool:            utils.RefPointer(policy.UsePool),
				Encrypted:          utils.RefPointer(policy.Encrypted),
				Protected:          utils.RefPointer(policy.Protected),
				RequireCheckIn:     utils.RefPointer(policy.RequireCheckIn),
				RequireHeartbeat:   utils.RefPointer(policy.RequireHeartbeat),
				MaxMachines:        utils.RefPointer(policy.MaxMachines),
				MaxUsers:           utils.RefPointer(policy.MaxUsers),
				MaxUses:            utils.RefPointer(policy.MaxUses),
				HeartbeatDuration:  utils.RefPointer(policy.HeartbeatDuration),
				Duration:           utils.RefPointer(policy.Duration),
				CheckInInterval:    utils.RefPointer(policy.CheckInInterval),
				HeartbeatBasis:     utils.RefPointer(policy.HeartbeatBasis),
				ExpirationStrategy: utils.RefPointer(policy.ExpirationStrategy),
				RenewalBasis:       utils.RefPointer(policy.RenewalBasis),
				OverageStrategy:    utils.RefPointer(policy.OverageStrategy),
				Metadata:           policy.Metadata,
			},
		})
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = policiesOutput
	return resp, nil
}

func (svc *PolicyService) Retrieve(ctx *gin.Context, input *models.PolicyRetrievalInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-product")
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrPolicyIDIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrPolicyIDIsInvalid]
			return resp, cerrors.ErrPolicyIDIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	respData := &models.PolicyRetrievalOutput{
		ID:         policy.ID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:               utils.RefPointer(policy.Name),
			Scheme:             utils.RefPointer(policy.Scheme),
			Strict:             utils.RefPointer(policy.Strict),
			RateLimited:        utils.RefPointer(policy.RateLimited),
			Floating:           utils.RefPointer(policy.Floating),
			UsePool:            utils.RefPointer(policy.UsePool),
			Encrypted:          utils.RefPointer(policy.Encrypted),
			Protected:          utils.RefPointer(policy.Protected),
			RequireCheckIn:     utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:   utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:        utils.RefPointer(policy.MaxMachines),
			MaxUsers:           utils.RefPointer(policy.MaxUsers),
			MaxUses:            utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:  utils.RefPointer(policy.HeartbeatDuration),
			Duration:           utils.RefPointer(policy.Duration),
			CheckInInterval:    utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:     utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy: utils.RefPointer(policy.ExpirationStrategy),
			RenewalBasis:       utils.RefPointer(policy.RenewalBasis),
			OverageStrategy:    utils.RefPointer(policy.OverageStrategy),
			Metadata:           policy.Metadata,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *PolicyService) Delete(ctx *gin.Context, input *models.PolicyDeletionInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("deleting policy [%s] and associated licenses", utils.DerefPointer(input.PolicyID)))
	err = svc.repo.DeletePolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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

func (svc *PolicyService) Update(ctx *gin.Context, input *models.PolicyUpdateInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "update-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	// Check if tenant exists
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

	// Query policy
	_, cSpan = input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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
	cSpan.End()

	// Update fields
	policy, err = svc.updatePolicyField(ctx, input, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}

	// Update existing policy
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-policy")
	svc.logger.GetLogger().Info("updating policy to database")
	policy.UpdatedAt = time.Now()
	err = svc.repo.UpdatePolicyByPK(ctx, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.PolicyRetrievalOutput{
		ID:         policy.ID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:               utils.RefPointer(policy.Name),
			Scheme:             utils.RefPointer(policy.Scheme),
			Strict:             utils.RefPointer(policy.Strict),
			RateLimited:        utils.RefPointer(policy.RateLimited),
			Floating:           utils.RefPointer(policy.Floating),
			UsePool:            utils.RefPointer(policy.UsePool),
			Encrypted:          utils.RefPointer(policy.Encrypted),
			Protected:          utils.RefPointer(policy.Protected),
			RequireCheckIn:     utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:   utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:        utils.RefPointer(policy.MaxMachines),
			MaxUsers:           utils.RefPointer(policy.MaxUsers),
			MaxUses:            utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:  utils.RefPointer(policy.HeartbeatDuration),
			Duration:           utils.RefPointer(policy.Duration),
			CheckInInterval:    utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:     utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy: utils.RefPointer(policy.ExpirationStrategy),
			RenewalBasis:       utils.RefPointer(policy.RenewalBasis),
			OverageStrategy:    utils.RefPointer(policy.OverageStrategy),
			Metadata:           policy.Metadata,
		},
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *PolicyService) Attach(ctx *gin.Context, input *models.PolicyAttachmentInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "attach-policy-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "query-entitlement")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying entitlements [%s]", input.EntitlementID))
	entitlementIDs := make([]uuid.UUID, 0)
	for _, id := range input.EntitlementID {
		entitlementIDs = append(entitlementIDs, uuid.MustParse(id))
	}

	entitlements, err := svc.repo.SelectEntitlementsByPK(ctx, entitlementIDs)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if len(entitlements) == 0 {
		svc.logger.GetLogger().Error("no entitlement record found")
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrEntitlementIDIsInvalid]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrEntitlementIDIsInvalid]
		return resp, cerrors.ErrEntitlementIDIsInvalid
	} else {
		for _, entitlement := range entitlements {
			exist, err := svc.repo.CheckPolicyEntitlementExistsByPolicyIDAndEntitlementID(ctx, policy.ID, entitlement.ID)
			if err != nil {
				svc.logger.GetLogger().Error(err.Error())
				cSpan.End()
			}
			if exist {
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrPolicyEntitlementAlreadyExist]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrPolicyEntitlementAlreadyExist]
				return resp, cerrors.ErrPolicyEntitlementAlreadyExist
			}
		}
	}

	_, cSpan = input.Tracer.Start(rootCtx, "insert-policy-entitlement")
	policyEntitlements := make([]entities.PolicyEntitlement, 0)
	policyEntitlementsOutput := make([]models.PolicyAttachmentOutput, 0)
	for _, entitlement := range entitlements {
		policyEntitlementID := uuid.New()
		now := time.Now()
		policyEntitlements = append(policyEntitlements, entities.PolicyEntitlement{
			ID:            policyEntitlementID,
			TenantName:    tenant.Name,
			PolicyID:      policy.ID,
			EntitlementID: entitlement.ID,
			Metadata:      nil,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		policyEntitlementsOutput = append(policyEntitlementsOutput, models.PolicyAttachmentOutput{
			ID:            policyEntitlementID.String(),
			TenantName:    tenant.Name,
			PolicyID:      policy.ID.String(),
			EntitlementID: entitlement.ID.String(),
			Metadata:      nil,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	err = svc.repo.InsertNewPolicyEntitlements(ctx, policyEntitlements)
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
	resp.Data = policyEntitlementsOutput

	return resp, nil
}

func (svc *PolicyService) Detach(ctx *gin.Context, input *models.PolicyDetachmentInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "detach-policy-handler")
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying policy [%s]", utils.DerefPointer(input.PolicyID)))
	_, err = svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "delete-policy-entitlement")
	svc.logger.GetLogger().Info(fmt.Sprintf("deleting policy entitlements [%s]", input.ID))
	policyEntitlementIDs := make([]uuid.UUID, 0)
	for _, id := range input.ID {
		policyEntitlementIDs = append(policyEntitlementIDs, uuid.MustParse(id))
	}

	err = svc.repo.DeletePolicyEntitlementsByPK(ctx, policyEntitlementIDs)
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

func (svc *PolicyService) ListEntitlements(ctx *gin.Context, input *models.PolicyEntitlementListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-policy-entilement")
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
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
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "listing-policy-entitlement")
	svc.logger.GetLogger().Info("listing policy entitlements")
	entitlements, total, err := svc.repo.SelectPolicyEntitlements(ctx, policy.ID, input.QueryCommonParam)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	outputs := make([]models.PolicyEntitlementListOutput, 0)
	for _, entitlement := range entitlements {
		outputs = append(outputs, models.PolicyEntitlementListOutput{
			ID:            entitlement.ID.String(),
			TenantName:    entitlement.TenantName,
			PolicyID:      entitlement.PolicyID.String(),
			EntitlementID: entitlement.EntitlementID.String(),
			Metadata:      entitlement.Metadata,
			CreatedAt:     entitlement.CreatedAt,
			UpdatedAt:     entitlement.UpdatedAt,
		})
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = total
	resp.Data = outputs
	return resp, nil
}
