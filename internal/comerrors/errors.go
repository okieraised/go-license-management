package comerrors

import "errors"

var (
	ErrInvalidDatabaseClient = errors.New("invalid database client")
)

var (
	ErrGenericBadRequest      = errors.New("bad request")
	ErrGenericUnauthorized    = errors.New("unauthorized request")
	ErrGenericPermission      = errors.New("invalid permission")
	ErrGenericInternalServer  = errors.New("internal server error")
	ErrGenericRequestTimedOut = errors.New("request timeout error")
)

var (
	ErrTenantNameIsEmpty      = errors.New("tenant name is empty")
	ErrTenantNameAlreadyExist = errors.New("tenant name already exists")
	ErrTenantNameIsInvalid    = errors.New("tenant name is invalid")
)

var (
	ErrAccountUsernameIsEmpty      = errors.New("account username is empty")
	ErrAccountPasswordIsEmpty      = errors.New("account password is empty")
	ErrAccountEmailIsEmpty         = errors.New("account email is empty")
	ErrAccountRoleIsEmpty          = errors.New("account role is empty")
	ErrAccountRoleIsInvalid        = errors.New("account role is invalid")
	ErrAccountUsernameAlreadyExist = errors.New("account username already exists")
	ErrAccountActionIsEmpty        = errors.New("account action is empty")
	ErrAccountActionIsInvalid      = errors.New("account action is invalid")
	ErrAccountUsernameIsInvalid    = errors.New("account username is invalid")
)

var (
	ErrProductNameIsEmpty                    = errors.New("product name is empty")
	ErrProductCodeIsEmpty                    = errors.New("product code is empty")
	ErrProductDistributionStrategyIsInvalid  = errors.New("product distribution strategy is invalid")
	ErrProductNameAlreadyExist               = errors.New("product name already exists")
	ErrProductCodeAlreadyExist               = errors.New("product code already exists")
	ErrProductIDIsEmpty                      = errors.New("product id is empty")
	ErrProductIDIsInvalid                    = errors.New("product id is invalid")
	ErrProductTokenExpirationFormatIsInvalid = errors.New("product token expiration format is invalid")
)

var (
	ErrEntitlementIDIsEmpty        = errors.New("entitlement id is empty")
	ErrEntitlementNameIsEmpty      = errors.New("entitlement name is empty")
	ErrEntitlementCodeIsEmpty      = errors.New("entitlement code is empty")
	ErrEntitlementCodeAlreadyExist = errors.New("entitlement code already exists")
	ErrEntitlementIDIsInvalid      = errors.New("entitlement id is invalid")
)

var (
	ErrPolicyNameIsEmpty                     = errors.New("policy name is empty")
	ErrPolicySchemeIsInvalid                 = errors.New("policy scheme is invalid")
	ErrPolicyIDIsEmpty                       = errors.New("policy id is empty")
	ErrPolicyIDIsInvalid                     = errors.New("policy id is invalid")
	ErrPolicyDurationIsLessThanZero          = errors.New("policy duration is less than zero")
	ErrPolicyMaxMachinesIsLessThanZero       = errors.New("policy max machines is less than zero")
	ErrPolicyMaxUsesIsLessThanZero           = errors.New("policy max uses is less than zero")
	ErrPolicyMaxUsersIsLessThanZero          = errors.New("policy max users is less than zero")
	ErrPolicyHeartbeatDurationIsLessThanZero = errors.New("policy heartbeat duration is less than zero")
	ErrPolicyInvalidExpirationStrategy       = errors.New("policy expiration strategy is invalid")
	ErrPolicyInvalidAuthenticationStrategy   = errors.New("policy authentication strategy is invalid")
	ErrPolicyInvalidExpirationBasis          = errors.New("policy expiration basis is invalid")
	ErrPolicyInvalidOverageStrategy          = errors.New("policy overage strategy is invalid")
	ErrPolicyInvalidRenewalBasis             = errors.New("policy renewal basis is invalid")
	ErrPolicyInvalidHeartbeatBasis           = errors.New("policy heartbeat basis is invalid")
	ErrPolicyInvalidCheckinIntervalBasis     = errors.New("policy checkin interval basis is invalid")
)

