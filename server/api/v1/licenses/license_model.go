package licenses

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/license_attribute"
	"go-license-management/internal/services/v1/licenses/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type LicenseRegistrationRequest struct {
	PolicyID    *string                `json:"policy_id" validate:"required" example:"test"`
	ProductID   *string                `json:"product_id" validate:"required" example:"test"`
	Name        *string                `json:"name" validate:"required" example:"test"`
	MaxMachines *int                   `json:"max_machines" validate:"optional" example:"1"`
	MaxUsers    *int                   `json:"max_users" validate:"optional" example:"1"`
	MaxUses     *int                   `json:"max_uses" validate:"optional" example:"1"`
	Expiry      *string                `json:"expiry" validate:"optional" example:"test"`
	Metadata    map[string]interface{} `json:"metadata" validate:"optional"`
}

func (req *LicenseRegistrationRequest) Validate() error {
	if req.ProductID == nil {
		return cerrors.ErrLicenseProductIDIsEmpty
	} else {
		_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
		if err != nil {
			return cerrors.ErrProductIDIsInvalid
		}
	}

	if req.PolicyID == nil {
		return cerrors.ErrLicensePolicyIDIsEmpty
	} else {
		_, err := uuid.Parse(utils.DerefPointer(req.PolicyID))
		if err != nil {
			return cerrors.ErrPolicyIDIsInvalid
		}
	}
	if req.Name == nil {
		return cerrors.ErrLicenseNameIsEmpty
	}

	if req.Expiry != nil {
		exp, err := time.Parse(time.RFC3339, utils.DerefPointer(req.Expiry))
		if err != nil {
			return cerrors.ErrLicenseExpiryFormatIsInvalid
		}
		if exp.Before(time.Now()) {
			return cerrors.ErrLicenseExpireDateIsInvalid
		}
	}

	if req.MaxMachines != nil {
		if utils.DerefPointer(req.MaxMachines) <= 0 {
			return cerrors.ErrLicenseMaxMachinesIsInvalid
		}
	}

	if req.MaxUses != nil {
		if utils.DerefPointer(req.MaxUses) <= 0 {
			return cerrors.ErrLicenseMaxUsesIsInvalid
		}
	}

	if req.MaxUsers != nil {
		if utils.DerefPointer(req.MaxUsers) <= 0 {
			return cerrors.ErrLicenseMaxUsersIsInvalid
		}
	}

	return nil
}

func (req *LicenseRegistrationRequest) ToLicenseRegistrationInput(ctx context.Context, tracer trace.Tracer, licenseURI license_attribute.LicenseCommonURI) *models.LicenseRegistrationInput {

	return &models.LicenseRegistrationInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: licenseURI,
		PolicyID:         req.PolicyID,
		ProductID:        req.ProductID,
		Name:             req.Name,
		MaxMachines:      req.MaxMachines,
		MaxUsers:         req.MaxUsers,
		MaxUses:          req.MaxUses,
		Expiry:           req.Expiry,
		Metadata:         req.Metadata,
	}
}

type LicenseUpdateRequest struct {
	PolicyID    *string                `json:"policy_id" validate:"required" example:"test"`
	ProductID   *string                `json:"product_id" validate:"required" example:"test"`
	Name        *string                `json:"name" validate:"required" example:"test"`
	MaxMachines *int                   `json:"max_machines" validate:"optional" example:"1"`
	MaxUsers    *int                   `json:"max_users" validate:"optional" example:"1"`
	MaxUses     *int                   `json:"max_uses" validate:"optional" example:"1"`
	Expiry      *string                `json:"expiry" validate:"optional" example:"test"`
	Metadata    map[string]interface{} `json:"metadata" validate:"optional"`
}

