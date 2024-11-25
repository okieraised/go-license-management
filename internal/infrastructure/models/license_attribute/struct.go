package license_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/utils"
)

type LicenseCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	LicenseID  *string `uri:"license_id"`
	LicenseKey *string `uri:"license_key"`
	Action     *string `uri:"action"`
}

func (req *LicenseCommonURI) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	if req.LicenseID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.LicenseID)); err != nil {
			return comerrors.ErrLicenseIDIsInvalid
		}
	}

	if req.Action != nil {
		if _, ok := constants.ValidLicenseActionMapper[utils.DerefPointer(req.Action)]; !ok {
			return comerrors.ErrLicenseActionIsInvalid
		}
	}

	return nil
}
