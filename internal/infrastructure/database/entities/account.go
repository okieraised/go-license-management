package entities

import (
	"github.com/uptrace/bun"
	"time"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:a" swaggerignore:"true"`

	Username            string                 `bun:"username,pk,type:varchar(128)"`
	TenantName          string                 `bun:"tenant_name,pk,type:varchar(256),notnull"`
	RoleName            string                 `bun:"role_name,type:varchar(256),notnull"`
	Email               string                 `bun:"email,type:varchar(256),notnull"`
	FirstName           string                 `bun:"first_name,type:varchar(128)"`
	LastName            string                 `bun:"last_name,type:varchar(128)"`
	Status              string                 `bun:"status,type:varchar(32),notnull"`
	PasswordDigest      string                 `bun:"password_digest,type:varchar(256)"`
	PasswordResetToken  string                 `bun:"password_reset_token,type:varchar(256),notnull"`
	Metadata            map[string]interface{} `bun:"metadata,type:jsonb"`
	PasswordResetSentAt time.Time              `bun:"password_reset_sent_at,nullzero"`
	BannedAt            time.Time              `bun:"banned_at,nullzero"`
	CreatedAt           time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt           time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Tenant              *Tenant                `bun:"rel:belongs-to,join:tenant_name=name"`
	Role                *Role                  `bun:"rel:belongs-to,join:role_name=name"`
}
