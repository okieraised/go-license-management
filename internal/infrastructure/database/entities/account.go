package entities

import (
	"github.com/uptrace/bun"
	"time"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts" swaggerignore:"true"`

	ID          string    `bun:",pk,type:varchar(128)" json:"id"`
	Name        string    `bun:"name,type:varchar(128)" json:"name"`
	Description string    `bun:"description,type:varchar(256)" json:"description"`
	CreatedAt   time.Time `bun:"created_at,type:timestamptz" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,type:timestamptz" json:"updated_at"`
}
