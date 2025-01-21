package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/authentications/models"
	"go-license-management/internal/services/v1/authentications/repository"
	"go-license-management/internal/utils"
	"go.uber.org/zap"
)

type AuthenticationService struct {
	repo   repository.IAuthentication
	logger *logging.Logger
}

func NewAuthenticationService(options ...func(*AuthenticationService)) *AuthenticationService {
	svc := &AuthenticationService{}

	for _, opt := range options {
		opt(svc)
	}
	logger := logging.NewECSLogger()
	svc.logger = logger

	return svc
}

func WithRepository(repo repository.IAuthentication) func(*AuthenticationService) {
	return func(c *AuthenticationService) {
		c.repo = repo
	}
}

// Login handles the login logic.
func (svc *AuthenticationService) Login(ctx *gin.Context, input *models.AuthenticationLoginInput) (*response.BaseOutput, error) {
	rootCtx, span := input.Tracer.Start(input.TracerCtx, "create-handler")
	defer span.End()

	resp := &response.BaseOutput{}
	svc.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	var token string
	var exp int64

	// Login admin
	if utils.DerefPointer(input.Username) == config.SuperAdminUsername {
		_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
		master, err := svc.repo.SelectMasterByPK(ctx, utils.DerefPointer(input.Username))
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

		// compare hash
		_, cSpan = input.Tracer.Start(rootCtx, "compare-hash")
		match := utils.CompareHashedPassword(master.PasswordDigest, utils.DerefPointer(input.Password))
		if !match {
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericUnauthorized]
			return resp, cerrors.ErrGenericUnauthorized
		}
		cSpan.End()

		// generate jwt
		_, cSpan = input.Tracer.Start(rootCtx, "generate-master-token")
		token, exp, err = svc.generateSuperadminJWT(ctx, master)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		cSpan.End()

	} else {
		// Login user
		_, cSpan := input.Tracer.Start(rootCtx, "query-tenant-by-name")
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

		// compare hash
		_, cSpan = input.Tracer.Start(rootCtx, "compare-hash")
		match := utils.CompareHashedPassword(account.PasswordDigest, utils.DerefPointer(input.Password))
		if !match {
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericUnauthorized]
			return resp, cerrors.ErrGenericUnauthorized
		}
		cSpan.End()

		// If account is inactive of banned
		if account.Status == constants.AccountStatusInactive {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountIsInactive]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountIsInactive]
			return resp, cerrors.ErrAccountIsInactive
		}

		if account.Status == constants.AccountStatusBanned {
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrAccountIsBanned]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrAccountIsBanned]
			return resp, cerrors.ErrAccountIsBanned
		}

		// generate jwt
		_, cSpan = input.Tracer.Start(rootCtx, "generate-account-token")
		token, exp, err = svc.generateJWT(ctx, tenant, account)
		if err != nil {
			svc.logger.GetLogger().Error(err.Error())
			cSpan.End()
			resp.Code = cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer]
			resp.Message = cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer]
			return resp, cerrors.ErrGenericInternalServer
		}
		cSpan.End()
	}

	resp.Code = cerrors.ErrCodeMapper[nil]
	resp.Message = cerrors.ErrMessageMapper[nil]
	resp.Data = models.AuthenticationLoginOutput{
		Access:   token,
		ExpireAt: exp,
	}

	return resp, nil
}
