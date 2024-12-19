package constants

const (
	// PolicySchemeED25519 signs license keys with your account's Ed25519 signing key,
	// using elliptic curve cryptography and SHA512. The given license key data will be base64url encoded.
	PolicySchemeED25519 = "ED25519"

	// PolicySchemeRSA2048PKCS1 signs license keys with your account's 2048-bit RSA private key using RSA
	// PKCS1 v1.5 padding, with a SHA256 digest. The given license key data will be base64url encoded.
	PolicySchemeRSA2048PKCS1 = "RSA2048PKCS1"
)

var ValidPolicySchemeMapper = map[string]bool{
	PolicySchemeED25519:      true,
	PolicySchemeRSA2048PKCS1: true,
}

const (
	// PolicyExpirationStrategyRestrictAccess - expired licenses can continue to access releases published prior to
	// their license expiry. Automatic upgrades are enabled, but only for releases published prior to their expiry.
	// Validation scopes take precedence over expiry check during license validation. This is the default.
	PolicyExpirationStrategyRestrictAccess = "restrict"

	// PolicyExpirationStrategyRevokeAccess - Expired licenses are no longer able to access any releases,
	// including past releases. Automatic upgrades are disabled. Expiry check takes precedence over
	// scopes during license validation.
	PolicyExpirationStrategyRevokeAccess = "revoke"

	// PolicyExpirationStrategyMaintainAccess - Expired licenses can continue to access releases published prior to their
	// license expiry. Automatic upgrades are enabled, but only for releases published prior to their expiry.
	// Validation scopes take precedence over expiry check during license validation. In addition,
	// validations with an EXPIRED code will return valid=true
	PolicyExpirationStrategyMaintainAccess = "maintain"

	// PolicyExpirationStrategyAllowAccess - Expired licenses can access any releases, including past releases and
	// future releases. Automatic upgrades are enabled. Validation scopes take precedence over expiry check during
	// license validation. In addition, validations with an EXPIRED code will return valid=true.
	PolicyExpirationStrategyAllowAccess = "allow"
)

var ValidPolicyExpirationStrategyMapper = map[string]bool{
	PolicyExpirationStrategyRestrictAccess: true,
	PolicyExpirationStrategyRevokeAccess:   true,
	PolicyExpirationStrategyMaintainAccess: true,
	PolicyExpirationStrategyAllowAccess:    true,
}

const (
	// PolicyCheckinIntervalDaily requires a license implementing the policy checkin at least once every day to remain valid.
	PolicyCheckinIntervalDaily = "daily"
	// PolicyCheckinIntervalWeekly requires a license implementing the policy checkin at least once every week to remain valid.
	PolicyCheckinIntervalWeekly = "weekly"
	// PolicyCheckinIntervalMonthly requires a license implementing the policy checkin at least once every month to remain valid.
	PolicyCheckinIntervalMonthly = "monthly"
	// PolicyCheckinIntervalYearly requires a license implementing the policy to check-in at least once every year to remain valid.
	PolicyCheckinIntervalYearly = "yearly"
)

var ValidPolicyCheckinIntervalMapper = map[string]bool{
	PolicyCheckinIntervalDaily:   true,
	PolicyCheckinIntervalWeekly:  true,
	PolicyCheckinIntervalMonthly: true,
	PolicyCheckinIntervalYearly:  true,
}

const (
	// PolicyExpirationBasisFromCreation - License expirations are set immediately upon creation.
	PolicyExpirationBasisFromCreation = "from_creation"
	// PolicyExpirationBasisFromFirstActivation - License expirations are set after their first license validation event.
	PolicyExpirationBasisFromFirstActivation = "from_first_activation"
	// PolicyExpirationBasisFromFirstValidation - License expirations are set after their first machine activation event.
	PolicyExpirationBasisFromFirstValidation = "from_first_validation"
	// PolicyExpirationBasisFromFirstUse - License expirations are set after their first usage increment event.
	PolicyExpirationBasisFromFirstUse = "from_first_use"
)

var ValidPolicyExpirationBasisMapper = map[string]bool{
	PolicyExpirationBasisFromCreation:        true,
	PolicyExpirationBasisFromFirstActivation: true,
	PolicyExpirationBasisFromFirstValidation: true,
	PolicyExpirationBasisFromFirstUse:        true,
}

