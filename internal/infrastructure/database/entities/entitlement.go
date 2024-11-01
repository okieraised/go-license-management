package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Entitlement struct {
	bun.BaseModel `bun:"table:entitlements,alias:e" swaggerignore:"true"`

	ID            uuid.UUID              `bun:"id,pk,type:uuid"`
	AccountID     uuid.UUID              `bun:"account_id,type:uuid,notnull"`
	EnvironmentID uuid.UUID              `bun:"environment_id,type:uuid,nullzero"`
	Name          string                 `bun:"name,type:varchar(256),notnull"`
	Code          string                 `bun:"code,type:varchar(256),notnull"`
	Metadata      map[string]interface{} `bun:"metadata,type:jsonb,nullzero"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
