package product_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/utils"
)

type ProductCommonURI struct {
	TenantName *string `uri:"tenant_name"`
	ProductID  *string `uri:"product_id"`
}

func (req *ProductCommonURI) Validate() error {
	if req.TenantName == nil {
		return cerrors.ErrTenantNameIsEmpty
	}

	if req.ProductID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.ProductID)); err != nil {
			return cerrors.ErrProductIDIsInvalid
		}
	}

	return nil
}

type ProductAttribute struct {
	Name                 *string                `json:"name" validate:"required" example:"test"`
	Code                 *string                `json:"code" validate:"required" example:"test"`
	DistributionStrategy *string                `json:"distribution_strategy" validate:"optional" example:"test"`
	Url                  *string                `json:"url" validate:"optional" example:"test"`
	Permissions          []string               `json:"permissions" validate:"optional" example:"test"`
	Platforms            []string               `json:"platforms" validate:"optional" example:"test"`
	Metadata             map[string]interface{} `json:"metadata" validate:"optional"`
}
