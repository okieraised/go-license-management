package license

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/server/v1/licenses/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type LicenseRegistrationRequest struct {
	Name         *string                `json:"name" validate:"required" example:"test"`
	Key          *string                `json:"key" validate:"optional" example:"test"`
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

	if req.Name == nil {
		return comerrors.ErrLicenseNameIsEmpty
	}

	return nil
}

func (req *LicenseRegistrationRequest) ToLicenseRegistrationInput(ctx context.Context, tracer trace.Tracer, tenantName string) *models.LicenseRegistrationInput {

	return &models.LicenseRegistrationInput{
		TracerCtx:  ctx,
		Tracer:     tracer,
		TenantName: utils.RefPointer(tenantName),
	}
}

type LicenseRetrievalRequest struct {
	TenantName *string `uri:"tenant_name" validate:"required" example:"test"`
	LicenseID  *string `uri:"license_id" validate:"required" example:"test"`
}

func (req *LicenseRetrievalRequest) Validate() error {

	return nil
}
