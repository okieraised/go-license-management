package service

import (
	"database/sql"
	"errors"
	"fmt"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/accounts/models"
	"go-license-management/internal/services/v1/accounts/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
	"time"
)

type AccountService struct {
	repo   repository.IAccount
	casbin *xormadapter.Adapter
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

func WithCasbinAdapter(casbinAdapter *xormadapter.Adapter) func(*AccountService) {
	return func(c *AccountService) {
		c.casbin = casbinAdapter
	}
}

// Create creates new user
func (svc *AccountService) Create(ctx *gin.Context, input *models.AccountRegistrationInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-account-by-name")
	exists, err := svc.repo.CheckAccountExistByPK(ctx, tenant.Name, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// If username/tenant combo already exists, return with error
	if exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("username [%s] already exists in tenant [%s]", utils.DerefPointer(input.Username), utils.DerefPointer(input.TenantName)))
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountUsernameAlreadyExist]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountUsernameAlreadyExist]
		return resp, cerrors.ErrAccountUsernameAlreadyExist
	}

	_, cSpan = input.Tracer.Start(rootCtx, "query-account-by-name")
	exists, err = svc.repo.CheckAccountEmailExistByPK(ctx, tenant.Name, utils.DerefPointer(input.Email))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// If email/tenant combo already exists, return with error
	if exists {
		svc.logger.GetLogger().Info(fmt.Sprintf("email [%s] has already been used in tenant [%s]", utils.DerefPointer(input.Username), utils.DerefPointer(input.TenantName)))
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountEmailAlreadyExist]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountEmailAlreadyExist]
		return resp, cerrors.ErrAccountEmailAlreadyExist
	}

	// Hashing password
	_, cSpan = input.Tracer.Start(rootCtx, "hash-password")
	hashed, err := utils.HashPassword(utils.DerefPointer(input.Password))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// Insert new account
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-account")
	now := time.Now()
	account := &entities.Account{
		Username:       utils.DerefPointer(input.Username),
		TenantName:     tenant.Name,
		Status:         constants.AccountStatusActive,
		RoleName:       utils.DerefPointer(input.Role),
		Email:          utils.DerefPointer(input.Email),
		FirstName:      utils.DerefPointer(input.FirstName),
		LastName:       utils.DerefPointer(input.LastName),
		PasswordDigest: hashed,
		Metadata:       input.Metadata,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = svc.repo.InsertNewAccount(ctx, account)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// Insert account to casbin
	_, cSpan = input.Tracer.Start(rootCtx, "insert-new-account-casbin")
	err = svc.casbin.AddPolicy("g", "g", []string{utils.DerefPointer(input.TenantName), utils.DerefPointer(input.Username), utils.DerefPointer(input.Role)})
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
	resp.Data = models.AccountRegistrationOutput{
		Username:  account.Username,
		RoleName:  account.RoleName,
		Email:     account.Email,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Status:    account.Status,
		Metadata:  account.Metadata,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	return resp, nil
}

func (svc *AccountService) List(ctx *gin.Context, input *models.AccountListInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	_, cSpan = input.Tracer.Start(rootCtx, "query-accounts")
	accounts, count, err := svc.repo.SelectAccountsByTenant(ctx, tenant.Name, input.QueryCommonParam)
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "convert-tenants-to-output")
	respData := make([]models.AccountRetrievalOutput, 0)
	for _, account := range accounts {
		respData = append(respData, models.AccountRetrievalOutput{
			Username:  account.Username,
			RoleName:  account.RoleName,
			Email:     account.Email,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Status:    account.Status,
			Metadata:  account.Metadata,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		})
	}
	cSpan.End()

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Count = count
	resp.Data = respData

	return resp, nil
}

