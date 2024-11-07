package license

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type LicenseRegistrationRequest struct {
	PolicyID     *string                `json:"policy_id" validate:"required" example:"test"`
	ProductID    *string                `json:"product_id" validate:"required" example:"test"`
	Name         *string                `json:"name" validate:"required" example:"test"`
	Expiry       *string                `json:"expiry" validate:"optional" example:"test"`
	MaxMachine   *int                   `json:"max_machine" validate:"optional" example:"test"`
	MaxProcesses *int                   `json:"max_processes" validate:"optional" example:"test"`
	MaxUsers     *int                   `json:"max_users" validate:"optional" example:"test"`
	MaxUses      *int                   `json:"max_uses" validate:"optional" example:"test"`
	MaxCores     *int                   `json:"max_cores" validate:"optional" example:"test"`
	Protected    *bool                  `json:"protected" validate:"optional" example:"test"`
	Suspended    *bool                  `json:"suspended" validate:"optional" example:"test"`
	Permissions  []string               `json:"permissions" validate:"optional" example:"test"`
	Metadata     map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

func (req *LicenseRegistrationRequest) Validate() error {
	if req.ProductID == nil {
		return comerrors.ErrLicenseProductIDIsEmpty
	} else {
		_, err := uuid.Parse(utils.DerefPointer(req.ProductID))
		if err != nil {
			return comerrors.ErrProductIDIsInvalid
		}
	}

	if req.PolicyID == nil {
		return comerrors.ErrLicensePolicyIDIsEmpty
	} else {
		_, err := uuid.Parse(utils.DerefPointer(req.PolicyID))
		if err != nil {
			return comerrors.ErrPolicyIDIsInvalid
		}
	}
	if req.Name == nil {
		return comerrors.ErrLicenseNameIsEmpty
	}

	if req.Protected == nil {
		req.Protected = utils.RefPointer(true)
	}
	if req.Suspended == nil {
		req.Suspended = utils.RefPointer(false)
	}

	if req.Expiry != nil {
		if req.Expiry != nil {
			_, err := time.Parse(constants.DateFormatISO8601Hyphen, utils.DerefPointer(req.Expiry))
			if err != nil {
				return comerrors.ErrLicenseExpiryFormatIsInvalid
			}
		}
	}

	if req.MaxMachine == nil {
		req.MaxMachine = utils.RefPointer(0)
	}
	if req.MaxProcesses == nil {
		req.MaxProcesses = utils.RefPointer(0)
	}
	if req.MaxUsers == nil {
		req.MaxUsers = utils.RefPointer(0)
	}
	if req.MaxUses == nil {
		req.MaxUses = utils.RefPointer(0)
	}
	if req.MaxCores == nil {
		req.MaxCores = utils.RefPointer(0)
	}

	return nil
}

func (req *LicenseRegistrationRequest) ToLicenseRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.LicenseRegistrationInput {

	return &models.LicenseRegistrationInput{
		TracerCtx:    ctx,
		Tracer:       tracer,
		TenantName:   utils.RefPointer(tenantName),
		PolicyID:     uuid.MustParse(utils.DerefPointer(req.PolicyID)),
		ProductID:    uuid.MustParse(utils.DerefPointer(req.ProductID)),
		Name:         req.Name,
		Expiry:       req.Expiry,
		MaxMachine:   req.MaxMachine,
		MaxProcesses: req.MaxProcesses,
		MaxUsers:     req.MaxUsers,
		MaxUses:      req.MaxUses,
		MaxCores:     req.MaxCores,
		Protected:    req.Protected,
		Suspended:    req.Suspended,
		Permissions:  req.Permissions,
		Metadata:     req.Metadata,
	}
}

type LicenseRetrievalRequest struct {
	TenantName *string `uri:"tenant_name" validate:"required" example:"test"`
	LicenseID  *string `uri:"license_id" validate:"required" example:"test"`
}

func (req *LicenseRetrievalRequest) Validate() error {

	return nil
}
