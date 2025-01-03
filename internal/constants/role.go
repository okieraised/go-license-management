package constants

const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleProduct    = "product"
	RoleLicense    = "license"
	RoleUser       = "user"
)

var ValidRoleMapper = map[string]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleProduct:    true,
	RoleLicense:    true,
	RoleUser:       true,
}

var ValidAccountCreationRoleMapper = map[string]bool{
	RoleAdmin:   true,
	RoleProduct: true,
	RoleLicense: true,
	RoleUser:    true,
}
