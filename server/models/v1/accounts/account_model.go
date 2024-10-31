package accounts

import "errors"

type AccountCreateModelRequest struct {
	Name        *string `json:"name" validate:"required" example:"test"`
	Description *string `json:"description"  validate:"optional" example:"test"`
}

func (req *AccountCreateModelRequest) Validate() error {
	if req.Name == nil {
		return errors.New("account name is empty")
	}

	return nil
}
