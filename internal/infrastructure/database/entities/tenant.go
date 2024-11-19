package entities

import (
	"github.com/uptrace/bun"
	"time"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenants,alias:tn" swaggerignore:"true"`

	Name              string    `bun:"name,pk,type:varchar(256),notnull"`
	Ed25519PublicKey  string    `bun:"ed25519_public_key,type:varchar(512),notnull"`
	Ed25519PrivateKey string    `bun:"ed25519_private_key,type:varchar(512),notnull"`
	CreatedAt         time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt         time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
