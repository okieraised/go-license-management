package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenants,alias:tn" swaggerignore:"true"`

	ID                uuid.UUID `bun:"id,pk,type:uuid"`
	Name              string    `bun:"name,type:varchar(256),nullzero"`
	Slug              string    `bun:"slug,type:varchar(256),nullzero"`
	Protected         bool      `bun:"protected,default:true"`
	PublicKey         string    `bun:"public_key,nullzero"`
	PrivateKey        string    `bun:"private_key,nullzero"`
	SecretKey         string    `bun:"secret_key,nullzero"`
	Ed25519PrivateKey string    `bun:"ed25519_private_key,nullzero"`
	Ed25519PublicKey  string    `bun:"ed25519_public_key,nullzero"`
}