func (svc *AccountService) Retrieve(ctx *gin.Context, input *models.AccountRetrievalInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	_, cSpan = input.Tracer.Start(rootCtx, "select-account")
	account, err := svc.repo.SelectAccountByPK(ctx, tenant.Name, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	respData := &models.AccountRetrievalOutput{
		Username:  account.Username,
		RoleName:  account.RoleName,
		Email:     account.Email,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Status:    account.Status,
		Metadata:  account.Metadata,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = respData

	return resp, nil
}

func (svc *AccountService) Delete(ctx *gin.Context, input *models.AccountDeletionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	_, cSpan = input.Tracer.Start(rootCtx, "delete-account")
	err = svc.repo.DeleteAccountByPK(ctx, tenant.Name, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
		resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
		return resp, cerrors.ErrGenericInternalServer
	}
	cSpan.End()

	// Remove user policy from casbin
	_, cSpan = input.Tracer.Start(rootCtx, "delete-account-casbin")
	err = svc.casbin.RemovePolicy("g", "g", []string{utils.DerefPointer(input.TenantName), utils.DerefPointer(input.Username)})
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

func (svc *AccountService) Update(ctx *gin.Context, input *models.AccountUpdateInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	// Query tenant
	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	// Query account
	_, cSpan = input.Tracer.Start(rootCtx, "query-account")
	account, err := svc.repo.SelectAccountByPK(ctx, tenant.Name, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountUsernameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountUsernameIsInvalid]
			return resp, cerrors.ErrAccountUsernameIsInvalid
		default:
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}

	}
	cSpan.End()

	// Update account
	if input.Password != nil {
		_, cSpan = input.Tracer.Start(rootCtx, "query-account")
		hashed, err := utils.HashPassword(utils.DerefPointer(input.Password))
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		cSpan.End()
		account.PasswordDigest = hashed
	}

	if input.Email != nil {
		account.Email = utils.DerefPointer(input.Email)
	}

	if input.FirstName != nil {
		account.FirstName = utils.DerefPointer(input.FirstName)
	}

	if input.LastName != nil {
		account.LastName = utils.DerefPointer(input.LastName)
	}

	if input.Metadata != nil {
		account.Metadata = input.Metadata
	}

	if input.Role != nil {
		if account.RoleName != utils.DerefPointer(input.Role) {
			svc.logger.GetLogger().Info("updating current account role")
			err = svc.casbin.UpdatePolicy(
				"g",
				"g",
				[]string{utils.DerefPointer(input.TenantName), utils.DerefPointer(input.Username), account.RoleName},
				[]string{utils.DerefPointer(input.TenantName), utils.DerefPointer(input.Username), utils.DerefPointer(input.Role)},
			)
			if err != nil {
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
			account.RoleName = utils.DerefPointer(input.Role)
		} else {
			svc.logger.GetLogger().Info("current account role is the same as the update role")
		}
	}

	_, cSpan = input.Tracer.Start(rootCtx, "update-account")
	account, err = svc.repo.UpdateAccountByPK(ctx, account)
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
	resp.Data = models.AccountUpdateOutput{
		Username:  account.Username,
		RoleName:  account.RoleName,
		Email:     account.Email,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Status:    account.Status,
		Metadata:  account.Metadata,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	return resp, nil
}

func (svc *AccountService) Action(ctx *gin.Context, input *models.AccountActionInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "list-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	)

	_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
	svc.logger.GetLogger().Info(fmt.Sprintf("checking existing tenant [%s]", utils.DerefPointer(input.TenantName)))
	tenant, err := svc.repo.SelectTenantByPK(ctx, utils.DerefPointer(input.TenantName))
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

	_, cSpan = input.Tracer.Start(rootCtx, "queryaccount-by-pk")
	account, err := svc.repo.SelectAccountByPK(ctx, tenant.Name, utils.DerefPointer(input.Username))
	if err != nil {
		svc.logger.GetLogger().Error(err.Error())
		cSpan.End()
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountUsernameIsInvalid]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountUsernameIsInvalid]
			return resp, cerrors.ErrAccountUsernameIsInvalid
		} else {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	}
	cSpan.End()

	_, cSpan = input.Tracer.Start(rootCtx, "account-action")
	var output *entities.Account
	switch utils.DerefPointer(input.Action) {
	case constants.AccountActionBan:
		output, err = svc.actionBan(ctx, account)
		if err != nil {
			cSpan.End()
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	case constants.AccountActionUnban:
		output, err = svc.actionUnban(ctx, account)
		if err != nil {
			cSpan.End()
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	case constants.AccountActionGenerateResetToken:
		output, err = svc.actionGenerateResetToken(ctx, account)
		if err != nil {
			cSpan.End()
			svc.logger.GetLogger().Error(err.Error())
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
	case constants.AccountActionResetPassword:
		output, err = svc.actionResetPassword(ctx, utils.DerefPointer(input.ResetToken), utils.DerefPointer(input.NewPassword), account)
		if err != nil {
			cSpan.End()
			svc.logger.GetLogger().Error(err.Error())
			switch {
			case errors.Is(err, cerrors.ErrAccountResetTokenIsInvalid),
				errors.Is(err, cerrors.ErrAccountResetTokenIsExpired):
				resp.Code = cerrors.ErrCodeMapper[err]
				resp.Message = cerrors.ErrMessageMapper[err]
				return resp, err
			default:
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		}
	case constants.AccountActionUpdatePassword:
		output, err = svc.actionUpdatePassword(ctx, utils.DerefPointer(input.CurrentPassword), utils.DerefPointer(input.NewPassword), account)
		if err != nil {
			cSpan.End()
			svc.logger.GetLogger().Error(err.Error())
			switch {
			case errors.Is(err, cerrors.ErrAccountPasswordNotMatch):
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountPasswordNotMatch]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountPasswordNotMatch]
				return resp, cerrors.ErrAccountPasswordNotMatch
			default:
				resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
				resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
				return resp, cerrors.ErrGenericInternalServer
			}
		}
	default:
	}
	cSpan.End()

	if utils.DerefPointer(input.Action) == constants.AccountActionGenerateResetToken {
		resp.Data = models.AccountPasswordTokenOutput{
			Token: output.PasswordResetToken,
		}
	} else {
		resp.Data = models.AccountUpdateOutput{
			Username:  output.Username,
			RoleName:  output.RoleName,
			Email:     output.Email,
			FirstName: output.FirstName,
			LastName:  output.LastName,
			Status:    output.Status,
			Metadata:  output.Metadata,
			CreatedAt: output.CreatedAt,
			UpdatedAt: output.UpdatedAt,
		}
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	return resp, nil
}
