package constants

const (
	LicenseActionValidate  = "validate"
	LicenseActionRevoke    = "revoke"
	LicenseActionSuspend   = "suspend"
	LicenseActionReinstate = "reinstate"
	LicenseActionRenew     = "renew"
	LicenseActionCheckout  = "checkout"
	LicenseActionCheckin   = "checkin"
)

var ValidLicenseActionMapper = map[string]interface{}{
	LicenseActionValidate:  true,
	LicenseActionRevoke:    true,
	LicenseActionSuspend:   true,
	LicenseActionReinstate: true,
	LicenseActionRenew:     true,
	LicenseActionCheckout:  true,
	LicenseActionCheckin:   true,
}

//The status of the license to filter by. One of: ACTIVE, INACTIVE, EXPIRED, SUSPENDED, or BANNED.

const (
	LicenseStatusNotActivated = "not_activated"
	LicenseStatusActive       = "active"
	LicenseStatusInactive     = "inactive"
	LicenseStatusSuspended    = "suspended"
	LicenseStatusExpired      = "expired"
	LicenseStatusBanned       = "banned"
)
