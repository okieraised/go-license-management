package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
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

	err = e.LoadPolicy()
	if err != nil {
		fmt.Println(err)
		return
	}

	//policies, err := e.GetPolicy()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(policies)

	//policies, err := e.GetFilteredPolicy(0, "user")
	//fmt.Println(policies)

	ok, err := e.Enforce("test", "user2", "product", "create")
	if err != nil {
		fmt.Printf("Error checking permission: %v\n", err)
		return
	}
	fmt.Println("result", ok)

	//// Modify the policy.
	//for _, record := range constants.CreateAdminPermission("test") {
	//	_, err := e.AddPolicy(record[1], record[2], record[3], record[4])
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	err = e.SavePolicy()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}

}
