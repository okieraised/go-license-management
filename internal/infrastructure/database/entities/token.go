package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductToken struct {
	bun.BaseModel `bun:"table:product_tokens,alias:pt"  swaggerignore:"true"`

	ID        uuid.UUID `bun:"id,pk,type:uuid"`
	TenantID  uuid.UUID `bun:"tenant_id,type:uuid,notnull"`
	ProductID uuid.UUID `bun:"product_id,type:uuid,notnull"`
	Token     string    `bun:"token,type:uuid"`
	Product   *Product  `bun:"rel:belongs-to,join:product_id=id,join:tenant_id=tenant_id"`
}

type LicenseToken struct {
	bun.BaseModel `bun:"table:product_tokens"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	TenantID      uuid.UUID `bun:"tenant_id,type:uuid,notnull"`
	LicenseID     uuid.UUID `bun:"license_id,type:uuid,notnull"`
	Token         string    `bun:"token,type:uuid"`
	License       *License  `bun:"rel:belongs-to,join:license_id=id"`
}