var (
	ErrLicenseNameIsEmpty           = errors.New("license name is empty")
	ErrLicenseProductIDIsEmpty      = errors.New("license product id is empty")
	ErrLicensePolicyIDIsEmpty       = errors.New("license policy id is empty")
	ErrLicenseExpiryFormatIsInvalid = errors.New("license expiry format is invalid")
	ErrLicenseIDIsEmpty             = errors.New("license id is empty")
	ErrLicenseIDIsInvalid           = errors.New("license id is invalid")
	ErrLicenseActionIsEmpty         = errors.New("license action is empty")
	ErrLicenseActionIsInvalid       = errors.New("license action is invalid")
	ErrLicenseIsSuspended           = errors.New("license is suspended")
	ErrLicenseIsExpired             = errors.New("license is expired")
	ErrLicenseIsBanned              = errors.New("license is banned")
)

var (
	ErrMachineIDIsEmpty                        = errors.New("machine id is empty")
	ErrMachineIDIsInvalid                      = errors.New("machine id is invalid")
	ErrMachineFingerprintIsEmpty               = errors.New("machine fingerprint is empty")
	ErrMachineLicenseIsEmpty                   = errors.New("machine license is empty")
	ErrMachineLicenseIsInvalid                 = errors.New("machine license is invalid")
	ErrMachineFingerprintAssociatedWithLicense = errors.New("machine fingerprint is already associated with a license")
	ErrMachineActionIsEmpty                    = errors.New("machine action is empty")
	ErrMachineActionIsInvalid                  = errors.New("machine action is invalid")
	ErrMachineActionCheckoutTTLIsInvalid       = errors.New("machine license TTL is invalid (must be >= 3600 or <= 31556952 seconds)")
)

var ErrCodeMapper = map[error]string{
	nil:                                        "00000",
	ErrGenericInternalServer:                   "50000",
	ErrInvalidDatabaseClient:                   "50001",
	ErrGenericRequestTimedOut:                  "50004",
	ErrGenericBadRequest:                       "40000",
	ErrGenericUnauthorized:                     "40001",
	ErrGenericPermission:                       "40003",
	ErrTenantNameIsEmpty:                       "42000",
	ErrTenantNameAlreadyExist:                  "42001",
	ErrTenantNameIsInvalid:                     "42002",
	ErrAccountUsernameIsEmpty:                  "43000",
	ErrAccountEmailIsEmpty:                     "43001",
	ErrAccountRoleIsEmpty:                      "43002",
	ErrAccountRoleIsInvalid:                    "43003",
	ErrAccountPasswordIsEmpty:                  "43004",
	ErrAccountUsernameAlreadyExist:             "43005",
	ErrAccountActionIsEmpty:                    "43006",
	ErrAccountActionIsInvalid:                  "43007",
	ErrAccountUsernameIsInvalid:                "43008",
	ErrProductNameIsEmpty:                      "44000",
	ErrProductCodeIsEmpty:                      "44001",
	ErrProductDistributionStrategyIsInvalid:    "44002",
	ErrProductNameAlreadyExist:                 "44003",
	ErrProductCodeAlreadyExist:                 "44004",
	ErrProductIDIsEmpty:                        "44005",
	ErrProductIDIsInvalid:                      "44006",
	ErrProductTokenExpirationFormatIsInvalid:   "44007",
	ErrEntitlementIDIsEmpty:                    "45000",
	ErrEntitlementNameIsEmpty:                  "45001",
	ErrEntitlementCodeIsEmpty:                  "45002",
	ErrEntitlementCodeAlreadyExist:             "45003",
	ErrEntitlementIDIsInvalid:                  "45004",
	ErrPolicyNameIsEmpty:                       "46000",
	ErrPolicySchemeIsInvalid:                   "46001",
	ErrPolicyIDIsEmpty:                         "46002",
	ErrPolicyIDIsInvalid:                       "46003",
	ErrPolicyDurationIsLessThanZero:            "46004",
	ErrPolicyMaxMachinesIsLessThanZero:         "46005",
	ErrPolicyMaxUsesIsLessThanZero:             "46006",
	ErrPolicyMaxUsersIsLessThanZero:            "46007",
	ErrPolicyHeartbeatDurationIsLessThanZero:   "46008",
	ErrPolicyInvalidExpirationStrategy:         "46009",
	ErrPolicyInvalidAuthenticationStrategy:     "46010",
	ErrPolicyInvalidExpirationBasis:            "46011",
	ErrPolicyInvalidOverageStrategy:            "46012",
	ErrPolicyInvalidRenewalBasis:               "46013",
	ErrPolicyInvalidHeartbeatBasis:             "46014",
	ErrPolicyInvalidCheckinIntervalBasis:       "46015",
	ErrLicenseNameIsEmpty:                      "47001",
	ErrLicenseProductIDIsEmpty:                 "47002",
	ErrLicensePolicyIDIsEmpty:                  "47003",
	ErrLicenseExpiryFormatIsInvalid:            "47004",
	ErrLicenseIDIsEmpty:                        "47005",
	ErrLicenseIDIsInvalid:                      "47006",
	ErrLicenseActionIsEmpty:                    "47007",
	ErrLicenseActionIsInvalid:                  "47008",
	ErrLicenseIsSuspended:                      "47009",
	ErrLicenseIsExpired:                        "47010",
	ErrLicenseIsBanned:                         "47011",
	ErrMachineIDIsEmpty:                        "48000",
	ErrMachineIDIsInvalid:                      "48001",
	ErrMachineFingerprintIsEmpty:               "48002",
	ErrMachineLicenseIsEmpty:                   "48003",
	ErrMachineLicenseIsInvalid:                 "48004",
	ErrMachineFingerprintAssociatedWithLicense: "48005",
	ErrMachineActionIsEmpty:                    "48006",
	ErrMachineActionIsInvalid:                  "48007",
	ErrMachineActionCheckoutTTLIsInvalid:       "48008",
}

