package accounts

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/account_attribute"
	"go-license-management/internal/server/v1/accounts/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type AccountRegistrationRequest struct {
	Username  *string                `json:"username" validate:"required" example:"test"`
	Password  *string                `json:"password" validate:"required" example:"test"`
	FirstName *string                `json:"first_name" validate:"optional" example:"test"`
	LastName  *string                `json:"lastName" validate:"optional" example:"test"`
	Email     *string                `json:"email" validate:"required" example:"test"`
	Role      *string                `json:"role" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *AccountRegistrationRequest) Validate() error {
	if req.Username == nil {
		return comerrors.ErrAccountUsernameIsEmpty
	}
	if req.Password == nil {
		return comerrors.ErrAccountPasswordIsEmpty
	}

	if req.Email == nil {
		return comerrors.ErrAccountEmailIsEmpty
	}
	if req.Role == nil {
		req.Role = utils.RefPointer(constants.RoleUser)
	} else {
		if _, ok := constants.ValidRoleMapper[utils.DerefPointer(req.Role)]; !ok {
			return comerrors.ErrAccountRoleIsInvalid
		}
	}
	if req.FirstName == nil {
		req.FirstName = utils.RefPointer("")
	}
	if req.LastName == nil {
		req.LastName = utils.RefPointer("")
	}

	return nil
}

func (req *AccountRegistrationRequest) ToAccountRegistrationInput(ctx context.Context, tracer trace.Tracer, accountURI account_attribute.AccountCommonURI) *models.AccountRegistrationInput {
	return &models.AccountRegistrationInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		AccountCommonURI: accountURI,
		Username:         req.Username,
		Password:         req.Password,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Role:             req.Role,
		Metadata:         req.Metadata,
	}
}

type AccountRetrievalRequest struct {
	account_attribute.AccountCommonURI
}

func (req *AccountRetrievalRequest) Validate() error {
	if req.Username == nil {
		return comerrors.ErrAccountUsernameIsEmpty
	}
	return req.AccountCommonURI.Validate()
}

func (req *AccountRetrievalRequest) ToAccountRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.AccountRetrievalInput {
	return &models.AccountRetrievalInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		AccountCommonURI: req.AccountCommonURI,
	}
}

type AccountDeletionRequest struct {
	account_attribute.AccountCommonURI
}

func (req *AccountDeletionRequest) Validate() error {
	if req.Username == nil {
		return comerrors.ErrAccountUsernameIsEmpty
	}
	return req.AccountCommonURI.Validate()
}

func (req *AccountDeletionRequest) ToAccountDeletionInput(ctx context.Context, tracer trace.Tracer) *models.AccountDeletionInput {
	return &models.AccountDeletionInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		AccountCommonURI: req.AccountCommonURI,
	}
}

type AccountUpdateRequest struct {
	Password  *string                `json:"password" validate:"required" example:"test"`
	FirstName *string                `json:"first_name" validate:"optional" example:"test"`
	LastName  *string                `json:"lastName" validate:"optional" example:"test"`
	Email     *string                `json:"email" validate:"required" example:"test"`
	Role      *string                `json:"role" validate:"required" example:"test"`
	Metadata  map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *AccountUpdateRequest) Validate() error {
	if req.Role == nil {
		req.Role = utils.RefPointer(constants.RoleUser)
	} else {
		if _, ok := constants.ValidRoleMapper[utils.DerefPointer(req.Role)]; !ok {
			return comerrors.ErrAccountRoleIsInvalid
		}
	}
	return nil
}

func (req *AccountUpdateRequest) ToAccountUpdateInput(ctx context.Context, tracer trace.Tracer, tenantName, username string) *models.AccountUpdateInput {
	return &models.AccountUpdateInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
		Username:   utils.RefPointer(username),
		Password:   req.Password,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Role:       req.Role,
		Metadata:   req.Metadata,
	}
}

type AccountListRequest struct {
	constants.QueryCommonParam
}

func (req *AccountListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *AccountListRequest) ToAccountListInput(ctx context.Context, tracer trace.Tracer, accountURI account_attribute.AccountCommonURI) *models.AccountListInput {
	return &models.AccountListInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		AccountCommonURI: accountURI,
		QueryCommonParam: req.QueryCommonParam,
	}
}
