package constants

const (
	RoleSuperAdmin = iota
	RoleAdmin
	RoleProduct
	RoleLicense
	RoleUser
)

var RoleMapper = map[int]string{
	RoleSuperAdmin: "superadmin",
	RoleAdmin:      "admin",
	RoleProduct:    "product",
	RoleLicense:    "license",
	RoleUser:       "user",
}
