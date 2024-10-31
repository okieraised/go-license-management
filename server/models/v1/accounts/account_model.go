package accounts

import "errors"

type UserRegistrationRequest struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func (req *UserRegistrationRequest) Validate() error {
	if req.Username == nil {
		return errors.New("username is empty")
	}

	if req.Password == nil {
		return errors.New("password is empty")
	}

	return nil
}
