package constants

const (
	// PolicySchemeED25519 signs license keys with your account's Ed25519 signing key,
	// using elliptic curve cryptography and SHA512.
	// The given license key data will be base64url encoded and then prefixed with key/ before signing,
	// and the signing data's signature will be base64url encoded and then appended onto the end of
	// the signing data, delimited by the . character, e.g. key/{URLBASE64URL_KEY}.{URLBASE64URL_SIGNATURE}.
	// This is our recommended signing scheme, but it may not be supported in your preferred programming language.
	PolicySchemeED25519 = "ED25519"

	// PolicySchemeRSA2048PKCS1 signs license keys with your account's 2048-bit RSA private key using RSA
	// PKCS1 v1.5 padding, with a SHA256 digest. The provided embedded dataset will be base64url encoded and then
	// prefixed with key/ before signing, and the signing data's signature will be base64url encoded and
	// then appended onto the end of the signing data, delimited by the . character,
	// e.g. key/{URLBASE64URL_KEY}.{URLBASE64URL_SIGNATURE}, resulting in the final key.
	PolicySchemeRSA2048PKCS1 = "RSA2048PKCS1"

	// PolicySchemeRSA2048JWTRS256 encodes a license claims payload into a JWT using the RS256 algorithm.
	// The license key must be a valid JWT claims payload (i.e. a JSON encoded string).
	// The JWT will be signed using your account's 2048-bit RSA private key and
	// can be verified using your account's public key. The resulting key will be a full JSON Web Token.
	// We do not modify your claim payload.
	PolicySchemeRSA2048JWTRS256 = "RSA2048JWTRS256"
)

var ValidPolicySchemeMapper = map[string]bool{
	PolicySchemeED25519:         true,
	PolicySchemeRSA2048PKCS1:    true,
	PolicySchemeRSA2048JWTRS256: true,
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

const (
	// PolicyRenewalBasisFromExpiry - License expiry is extended from the license's current expiry value,
	//i.e. license.expiry = license.expiry + policy.duration. This is the default.
	PolicyRenewalBasisFromExpiry = "from_expiry"

	// PolicyRenewalFromNow - License expiry is extended from the current time, i.e. license.expiry = time.now + policy.duration.
	PolicyRenewalFromNow = "from_now"

	// PolicyRenewalFromNowIfExpired - Conditionally extend license expiry from the current time if the license is expired, otherwise extend from the license's current expiry value.
	PolicyRenewalFromNowIfExpired = "from_now_if_expired"
)

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

const (
	// PolicyHeartbeatBasisFromCreation - Machine heartbeat is started immediately upon creation.
	PolicyHeartbeatBasisFromCreation = "from_creation"

	// PolicyHeartbeatBasisFromFirstPing - Machine heartbeat is started after their first heartbeat ping event.
	PolicyHeartbeatBasisFromFirstPing = "from_first_ping"
)

const (
	// PolicyOverageStrategyNoOverage - Do not allow overages. Attempts to exceed limits will fail. This is the default.
	PolicyOverageStrategyNoOverage = "no_overage"
	// PolicyOverageStrategyAlwaysAllow - The license may exceed its limits, and doing so will not affect the license validity.
	PolicyOverageStrategyAlwaysAllow = "always_allow"
)

const (
	// PolicyProcessLeasingStrategyPerMachine - Processes are counted per-machine. This is the default.
	PolicyProcessLeasingStrategyPerMachine = "per_machine"
	// PolicyProcessLeasingStrategyPerUser - Processes are counted per-user, per-license.
	PolicyProcessLeasingStrategyPerUser = "per_user"
	// PolicyProcessLeasingStrategyPerLicense - Processes are counted per-license.
	PolicyProcessLeasingStrategyPerLicense = "per_license"
)

const (
	// PolicyTransferStrategyResetExpiry resets the transferred license's expiry from the time of transfer.
	PolicyTransferStrategyResetExpiry = "reset_expiry"
	// PolicyTransferStrategyKeepExpiry keeps the license's current expiry.
	PolicyTransferStrategyKeepExpiry = "keep_expiry"
)
