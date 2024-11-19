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

//var AdminPermissions = map[string]bool{
//	"account.analytics.read":      true,
//	"account.billing.read":        true,
//	"account.billing.update":      true,
//	"account.plan.read":           true,
//	"account.plan.update":         true,
//	"account.read":                true,
//	"account.subscription.read":   true,
//	"account.subscription.update": true,
//	"account.update":              true,
//	"admin.create":                true,
//	"admin.delete":                true,
//	"admin.invite":                true,
//	"admin.read":                  true,
//	"admin.update":                true,
//	"arch.read":                   true,
//	"artifact.create":             true,
//	"artifact.delete":             true,
//	"artifact.read":               true,
//	"artifact.update":             true,
//	"channel.read":                true,
//	"constraint.read":             true,
//	"engine.read":                 true,
//	"entitlement.create":          true,
//	"entitlement.delete":          true,
//	"entitlement.read":            true,
//	"entitlement.update":          true,
//	"environment.create":          true,
//	"environment.delete":          true,
//	"environment.read":            true,
//	"environment.tokens.generate": true,
//	"environment.update":          true,
//	"event-log.read":              true,
//	"group.create":                true,
//	"group.delete":                true,
//	"group.licenses.read":         true,
//	"group.machines.read":         true,
//	"group.owners.attach":         true,
//	"group.owners.detach":         true,
//	"group.owners.read":           true,
//	"group.read":                  true,
//	"group.update":                true,
//	"group.users.read":            true,
//	"key.create":                  true,
//	"key.delete":                  true,
//	"key.read":                    true,
//	"key.update":                  true,
//	"license.check-in":            true,
//	"license.check-out":           true,
//	"license.create":              true,
//	"license.delete":              true,
//	"license.entitlements.attach": true,
//	"license.entitlements.detach": true,
//	"license.group.update":        true,
//	"license.owner.update":        true,
//	"license.policy.update":       true,
//	"license.read":                true,
//	"license.reinstate":           true,
//	"license.renew":               true,
//	"license.revoke":              true,
//	"license.suspend":             true,
//	"license.tokens.generate":     true,
//	"license.update":              true,
//	"license.usage.decrement":     true,
//	"license.usage.increment":     true,
//	"license.usage.reset":         true,
//	"license.users.attach":        true,
//	"license.users.detach":        true,
//	"license.validate":            true,
//	"machine.check-out":           true,
//	"machine.create":              true,
//	"machine.delete":              true,
//	"machine.group.update":        true,
//	"machine.heartbeat.ping":      true,
//	"machine.heartbeat.reset":     true,
//	"machine.owner.update":        true,
//	"machine.proofs.generate":     true,
//	"machine.read":                true,
//	"machine.update":              true,
//	"metric.read":                 true,
//	"package.create":              true,
//	"package.delete":              true,
//	"package.read":                true,
//	"package.update":              true,
//	"platform.read":               true,
//	"policy.create":               true,
//	"policy.delete":               true,
//	"policy.entitlements.attach":  true,
//	"policy.entitlements.detach":  true,
//	"policy.pool.pop":             true,
//	"policy.read":                 true,
//	"policy.update":               true,
//	"process.create":              true,
//	"process.delete":              true,
//	"process.heartbeat.ping":      true,
//	"process.read":                true,
//	"process.update":              true,
//	"product.create":              true,
//	"product.delete":              true,
//	"product.read":                true,
//	"product.tokens.generate":     true,
//	"product.update":              true,
//	"release.constraints.attach":  true,
//	"release.constraints.detach":  true,
//	"release.create":              true,
//	"release.delete":              true,
//	"release.download":            true,
//	"release.package.update":      true,
//	"release.publish":             true,
//	"release.read":                true,
//	"release.update":              true,
//	"release.upgrade":             true,
//	"release.upload":              true,
//	"release.yank":                true,
//	"request-log.read":            true,
//	"token.generate":              true,
//	"token.read":                  true,
//	"token.regenerate":            true,
//	"token.revoke":                true,
//	"user.ban":                    true,
//	"user.create":                 true,
//	"user.delete":                 true,
//	"user.group.update":           true,
//	"user.invite":                 true,
//	"user.password.reset":         true,
//	"user.password.update":        true,
//	"user.read":                   true,
//	"user.second-factors.create":  true,
//	"user.second-factors.delete":  true,
//	"user.second-factors.read":    true,
//	"user.second-factors.update":  true,
//	"user.tokens.generate":        true,
//	"user.unban":                  true,
//	"user.update":                 true,
//}
//
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
