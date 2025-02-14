package models

import (
	"context"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/license_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type LicenseRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
	PolicyID    *string                `json:"policy_id" validate:"required" example:"test"`
	ProductID   *string                `json:"product_id" validate:"required" example:"test"`
	Name        *string                `json:"name" validate:"required" example:"test"`
	MaxMachines *int                   `json:"max_machines" validate:"optional" example:"test"`
	MaxUsers    *int                   `json:"max_users" validate:"optional" example:"test"`
	MaxUses     *int                   `json:"max_uses" validate:"optional" example:"test"`
	Expiry      *string                `json:"expiry" validate:"optional" example:"test"`
	Metadata    map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type LicenseUpdateInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
	PolicyID    *string                `json:"policy_id" validate:"required" example:"test"`
	ProductID   *string                `json:"product_id" validate:"required" example:"test"`
	Name        *string                `json:"name" validate:"required" example:"test"`
	MaxMachines *int                   `json:"max_machines" validate:"optional" example:"test"`
	MaxUsers    *int                   `json:"max_users" validate:"optional" example:"test"`
	MaxUses     *int                   `json:"max_uses" validate:"optional" example:"test"`
	Expiry      *string                `json:"expiry" validate:"optional" example:"test"`
	Metadata    map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}

type LicenseRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
}

type LicenseInfoOutput struct {
	LicenseID      string                 `json:"license_id"`
	ProductID      string                 `json:"product_id"`
	PolicyID       string                 `json:"policy_id"`
	Name           string                 `json:"name"`
	LicenseKey     string                 `json:"license_key"`
	MD5Checksum    string                 `json:"md5_checksum"`
	Sha1Checksum   string                 `json:"sha1_checksum"`
	Sha256Checksum string                 `json:"sha256_checksum"`
	Status         string                 `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
	Expiry         time.Time              `json:"expiry"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	LicensePolicy  LicensePolicyOutput    `json:"license_policy"`
	LicenseProduct LicenseProductOutput   `json:"license_product"`
}

type LicenseProductOutput struct {
	Name                 string                 `json:"name,type:varchar(256)"`
	DistributionStrategy string                 `json:"distribution_strategy,type:varchar(128)"`
	Code                 string                 `json:"code,type:varchar(128),unique"`
	URL                  string                 `json:"url,type:varchar(1024)"`
	Platforms            []string               `json:"platform,type:jsonb"`
	Metadata             map[string]interface{} `json:"metadata,type:jsonb"`
	CreatedAt            time.Time              `json:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `json:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type LicensePolicyOutput struct {
	PolicyPublicKey    string `json:"policy_public_key"`
	PolicyScheme       string `json:"policy_scheme"`
	ExpirationStrategy string `json:"expiration_strategy"`
	CheckInInterval    string `json:"check_in_interval"`
	OverageStrategy    string `json:"overage_strategy"`
	HeartbeatBasis     string `json:"heartbeat_basis"`
	RenewalBasis       string `json:"renewal_basis"`
	RequireCheckIn     bool   `json:"require_check_in"`
	Concurrent         bool   `json:"concurrent"`
	RequireHeartbeat   bool   `json:"require_heartbeat"`
	Strict             bool   `json:"strict"`
	Floating           bool   `json:"floating"`
	UsePool            bool   `json:"use_pool"`
	RateLimited        bool   `json:"rate_limited"`
	Encrypted          bool   `json:"encrypted"`
	Protected          bool   `json:"protected"`
	Duration           int64  `json:"duration"`
	MaxMachines        int    `json:"max_machines"`
	MaxUses            int    `json:"max_uses"`
	MaxUsers           int    `json:"max_users"`
	HeartbeatDuration  int    `json:"heartbeat_duration"`
}

type LicenseListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
	constants.QueryCommonParam
}

type LicenseDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
}

type LicenseActionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	license_attribute.LicenseCommonURI
	LicenseKey *string `json:"license_key"`
	Nonce      *int    `json:"nonce"`
	Increment  *int    `json:"increment"`
	Decrement  *int    `json:"decrement"`
}

type LicenseValidationOutput struct {
	Valid bool   `json:"valid"`
	Code  string `json:"code"`
}

type LicenseActionCheckoutOutput struct {
	Certificate string    `json:"certificate"`
	TTL         int       `json:"ttl"`
	ExpiryAt    time.Time `json:"expiry_at"`
	IssuedAt    time.Time `json:"issued_at"`
}
