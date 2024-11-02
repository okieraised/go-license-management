package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Environment struct {
	bun.BaseModel `bun:"table:environments,alias:e"  swaggerignore:"true"`

	ID                uuid.UUID `bun:"id,pk,type:uuid"`
	AccountID         uuid.UUID `bun:"type:uuid,notnull"`
	Name              string    `bun:",notnull"` // Environment name, not nullable
	CreatedAt         time.Time `bun:",notnull"` // Not nullable timestamp
	UpdatedAt         time.Time `bun:",notnull"` // Not nullable timestamp
	Code              string    `bun:",notnull"` // Code, not nullable
	IsolationStrategy string    `bun:",notnull"` // Isolation strategy, not nullable
}
