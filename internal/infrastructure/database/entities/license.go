package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type License struct {
	bun.BaseModel `bun:"table:licenses,alias:l" swaggerignore:"true"`

	ID                          uuid.UUID              `bun:"id,pk,type:uuid"`
	UserID                      uuid.UUID              `bun:"user_id,type:uuid,nullzero"`
	PolicyID                    uuid.UUID              `bun:"policy_id,type:uuid,notnull"`
	AccountID                   uuid.UUID              `bun:"account_id,type:uuid,notnull"`
	GroupID                     uuid.UUID              `bun:"group_id,type:uuid,nullzero"`
	ProductID                   uuid.UUID              `bun:"product_id,type:uuid,notnull"`
	EnvironmentID               uuid.UUID              `bun:"environment_id,type:uuid,nullzero"`
	Key                         string                 `bun:"key,type:varchar(1028),notnull"`
	Name                        string                 `bun:"name,type:varchar(256),notnull"`
	LastValidatedChecksum       string                 `bun:"last_validated_checksum,type:varchar(1028),notnull"`
	LastValidatedVersion        string                 `bun:"last_validated_version,type:varchar(1028),notnull"`
	Suspended                   bool                   `bun:"suspended,default:false"`
	Protected                   bool                   `bun:"protected,default:true"`
	Uses                        int                    `bun:"uses,type:integer,default:0"`
	MachinesCount               int                    `bun:"machines_count,type:integer,default:0"`
	MachinesCoreCount           int                    `bun:"machines_core_count,nullzero"`
	MaxMachinesOverride         int                    `bun:"max_machine_override,nullzero"`
	MaxCoresOverride            int                    `bun:"max_cores_override,nullzero"`
	MaxUsesOverride             int                    `bun:"max_uses_override,nullzero"`
	MaxProcessesOverride        int                    `bun:"max_processes_override,nullzero"`
	MaxUsersOverride            int                    `bun:"max_user_override,nullzero"`
	LicenseUsersCount           int                    `bun:"license_users_count,default:0,notnull"`
	Metadata                    map[string]interface{} `bun:"metadata,type:jsonb"`
	CreatedAt                   time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                   time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Expiry                      time.Time              `bun:"expiry,nullzero"`
	LastCheckInAt               time.Time              `bun:"last_checked_in_at,nullzero"`
	LastExpirationEventSentAt   time.Time              `bun:"last_expiration_event_sent_at,nullzero"`
	LastCheckInEventSentAt      time.Time              `bun:"last_checked_in_event_sent_at,nullzero"`
	LastExpiringSoonEventSentAt time.Time              `bun:"last_expiring_soon_event_sent_at,nullzero"`
	LastCheckInSoonEventSentAt  time.Time              `bun:"last_checked_in_soon_event_sent_at,nullzero"`
	LastCheckOutAt              time.Time              `bun:"last_checkout_at,nullzero"`
	LastValidatedAt             time.Time              `bun:"last_validated_at,nullzero"`
}

type LicenseUser struct {
	bun.BaseModel `bun:"table:license_users,alias:lu" swaggerignore:"true"`

	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull"`
	EnvironmentID uuid.UUID `bun:"environment_id,type:uuid,nullzero"`
	LicenseID     uuid.UUID `bun:"license_id,type:uuid,notnull"`
	UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	License       *License  `bun:"rel:belongs-to,join:license_id=id"`
}

type LicenseEntitlement struct {
	bun.BaseModel `bun:"table:license_entitlements,alias:le" swaggerignore:"true"`

	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull"`
	LicenseID     uuid.UUID `bun:"license_id,type:uuid,notnull"`
	EntitlementID uuid.UUID `bun:"entitlement_id,type:uuid,notnull"`
	EnvironmentID uuid.UUID `bun:"environment_id,type:uuid,nullzero"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	License       *License  `bun:"rel:belongs-to,join:license_id=id"`
}
