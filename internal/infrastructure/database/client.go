package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/utils"
	"time"
)

var client *bun.DB

const (
	defaultTimeout             = 10 * time.Second
	defaultReadTimeout         = 10 * time.Second
	defaultWriteTimeout        = 10 * time.Second
	defaultConnMaxIdleTime     = 30 * 60 * time.Second
	defaultConnMaxLifetime     = 60 * 60 * time.Second
	defaultMaxIdleConn     int = 100
	defaultMaxOpenConn     int = 100
)

func GetInstance() *bun.DB {
	return client
}

type DBConfig struct {
	driver       string
	host         string
	port         string
	dbName       string
	username     string
	password     string
	maxIdleConn  int
	maxOpenConn  int
	maxLifetime  time.Duration
	maxIdleTime  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	timeout      time.Duration
}

func WithDriver(driver string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.driver = driver
	}
}

func WithHost(host string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.host = host
	}
}

func WithPort(port string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.port = port
	}
}

func WithDBName(name string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.dbName = name
	}
}

func WithUsername(username string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.username = username
	}
}

func WithPassword(password string) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.password = password
	}
}

func WithMaxIdleConn(maxIdleConn int) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.maxIdleConn = maxIdleConn
	}
}

func WithMaxOpenConn(maxOpenConn int) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.maxOpenConn = maxOpenConn
	}
}

func WithMaxLifetime(maxLifetime time.Duration) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.maxLifetime = maxLifetime
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.maxIdleTime = maxIdleTime
	}
}

func WithReadTimeout(readTimeout time.Duration) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.readTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.writeTimeout = writeTimeout
	}
}

func WithTimeout(timeout time.Duration) func(*DBConfig) {
	return func(cfg *DBConfig) {
		cfg.timeout = timeout
	}
}

func NewDBConfig(opts ...func(*DBConfig)) *DBConfig {
	cfg := new(DBConfig)
	cfg.maxIdleConn = defaultMaxIdleConn
	cfg.maxOpenConn = defaultMaxOpenConn
	cfg.maxLifetime = defaultConnMaxLifetime
	cfg.maxIdleTime = defaultConnMaxIdleTime
	cfg.readTimeout = defaultReadTimeout
	cfg.writeTimeout = defaultWriteTimeout
	cfg.timeout = defaultTimeout

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func NewDatabaseClient(cfg *DBConfig) (*bun.DB, error) {

	if cfg.host == "" || cfg.username == "" || cfg.password == "" || cfg.dbName == "" {
		return nil, errors.New("one or more required connection parameters are empty")
	}

	switch cfg.driver {
	case "postgresql":
		pgconn := pgdriver.NewConnector(
			pgdriver.WithNetwork("tcp"),
			pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.host, cfg.port)),
			pgdriver.WithUser(cfg.username),
			pgdriver.WithPassword(cfg.password),
			pgdriver.WithDatabase(cfg.dbName),
			pgdriver.WithTimeout(cfg.timeout),
			pgdriver.WithReadTimeout(cfg.readTimeout),
			pgdriver.WithWriteTimeout(cfg.writeTimeout),
			pgdriver.WithInsecure(true),
		)
		postgresDB := sql.OpenDB(pgconn)
		postgresDB.SetMaxIdleConns(cfg.maxIdleConn)
		postgresDB.SetMaxOpenConns(cfg.maxOpenConn)
		postgresDB.SetConnMaxIdleTime(cfg.maxIdleTime)
		postgresDB.SetConnMaxLifetime(cfg.maxLifetime)
		client = bun.NewDB(postgresDB, pgdialect.New())

	case "mysql":
		mysqlCfg := &mysql.Config{
			User:                 cfg.username,
			Passwd:               cfg.password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%s", cfg.host, cfg.port),
			DBName:               cfg.dbName,
			AllowNativePasswords: true,
			Timeout:              cfg.timeout,
			ReadTimeout:          cfg.readTimeout,
			WriteTimeout:         cfg.writeTimeout,
		}

		mysqlDB, err := sql.Open("mysql", mysqlCfg.FormatDSN())
		if err != nil {
			return nil, err
		}
		mysqlDB.SetMaxOpenConns(cfg.maxOpenConn)
		mysqlDB.SetMaxIdleConns(cfg.maxIdleConn)
		mysqlDB.SetConnMaxIdleTime(cfg.maxIdleTime)
		mysqlDB.SetConnMaxLifetime(cfg.maxLifetime)
		client = bun.NewDB(mysqlDB, mysqldialect.New())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.driver)
	}

	return client, nil
}

func SeedingDatabase() error {
	logging.GetInstance().GetLogger().Info("started populating license database")
	roles := make([]entities.Role, 0)
	for k, _ := range constants.ValidRoleMapper {
		roles = append(roles, entities.Role{
			Name:      k,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	_, err := GetInstance().NewInsert().Model(&roles).On("CONFLICT DO NOTHING").Exec(context.Background())
	if err != nil {
		return err
	}

	superadminPassword := constants.DefaultSuperAdminPassword

	if viper.GetString(config.SuperAdminPassword) != "" {
		superadminPassword = viper.GetString(config.SuperAdminPassword)
	} else {
		viper.Set(config.SuperAdminPassword, superadminPassword)
	}

	digest, err := utils.HashPassword(superadminPassword)
	if err != nil {
		return err
	}
	privateKey, publicKey, err := utils.NewEd25519KeyPair()
	if err != nil {
		return err
	}
	superadmin := entities.Master{
		Username:          config.SuperAdminUsername,
		RoleName:          constants.RoleSuperAdmin,
		PasswordDigest:    digest,
		Ed25519PublicKey:  publicKey,
		Ed25519PrivateKey: privateKey,
	}

	_, err = GetInstance().NewInsert().Model(&superadmin).Exec(context.Background())
	if err != nil {
		return err
	}

	logging.GetInstance().GetLogger().Info("completed populating license database")
	return nil
}

func CreateSchemaIfNotExists() error {
	logging.GetInstance().GetLogger().Info("started initializing database schemas")
	_, err := GetInstance().NewDropTable().
		Model((*entities.Master)(nil)).IfExists().Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		IfNotExists().
		Model((*entities.Master)(nil)).
		WithForeignKeys().Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		IfNotExists().
		Model((*entities.Tenant)(nil)).
		WithForeignKeys().Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		Model((*entities.Role)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		Model((*entities.Account)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		ForeignKey(`("role_name") REFERENCES "roles" ("name") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		Model((*entities.Product)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().
		NewCreateTable().
		Model((*entities.ProductToken)(nil)).
		IfNotExists().
		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.Entitlement)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.Policy)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.PolicyEntitlement)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		ForeignKey(`("policy_id") REFERENCES "policies" ("id") ON DELETE CASCADE`).
		ForeignKey(`("entitlement_id") REFERENCES "entitlements" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.License)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		ForeignKey(`("policy_id") REFERENCES "policies" ("id") ON DELETE CASCADE`).
		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.Machine)(nil)).
		IfNotExists().
		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
		ForeignKey(`("license_id") REFERENCES "licenses" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = GetInstance().NewCreateTable().
		Model((*entities.Key)(nil)).
		IfNotExists().
		ForeignKey(`("policy_id") REFERENCES "policies" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	if err != nil {
		return err
	}
	logging.GetInstance().GetLogger().Info("completed initializing database schemas")

	return nil
}
