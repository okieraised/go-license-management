package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
)

func main() {
	a, err := xormadapter.NewAdapter("postgres",
		"dbname=rbac_rules  user=postgres password=postgres host=127.0.0.1 port=5432 sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}

	e, err := casbin.NewEnforcer("/Users/tripham/Desktop/go-license-management/conf/rbac_model.conf", a)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(e.GetPolicy())

	subject := "alice"  // the user who wants to access the resource
	domain := "domain1" // the domain in which access is requested
	object := "data1"   // the resource to access
	action := "read"    // the action the user wants to perform

	// Check if the subject has permission
	allowed, err := e.Enforce(domain, subject, object, action)
	if err != nil {
		fmt.Println(err)
		return
	}
	if allowed {
		fmt.Printf("Access granted for %s to %s %s in %s\n", subject, action, object, domain)
	} else {
		fmt.Printf("Access denied for %s to %s %s in %s\n", subject, action, object, domain)
	}

	// Modify the policy.
	//policy, err := e.AddPolicy(domain, subject, object, action)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(policy)
	//// e.RemovePolicy(...)
	//
	//// Save the policy back to DB.
	//err = e.SavePolicy()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}
