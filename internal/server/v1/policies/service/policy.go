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
	exists, err := svc.repo.CheckProductExistByID(ctx, tenant.ID, productID)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	if !exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("product id [%s] does not exist in tenant [%s]", productID.String(), tenant.ID.String()))
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
	case constants.PolicySchemeRSA2048PKCS1, constants.PolicySchemeRSA2048JWTRS256:
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
	entitlement := &entities.Policy{
		ID:                            policyID,
		TenantID:                      tenant.ID,
		ProductID:                     productID,
		Duration:                      0,
		LockVersion:                   0,
		MaxMachines:                   0,
		CheckInIntervalCount:          0,
		MaxUses:                       0,
		MaxProcesses:                  0,
		HeartbeatDuration:             0,
		MaxCores:                      0,
		MaxUsers:                      0,
		Strict:                        utils.DerefPointer(input.Strict),
		Floating:                      utils.DerefPointer(input.Floating),
		UsePool:                       utils.DerefPointer(input.UsePool),
		Encrypted:                     false,
		Protected:                     utils.DerefPointer(input.Protected),
		RequireCheckIn:                utils.DerefPointer(input.RequireCheckIn),
		RequireProductScope:           utils.DerefPointer(input.RequireProductScope),
		RequirePolicyScope:            utils.DerefPointer(input.RequirePolicyScope),
		RequireMachineScope:           utils.DerefPointer(input.RequireMachineScope),
		RequireFingerprintScope:       utils.DerefPointer(input.RequireFingerprintScope),
		Concurrent:                    false,
		RequireHeartbeat:              utils.DerefPointer(input.RequireHeartbeat),
		RequireChecksumScope:          false,
		RequireVersionScope:           utils.DerefPointer(input.RequireVersionScope),
		RequireComponentsScope:        utils.DerefPointer(input.RequireComponentsScope),
		RequireAccountScope:           false,
		PublicKey:                     publicKey,
		PrivateKey:                    privateKey,
		Name:                          utils.DerefPointer(input.Name),
		Scheme:                        scheme,
		FingerprintUniquenessStrategy: "",
		FingerprintMatchingStrategy:   "",
		LeasingStrategy:               "",
		ExpirationStrategy:            utils.DerefPointer(input.ExpirationStrategy),
		ExpirationBasis:               utils.DerefPointer(input.ExpirationBasis),
		AuthenticationStrategy:        utils.DerefPointer(input.AuthenticationStrategy),
		HeartbeatCullStrategy:         "",
		HeartbeatResurrectionStrategy: "",
		CheckInInterval:               "",
		TransferStrategy:              utils.DerefPointer(input.TransferStrategy),
		OverageStrategy:               utils.DerefPointer(input.OverageStrategy),
		HeartbeatBasis:                "",
		MachineUniquenessStrategy:     "",
		MachineMatchingStrategy:       "",
		ComponentUniquenessStrategy:   "",
		ComponentMatchingStrategy:     "",
		RenewalBasis:                  "",
		MachineLeasingStrategy:        "",
		ProcessLeasingStrategy:        "",
		Metadata:                      input.Metadata,
		CreatedAt:                     now,
		UpdatedAt:                     now,
	}
	err = svc.repo.InsertNewPolicy(ctx, entitlement)
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
		"policy_id": policyID.String(),
	}
	return resp, nil
}

func (svc *PolicyService) List(ctx *gin.Context, input *models.PolicyListInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Retrieve(ctx *gin.Context, input *models.PolicyRetrievalInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *PolicyService) Delete(ctx *gin.Context, input *models.PolicyDeletionInput) (*response.BaseOutput, error) {
	return nil, nil
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
