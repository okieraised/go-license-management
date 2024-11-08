package postgres

import (
	"context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"testing"
	"time"
)

func TestNewPostgresClient(t *testing.T) {
	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	//dbClient.
}

func TestNewPostgresClient_CheckDBExists(t *testing.T) {
	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	err = CreateDatabase(context.Background(), viper.GetString(config.PostgresDatabase))
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateTenantSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Tenant)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateRoleSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Role)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)

	roles := make([]entities.Role, 0)
	for k, _ := range constants.ValidRoleMapper {
		roles = append(roles, entities.Role{
			Name:      k,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	_, err = dbClient.NewInsert().Model(&roles).Exec(context.Background())
	assert.NoError(t, err)

}

func TestNewPostgresClient_CreateAccountsSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Account)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateProductsSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Product)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateEntitlementsSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Entitlement)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreatePolicySchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Policy)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateLicenseSchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.License)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}

func TestNewPostgresClient_CreateKeySchema(t *testing.T) {

	viper.Set(config.PostgresHost, "127.0.0.1:5432")
	viper.Set(config.PostgresDatabase, "licenses")
	viper.Set(config.PostgresUsername, "postgres")
	viper.Set(config.PostgresPassword, "123qweA#")

	dbClient, err := NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dbClient)

	_, err = dbClient.NewCreateTable().Model((*entities.Key)(nil)).WithForeignKeys().Exec(context.Background())
	assert.NoError(t, err)
}