func (req *LicenseUpdateRequest) Validate() error {
	if req.ProductID != nil {
		_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
		if err != nil {
			return cerrors.ErrProductIDIsInvalid
		}
	}

	if req.PolicyID != nil {
		_, err := uuid.Parse(utils.DerefPointer(req.PolicyID))
		if err != nil {
			return cerrors.ErrPolicyIDIsInvalid
		}
	}

	if req.Expiry != nil {
		exp, err := time.Parse(time.RFC3339, utils.DerefPointer(req.Expiry))
		if err != nil {
			return cerrors.ErrLicenseExpiryFormatIsInvalid
		}
		if exp.Before(time.Now()) {
			return cerrors.ErrLicenseExpireDateIsInvalid
		}
	}

	if req.MaxMachines != nil {
		if utils.DerefPointer(req.MaxMachines) <= 0 {
			return cerrors.ErrLicenseMaxMachinesIsInvalid
		}
	}

	if req.MaxUses != nil {
		if utils.DerefPointer(req.MaxUses) <= 0 {
			return cerrors.ErrLicenseMaxUsesIsInvalid
		}
	}

	if req.MaxUsers != nil {
		if utils.DerefPointer(req.MaxUsers) <= 0 {
			return cerrors.ErrLicenseMaxUsersIsInvalid
		}
	}

	return nil
}

func (req *LicenseUpdateRequest) ToLicenseUpdateInput(ctx context.Context, tracer trace.Tracer, licenseURI license_attribute.LicenseCommonURI) *models.LicenseUpdateInput {
	return &models.LicenseUpdateInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: licenseURI,
		PolicyID:         req.PolicyID,
		ProductID:        req.ProductID,
		Name:             req.Name,
		MaxMachines:      req.MaxMachines,
		MaxUsers:         req.MaxUsers,
		MaxUses:          req.MaxUses,
		Expiry:           req.Expiry,
		Metadata:         req.Metadata,
	}
}

type LicenseRetrievalRequest struct {
	license_attribute.LicenseCommonURI
}

func (req *LicenseRetrievalRequest) Validate() error {
	return req.LicenseCommonURI.Validate()
}

func (req *LicenseRetrievalRequest) ToLicenseRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.LicenseRetrievalInput {
	return &models.LicenseRetrievalInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: req.LicenseCommonURI,
	}
}

type LicenseListRequest struct {
	constants.QueryCommonParam
}

func (req *LicenseListRequest) Validate() error {
	req.QueryCommonParam.Validate()
	return nil
}

func (req *LicenseListRequest) ToLicenseListInput(ctx context.Context, tracer trace.Tracer, licenseURI license_attribute.LicenseCommonURI) *models.LicenseListInput {
	return &models.LicenseListInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: licenseURI,
		QueryCommonParam: req.QueryCommonParam,
	}
}

type LicenseDeletionRequest struct {
	license_attribute.LicenseCommonURI
}

func (req *LicenseDeletionRequest) Validate() error {
	return req.LicenseCommonURI.Validate()
}

func (req *LicenseDeletionRequest) ToLicenseDeletionInput(ctx context.Context, tracer trace.Tracer) *models.LicenseDeletionInput {
	return &models.LicenseDeletionInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: req.LicenseCommonURI,
	}
}

type LicenseActionsRequest struct {
	LicenseKey *string `json:"license_key"`
	Nonce      *int    `json:"nonce"`
	Increment  *int    `json:"increment"`
	Decrement  *int    `json:"decrement"`
}

func (req *LicenseActionsRequest) Validate() error {

	if req.LicenseKey == nil {
		return cerrors.ErrLicenseKeyIsEmpty
	}

	if req.Decrement != nil {
		if utils.DerefPointer(req.Decrement) <= 0 {
			return cerrors.ErrLicenseDecrementIsInvalid
		}
	} else {
		req.Decrement = utils.RefPointer(1)
	}

	if req.Increment != nil {
		if utils.DerefPointer(req.Increment) <= 0 {
		}
	} else {
		req.Decrement = utils.RefPointer(1)
	}

	return nil
}

func (req *LicenseActionsRequest) ToLicenseActionsInput(ctx context.Context, tracer trace.Tracer, licenseUri license_attribute.LicenseCommonURI) *models.LicenseActionInput {
	return &models.LicenseActionInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		LicenseCommonURI: licenseUri,
		LicenseKey:       req.LicenseKey,
		Nonce:            req.Nonce,
		Increment:        req.Increment,
		Decrement:        req.Decrement,
	}
}
