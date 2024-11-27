package entities

import (
	"github.com/google/uuid"
	"time"
)

type Machine struct {
	ID                   uuid.UUID              `bun:"id,pk,type:uuid"`
	LicenseID            uuid.UUID              `bun:"license_id,type:uuid,notnull"`
	LicenseKey           string                 `bun:"license_key,type:varchar(256),notnull"`
	TenantName           string                 `bun:"tenant_name,type:varchar(256),notnull"`
	Fingerprint          string                 `bun:"fingerprint"`
	IP                   string                 `bun:"ip,type:varchar(64)"`
	Hostname             string                 `bun:"hostname,type:varchar(128)"`
	Platform             string                 `bun:"platform,type:varchar(128)"`
	Name                 string                 `bun:"name,type:varchar(128)"`
	Metadata             map[string]interface{} `bun:"metadata,type:jsonb"`
	Cores                int                    `bun:"cores,type:integer"`
	LastHeartbeatAt      time.Time              `bun:"last_heartbeat_at"`
	LastDeathEventSentAt time.Time              `bun:"last_death_event_sent_at"`
	LastCheckOutAt       time.Time              `bun:"last_check_out_at"`
	CreatedAt            time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant               *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
	License              *License               `bun:"rel:belongs-to,join:license_id=id"`
}
