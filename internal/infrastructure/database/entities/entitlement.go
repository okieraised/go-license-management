package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Entitlement struct {
	bun.BaseModel `bun:"table:entitlements,alias:e" swaggerignore:"true"`

	ID         uuid.UUID              `bun:"id,pk,type:uuid"`
	TenantName string                 `bun:"tenant_name,type:varchar(256),notnull"`
	Name       string                 `bun:"name,type:varchar(256),notnull"`
	Code       string                 `bun:"code,type:varchar(256),unique,notnull"`
	Metadata   map[string]interface{} `bun:"metadata,type:jsonb,nullzero"`
	CreatedAt  time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt  time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant     *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
}
