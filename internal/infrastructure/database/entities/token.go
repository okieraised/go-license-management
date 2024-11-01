package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Token struct {
	bun.BaseModel `bun:"table:tokens,alias:t" swaggerignore:"true"`

	ID               uuid.UUID `bun:"id,pk,type:uuid"`
	BearerID         uuid.UUID `bun:"bearer_id,type:uuid,nullzero"`
	AccountID        uuid.UUID `bun:"account_id,type:uuid,nullzero"`
	Name             string    `bun:"name,type:varchar(128)"`
	Digest           string    `bun:"digest,type:varchar(256),unique,notnull"`
	BearerType       string    `bun:"bearer_type,type:varchar(256)"`
	MaxActivations   int       `bun:"max_activations,type:integer"`
	MaxDeactivations int       `bun:"max_deactivations,type:integer"`
	Activations      int       `bun:"activations,default:0"`
	Deactivations    int       `bun:"deactivations,default:0"`
	CreatedAt        time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Expiry           time.Time `bun:"expiry,nullzero"`
	EnvironmentID    uuid.UUID `bun:"type:uuid,nullzero"`
}

type TokenPermission struct {
	bun.BaseModel `bun:"table:token_permissions,alias:tp" swaggerignore:"true"`

	ID           uuid.UUID `bun:"id,pk,type:uuid"` // primary key with UUID type
	PermissionID uuid.UUID `bun:"permission_id,type:uuid,notnull"`
	TokenID      uuid.UUID `bun:"token_id,type:uuid,notnull"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
