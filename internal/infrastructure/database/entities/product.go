package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p" swaggerignore:"true"`

	ID                   uuid.UUID              `bun:"id,pk,type:uuid"` // primary key with UUID type
	AccountID            uuid.UUID              `bun:"account_id,type:uuid,notnull"`
	EnvironmentID        uuid.UUID              `bun:"type:uuid,nullzero"`
	Name                 string                 `bun:"name,type:varchar(256)"`
	DistributionStrategy string                 `bun:"distribution_strategy,type:varchar(128)"`
	Code                 string                 `bun:"code,type:varchar(128)"`
	Platforms            map[string]interface{} `bun:"platform,type:jsonb"`
	Metadata             map[string]interface{} `bun:"metadata,type:jsonb"`
	URL                  string                 `bun:"url,type:varchar(1024)"`
	CreatedAt            time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt            time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
