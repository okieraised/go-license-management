package database

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/utils"
	"time"
)

var client *bun.DB

const (
	defaultTimeout          = 10 * time.Second
	defaultReadTimeout      = 10 * time.Second
	defaultWriteTimeout     = 10 * time.Second
	defaultMaxIdleConn  int = 100
	defaultMaxOpenConn  int = 100
)

func GetInstance() *bun.DB {
	return client
}

func NewDatabaseClient(driver, host, port, dbname, userName, password string) (*bun.DB, error) {

	switch driver {
	case "postgres":
	case "mysql":
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
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

	superadminPassword := "superadmin"

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
