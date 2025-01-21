package license_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/utils"
)

type LicenseCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	LicenseID  *string `uri:"license_id"`
	Action     *string `uri:"action"`
}

func (req *LicenseCommonURI) Validate() error {
	if req.TenantName == nil {
		return cerrors.ErrTenantNameIsEmpty
	}

	if req.LicenseID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.LicenseID)); err != nil {
			return cerrors.ErrLicenseIDIsInvalid
		}
	}

	if req.Action != nil {
		if _, ok := constants.ValidLicenseActionMapper[utils.DerefPointer(req.Action)]; !ok {
			return cerrors.ErrLicenseActionIsInvalid
		}
	}

	return nil
}

// LicenseFileContent contains information about the license file
type LicenseFileContent struct {
	Enc string `json:"enc"`
	Sig string `json:"sig"`
	Alg string `json:"alg"`
}
