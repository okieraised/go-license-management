package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:a" swaggerignore:"true"`

	ID                             uuid.UUID `bun:"id,pk,type:uuid"`
	PlanID                         uuid.UUID `bun:"plan_id,type:uuid,nullzero"`
	Name                           string    `bun:"name,type:varchar(256),nullzero"`
	Slug                           string    `bun:"slug,type:varchar(256),nullzero"`
	Protected                      bool      `bun:"protected,default:false"`
	PublicKey                      string    `bun:"public_key,nullzero"`
	PrivateKey                     string    `bun:"private_key,nullzero"`
	SecretKey                      string    `bun:"secret_key,nullzero"`
	Ed25519PrivateKey              string    `bun:"ed25519_private_key,nullzero"`
	Ed25519PublicKey               string    `bun:"ed25519_public_key,nullzero"`
	Domain                         string    `bun:"domain,type:varchar(256),nullzero"`
	Subdomain                      string    `bun:"subdomain,type:varchar(256),nullzero"`
	CNAME                          string    `bun:"cname,type:varchar(256),nullzero"`
	Backend                        string    `bun:"backend,type:varchar(256),nullzero"`
	LastLowActivityLifelineSentAt  time.Time `bun:"last_low_activity_lifeline_sent_at,nullzero"`
	LastTrialWillEndSentAt         time.Time `bun:"last_trial_will_end_sent_at,nullzero"`
	LastLicenseLimitExceededSentAt time.Time `bun:"last_license_limit_exceeded_sent_at,nullzero"`
	LastRequestLimitExceededSentAt time.Time `bun:"last_request_limit_exceeded_sent_at,nullzero"`
	LastPromptForReviewSentAt      time.Time `bun:"last_prompt_for_review_sent_at,nullzero"`
	CreatedAt                      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
