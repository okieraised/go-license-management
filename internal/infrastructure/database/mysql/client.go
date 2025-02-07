package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"time"
)

const (
	defaultTimeout          = 10 * time.Second
	defaultReadTimeout      = 10 * time.Second
	defaultWriteTimeout     = 10 * time.Second
	defaultMaxIdleConn  int = 100
	defaultMaxOpenConn  int = 100
)

var mysqlClient *bun.DB

func GetInstance() *bun.DB {
	return mysqlClient
}

func NewMysqlClient(host, port, dbname, userName, password string) (*bun.DB, error) {

	config := &mysql.Config{
		User:                 userName,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbname,
		AllowNativePasswords: true,
		Timeout:              defaultTimeout,
		ReadTimeout:          defaultReadTimeout,
		WriteTimeout:         defaultWriteTimeout,
	}

	mysqlDB, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	mysqlDB.SetMaxOpenConns(defaultMaxOpenConn)
	mysqlDB.SetMaxIdleConns(defaultMaxIdleConn)
	mysqlDB.SetConnMaxIdleTime(30 * time.Minute)
	mysqlDB.SetConnMaxLifetime(60 * time.Minute)

	mysqlClient = bun.NewDB(mysqlDB, mysqldialect.New())

	return mysqlClient, nil
}
