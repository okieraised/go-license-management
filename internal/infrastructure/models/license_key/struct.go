package license_key

import (
	"time"
)

type LicenseKeyContent struct {
	TenantName *string                `json:"tenant_name"`
	ProductID  *string                `json:"product_id"`
	PolicyID   *string                `json:"policy_id"`
	LicenseID  *string                `json:"license_id"`
	Metadata   map[string]interface{} `json:"metadata"`
	Expiry     time.Time              `json:"expiry"`
	CreatedAt  time.Time              `json:"created_at"`
}
