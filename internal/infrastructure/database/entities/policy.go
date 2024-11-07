package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Policy struct {
	bun.BaseModel `bun:"table:policies,alias:p" swaggerignore:"true"`

	ID                            uuid.UUID              `bun:"id,pk,type:uuid"`
	TenantID                      uuid.UUID              `bun:"tenant_id,type:uuid,notnull"`
	ProductID                     uuid.UUID              `bun:"product_id,type:uuid"`
	Duration                      int64                  `bun:"duration,nullzero"`
	LockVersion                   int                    `bun:"lock_version,nullzero"`
	MaxMachines                   int                    `bun:"max_machines,nullzero"`
	CheckInIntervalCount          int                    `bun:"check_in_interval_count,nullzero"`
	MaxUses                       int                    `bun:"max_uses,nullzero"`
	MaxProcesses                  int                    `bun:"max_processes,nullzero"`
	HeartbeatDuration             int                    `bun:"heartbeat_duration,nullzero"`
	MaxCores                      int                    `bun:"max_cores,nullzero"`
	MaxUsers                      int                    `bun:"max_users,nullzero"`
	Strict                        bool                   `bun:"strict,default:false"`
	Floating                      bool                   `bun:"floating,default:false"`
	UsePool                       bool                   `bun:"use_pool,default:false"`
	Encrypted                     bool                   `bun:"encrypted,default:false"`
	Protected                     bool                   `bun:"protected,nullzero"`
	RequireCheckIn                bool                   `bun:"require_check_in,default:false"`
	RequireProductScope           bool                   `bun:"require_product_scope,default:false"`
	RequirePolicyScope            bool                   `bun:"require_policy_scope,default:false"`
	RequireMachineScope           bool                   `bun:"require_machine_scope,default:false"`
	RequireFingerprintScope       bool                   `bun:"require_fingerprint_scope,default:false"`
	Concurrent                    bool                   `bun:"concurrent,default:true"`
	RequireHeartbeat              bool                   `bun:"require_heartbeat,default:false,notnull"`
	RequireChecksumScope          bool                   `bun:"require_checksum_scope,default:false,notnull"`
	RequireVersionScope           bool                   `bun:"require_version_scope,default:false,notnull"`
	RequireComponentsScope        bool                   `bun:"require_components_scope,default:false,notnull"`
	RequireAccountScope           bool                   `bun:"require_account_scope,default:false,notnull"`
	PublicKey                     string                 `bun:"public_key,type:varchar(1024),notnull"`
	PrivateKey                    string                 `bun:"private_key,type:varchar(1024),notnull"`
	Name                          string                 `bun:"name,type:varchar(256),nullzero"`
	Scheme                        string                 `bun:"scheme,type:varchar(128),nullzero"`
	FingerprintUniquenessStrategy string                 `bun:"fingerprint_uniqueness_strategy,type:varchar(64),nullzero"`
	FingerprintMatchingStrategy   string                 `bun:"fingerprint_matching_strategy,type:varchar(64),nullzero"`
	LeasingStrategy               string                 `bun:"leasing_strategy,type:varchar(64),nullzero"`
	ExpirationStrategy            string                 `bun:"expiration_strategy,type:varchar(64),nullzero"`
	ExpirationBasis               string                 `bun:"expiration_basis,type:varchar(64),nullzero"`
	AuthenticationStrategy        string                 `bun:"authentication_strategy,type:varchar(64),nullzero"`
	HeartbeatCullStrategy         string                 `bun:"heartbeat_cull_strategy,type:varchar(64),nullzero"`
	HeartbeatResurrectionStrategy string                 `bun:"heartbeat_resurrection_strategy,type:varchar(64),nullzero"`
	CheckInInterval               string                 `bun:"check_in_interval,type:varchar(64),nullzero"`
	TransferStrategy              string                 `bun:"transfer_strategy,type:varchar(64),nullzero"`
	OverageStrategy               string                 `bun:"overage_strategy,type:varchar(64),nullzero"`
	HeartbeatBasis                string                 `bun:"heartbeat_basis,type:varchar(64),nullzero"`
	MachineUniquenessStrategy     string                 `bun:"machine_uniqueness_strategy,type:varchar(64),nullzero"`
	MachineMatchingStrategy       string                 `bun:"machine_matching_strategy,type:varchar(64),nullzero"`
	ComponentUniquenessStrategy   string                 `bun:"component_uniqueness_strategy,type:varchar(64),nullzero"`
	ComponentMatchingStrategy     string                 `bun:"component_matching_strategy,type:varchar(64),nullzero"`
	RenewalBasis                  string                 `bun:"renewal_basis,type:varchar(64),nullzero"`
	MachineLeasingStrategy        string                 `bun:"machine_leasing_strategy,type:varchar(64),nullzero"`
	ProcessLeasingStrategy        string                 `bun:"process_leasing_strategy,type:varchar(64),nullzero"`
	Metadata                      map[string]interface{} `bun:"type:jsonb,nullzero"`
	CreatedAt                     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant                        *Tenant                `bun:"rel:belongs-to,join:tenant_id=id"`
	Product                       *Product               `bun:"rel:belongs-to,join:product_id=id,join:tenant_id=tenant_id"`
}
