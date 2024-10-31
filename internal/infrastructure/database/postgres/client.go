package postgres

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

const (
	defaultTimeout          = 10 * time.Second
	defaultReadTimeout      = 10 * time.Second
	defaultWriteTimeout     = 10 * time.Second
	defaultMaxIdleConn  int = 100
	defaultMaxOpenConn  int = 100
)

var postgresClient *bun.DB

func GetInstance() *bun.DB {
	return postgresClient
}

func NewPostgresClient(host, dbname, userName, password string) (*bun.DB, error) {

	if host == "" || userName == "" || password == "" || dbname == "" {
		return nil, errors.New("one or more required connection parameters are empty")
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(host),
		pgdriver.WithUser(userName),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(dbname),
		pgdriver.WithTimeout(defaultTimeout),
		pgdriver.WithReadTimeout(defaultReadTimeout),
		pgdriver.WithWriteTimeout(defaultWriteTimeout),
		pgdriver.WithInsecure(true),
	)
	postgresDB := sql.OpenDB(pgconn)
	postgresDB.SetMaxIdleConns(defaultMaxIdleConn)
	postgresDB.SetMaxOpenConns(defaultMaxOpenConn)
	postgresDB.SetConnMaxIdleTime(30 * time.Minute)
	postgresDB.SetConnMaxLifetime(60 * time.Minute)

	postgresClient = bun.NewDB(postgresDB, pgdialect.New())

	return postgresClient, nil
}