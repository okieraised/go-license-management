package models

import "github.com/uptrace/bun"

type DataSource struct {
	database *bun.DB
}

func (ds *DataSource) SetDatabase(db *bun.DB) {
	ds.database = db
}

func (ds *DataSource) GetDatabase() *bun.DB {
	return ds.database
}
