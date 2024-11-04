package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenants,alias:tn" swaggerignore:"true"`

	ID                             uuid.UUID `bun:"id,pk,type:uuid"`
	Name                           string    `bun:"name,type:varchar(256),unique,notnull"`
	Protected                      bool      `bun:"protected,default:false"`
	Ed25519PublicKey               string    `bun:"ed25519_public_key,type:varchar(512),notnull"`
	Ed25519PrivateKey              string    `bun:"ed25519_private_key,type:varchar(512),notnull"`
	LastLowActivityLifelineSentAt  time.Time `bun:"last_low_activity_lifeline_sent_at,nullzero"`
	LastTrialWillEndSentAt         time.Time `bun:"last_trial_will_end_sent_at,nullzero"`
	LastLicenseLimitExceededSentAt time.Time `bun:"last_license_limit_exceeded_sent_at,nullzero"`
	LastRequestLimitExceededSentAt time.Time `bun:"last_request_limit_exceeded_sent_at,nullzero"`
	LastPromptForReviewSentAt      time.Time `bun:"last_prompt_for_review_sent_at,nullzero"`
	CreatedAt                      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
