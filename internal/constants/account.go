package constants

const (
	AccountStatusActive   = "active"
	AccountStatusInactive = "inactive"
	AccountStatusBanned   = "banned"
)

const (
	AccountActionUpdatePassword     = "update-password"
	AccountActionResetPassword      = "reset-password"
	AccountActionGenerateResetToken = "password-token"
	AccountActionBan                = "ban"
	AccountActionUnban              = "unban"
)

var ValidAccountActionMapper = map[string]bool{
	AccountActionUpdatePassword:     true,
	AccountActionResetPassword:      true,
	AccountActionGenerateResetToken: true,
	AccountActionBan:                true,
	AccountActionUnban:              true,
}
