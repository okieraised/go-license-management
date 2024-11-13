package license_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/utils"
)

type LicenseCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	LicenseID  *string `uri:"license_id"`
}

func (req *LicenseCommonURI) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	if req.LicenseID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.LicenseID)); err != nil {
			return comerrors.ErrLicenseIDIsEmpty
		}
	}

	return nil
}
