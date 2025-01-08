package constants

const (
	// PolicySchemeED25519 signs license keys with your account's
	// Ed25519 signing key,
	PolicySchemeED25519 = "ED25519"

	// PolicySchemeRSA2048PKCS1 signs license keys with your account's
	// 2048-bit RSA private key using RSA PKCS1 v1.5 padding
	PolicySchemeRSA2048PKCS1 = "RSA2048PKCS1"
)

var ValidPolicySchemeMapper = map[string]bool{
	PolicySchemeED25519:      true,
	PolicySchemeRSA2048PKCS1: true,
}

const (
	// PolicyExpirationStrategyRestrictAccess - expired licenses can
	// continue to access releases published prior to
	// their license expiry. This is the default.
	PolicyExpirationStrategyRestrictAccess = "restrict"

	// PolicyExpirationStrategyRevokeAccess - Expired licenses are
	// no longer able to access any releases.
	PolicyExpirationStrategyRevokeAccess = "revoke"

	// PolicyExpirationStrategyMaintainAccess - Expired licenses can continue
	// to access releases published prior to their license expiry.
	PolicyExpirationStrategyMaintainAccess = "maintain"

	// PolicyExpirationStrategyAllowAccess - Expired licenses can access any releases.
	PolicyExpirationStrategyAllowAccess = "allow"
)

var ValidPolicyExpirationStrategyMapper = map[string]bool{
	PolicyExpirationStrategyRestrictAccess: true,
	PolicyExpirationStrategyRevokeAccess:   true,
	PolicyExpirationStrategyMaintainAccess: true,
	PolicyExpirationStrategyAllowAccess:    true,
}

const (
	// PolicyCheckinIntervalDaily requires a license implementing
	//the policy checkin at least once every day to remain valid.
	PolicyCheckinIntervalDaily = "daily"
	// PolicyCheckinIntervalWeekly requires a license implementing
	// the policy checkin at least once every week to remain valid.
	PolicyCheckinIntervalWeekly = "weekly"
	// PolicyCheckinIntervalMonthly requires a license implementing
	// the policy checkin at least once every month to remain valid.
	PolicyCheckinIntervalMonthly = "monthly"
	// PolicyCheckinIntervalYearly requires a license implementing
	// the policy to check in at least once every year to remain valid.
	PolicyCheckinIntervalYearly = "yearly"
)

var ValidPolicyCheckinIntervalMapper = map[string]bool{
	PolicyCheckinIntervalDaily:   true,
	PolicyCheckinIntervalWeekly:  true,
	PolicyCheckinIntervalMonthly: true,
	PolicyCheckinIntervalYearly:  true,
}

const (
	// PolicyRenewalBasisFromExpiry - License expiry is extended from the license's current expiry value,
	// i.e. license.expiry = license.expiry + policy.duration. This is the default.
	PolicyRenewalBasisFromExpiry = "from_expiry"
	// PolicyRenewalFromNow - License expiry is extended from the current time,
	// i.e. license.expiry = time.now + policy.duration.
	PolicyRenewalFromNow = "from_now"
	// PolicyRenewalFromNowIfExpired - Conditionally extend license expiry from
	// the current time if the license is expired, otherwise extend from the license's current expiry value.
	PolicyRenewalFromNowIfExpired = "from_now_if_expired"
)

var ValidPolicyRenewalBasisMapper = map[string]bool{
	PolicyRenewalBasisFromExpiry:  true,
	PolicyRenewalFromNow:          true,
	PolicyRenewalFromNowIfExpired: true,
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
