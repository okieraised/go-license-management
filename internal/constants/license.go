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

const (
	LicenseValidationStatusValid               = "valid"                 // The validated license resource or license key is valid.
	LicenseValidationStatusSuspended           = "suspended"             // The validated license has been suspended.
	LicenseValidationStatusExpired             = "expired"               // The validated license is expired.
	LicenseValidationStatusBanned              = "banned"                // The user that owns the validated license has been banned.
	LicenseValidationStatusOverdue             = "overdue"               // The validated license is overdue for check-in.
	LicenseValidationStatusNoMachine           = "no_machine"            // Not activated. The validated license does not meet its node-locked policy's requirement of exactly 1 associated machine.
	LicenseValidationStatusNoMachines          = "no_machines"           // The validated license does not meet its floating policy's requirement of at least 1 associated machine.
	LicenseValidationStatusTooManyMachine      = "too_many_machines"     // The validated license has exceeded its policy's machine limit.
	LicenseValidationStatusHeartbeatNotStarted = "heartbeat_not_started" // The validated machine or fingerprint scope requires a heartbeat but one is not started.
	LicenseValidationStatusHeartbeatDead       = "heartbeat_dead"        // The validated machine or fingerprint scope belongs to a dead machine.
)
