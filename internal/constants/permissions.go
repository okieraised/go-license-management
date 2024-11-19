package constants

const (
	TenantCreate = "tenant.create"
	TenantUpdate = "tenant.update"
	TenantDelete = "tenant.delete"
	TenantRead   = "tenant.read"
)

const (
	AdminCreate = "admin.create"
	AdminDelete = "admin.delete"
	AdminRead   = "admin.read"
	AdminUpdate = "admin.update"
)

const (
	UserBan            = "user.ban"
	UserCreate         = "user.create"
	UserDelete         = "user.delete"
	UserPasswordReset  = "user_password.reset"
	UserPasswordUpdate = "user_password.update"
	UserRead           = "user.read"
	UserUnban          = "user.unban"
	UserUpdate         = "user.update"
)

const (
	EntitlementCreate = "entitlement.create"
	EntitlementDelete = "entitlement.delete"
	EntitlementRead   = "entitlement.read"
	EntitlementUpdate = "entitlement.update"
)

const (
	ProductCreate         = "product.create"
	ProductDelete         = "product.delete"
	ProductRead           = "product.read"
	ProductTokensGenerate = "product_tokens.generate"
	ProductUpdate         = "product.update"
)

const (
	PolicyCreate             = "policy.create"
	PolicyDelete             = "policy.delete"
	PolicyRead               = "policy.read"
	PolicyUpdate             = "policy.update"
	PolicyEntitlementsAttach = "policy_entitlements.attach"
	PolicyEntitlementsDetach = "policy_entitlements.detach"
)

const (
	LicenseCheckIn            = "license.check-in"
	LicenseCheckOut           = "license.check-out"
	LicenseCreate             = "license.create"
	LicenseDelete             = "license.delete"
	LicenseRead               = "license.read"
	LicenseReinstate          = "license.reinstate"
	LicenseRenew              = "license.renew"
	LicenseRevoke             = "license.revoke"
	LicenseSuspend            = "license.suspend"
	LicenseValidate           = "license.validate"
	LicenseUpdate             = "license.update"
	LicenseUsageDecrement     = "license-usage.decrement"
	LicenseUsageIncrement     = "license-usage.increment"
	LicenseTokensGenerate     = "license-tokens.generate"
	LicenseUsageReset         = "license-usage.reset"
	LicenseEntitlementsAttach = "license-entitlements.attach"
	LicenseEntitlementsDetach = "license-entitlements.detach"
	LicensePolicyUpdate       = "license-policy.update"
	LicenseUsersAttach        = "license-users.attach"
	LicenseUsersDetach        = "license-users.detach"
)

const (
	MachineCreate         = "machine.create"
	MachineDelete         = "machine.delete"
	MachineRead           = "machine.read"
	MachineUpdate         = "machine.update"
	MachineCheckOut       = "machine.check-out"
	MachineHeartbeatPing  = "machine-heartbeat.ping"
	MachineHeartbeatReset = "machine-heartbeat.reset"
)

var SuperAdminPermissionMapper = map[string]bool{
	TenantCreate:              true,
	TenantUpdate:              true,
	TenantDelete:              true,
	TenantRead:                true,
	AdminCreate:               true,
	AdminDelete:               true,
	AdminRead:                 true,
	AdminUpdate:               true,
	UserBan:                   true,
	UserCreate:                true,
	UserDelete:                true,
	UserPasswordReset:         true,
	UserPasswordUpdate:        true,
	UserRead:                  true,
	UserUnban:                 true,
	UserUpdate:                true,
	EntitlementCreate:         true,
	EntitlementDelete:         true,
	EntitlementRead:           true,
	EntitlementUpdate:         true,
	ProductCreate:             true,
	ProductDelete:             true,
	ProductRead:               true,
	ProductTokensGenerate:     true,
	ProductUpdate:             true,
	PolicyCreate:              true,
	PolicyDelete:              true,
	PolicyRead:                true,
	PolicyUpdate:              true,
	PolicyEntitlementsAttach:  true,
	PolicyEntitlementsDetach:  true,
	LicenseCheckIn:            true,
	LicenseCheckOut:           true,
	LicenseCreate:             true,
	LicenseDelete:             true,
	LicenseRead:               true,
	LicenseReinstate:          true,
	LicenseRenew:              true,
	LicenseRevoke:             true,
	LicenseSuspend:            true,
	LicenseValidate:           true,
	LicenseUpdate:             true,
	LicenseUsageDecrement:     true,
	LicenseUsageIncrement:     true,
	LicenseTokensGenerate:     true,
	LicenseUsageReset:         true,
	LicenseEntitlementsAttach: true,
	LicenseEntitlementsDetach: true,
	LicensePolicyUpdate:       true,
	LicenseUsersAttach:        true,
	LicenseUsersDetach:        true,
	MachineCreate:             true,
	MachineDelete:             true,
	MachineRead:               true,
	MachineUpdate:             true,
	MachineCheckOut:           true,
	MachineHeartbeatPing:      true,
	MachineHeartbeatReset:     true,
}

//func CreateAdminPermission(domain string) [][]string {
//
//	result := make([][]string, 0)
//	for key, _ := range AdminPermissions {
//		parts := strings.Split(key, ".")
//
//		object := ""
//		perm := ""
//
//		if len(parts) == 3 {
//			object = parts[0] + "_" + parts[1]
//			perm = parts[2]
//		} else {
//			object = parts[0]
//			perm = parts[1]
//		}
//
//		result = append(result, []string{"p", domain, "admin", object, perm})
//	}
//
//	return result
//}
