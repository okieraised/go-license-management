package entities

import (
	"github.com/google/uuid"
	"time"
)

type Machine struct {
	ID                   uuid.UUID              `bun:"id,pk,type:uuid"`
	AccountID            uuid.UUID              `bun:"account_id,notnull"`
	LicenseID            uuid.UUID              `bun:"license_id,notnull"`
	GroupID              uuid.UUID              `bun:"group_id"`
	EnvironmentID        uuid.UUID              `bun:"environment_id"`
	OwnerID              uuid.UUID              `bun:"owner_id"`
	Fingerprint          string                 `bun:"fingerprint"`
	IP                   string                 `bun:"ip,type:varchar(64)"`
	Hostname             string                 `bun:"hostname,type:varchar(128)"`
	Platform             string                 `bun:"platform,type:varchar(128)"`
	Name                 string                 `bun:"name,type:varchar(128)"`
	HeartbeatJID         string                 `bun:"heartbeat_jid,type:varchar(128)"`
	Metadata             map[string]interface{} `bun:"metadata,type:jsonb"`
	Cores                int                    `bun:"cores,type:integer"`
	MaxProcessesOverride int                    `bun:"max_processes_override"`
	LastHeartbeatAt      time.Time              `bun:"last_heartbeat_at"`
	LastDeathEventSentAt time.Time              `bun:"last_death_event_sent_at"`
	LastCheckOutAt       time.Time              `bun:"last_check_out_at"`
	CreatedAt            time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type MachineProcess struct {
	ID                   uuid.UUID              `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AccountID            uuid.UUID              `bun:"account_id,notnull"`
	MachineID            uuid.UUID              `bun:"machine_id,notnull"`
	EnvironmentID        uuid.UUID              `bun:"environment_id"`
	PID                  string                 `bun:"pid,notnull"`
	HeartbeatJID         string                 `bun:"heartbeat_jid"`
	Metadata             map[string]interface{} `bun:"metadata,type:jsonb"`
	LastHeartbeatAt      time.Time              `bun:"last_heartbeat_at,notnull"`
	LastDeathEventSentAt time.Time              `bun:"last_death_event_sent_at"`
	CreatedAt            time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type MachineComponent struct {
	ID            uuid.UUID              `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AccountID     uuid.UUID              `bun:"account_id,notnull"`
	MachineID     uuid.UUID              `bun:"machine_id,notnull"`
	EnvironmentID uuid.UUID              `bun:"environment_id"`
	Fingerprint   string                 `bun:"fingerprint,notnull"`
	Name          string                 `bun:"name,notnull"`
	Metadata      map[string]interface{} `bun:"metadata,type:jsonb"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
