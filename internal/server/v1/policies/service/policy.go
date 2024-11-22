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
	"go-license-management/internal/infrastructure/models/policy_attribute"
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

	// Check if productID exists
	_, cSpan = input.Tracer.Start(rootCtx, "check-product-id")
	productID := uuid.MustParse(utils.DerefPointer(input.ProductID))
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying product [%s]", productID))
	exists, err := svc.repo.CheckProductExistByID(ctx, productID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}

	if !exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("product id [%s] does not exist in tenant [%s]", productID.String(), tenant.Name))
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrProductIDIsInvalid]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrProductIDIsInvalid]
		return resp, comerrors.ErrProductIDIsInvalid
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
		svc.logger.GetLogger().Error(fmt.Sprintf("invalid supported sheme [%s]", scheme))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrPolicySchemeIsInvalid]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrPolicySchemeIsInvalid]
		return resp, comerrors.ErrPolicySchemeIsInvalid
	}

	// Insert new policy
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-policy")
	svc.logger.GetLogger().Info("inserting new policy to database")
	policyID := uuid.New()
	now := time.Now()
	policy := &entities.Policy{
		ID:                     policyID,
		ProductID:              productID,
		TenantName:             tenant.Name,
		PublicKey:              publicKey,
		PrivateKey:             privateKey,
		Name:                   utils.DerefPointer(input.Name),
		Scheme:                 scheme,
		ExpirationStrategy:     utils.DerefPointer(input.ExpirationStrategy),
		ExpirationBasis:        utils.DerefPointer(input.ExpirationBasis),
		AuthenticationStrategy: utils.DerefPointer(input.AuthenticationStrategy),
		CheckInInterval:        utils.DerefPointer(input.CheckInInterval),
		OverageStrategy:        utils.DerefPointer(input.OverageStrategy),
		HeartbeatBasis:         utils.DerefPointer(input.HeartbeatBasis),
		RenewalBasis:           utils.DerefPointer(input.RenewalBasis),
		Duration:               utils.DerefPointer(input.Duration),
		MaxMachines:            utils.DerefPointer(input.MaxMachines),
		MaxUses:                utils.DerefPointer(input.MaxUses),
		MaxUsers:               utils.DerefPointer(input.MaxUsers),
		HeartbeatDuration:      utils.DerefPointer(input.HeartbeatDuration),
		Strict:                 utils.DerefPointer(input.Strict),
		Floating:               utils.DerefPointer(input.Floating),
		UsePool:                utils.DerefPointer(input.UsePool),
		RateLimited:            utils.DerefPointer(input.RateLimited),
		Encrypted:              utils.DerefPointer(input.Encrypted),
		Protected:              utils.DerefPointer(input.Protected),
		RequireCheckIn:         utils.DerefPointer(input.RequireCheckIn),
		RequireHeartbeat:       utils.DerefPointer(input.RequireHeartbeat),
		Metadata:               input.Metadata,
		CreatedAt:              now,
		UpdatedAt:              now,
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

	respData := models.PolicyRetrievalOutput{
		ID:         policyID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:                   utils.RefPointer(policy.Name),
			Scheme:                 utils.RefPointer(policy.Scheme),
			Strict:                 utils.RefPointer(policy.Strict),
			RateLimited:            utils.RefPointer(policy.RateLimited),
			Floating:               utils.RefPointer(policy.Floating),
			UsePool:                utils.RefPointer(policy.UsePool),
			Encrypted:              utils.RefPointer(policy.Encrypted),
			Protected:              utils.RefPointer(policy.Protected),
			RequireCheckIn:         utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:       utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:            utils.RefPointer(policy.MaxMachines),
			MaxUsers:               utils.RefPointer(policy.MaxUsers),
			MaxUses:                utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:      utils.RefPointer(policy.HeartbeatDuration),
			Duration:               utils.RefPointer(policy.Duration),
			CheckInInterval:        utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:         utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy:     utils.RefPointer(policy.ExpirationStrategy),
			ExpirationBasis:        utils.RefPointer(policy.ExpirationBasis),
			RenewalBasis:           utils.RefPointer(policy.RenewalBasis),
			AuthenticationStrategy: utils.RefPointer(policy.AuthenticationStrategy),
			OverageStrategy:        utils.RefPointer(policy.OverageStrategy),
			Metadata:               policy.Metadata,
		},
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
			ID:         policy.ID.String(),
			TenantName: policy.TenantName,
			PublicKey:  policy.PublicKey,
			CreatedAt:  policy.CreatedAt,
			UpdatedAt:  policy.UpdatedAt,
			PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
				Name:                   utils.RefPointer(policy.Name),
				Scheme:                 utils.RefPointer(policy.Scheme),
				Strict:                 utils.RefPointer(policy.Strict),
				RateLimited:            utils.RefPointer(policy.RateLimited),
				Floating:               utils.RefPointer(policy.Floating),
				UsePool:                utils.RefPointer(policy.UsePool),
				Encrypted:              utils.RefPointer(policy.Encrypted),
				Protected:              utils.RefPointer(policy.Protected),
				RequireCheckIn:         utils.RefPointer(policy.RequireCheckIn),
				RequireHeartbeat:       utils.RefPointer(policy.RequireHeartbeat),
				MaxMachines:            utils.RefPointer(policy.MaxMachines),
				MaxUsers:               utils.RefPointer(policy.MaxUsers),
				MaxUses:                utils.RefPointer(policy.MaxUses),
				HeartbeatDuration:      utils.RefPointer(policy.HeartbeatDuration),
				Duration:               utils.RefPointer(policy.Duration),
				CheckInInterval:        utils.RefPointer(policy.CheckInInterval),
				HeartbeatBasis:         utils.RefPointer(policy.HeartbeatBasis),
				ExpirationStrategy:     utils.RefPointer(policy.ExpirationStrategy),
				ExpirationBasis:        utils.RefPointer(policy.ExpirationBasis),
				RenewalBasis:           utils.RefPointer(policy.RenewalBasis),
				AuthenticationStrategy: utils.RefPointer(policy.AuthenticationStrategy),
				OverageStrategy:        utils.RefPointer(policy.OverageStrategy),
				Metadata:               policy.Metadata,
			},
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
		ID:         policy.ID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:                   utils.RefPointer(policy.Name),
			Scheme:                 utils.RefPointer(policy.Scheme),
			Strict:                 utils.RefPointer(policy.Strict),
			RateLimited:            utils.RefPointer(policy.RateLimited),
			Floating:               utils.RefPointer(policy.Floating),
			UsePool:                utils.RefPointer(policy.UsePool),
			Encrypted:              utils.RefPointer(policy.Encrypted),
			Protected:              utils.RefPointer(policy.Protected),
			RequireCheckIn:         utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:       utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:            utils.RefPointer(policy.MaxMachines),
			MaxUsers:               utils.RefPointer(policy.MaxUsers),
			MaxUses:                utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:      utils.RefPointer(policy.HeartbeatDuration),
			Duration:               utils.RefPointer(policy.Duration),
			CheckInInterval:        utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:         utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy:     utils.RefPointer(policy.ExpirationStrategy),
			ExpirationBasis:        utils.RefPointer(policy.ExpirationBasis),
			RenewalBasis:           utils.RefPointer(policy.RenewalBasis),
			AuthenticationStrategy: utils.RefPointer(policy.AuthenticationStrategy),
			OverageStrategy:        utils.RefPointer(policy.OverageStrategy),
			Metadata:               policy.Metadata,
		},
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("deleting policy [%s] and associated licenses", utils.DerefPointer(input.PolicyID)))
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
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "update-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	// Check if tenant exists
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

	// Query policy
	_, cSpan = input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("verifying policy [%s]", utils.DerefPointer(input.PolicyID)))
	policy, err := svc.repo.SelectPolicyByPK(ctx, uuid.MustParse(utils.DerefPointer(input.PolicyID)))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
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

	// Update fields
	policy, err = svc.updatePolicyField(ctx, input, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}

	// Update existing policy
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-policy")
	svc.logger.GetLogger().Info("updating policy to database")
	policy.UpdatedAt = time.Now()
	err = svc.repo.UpdatePolicyByPK(ctx, policy)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := models.PolicyRetrievalOutput{
		ID:         policy.ID.String(),
		TenantName: policy.TenantName,
		PublicKey:  policy.PublicKey,
		CreatedAt:  policy.CreatedAt,
		UpdatedAt:  policy.UpdatedAt,
		PolicyAttributeModel: policy_attribute.PolicyAttributeModel{
			Name:                   utils.RefPointer(policy.Name),
			Scheme:                 utils.RefPointer(policy.Scheme),
			Strict:                 utils.RefPointer(policy.Strict),
			RateLimited:            utils.RefPointer(policy.RateLimited),
			Floating:               utils.RefPointer(policy.Floating),
			UsePool:                utils.RefPointer(policy.UsePool),
			Encrypted:              utils.RefPointer(policy.Encrypted),
			Protected:              utils.RefPointer(policy.Protected),
			RequireCheckIn:         utils.RefPointer(policy.RequireCheckIn),
			RequireHeartbeat:       utils.RefPointer(policy.RequireHeartbeat),
			MaxMachines:            utils.RefPointer(policy.MaxMachines),
			MaxUsers:               utils.RefPointer(policy.MaxUsers),
			MaxUses:                utils.RefPointer(policy.MaxUses),
			HeartbeatDuration:      utils.RefPointer(policy.HeartbeatDuration),
			Duration:               utils.RefPointer(policy.Duration),
			CheckInInterval:        utils.RefPointer(policy.CheckInInterval),
			HeartbeatBasis:         utils.RefPointer(policy.HeartbeatBasis),
			ExpirationStrategy:     utils.RefPointer(policy.ExpirationStrategy),
			ExpirationBasis:        utils.RefPointer(policy.ExpirationBasis),
			RenewalBasis:           utils.RefPointer(policy.RenewalBasis),
			AuthenticationStrategy: utils.RefPointer(policy.AuthenticationStrategy),
			OverageStrategy:        utils.RefPointer(policy.OverageStrategy),
			Metadata:               policy.Metadata,
		},
	}

	resp.Code = comerrors.ErrCodeMapper[nil]
	resp.Message = comerrors.ErrMessageMapper[nil]
	resp.Data = respData
	return resp, nil
}

func (svc *PolicyService) Attach(ctx *gin.Context, input *models.PolicyAttachmentInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-policy")
	svc.logger.GetLogger().Info(fmt.Sprintf("querying policy [%s]", utils.DerefPointer(input.PolicyID)))
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

func (svc *PolicyService) Detach(ctx *gin.Context, input *models.PolicyDetachmentInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) ListEntitlements(ctx *gin.Context, input *models.PolicyEntitlementListInput) (*response.BaseOutput, error) {
	return nil, nil
}
