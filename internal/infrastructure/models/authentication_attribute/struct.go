package authentication_attribute

import (
	"go-license-management/internal/comerrors"
)

type AuthenticationCommonURI struct {
	TenantName *string `uri:"tenant_name"`
}

func (req *AuthenticationCommonURI) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	return nil
}
