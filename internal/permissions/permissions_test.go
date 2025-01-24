package permissions

import (
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAdminPermission(t *testing.T) {
	adminPolicies := CreateAdminPermission()

	a, err := xormadapter.NewAdapter("postgres",
		"user=postgres password=123qweA# host=127.0.0.1 port=5432 sslmode=disable")
	assert.NoError(t, err)

	e, err := casbin.NewEnforcer("../../conf/rbac_model.conf", a)
	assert.NoError(t, err)

	err = e.LoadPolicy()
	assert.NoError(t, err)

	// Modify the policy.
	for _, record := range adminPolicies {
		_, err = e.AddPolicy(record[1], record[2], record[3])
		assert.NoError(t, err)

		err = e.SavePolicy()
		assert.NoError(t, err)
	}
}

func TestCreateSuperAdminPermission(t *testing.T) {
	superadminPolicies := CreateSuperAdminPermission()

	a, err := xormadapter.NewAdapter("postgres",
		"user=postgres password=123qweA# host=127.0.0.1 port=5432 sslmode=disable")
	assert.NoError(t, err)

	e, err := casbin.NewEnforcer("../../conf/rbac_model.conf", a)
	assert.NoError(t, err)

	err = e.LoadPolicy()
	assert.NoError(t, err)

	// Modify the policy.
	for _, record := range superadminPolicies {
		_, err = e.AddPolicy(record[1], record[2], record[3])
		assert.NoError(t, err)

		err = e.SavePolicy()
		assert.NoError(t, err)
	}
}

func TestCreateUserPermission(t *testing.T) {
	userPolicies := CreateUserPermission()

	a, err := xormadapter.NewAdapter("postgres",
		"user=postgres password=123qweA# host=127.0.0.1 port=5432 sslmode=disable")
	assert.NoError(t, err)

	e, err := casbin.NewEnforcer("../../conf/rbac_model.conf", a)
	assert.NoError(t, err)

	err = e.LoadPolicy()
	assert.NoError(t, err)

	// Modify the policy.
	for _, record := range userPolicies {
		_, err = e.AddPolicy(record[1], record[2], record[3])
		assert.NoError(t, err)

		err = e.SavePolicy()
		assert.NoError(t, err)
	}
}
