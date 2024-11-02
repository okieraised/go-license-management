package accounts

import (
	"context"
	"errors"
	"go-license-management/internal/server/v1/accounts/models"
	"go.opentelemetry.io/otel/trace"
)

type AccountCreateModelRequest struct {
	Name        *string `json:"name" validate:"required" example:"test"`
	Description *string `json:"description"  validate:"optional" example:"test"`
	Slug        *string `json:"slug"`
	Protected   *bool   `json:"protected,default:true"`
}

func (req *AccountCreateModelRequest) Validate() error {
	if req.Name == nil {
		return errors.New("account name is empty")
	}

	if req.Description == nil {
		return errors.New("account description is empty")
	}

	return nil
}

func (req *AccountCreateModelRequest) ToAccountRegistrationInput(ctx context.Context, tracer trace.Tracer) *models.AccountRegistrationInput {
	return &models.AccountRegistrationInput{
		TracerCtx:   ctx,
		Tracer:      tracer,
		Name:        req.Name,
		Description: req.Description,
	}
}
