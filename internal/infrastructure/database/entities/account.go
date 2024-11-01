package entities

import (
	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts" swaggerignore:"true"`

	ID          string `bun:",pk,type:uuid" json:"id"`
	Name        string `bun:"name,type:varchar(128)" json:"name"`
	Description string `bun:"description,type:varchar(256)" json:"description"`
	CreatedAt   int64  `bun:"created_at,type:bigint" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at,type:bigint" json:"updated_at"`
}
