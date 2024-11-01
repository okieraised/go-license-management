package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Policy struct {
	bun.BaseModel `bun:"table:policies,alias:p" swaggerignore:"true"`

	ID                            uuid.UUID              `bun:"id,pk,type:uuid"`
	ProductID                     uuid.UUID              `bun:"product_id,type:uuid,notnull"`
	AccountID                     uuid.UUID              `bun:"account_id,type:uuid,notnull"`
	EnvironmentID                 uuid.UUID              `bun:"environment_id,type:uuid,nullzero"`
	Duration                      int64                  `bun:"duration,nullzero"`
	LockVersion                   int                    `bun:"lock_version,default:0,notnull"`
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
	RequireEnvironmentScope       bool                   `bun:"require_environment_scope,default:false,notnull"`
	RequireChecksumScope          bool                   `bun:"require_checksum_scope,default:false,notnull"`
	RequireVersionScope           bool                   `bun:"require_version_scope,default:false,notnull"`
	RequireComponentsScope        bool                   `bun:"require_components_scope,default:false,notnull"`
	RequireUserScope              bool                   `bun:"require_user_scope,default:false,notnull"`
	Name                          string                 `bun:"name,type:varchar(256),nullzero"`
	Scheme                        string                 `bun:"scheme,nullzero"`
	FingerprintUniquenessStrategy string                 `bun:"fingerprint_uniqueness_strategy,nullzero"`
	FingerprintMatchingStrategy   string                 `bun:"fingerprint_matching_strategy,nullzero"`
	LeasingStrategy               string                 `bun:"leasing_strategy,nullzero"`
	ExpirationStrategy            string                 `bun:"expiration_strategy,nullzero"`
	ExpirationBasis               string                 `bun:"expiration_basis,nullzero"`
	AuthenticationStrategy        string                 `bun:"authentication_strategy,nullzero"`
	HeartbeatCullStrategy         string                 `bun:"heartbeat_cull_strategy,nullzero"`
	HeartbeatResurrectionStrategy string                 `bun:"heartbeat_resurrection_strategy,nullzero"`
	CheckInInterval               string                 `bun:"check_in_interval,nullzero"`
	TransferStrategy              string                 `bun:"transfer_strategy,nullzero"`
	OverageStrategy               string                 `bun:"overage_strategy,nullzero"`
	HeartbeatBasis                string                 `bun:"heartbeat_basis,nullzero"`
	MachineUniquenessStrategy     string                 `bun:"machine_uniqueness_strategy,nullzero"`
	MachineMatchingStrategy       string                 `bun:"machine_matching_strategy,nullzero"`
	ComponentUniquenessStrategy   string                 `bun:"component_uniqueness_strategy,nullzero"`
	ComponentMatchingStrategy     string                 `bun:"component_matching_strategy,nullzero"`
	RenewalBasis                  string                 `bun:"renewal_basis,nullzero"`
	MachineLeasingStrategy        string                 `bun:"machine_leasing_strategy,nullzero"`
	ProcessLeasingStrategy        string                 `bun:"process_leasing_strategy,nullzero"`
	Metadata                      map[string]interface{} `bun:"type:jsonb,nullzero"`
	CreatedAt                     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
