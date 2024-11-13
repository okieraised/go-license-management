package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/models/license_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type LicenseRegistrationInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string                `json:"tenant_name" validate:"required" example:"test"`
	PolicyID   uuid.UUID              `json:"policy_id" validate:"required" example:"test"`
	ProductID  uuid.UUID              `json:"product_id" validate:"required" example:"test"`
	Name       *string                `json:"name" validate:"required" example:"test"`
	Expiry     *string                `json:"expiry" validate:"optional" example:"test"`
	Metadata   map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type LicenseRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
}

type LicenseInfoOutput struct {
	ID                            string                 `json:"id"`
	TenantID                      string                 `json:"tenant_id"`
	ProductID                     string                 `json:"product_id"`
	PolicyID                      string                 `json:"policy_id"`
	Name                          string                 `json:"name"`
	Key                           string                 `json:"key"`
	MD5                           string                 `json:"md5"`
	Sha1                          string                 `json:"sha1"`
	Sha256                        string                 `json:"sha256"`
	PublicKey                     string                 `json:"public_key"`
	Scheme                        string                 `json:"scheme"`
	ExpirationStrategy            string                 `json:"expiration_strategy"`
	ExpirationBasis               string                 `json:"expiration_basis"`
	AuthenticationStrategy        string                 `json:"authentication_strategy"`
	HeartbeatCullStrategy         string                 `json:"heartbeat_cull_strategy"`
	HeartbeatResurrectionStrategy string                 `json:"heartbeat_resurrection_strategy"`
	CheckInInterval               string                 `json:"check_in_interval"`
	TransferStrategy              string                 `json:"transfer_strategy"`
	OverageStrategy               string                 `json:"overage_strategy"`
	HeartbeatBasis                string                 `json:"heartbeat_basis"`
	RenewalBasis                  string                 `json:"renewal_basis"`
	RequireCheckIn                bool                   `json:"require_check_in"`
	Concurrent                    bool                   `json:"concurrent"`
	RequireHeartbeat              bool                   `json:"require_heartbeat"`
	Strict                        bool                   `json:"strict"`
	Floating                      bool                   `json:"floating"`
	UsePool                       bool                   `json:"use_pool"`
	RateLimited                   bool                   `json:"rate_limited"`
	Encrypted                     bool                   `json:"encrypted"`
	Protected                     bool                   `json:"protected"`
	Duration                      int64                  `json:"duration"`
	MaxMachines                   int                    `json:"max_machines"`
	MaxUses                       int                    `json:"max_uses"`
	MaxUsers                      int                    `json:"max_users"`
	HeartbeatDuration             int                    `json:"heartbeat_duration"`
	Metadata                      map[string]interface{} `json:"metadata"`
	CreatedAt                     time.Time              `json:"created_at"`
	UpdatedAt                     time.Time              `json:"updated_at"`
}

type LicenseListInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
	LicenseID  *string `json:"license_id" validate:"required" example:"test"`
}

type LicenseDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
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
