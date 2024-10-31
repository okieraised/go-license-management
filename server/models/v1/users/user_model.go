package users

import "errors"

type UserRegistrationRequest struct {
	Username *string `json:"username,omitempty" validate:"required" example:"test"`
	Password *string `json:"password,omitempty" validate:"required" example:"test"`
	Email    *string `json:"email,omitempty" validate:"required" example:"test"`
	Role     *string `json:"role,omitempty" validate:"required" example:"test"`
}

func (req *UserRegistrationRequest) Validate() error {
	if req.Username == nil {
		return errors.New("username is empty")
	}

	if req.Password == nil {
		return errors.New("password is empty")
	}

	if req.Email == nil {
		return errors.New("email is empty")
	}

	if req.Role == nil {
		return errors.New("user role is empty")
	}

	return nil
}
