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
