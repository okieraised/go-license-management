package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	_ "github.com/lib/pq"
)

func main() {
	a, err := xormadapter.NewAdapter("postgres", "user=postgres password=123qweA# host=127.0.0.1 port=5432 sslmode=disable")
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

	fmt.Println(e.GetAllDomains())

	policies, err := e.GetPolicy()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("policies", policies)

	ok, err := e.Enforce("test", "user1", "product", "create")
	if err != nil {
		fmt.Printf("Error checking permission: %v\n", err)
		return
	}
	fmt.Println("result", ok)

}
