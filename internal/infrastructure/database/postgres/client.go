package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/utils"
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

func NewPostgresClient(host, port, dbname, userName, password string) (*bun.DB, error) {

	if host == "" || userName == "" || password == "" || dbname == "" {
		return nil, errors.New("one or more required connection parameters are empty")
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(host+":"+port),
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

	superadminUsername := "superadmin"
	superadminPassword := "superadmin"

	if viper.GetString(config.SuperAdminUsername) != "" {
		superadminUsername = viper.GetString(config.SuperAdminUsername)
	}

	if viper.GetString(config.SuperAdminPassword) != "" {
		superadminPassword = viper.GetString(config.SuperAdminPassword)
	}

	digest, err := utils.HashPassword(superadminPassword)
	if err != nil {
		return err
	}

	superadmin := entities.Master{
		Username:       superadminUsername,
		RoleName:       constants.RoleSuperAdmin,
		PasswordDigest: digest,
	}

	_, err = GetInstance().NewInsert().Model(&superadmin).Exec(context.Background())
	if err != nil {
		return err
	}

	logging.GetInstance().GetLogger().Info("completed populating license database")
	return nil
}

func checkDatabaseExists(ctx context.Context, dbName string) (bool, error) {
	var exists bool
	query := "SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower(?);"
	err := postgresClient.NewRaw(query, dbName).Scan(ctx, &exists)
	return exists, err
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
