package models

import (
	"context"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type PolicyRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
	policy_attribute.PolicyAttributeModel
	ProductID *string `json:"product_id" validate:"required" example:"test"`
}

type PolicyRegistrationOutput struct {
	ID                            string                 `json:"id"`
	TenantID                      string                 `json:"tenant_id"`
	ProductID                     string                 `json:"product_id"`
	Duration                      int64                  `json:"duration"`
	MaxMachines                   int                    `json:"max_machines"`
	MaxUses                       int                    `json:"max_uses"`
	MaxUsers                      int                    `json:"max_users"`
	CheckInIntervalCount          int                    `json:"check_in_interval_count"`
	HeartbeatDuration             int                    `json:"heartbeat_duration"`
	Strict                        bool                   `json:"strict"`
	Floating                      bool                   `json:"floating"`
	UsePool                       bool                   `json:"use_pool"`
	RateLimited                   bool                   `json:"rate_limited"`
	Encrypted                     bool                   `json:"encrypted"`
	Protected                     bool                   `json:"protected"`
	RequireCheckIn                bool                   `json:"require_check_in"`
	Concurrent                    bool                   `json:"concurrent"`
	RequireHeartbeat              bool                   `json:"require_heartbeat"`
	PublicKey                     string                 `json:"public_key"`
	Name                          string                 `json:"name"`
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
	Metadata                      map[string]interface{} `json:"metadata"`
	CreatedAt                     time.Time              `json:"created_at"`
	UpdatedAt                     time.Time              `json:"updated_at"`
}

type PolicyListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
	constants.QueryCommonParam
}

type PolicyRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}

type PolicyRetrievalOutput struct {
	ID                            string                 `json:"id"`
	TenantName                    string                 `json:"tenant_name"`
	ProductID                     string                 `json:"product_id"`
	PublicKey                     string                 `json:"public_key"`
	Name                          string                 `json:"name"`
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
	Metadata                      map[string]interface{} `json:"metadata"`
	Duration                      int64                  `json:"duration"`
	MaxMachines                   int                    `json:"max_machines"`
	MaxUses                       int                    `json:"max_uses"`
	MaxUsers                      int                    `json:"max_users"`
	CheckInIntervalCount          int                    `json:"check_in_interval_count"`
	HeartbeatDuration             int                    `json:"heartbeat_duration"`
	Strict                        bool                   `json:"strict"`
	Floating                      bool                   `json:"floating"`
	UsePool                       bool                   `json:"use_pool"`
	RateLimited                   bool                   `json:"rate_limited"`
	Encrypted                     bool                   `json:"encrypted"`
	Protected                     bool                   `json:"protected"`
	RequireCheckIn                bool                   `json:"require_check_in"`
	Concurrent                    bool                   `json:"concurrent"`
	RequireHeartbeat              bool                   `json:"require_heartbeat"`
	CreatedAt                     time.Time              `json:"created_at"`
	UpdatedAt                     time.Time              `json:"updated_at"`
}

type PolicyUpdateInput struct {
	TracerCtx  context.Context
	Tracer     trace.Tracer
	TenantName *string `json:"tenant_name" validate:"required" example:"test"`
}

type PolicyDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	policy_attribute.PolicyCommonURI
}
