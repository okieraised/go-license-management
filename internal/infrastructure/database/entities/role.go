package entities

import (
	"github.com/uptrace/bun"
	"time"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r" swaggerignore:"true"`
	Name          string    `bun:"name,pk,type:varchar(256),notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
