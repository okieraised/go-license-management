package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p" swaggerignore:"true"`

	ID                   uuid.UUID              `bun:"id,pk,type:uuid"`
	TenantID             uuid.UUID              `bun:"tenant_id,pk,type:uuid,notnull"`
	AccountName          string                 `bun:"account_name,type:varchar(128),notnull"`
	Name                 string                 `bun:"name,type:varchar(256)"`
	DistributionStrategy string                 `bun:"distribution_strategy,type:varchar(128)"`
	Code                 string                 `bun:"code,type:varchar(128),unique"`
	Platforms            map[string]interface{} `bun:"platform,type:jsonb"`
	Metadata             map[string]interface{} `bun:"metadata,type:jsonb"`
	URL                  string                 `bun:"url,type:varchar(1024)"`
	CreatedAt            time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant               *Tenant                `bun:"rel:belongs-to,join:tenant_id=id"`
	Account              *Account               `bun:"rel:belongs-to,join:account_name=username,join:tenant_id=tenant_id"`
}
