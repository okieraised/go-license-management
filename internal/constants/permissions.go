package constants

import (
	"strings"
)

const (
	AccountAnalyticsRead      = "account.analytics.read"
	AccountBillingRead        = "account.billing.read"
	AccountBillingUpdate      = "account.billing.update"
	AccountPlanRead           = "account.plan.read"
	AccountPlanUpdate         = "account.plan.update"
	AccountRead               = "account.read"
	AccountSubscriptionRead   = "account.subscription.read"
	AccountSubscriptionUpdate = "account.subscription.update"
	AccountUpdate             = "account.update"
	AdminCreate               = "admin.create"
	AdminDelete               = "admin.delete"
	AdminInvite               = "admin.invite"
	AdminRead                 = "admin.read"
	AdminUpdate               = "admin.update"
	ArchRead                  = "arch.read"
	ArtifactCreate            = "artifact.create"
	ArtifactDelete            = "artifact.delete"
	ArtifactRead              = "artifact.read"
	ArtifactUpdate            = "artifact.update"
	ChannelRead               = "channel.read"
	ConstraintRead            = "constraint.read"
	EngineRead                = "engine.read"
	EntitlementCreate         = "entitlement.create"
	EntitlementDelete         = "entitlement.delete"
	EntitlementRead           = "entitlement.read"
	EntitlementUpdate         = "entitlement.update"
	EnvironmentCreate         = "environment.create"
	EnvironmentDelete         = "environment.delete"
	EnvironmentRead           = "environment.read"
	EnvironmentTokensGenerate = "environment.tokens.generate"
	EnvironmentUpdate         = "environment.update"
	EventLogRead              = "event-log.read"
	GroupCreate               = "group.create"
	GroupDelete               = "group.delete"
	GroupLicensesRead         = "group.licenses.read"
	GroupMachinesRead         = "group.machines.read"
	GroupOwnersAttach         = "group.owners.attach"
	GroupOwnersDetach         = "group.owners.detach"
	GroupOwnersRead           = "group.owners.read"
	GroupRead                 = "group.read"
	GroupUpdate               = "group.update"
	GroupUsersRead            = "group.users.read"
	KeyCreate                 = "key.create"
	KeyDelete                 = "key.delete"
	KeyRead                   = "key.read"
	KeyUpdate                 = "key.update"
	LicenseCheckIn            = "license.check-in"
	LicenseCheckOut           = "license.check-out"
	LicenseCreate             = "license.create"
	LicenseDelete             = "license.delete"
	LicenseEntitlementsAttach = "license.entitlements.attach"
	LicenseEntitlementsDetach = "license.entitlements.detach"
	LicenseGroupUpdate        = "license.group.update"
	LicenseOwnerUpdate        = "license.owner.update"
	LicensePolicyUpdate       = "license.policy.update"
	LicenseRead               = "license.read"
	LicenseReinstate          = "license.reinstate"
	LicenseRenew              = "license.renew"
	LicenseRevoke             = "license.revoke"
	LicenseSuspend            = "license.suspend"
	LicenseTokensGenerate     = "license.tokens.generate"
	LicenseUpdate             = "license.update"
	LicenseUsageDecrement     = "license.usage.decrement"
	LicenseUsageIncrement     = "license.usage.increment"
	LicenseUsageReset         = "license.usage.reset"
	LicenseUsersAttach        = "license.users.attach"
	LicenseUsersDetach        = "license.users.detach"
	LicenseValidate           = "license.validate"
	MachineCheckOut           = "machine.check-out"
	MachineCreate             = "machine.create"
	MachineDelete             = "machine.delete"
	MachineGroupUpdate        = "machine.group.update"
	MachineHeartbeatPing      = "machine.heartbeat.ping"
	MachineHeartbeatReset     = "machine.heartbeat.reset"
	MachineOwnerUpdate        = "machine.owner.update"
	MachineProofsGenerate     = "machine.proofs.generate"
	MachineRead               = "machine.read"
	MachineUpdate             = "machine.update"
	MetricRead                = "metric.read"
	PackageCreate             = "package.create"
	PackageDelete             = "package.delete"
	PackageRead               = "package.read"
	PackageUpdate             = "package.update"
	PlatformRead              = "platform.read"
	PolicyCreate              = "policy.create"
	PolicyDelete              = "policy.delete"
	PolicyEntitlementsAttach  = "policy.entitlements.attach"
	PolicyEntitlementsDetach  = "policy.entitlements.detach"
	PolicyPoolPop             = "policy.pool.pop"
	PolicyRead                = "policy.read"
	PolicyUpdate              = "policy.update"
	ProcessCreate             = "process.create"
	ProcessDelete             = "process.delete"
	ProcessHeartbeatPing      = "process.heartbeat.ping"
	ProcessRead               = "process.read"
	ProcessUpdate             = "process.update"
	ProductCreate             = "product.create"
	ProductDelete             = "product.delete"
	ProductRead               = "product.read"
	ProductTokensGenerate     = "product.tokens.generate"
	ProductUpdate             = "product.update"
	ReleaseConstraintsAttach  = "release.constraints.attach"
	ReleaseConstraintsDetach  = "release.constraints.detach"
	ReleaseCreate             = "release.create"
	ReleaseDelete             = "release.delete"
	ReleaseDownload           = "release.download"
	ReleasePackageUpdate      = "release.package.update"
	ReleasePublish            = "release.publish"
	ReleaseRead               = "release.read"
	ReleaseUpdate             = "release.update"
	ReleaseUpgrade            = "release.upgrade"
	ReleaseUpload             = "release.upload"
	ReleaseYank               = "release.yank"
	RequestLogRead            = "request-log.read"
	TokenGenerate             = "token.generate"
	TokenRead                 = "token.read"
	TokenRegenerate           = "token.regenerate"
	TokenRevoke               = "token.revoke"
	UserBan                   = "user.ban"
	UserCreate                = "user.create"
	UserDelete                = "user.delete"
	UserGroupUpdate           = "user.group.update"
	UserInvite                = "user.invite"
	UserPasswordReset         = "user.password.reset"
	UserPasswordUpdate        = "user.password.update"
	UserRead                  = "user.read"
	UserSecondFactorsCreate   = "user.second-factors.create"
	UserSecondFactorsDelete   = "user.second-factors.delete"
	UserSecondFactorsRead     = "user.second-factors.read"
	UserSecondFactorsUpdate   = "user.second-factors.update"
	UserTokensGenerate        = "user.tokens.generate"
	UserUnban                 = "user.unban"
	UserUpdate                = "user.update"
)

