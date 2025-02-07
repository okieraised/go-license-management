package sqlite

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var sqliteClient *bun.DB

func GetInstance() *bun.DB {
	return sqliteClient
}

func NewSqliteClient() (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	sqliteClient = bun.NewDB(sqldb, sqlitedialect.New())

	return sqliteClient, nil
}
