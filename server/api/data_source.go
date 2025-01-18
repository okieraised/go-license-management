package api

import (
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/uptrace/bun"
)

type DataSource struct {
	database *bun.DB
	casbin   *xormadapter.Adapter
}

func (ds *DataSource) SetDatabase(db *bun.DB) {
	ds.database = db
}

func (ds *DataSource) GetDatabase() *bun.DB {
	return ds.database
}

func (ds *DataSource) SetCasbin(casbin *xormadapter.Adapter) {
	ds.casbin = casbin
}

func (ds *DataSource) GetCasbin() *xormadapter.Adapter {
	return ds.casbin
}
