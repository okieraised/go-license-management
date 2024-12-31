package casbin_adapter

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	_ "github.com/lib/pq"
)

func init() {
	enforcerModel = model.NewModel()
	enforcerModel.AddDef("r", "r", "dom, sub, obj, act")
	enforcerModel.AddDef("p", "p", "sub, obj, act")
	enforcerModel.AddDef("g", "g", "_, _, _")
	enforcerModel.AddDef("e", "e", "some(where (p.eft == allow)) && !some(where (p.eft == deny))")
	enforcerModel.AddDef("m", "m", "g(r.dom, r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == \"superadmin\"")
}

var enforcerModel model.Model

func GetEnforcerModel() model.Model {
	return enforcerModel
}

var adapter *xormadapter.Adapter

func GetAdapter() *xormadapter.Adapter {
	return adapter
}

func NewCasbinAdapter(userName, password, host, port string) (*xormadapter.Adapter, error) {
	var err error

	if host == "" || userName == "" || password == "" || port == "" {
		return nil, errors.New("one or more required connection parameters are empty")
	}

	adapter, err = xormadapter.NewAdapter(
		"postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
			userName, password, host, port,
		),
	)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}
