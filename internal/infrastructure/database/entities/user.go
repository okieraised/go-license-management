package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u" swaggerignore:"true"`

	ID                  uuid.UUID              `bun:",pk,type:uuid"`
	Email               string                 `bun:"email,type:varchar(256),notnull"`
	PasswordDigest      string                 `bun:"password_digest,type:varchar(256)"`
	CreatedAt           time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt           time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	PasswordResetToken  string                 `bun:"password_reset_token,type:varchar(256),notnull"`
	PasswordResetSentAt time.Time              `bun:"password_reset_sent_at,nullzero"`
	Metadata            map[string]interface{} `bun:"metadata,type:jsonb"`
	AccountID           uuid.UUID              `bun:"type:uuid,notnull"`
	FirstName           string                 `bun:"first_name,type:varchar(256)"`
	LastName            string                 `bun:"last_name,type:varchar(256)"`
	BannedAt            time.Time              `bun:"banned_at,nullzero"`
	GroupID             uuid.UUID              `bun:"group_id,type:uuid"`
}
