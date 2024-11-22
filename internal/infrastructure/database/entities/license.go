package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type License struct {
	bun.BaseModel `bun:"table:licenses,alias:l" swaggerignore:"true"`

	ID                          uuid.UUID              `bun:"id,pk,type:uuid"`
	PolicyID                    uuid.UUID              `bun:"policy_id,type:uuid,notnull"`
	ProductID                   uuid.UUID              `bun:"product_id,type:uuid,notnull"`
	TenantName                  string                 `bun:"tenant_name,type:varchar(256),notnull"`
	Key                         string                 `bun:"key,type:varchar(1028),notnull"`
	Name                        string                 `bun:"name,type:varchar(256),notnull"`
	LastValidatedChecksum       string                 `bun:"last_validated_checksum,type:varchar(1028),notnull"`
	LastValidatedVersion        string                 `bun:"last_validated_version,type:varchar(1028),notnull"`
	Status                      string                 `bun:"status,type:varchar(64),notnull"`
	Suspended                   bool                   `bun:"suspended,default:false"`
	Uses                        int                    `bun:"uses,type:integer,default:0"`
	MachinesCount               int                    `bun:"machines_count,type:integer,default:0"`
	Users                       int                    `bun:"users,default:0,notnull"`
	MaxMachines                 int                    `bun:"max_machines"`
	MaxUses                     int                    `bun:"max_uses"`
	MaxUsers                    int                    `bun:"max_users"`
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
	Tenant                      *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
	Product                     *Product               `bun:"rel:belongs-to,join:product_id=id"`
	Policy                      *Policy                `bun:"rel:belongs-to,join:policy_id=id"`
}
