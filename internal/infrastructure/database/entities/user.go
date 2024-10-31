package entities

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users" swaggerignore:"true"`

	ID        string    `bun:",pk,type:varchar(128)" json:"id"`
	Username  string    `bun:"username,type:varchar(256)" json:"username"`
	Email     string    `bun:"email,type:varchar(128)" json:"email"`
	Hash      string    `bun:"hash,type:text"  json:"hash"`
	CreatedAt time.Time `bun:",type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",type:timestamptz" json:"updated_at"`
}
