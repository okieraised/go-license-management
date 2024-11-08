package product_attribute

type ProductAttribute struct {
	Name                 *string                `json:"name" validate:"required" example:"test"`
	Code                 *string                `json:"code" validate:"required" example:"test"`
	DistributionStrategy *string                `json:"distribution_strategy" validate:"optional" example:"test"`
	Url                  *string                `json:"url" validate:"optional" example:"test"`
	Permissions          []string               `json:"permissions" validate:"optional" example:"test"`
	Platforms            []string               `json:"platforms" validate:"optional" example:"test"`
	Metadata             map[string]interface{} `json:"metadata" validate:"optional" example:"test"`
}
