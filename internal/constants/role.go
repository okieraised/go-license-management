package constants

const (
	// RoleSuperAdmin is the role with the highest privileges.
	RoleSuperAdmin = "superadmin"
	// RoleAdmin is an admin role with privileges to perform most action (except those related to tenant management).
	// Should only be used server-side
	RoleAdmin = "admin"
	// RoleUser is the default role when creating account. Can be used client-side to communicate with server
	RoleUser = "user"
)

var ValidRoleMapper = map[string]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleUser:       true,
}

var ValidAccountCreationRoleMapper = map[string]bool{
	RoleAdmin: true,
	RoleUser:  true,
}
