package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Policy struct {
	bun.BaseModel `bun:"table:policies,alias:p" swaggerignore:"true"`

	ID                     uuid.UUID              `bun:"id,pk,type:uuid"`
	ProductID              uuid.UUID              `bun:"product_id,type:uuid"`
	TenantName             string                 `bun:"tenant_name,type:varchar(256),notnull"`
	PublicKey              string                 `bun:"public_key,type:varchar(4096),notnull"`
	PrivateKey             string                 `bun:"private_key,type:varchar(4096),notnull"`
	Name                   string                 `bun:"name,type:varchar(256),nullzero"`
	Scheme                 string                 `bun:"scheme,type:varchar(128),nullzero"`
	ExpirationStrategy     string                 `bun:"expiration_strategy,type:varchar(64),nullzero"`
	ExpirationBasis        string                 `bun:"expiration_basis,type:varchar(64),nullzero"`
	AuthenticationStrategy string                 `bun:"authentication_strategy,type:varchar(64),nullzero"`
	CheckInInterval        string                 `bun:"check_in_interval,type:varchar(64),nullzero"`
	OverageStrategy        string                 `bun:"overage_strategy,type:varchar(64),nullzero"`
	HeartbeatBasis         string                 `bun:"heartbeat_basis,type:varchar(64),nullzero"`
	RenewalBasis           string                 `bun:"renewal_basis,type:varchar(64),nullzero"`
	Duration               int64                  `bun:"duration,nullzero"`
	MaxMachines            int                    `bun:"max_machines,nullzero"`
	MaxUses                int                    `bun:"max_uses,nullzero"`
	MaxUsers               int                    `bun:"max_users,nullzero"`
	HeartbeatDuration      int                    `bun:"heartbeat_duration,nullzero"`
	Strict                 bool                   `bun:"strict,default:false"`
	Floating               bool                   `bun:"floating,default:false"`
	UsePool                bool                   `bun:"use_pool,default:false"`
	RateLimited            bool                   `bun:"rate_limited,default:false"`
	Encrypted              bool                   `bun:"encrypted,default:false"`
	Protected              bool                   `bun:"protected,default:false"`
	RequireCheckIn         bool                   `bun:"require_check_in,default:false"`
	RequireHeartbeat       bool                   `bun:"require_heartbeat,default:false,notnull"`
	Metadata               map[string]interface{} `bun:"type:jsonb,nullzero"`
	CreatedAt              time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt              time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant                 *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
	Product                *Product               `bun:"rel:belongs-to,join:product_id=id"`
}

type PolicyEntitlement struct {
	bun.BaseModel `bun:"table:policy_entitlements,alias:pe" swaggerignore:"true"`

	ID            uuid.UUID              `bun:"id,pk,type:uuid"`
	TenantName    string                 `bun:"tenant_name,type:varchar(256),notnull"`
	PolicyID      uuid.UUID              `bun:"policy_id,type:uuid,notnull"`
	EntitlementID uuid.UUID              `bun:"entitlement_id,type:uuid,notnull"`
	Metadata      map[string]interface{} `bun:"type:jsonb,nullzero"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant        *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
	Policy        *Policy                `bun:"rel:belongs-to,join:policy_id=id"`
	Entitlement   *Entitlement           `bun:"rel:belongs-to,join:entitlement_id=name"`
}