const (
	// PolicyRenewalBasisFromExpiry - License expiry is extended from the license's current expiry value,
	// i.e. license.expiry = license.expiry + policy.duration. This is the default.
	PolicyRenewalBasisFromExpiry = "from_expiry"
	// PolicyRenewalFromNow - License expiry is extended from the current time, i.e. license.expiry = time.now + policy.duration.
	PolicyRenewalFromNow = "from_now"
	// PolicyRenewalFromNowIfExpired - Conditionally extend license expiry from the current time if the license is expired, otherwise extend from the license's current expiry value.
	PolicyRenewalFromNowIfExpired = "from_now_if_expired"
)

var ValidPolicyRenewalBasisMapper = map[string]bool{
	PolicyRenewalBasisFromExpiry:  true,
	PolicyRenewalFromNow:          true,
	PolicyRenewalFromNowIfExpired: true,
}

const (
	// PolicyAuthenticationStrategyToken - Allow licenses to authenticate using a license token. This is the default.
	PolicyAuthenticationStrategyToken = "auth_token"

	// PolicyAuthenticationStrategyLicense - Allow licenses to authenticate using a license key.
	PolicyAuthenticationStrategyLicense = "auth_license"

	// PolicyAuthenticationStrategyMixes - Allow both license token and license key authentication.
	PolicyAuthenticationStrategyMixes = "auth_mixed"

	// PolicyAuthenticationStrategyNone - Disable the ability for licenses to authenticate with the API.
	PolicyAuthenticationStrategyNone = "auth_none"
)

var ValidPolicyAuthenticationStrategyMap = map[string]bool{
	PolicyAuthenticationStrategyToken:   true,
	PolicyAuthenticationStrategyLicense: true,
	PolicyAuthenticationStrategyMixes:   true,
	PolicyAuthenticationStrategyNone:    true,
}

const (
	// PolicyHeartbeatBasisFromCreation - Machine heartbeat is started immediately upon creation.
	PolicyHeartbeatBasisFromCreation = "from_creation"

	// PolicyHeartbeatBasisFromFirstPing - Machine heartbeat is started after their first heartbeat ping event.
	PolicyHeartbeatBasisFromFirstPing = "from_first_ping"
)

var ValidPolicyHeartbeatBasisMapper = map[string]bool{
	PolicyHeartbeatBasisFromCreation:  true,
	PolicyHeartbeatBasisFromFirstPing: true,
}

const (
	// PolicyOverageStrategyNoOverage - Do not allow overages. Attempts to exceed limits will fail. This is the default.
	PolicyOverageStrategyNoOverage = "no_overage"
	// PolicyOverageStrategyAlwaysAllow - The license may exceed its limits, and doing so will not affect the license validity.
	PolicyOverageStrategyAlwaysAllow = "always_allow"
)

var ValidPolicyOverageStrategyMapper = map[string]bool{
	PolicyOverageStrategyNoOverage:   true,
	PolicyOverageStrategyAlwaysAllow: true,
}

const (
	// PolicyTransferStrategyResetExpiry resets the transferred license's expiry from the time of transfer.
	PolicyTransferStrategyResetExpiry = "reset_expiry"
	// PolicyTransferStrategyKeepExpiry keeps the license's current expiry.
	PolicyTransferStrategyKeepExpiry = "keep_expiry"
)

const (
	// PolicyHeartbeatCullPolicyDeactivateDead - Automatically deactivate machines that fail to maintain their heartbeat pings. This is the default.
	PolicyHeartbeatCullPolicyDeactivateDead = "deactivate_dead"
	// PolicyHeartbeatCullPolicyKeepDead - Mark machines that fail to maintain their heartbeat pings as dead, but do not deactivate.
	PolicyHeartbeatCullPolicyKeepDead = "keep_dead"
)

const (
	// PolicyHeartbeatResurrectionPolicyNoRevive -  Do not allow dead machines and processes to be revived. This is the default.
	PolicyHeartbeatResurrectionPolicyNoRevive = "no_revive"
	// PolicyHeartbeatResurrectionPolicyOneMinute - A machine or process can be revived if it sends a ping within 1 minute from its time of death.
	PolicyHeartbeatResurrectionPolicyOneMinute = "1_minute"
	// PolicyHeartbeatResurrectionPolicyOneHour - A machine or process can be revived if it sends a ping within 1 hour from its time of death.
	PolicyHeartbeatResurrectionPolicyOneHour = "1_hour"
	// PolicyHeartbeatResurrectionPolicyAlwaysRevive - A machine or process can always be revived.
	PolicyHeartbeatResurrectionPolicyAlwaysRevive = "always_revive"
)
