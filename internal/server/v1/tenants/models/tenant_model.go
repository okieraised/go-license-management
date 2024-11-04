package models

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type TenantRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
	Protected *bool   `json:"protected" validate:"optional" example:"true"`
}

type TenantListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}

type TenantRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}

type TenantRetrievalOutput struct {
	ID                             string    `json:"id"`
	Name                           string    `json:"name"`
	Protected                      bool      `json:"protected"`
	Ed25519PublicKey               string    `json:"ed25519_public_key"`
	LastLowActivityLifelineSentAt  time.Time `json:"last_low_activity_lifeline_sent_at"`
	LastTrialWillEndSentAt         time.Time `json:"last_trial_will_end_sent_at"`
	LastLicenseLimitExceededSentAt time.Time `json:"last_license_limit_exceeded_sent_at"`
	LastRequestLimitExceededSentAt time.Time `json:"last_request_limit_exceeded_sent_at"`
	LastPromptForReviewSentAt      time.Time `json:"last_prompt_for_review_sent_at"`
	CreatedAt                      time.Time `json:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at"`
}

type TenantDeletionInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	Name      *string `json:"name,omitempty" validate:"required" example:"test"`
}
