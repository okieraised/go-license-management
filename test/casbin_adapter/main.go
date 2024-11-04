package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"go-license-management/internal/constants"
)

func main() {
	a, err := xormadapter.NewAdapter("postgres",
		"dbname=rbac_rules  user=postgres password=123qweA# host=127.0.0.1 port=5432 sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}

	e, err := casbin.NewEnforcer("conf/rbac_model.conf", a)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Modify the policy.
	for _, record := range constants.CreateAdminPermission("test") {
		_, err := e.AddPolicy(record[1], record[2], record[3], record[4])
		if err != nil {
			fmt.Println(err)
			return
		}

		err = e.SavePolicy()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
