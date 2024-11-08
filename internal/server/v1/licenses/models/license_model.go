package models

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type LicenseRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string                `json:"tenant_name" validate:"required" example:"test"`
	PolicyID   uuid.UUID              `json:"policy_id" validate:"required" example:"test"`
	ProductID  uuid.UUID              `json:"product_id" validate:"required" example:"test"`
	Name       *string                `json:"name" validate:"required" example:"test"`
	Expiry     *string                `json:"expiry" validate:"optional" example:"test"`
	Protected  *bool                  `json:"protected" validate:"optional" example:"test"`
	Suspended  *bool                  `json:"suspended" validate:"optional" example:"test"`
	Metadata   map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type LicenseRegistrationOutput struct {
	ID                      string `json:"id"`
	TenantID                string `json:"tenant_id"`
	ProductID               string `json:"product_id"`
	PolicyID                string `json:"policy_id"`
	Name                    string `json:"name"`
	Key                     string `json:"key"`
	MD5                     string `json:"md5"`
	Sha1                    string `json:"sha1"`
	Sha256                  string `json:"sha256"`
	PublicKey               string `json:"public_key"`
	Scheme                  string `json:"scheme"`
	RequireCheckIn          bool   `json:"require_check_in"`
	RequireProductScope     bool   `json:"require_product_scope"`
	RequirePolicyScope      bool   `json:"require_policy_scope"`
	RequireMachineScope     bool   `json:"require_machine_scope"`
	RequireFingerprintScope bool   `json:"require_fingerprint_scope"`
	Concurrent              bool   `json:"concurrent"`
	RequireHeartbeat        bool   `json:"require_heartbeat"`
	RequireChecksumScope    bool   `json:"require_checksum_scope"`
	RequireVersionScope     bool   `json:"require_version_scope"`
	RequireComponentsScope  bool   `json:"require_components_scope"`
	RequireUserScope        bool   `json:"require_user_scope"`
}

type LicenseListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	LicenseID  *string `json:"license_id" validate:"required" example:"test"`
}

type LicenseRetrievalInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseDeletionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type LicenseActionInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}
