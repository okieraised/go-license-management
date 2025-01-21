package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"time"
)

func (svc *AccountService) actionBan(ctx *gin.Context, account *entities.Account) (*entities.Account, error) {
	svc.logger.GetLogger().Info(fmt.Sprintf("banning account [%s] in tenant [%s]", account.Username, account.Tenant))
	account.BannedAt = time.Now()
	account.Status = constants.AccountStatusBanned

	account, err := svc.repo.UpdateAccountByPK(ctx, account)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (svc *AccountService) actionUnban(ctx *gin.Context, account *entities.Account) (*entities.Account, error) {
	svc.logger.GetLogger().Info(fmt.Sprintf("unbanning account [%s] in tenant [%s]", account.Username, account.Tenant))
	account.BannedAt = time.Now()
	account.Status = constants.AccountStatusActive

	account, err := svc.repo.UpdateAccountByPK(ctx, account)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (svc *AccountService) actionGenerateResetToken(ctx *gin.Context, account *entities.Account) (*entities.Account, error) {
	svc.logger.GetLogger().Info(fmt.Sprintf("generate reset token for account [%s] in tenant [%s]", account.Username, account.TenantName))

	account.PasswordResetSentAt = time.Now()
	account.PasswordResetToken = utils.RandStringBytesMaskImprSrcSB(32)
	account, err := svc.repo.UpdateAccountByPK(ctx, account)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (svc *AccountService) actionResetPassword(ctx *gin.Context, token, newPass string, account *entities.Account) (*entities.Account, error) {
	svc.logger.GetLogger().Info(fmt.Sprintf("reset account [%s] in tenant [%s]", account.Username, account.TenantName))

	if token != account.PasswordResetToken {
		return account, cerrors.ErrAccountResetTokenIsInvalid
	}

	if time.Now().After(account.PasswordResetSentAt.Add(24 * time.Hour)) {
		return account, cerrors.ErrAccountResetTokenIsExpired
	}

	newHash, err := utils.HashPassword(newPass)
	if err != nil {
		return account, err
	}

	account.PasswordDigest = newHash
	account, err = svc.repo.UpdateAccountByPK(ctx, account)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (svc *AccountService) actionUpdatePassword(ctx *gin.Context, currentPass, newPass string, account *entities.Account) (*entities.Account, error) {
	svc.logger.GetLogger().Info(fmt.Sprintf("updating account [%s] in tenant [%s]", account.Username, account.TenantName))

	if !utils.CompareHashedPassword(account.PasswordDigest, currentPass) {
		return account, cerrors.ErrAccountPasswordNotMatch
	}

	newHash, err := utils.HashPassword(newPass)
	if err != nil {
		return account, err
	}

	account.PasswordDigest = newHash
	account, err = svc.repo.UpdateAccountByPK(ctx, account)
	if err != nil {
		return account, err
	}

	return account, nil

}
