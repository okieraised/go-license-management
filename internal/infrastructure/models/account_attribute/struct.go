package account_attribute

import (
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/utils"
)

type AccountCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	Username   *string `uri:"username"`
	Action     *string `uri:"action"`
}

func (req *AccountCommonURI) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	if req.Action != nil {
		if _, ok := constants.ValidAccountActionMapper[utils.DerefPointer(req.Action)]; !ok {
			return comerrors.ErrAccountActionIsInvalid
		}
	}

	return nil
}
