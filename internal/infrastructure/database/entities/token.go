package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type ProductToken struct {
	bun.BaseModel `bun:"table:product_tokens,alias:pt" swaggerignore:"true"`

	ID         uuid.UUID `bun:"id,pk,type:uuid"`
	ProductID  uuid.UUID `bun:"product_id,type:uuid,notnull"`
	TenantName string    `bun:"tenant_name,type:varchar(256),notnull"`
	Token      string    `bun:"token,type:varchar(128)"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	Product    *Product  `bun:"rel:belongs-to,join:product_id=id"`
}

type LicenseToken struct {
	bun.BaseModel `bun:"table:license_tokens,alias:lt" swaggerignore:"true"`

	ID         uuid.UUID `bun:"id,pk,type:uuid"`
	TenantName string    `bun:"tenant_name,type:varchar(256),notnull"`
	LicenseID  uuid.UUID `bun:"license_id,type:uuid,notnull"`
	Token      string    `bun:"token,type:varchar(128)"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	License    *License  `bun:"rel:belongs-to,join:license_id=id"`
}