var AdminPermissions = map[string]bool{
	"account.analytics.read":      true,
	"account.billing.read":        true,
	"account.billing.update":      true,
	"account.plan.read":           true,
	"account.plan.update":         true,
	"account.read":                true,
	"account.subscription.read":   true,
	"account.subscription.update": true,
	"account.update":              true,
	"admin.create":                true,
	"admin.delete":                true,
	"admin.invite":                true,
	"admin.read":                  true,
	"admin.update":                true,
	"arch.read":                   true,
	"artifact.create":             true,
	"artifact.delete":             true,
	"artifact.read":               true,
	"artifact.update":             true,
	"channel.read":                true,
	"constraint.read":             true,
	"engine.read":                 true,
	"entitlement.create":          true,
	"entitlement.delete":          true,
	"entitlement.read":            true,
	"entitlement.update":          true,
	"environment.create":          true,
	"environment.delete":          true,
	"environment.read":            true,
	"environment.tokens.generate": true,
	"environment.update":          true,
	"event-log.read":              true,
	"group.create":                true,
	"group.delete":                true,
	"group.licenses.read":         true,
	"group.machines.read":         true,
	"group.owners.attach":         true,
	"group.owners.detach":         true,
	"group.owners.read":           true,
	"group.read":                  true,
	"group.update":                true,
	"group.users.read":            true,
	"key.create":                  true,
	"key.delete":                  true,
	"key.read":                    true,
	"key.update":                  true,
	"license.check-in":            true,
	"license.check-out":           true,
	"license.create":              true,
	"license.delete":              true,
	"license.entitlements.attach": true,
	"license.entitlements.detach": true,
	"license.group.update":        true,
	"license.owner.update":        true,
	"license.policy.update":       true,
	"license.read":                true,
	"license.reinstate":           true,
	"license.renew":               true,
	"license.revoke":              true,
	"license.suspend":             true,
	"license.tokens.generate":     true,
	"license.update":              true,
	"license.usage.decrement":     true,
	"license.usage.increment":     true,
	"license.usage.reset":         true,
	"license.users.attach":        true,
	"license.users.detach":        true,
	"license.validate":            true,
	"machine.check-out":           true,
	"machine.create":              true,
	"machine.delete":              true,
	"machine.group.update":        true,
	"machine.heartbeat.ping":      true,
	"machine.heartbeat.reset":     true,
	"machine.owner.update":        true,
	"machine.proofs.generate":     true,
	"machine.read":                true,
	"machine.update":              true,
	"metric.read":                 true,
	"package.create":              true,
	"package.delete":              true,
	"package.read":                true,
	"package.update":              true,
	"platform.read":               true,
	"policy.create":               true,
	"policy.delete":               true,
	"policy.entitlements.attach":  true,
	"policy.entitlements.detach":  true,
	"policy.pool.pop":             true,
	"policy.read":                 true,
	"policy.update":               true,
	"process.create":              true,
	"process.delete":              true,
	"process.heartbeat.ping":      true,
	"process.read":                true,
	"process.update":              true,
	"product.create":              true,
	"product.delete":              true,
	"product.read":                true,
	"product.tokens.generate":     true,
	"product.update":              true,
	"release.constraints.attach":  true,
	"release.constraints.detach":  true,
	"release.create":              true,
	"release.delete":              true,
	"release.download":            true,
	"release.package.update":      true,
	"release.publish":             true,
	"release.read":                true,
	"release.update":              true,
	"release.upgrade":             true,
	"release.upload":              true,
	"release.yank":                true,
	"request-log.read":            true,
	"token.generate":              true,
	"token.read":                  true,
	"token.regenerate":            true,
	"token.revoke":                true,
	"user.ban":                    true,
	"user.create":                 true,
	"user.delete":                 true,
	"user.group.update":           true,
	"user.invite":                 true,
	"user.password.reset":         true,
	"user.password.update":        true,
	"user.read":                   true,
	"user.second-factors.create":  true,
	"user.second-factors.delete":  true,
	"user.second-factors.read":    true,
	"user.second-factors.update":  true,
	"user.tokens.generate":        true,
	"user.unban":                  true,
	"user.update":                 true,
}

func CreateAdminPermission(domain string) [][]string {

	result := make([][]string, 0)
	for key, _ := range AdminPermissions {
		parts := strings.Split(key, ".")

		object := ""
		perm := ""

		if len(parts) == 3 {
			object = parts[0] + "_" + parts[1]
			perm = parts[2]
		} else {
			object = parts[0]
			perm = parts[1]
		}

		result = append(result, []string{"p", domain, "admin", object, perm})
	}

	return result
}
