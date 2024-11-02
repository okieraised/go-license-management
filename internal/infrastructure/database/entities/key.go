package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Key struct {
	bun.BaseModel `bun:"table:keys,alias:k"`

	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	Key           string    `bun:"key,nullzero"`
	PolicyID      uuid.UUID `bun:"policy_id,type:uuid"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid"`
	EnvironmentID uuid.UUID `bun:"environment_id,type:uuid"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}