package entitlement_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/utils"
)

type EntitlementCommonURI struct {
	TenantName    *string `uri:"tenant_name" validate:"required" example:"test"`
	EntitlementID *string `uri:"entitlement_id" validate:"required" example:"test"`
}

func (req *EntitlementCommonURI) Validate() error {
	if req.TenantName == nil {
		return cerrors.ErrTenantNameIsEmpty
	}

	if req.EntitlementID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.EntitlementID)); err != nil {
			return cerrors.ErrEntitlementIDIsInvalid
		}
	}
	return nil
}
