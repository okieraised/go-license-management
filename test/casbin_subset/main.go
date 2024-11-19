package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func main() {
	enforcer, err := casbin.NewEnforcer()
	if err != nil {
		fmt.Println(err)
		return
	}

	adapter := fileadapter.NewFilteredAdapter("conf/rbac_policy.csv")
	err = enforcer.InitWithAdapter("conf/rbac_model.conf", adapter)
	if err != nil {
		fmt.Println(err)
		return
	}

	e, err := casbin.NewEnforcer("conf/rbac_model.conf", "conf/rbac_policy.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(e.GetPolicy())

	filter := &fileadapter.Filter{
		P: []string{"admin"},
		//G: []string{"", "", "admin"},
	}
	err = enforcer.LoadFilteredPolicy(filter)
	if err != nil {
		fmt.Println(err)
		return
	}
}
