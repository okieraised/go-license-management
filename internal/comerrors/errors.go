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
)

var (
	ErrProductNameIsEmpty                   = errors.New("product name is empty")
	ErrProductCodeIsEmpty                   = errors.New("product code is empty")
	ErrProductDistributionStrategyIsInvalid = errors.New("product distribution strategy is invalid")
	ErrProductNameAlreadyExist              = errors.New("product name already exists")
	ErrProductCodeAlreadyExist              = errors.New("product code already exists")
)

var (
	ErrEntitlementIDIsEmpty        = errors.New("entitlement id is empty")
	ErrEntitlementNameIsEmpty      = errors.New("entitlement name is empty")
	ErrEntitlementCodeIsEmpty      = errors.New("entitlement code is empty")
	ErrEntitlementCodeAlreadyExist = errors.New("entitlement code already exists")
	ErrEntitlementIDIsInvalid      = errors.New("entitlement id is invalid")
)

var ErrCodeMapper = map[error]string{
	nil:                                     "00000",
	ErrGenericInternalServer:                "50000",
	ErrInvalidDatabaseClient:                "50001",
	ErrGenericRequestTimedOut:               "50004",
	ErrGenericBadRequest:                    "40000",
	ErrGenericUnauthorized:                  "40001",
	ErrGenericPermission:                    "40003",
	ErrTenantNameIsEmpty:                    "42000",
	ErrTenantNameAlreadyExist:               "42001",
	ErrTenantNameIsInvalid:                  "42002",
	ErrAccountUsernameIsEmpty:               "43000",
	ErrAccountEmailIsEmpty:                  "43001",
	ErrAccountRoleIsEmpty:                   "43002",
	ErrAccountRoleIsInvalid:                 "43003",
	ErrAccountPasswordIsEmpty:               "43004",
	ErrAccountUsernameAlreadyExist:          "43005",
	ErrProductNameIsEmpty:                   "44000",
	ErrProductCodeIsEmpty:                   "44001",
	ErrProductDistributionStrategyIsInvalid: "44002",
	ErrProductNameAlreadyExist:              "44003",
	ErrProductCodeAlreadyExist:              "44004",
	ErrEntitlementIDIsEmpty:                 "45000",
	ErrEntitlementNameIsEmpty:               "45001",
	ErrEntitlementCodeIsEmpty:               "45002",
	ErrEntitlementCodeAlreadyExist:          "45003",
	ErrEntitlementIDIsInvalid:               "45004",
}

var ErrMessageMapper = map[error]string{
	nil:                                     "OK",
	ErrGenericInternalServer:                ErrGenericInternalServer.Error(),
	ErrGenericRequestTimedOut:               ErrGenericRequestTimedOut.Error(),
	ErrInvalidDatabaseClient:                ErrInvalidDatabaseClient.Error(),
	ErrGenericBadRequest:                    ErrGenericBadRequest.Error(),
	ErrGenericUnauthorized:                  ErrGenericUnauthorized.Error(),
	ErrGenericPermission:                    ErrGenericPermission.Error(),
	ErrTenantNameIsEmpty:                    ErrTenantNameIsEmpty.Error(),
	ErrTenantNameAlreadyExist:               ErrTenantNameAlreadyExist.Error(),
	ErrTenantNameIsInvalid:                  ErrTenantNameIsInvalid.Error(),
	ErrAccountUsernameIsEmpty:               ErrAccountUsernameIsEmpty.Error(),
	ErrAccountEmailIsEmpty:                  ErrAccountEmailIsEmpty.Error(),
	ErrAccountRoleIsEmpty:                   ErrAccountRoleIsEmpty.Error(),
	ErrAccountRoleIsInvalid:                 ErrAccountRoleIsInvalid.Error(),
	ErrAccountPasswordIsEmpty:               ErrAccountPasswordIsEmpty.Error(),
	ErrAccountUsernameAlreadyExist:          ErrAccountUsernameAlreadyExist.Error(),
	ErrProductNameIsEmpty:                   ErrProductNameIsEmpty.Error(),
	ErrProductCodeIsEmpty:                   ErrProductCodeIsEmpty.Error(),
	ErrProductDistributionStrategyIsInvalid: ErrProductDistributionStrategyIsInvalid.Error(),
	ErrProductNameAlreadyExist:              ErrProductNameAlreadyExist.Error(),
	ErrProductCodeAlreadyExist:              ErrProductCodeAlreadyExist.Error(),
	ErrEntitlementIDIsEmpty:                 ErrEntitlementIDIsEmpty.Error(),
	ErrEntitlementNameIsEmpty:               ErrEntitlementNameIsEmpty.Error(),
	ErrEntitlementCodeIsEmpty:               ErrEntitlementCodeIsEmpty.Error(),
	ErrEntitlementCodeAlreadyExist:          ErrEntitlementCodeAlreadyExist.Error(),
	ErrEntitlementIDIsInvalid:               ErrEntitlementIDIsInvalid.Error(),
}