var ErrMessageMapper = map[error]string{
	nil:                                        "OK",
	ErrGenericInternalServer:                   ErrGenericInternalServer.Error(),
	ErrGenericRequestTimedOut:                  ErrGenericRequestTimedOut.Error(),
	ErrInvalidDatabaseClient:                   ErrInvalidDatabaseClient.Error(),
	ErrGenericBadRequest:                       ErrGenericBadRequest.Error(),
	ErrGenericUnauthorized:                     ErrGenericUnauthorized.Error(),
	ErrGenericPermission:                       ErrGenericPermission.Error(),
	ErrTenantNameIsEmpty:                       ErrTenantNameIsEmpty.Error(),
	ErrTenantNameAlreadyExist:                  ErrTenantNameAlreadyExist.Error(),
	ErrTenantNameIsInvalid:                     ErrTenantNameIsInvalid.Error(),
	ErrAccountUsernameIsEmpty:                  ErrAccountUsernameIsEmpty.Error(),
	ErrAccountEmailIsEmpty:                     ErrAccountEmailIsEmpty.Error(),
	ErrAccountRoleIsEmpty:                      ErrAccountRoleIsEmpty.Error(),
	ErrAccountRoleIsInvalid:                    ErrAccountRoleIsInvalid.Error(),
	ErrAccountPasswordIsEmpty:                  ErrAccountPasswordIsEmpty.Error(),
	ErrAccountUsernameAlreadyExist:             ErrAccountUsernameAlreadyExist.Error(),
	ErrAccountActionIsEmpty:                    ErrAccountActionIsEmpty.Error(),
	ErrAccountActionIsInvalid:                  ErrAccountActionIsInvalid.Error(),
	ErrAccountUsernameIsInvalid:                ErrAccountUsernameIsInvalid.Error(),
	ErrProductNameIsEmpty:                      ErrProductNameIsEmpty.Error(),
	ErrProductCodeIsEmpty:                      ErrProductCodeIsEmpty.Error(),
	ErrProductDistributionStrategyIsInvalid:    ErrProductDistributionStrategyIsInvalid.Error(),
	ErrProductNameAlreadyExist:                 ErrProductNameAlreadyExist.Error(),
	ErrProductCodeAlreadyExist:                 ErrProductCodeAlreadyExist.Error(),
	ErrProductIDIsEmpty:                        ErrProductIDIsEmpty.Error(),
	ErrProductIDIsInvalid:                      ErrProductIDIsInvalid.Error(),
	ErrProductTokenExpirationFormatIsInvalid:   ErrProductTokenExpirationFormatIsInvalid.Error(),
	ErrEntitlementIDIsEmpty:                    ErrEntitlementIDIsEmpty.Error(),
	ErrEntitlementNameIsEmpty:                  ErrEntitlementNameIsEmpty.Error(),
	ErrEntitlementCodeIsEmpty:                  ErrEntitlementCodeIsEmpty.Error(),
	ErrEntitlementCodeAlreadyExist:             ErrEntitlementCodeAlreadyExist.Error(),
	ErrEntitlementIDIsInvalid:                  ErrEntitlementIDIsInvalid.Error(),
	ErrPolicyNameIsEmpty:                       ErrPolicyNameIsEmpty.Error(),
	ErrPolicySchemeIsInvalid:                   ErrPolicySchemeIsInvalid.Error(),
	ErrPolicyIDIsEmpty:                         ErrPolicyIDIsEmpty.Error(),
	ErrPolicyIDIsInvalid:                       ErrPolicyIDIsInvalid.Error(),
	ErrPolicyDurationIsLessThanZero:            ErrPolicyDurationIsLessThanZero.Error(),
	ErrPolicyMaxMachinesIsLessThanZero:         ErrPolicyMaxMachinesIsLessThanZero.Error(),
	ErrPolicyMaxUsesIsLessThanZero:             ErrPolicyMaxUsesIsLessThanZero.Error(),
	ErrPolicyMaxUsersIsLessThanZero:            ErrPolicyMaxUsersIsLessThanZero.Error(),
	ErrPolicyHeartbeatDurationIsLessThanZero:   ErrPolicyHeartbeatDurationIsLessThanZero.Error(),
	ErrPolicyInvalidExpirationStrategy:         ErrPolicyInvalidExpirationStrategy.Error(),
	ErrPolicyInvalidAuthenticationStrategy:     ErrPolicyInvalidAuthenticationStrategy.Error(),
	ErrPolicyInvalidExpirationBasis:            ErrPolicyInvalidExpirationBasis.Error(),
	ErrPolicyInvalidOverageStrategy:            ErrPolicyInvalidOverageStrategy.Error(),
	ErrPolicyInvalidRenewalBasis:               ErrPolicyInvalidRenewalBasis.Error(),
	ErrPolicyInvalidHeartbeatBasis:             ErrPolicyInvalidHeartbeatBasis.Error(),
	ErrPolicyInvalidCheckinIntervalBasis:       ErrPolicyInvalidCheckinIntervalBasis.Error(),
	ErrLicenseNameIsEmpty:                      ErrLicenseNameIsEmpty.Error(),
	ErrLicenseProductIDIsEmpty:                 ErrLicenseProductIDIsEmpty.Error(),
	ErrLicensePolicyIDIsEmpty:                  ErrLicensePolicyIDIsEmpty.Error(),
	ErrLicenseExpiryFormatIsInvalid:            ErrLicenseExpiryFormatIsInvalid.Error(),
	ErrLicenseIDIsEmpty:                        ErrLicenseIDIsEmpty.Error(),
	ErrLicenseIDIsInvalid:                      ErrLicenseIDIsInvalid.Error(),
	ErrLicenseActionIsEmpty:                    ErrLicenseActionIsEmpty.Error(),
	ErrLicenseActionIsInvalid:                  ErrLicenseActionIsInvalid.Error(),
	ErrLicenseIsSuspended:                      ErrLicenseIsSuspended.Error(),
	ErrLicenseIsExpired:                        ErrLicenseIsExpired.Error(),
	ErrLicenseIsBanned:                         ErrLicenseIsBanned.Error(),
	ErrMachineIDIsEmpty:                        ErrMachineIDIsEmpty.Error(),
	ErrMachineIDIsInvalid:                      ErrMachineIDIsInvalid.Error(),
	ErrMachineFingerprintIsEmpty:               ErrMachineFingerprintIsEmpty.Error(),
	ErrMachineLicenseIsEmpty:                   ErrMachineLicenseIsEmpty.Error(),
	ErrMachineLicenseIsInvalid:                 ErrMachineLicenseIsInvalid.Error(),
	ErrMachineFingerprintAssociatedWithLicense: ErrMachineFingerprintAssociatedWithLicense.Error(),
	ErrMachineActionIsEmpty:                    ErrMachineActionIsEmpty.Error(),
	ErrMachineActionIsInvalid:                  ErrMachineActionIsInvalid.Error(),
	ErrMachineActionCheckoutTTLIsInvalid:       ErrMachineActionCheckoutTTLIsInvalid.Error(),
}
