package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
)

func main() {
	e, err := casbin.NewEnforcer("conf/rbac_model.conf", "conf/rbac_policy.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	subject := "alice"
	domain := "domain1"
	object := "data1"
	action := "write"

	// Check if the subject has permission
	ok, err := e.Enforce(domain, subject, object, action)
	if err != nil {
		fmt.Printf("Error checking permission: %v\n", err)
		return
	}
	if ok {
		fmt.Printf("Access granted for %s to %s %s in %s\n", subject, action, object, domain)
	} else {
		fmt.Printf("Access denied for %s to %s %s in %s\n", subject, action, object, domain)
	}
}
