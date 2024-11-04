package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/accounts/models"
	"go-license-management/internal/server/v1/accounts/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type AccountService struct {
	repo   repository.IAccount
	logger *logging.Logger
}

func NewAccountService(options ...func(*AccountService)) *AccountService {
	svc := &AccountService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IAccount) func(*AccountService) {
	return func(c *AccountService) {
		c.repo = repo
	}
}

func (svc *AccountService) Create(ctx *gin.Context, input *models.AccountRegistrationInput) (*response.BaseOutput, error) {
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-account-by-name")
	exists, err := svc.repo.CheckAccountExistByPK(ctx, tenant.ID, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// If username/tenant combo already exists, return with error
	if exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("username [%s] already exists in tenant [%s]", utils.DerefPointer(input.Username), utils.DerefPointer(input.TenantName)))
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrAccountUsernameAlreadyExist]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrAccountUsernameAlreadyExist]
		return resp, comerrors.ErrAccountUsernameAlreadyExist
	}

	// Hashing password
	fmt.Println("utils.DerefPointer(input.Password)", utils.DerefPointer(input.Password))
	_, cSpan = input.Tracer.Start(rootCtx, "hash-password")
	hashed, err := utils.HashPassword(utils.DerefPointer(input.Password))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer]
		resp.Message = comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer]
		return resp, comerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// Insert new account
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-account")

	now := time.Now()
	account := &entities.Account{
		Username:           utils.DerefPointer(input.Username),
		TenantID:           tenant.ID,
		RoleName:           utils.DerefPointer(input.Role),
		Email:              utils.DerefPointer(input.Email),
		FirstName:          utils.DerefPointer(input.FirstName),
		LastName:           utils.DerefPointer(input.LastName),
		PasswordDigest:     hashed,
		PasswordResetToken: "",
		Metadata:           input.Metadata,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = svc.repo.InsertNewAccount(ctx, account)
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

func (svc *AccountService) List(ctx *gin.Context, input *models.AccountListInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *AccountService) Retrieve(ctx *gin.Context, input *models.AccountRetrievalInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *AccountService) Delete(ctx *gin.Context, input *models.AccountDeletionInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}

func (svc *AccountService) Update(ctx *gin.Context, input *models.AccountUpdateInput) (*response.BaseOutput, error) {

	return &response.BaseOutput{}, nil
}
