package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Group struct {
	bun.BaseModel `bun:"table:groups,alias:g" swaggerignore:"true"`

	ID            uuid.UUID              `bun:"id,pk,type:uuid"`
	AccountID     uuid.UUID              `bun:"account_id,type:uuid,notnull"`
	EnvironmentID uuid.UUID              `bun:"environment_id,type:uuid,nullzero"`
	Name          string                 `bun:"name,type:varchar(256),nullzero"`
	MaxUsers      int                    `bun:"max_users,type:integer,nullzero"`
	MaxLicenses   int                    `bun:"max_licenses,type:integer,nullzero"`
	MaxMachines   int                    `bun:"max_machines,type:integer,nullzero"`
	Metadata      map[string]interface{} `bun:"metadata,type:jsonb,nullzero"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type GroupOwner struct {
	bun.BaseModel `bun:"table:group_owners,alias:go" swaggerignore:"true"`

	ID            uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull"`
	GroupID       uuid.UUID `bun:"group_id,type:uuid,notnull"`
	UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
	EnvironmentID uuid.UUID `bun:"environment_id,type:uuid,nullzero"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
