package postgres

//func TestNewPostgresClient(t *testing.T) {
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	//dbClient.
//}
//
//func TestNewPostgresClient_CreateTenantSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().
//		IfNotExists().
//		Model((*entities.Tenant)(nil)).
//		WithForeignKeys().Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateRoleSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Role)(nil)).WithForeignKeys().Exec(context.Background())
//	assert.NoError(t, err)
//
//	roles := make([]entities.Role, 0)
//	for k, _ := range constants.ValidRoleMapper {
//		roles = append(roles, entities.Role{
//			Name:      k,
//			CreatedAt: time.Now(),
//			UpdatedAt: time.Now(),
//		})
//	}
//
//	_, err = dbClient.NewInsert().Model(&roles).Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateAccountsSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Account)(nil)).
//		IfNotExists().
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		ForeignKey(`("role_name") REFERENCES "roles" ("name") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateProductsSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Product)(nil)).
//		IfNotExists().
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateProductsToken(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.ProductToken)(nil)).
//		IfNotExists().
//		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateEntitlementsSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Entitlement)(nil)).
//		IfNotExists().
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreatePolicySchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Policy)(nil)).
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreatePolicyEntitlementSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.PolicyEntitlement)(nil)).
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		ForeignKey(`("policy_id") REFERENCES "policies" ("id") ON DELETE CASCADE`).
//		ForeignKey(`("entitlement_id") REFERENCES "entitlements" ("id") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateLicenseSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.License)(nil)).
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		ForeignKey(`("policy_id") REFERENCES "policies" ("id") ON DELETE CASCADE`).
//		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
////func TestNewPostgresClient_CreateLicenseTokenSchema(t *testing.T) {
////
////	viper.Set(config.PostgresHost, "127.0.0.1")
////	viper.Set(config.PostgresPort, "5432")
////	viper.Set(config.PostgresDatabase, "licenses")
////	viper.Set(config.PostgresUsername, "postgres")
////	viper.Set(config.PostgresPassword, "123qweA#")
////
////	dbClient, err := NewPostgresClient(
////		viper.GetString(config.PostgresHost),
////		viper.GetString(config.PostgresPort),
////		viper.GetString(config.PostgresDatabase),
////		viper.GetString(config.PostgresUsername),
////		viper.GetString(config.PostgresPassword),
////	)
////	assert.NoError(t, err)
////	assert.NotNil(t, dbClient)
////
////	_, err = dbClient.NewCreateTable().Model((*entities.LicenseToken)(nil)).
////		ForeignKey(`("license_id") REFERENCES "licenses" ("id") ON DELETE CASCADE`).
////		Exec(context.Background())
////	assert.NoError(t, err)
////}
//
//func TestNewPostgresClient_CreateKeySchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Key)(nil)).WithForeignKeys().Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestNewPostgresClient_CreateMachineSchema(t *testing.T) {
//
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	_, err = dbClient.NewCreateTable().Model((*entities.Machine)(nil)).
//		ForeignKey(`("tenant_name") REFERENCES "tenants" ("name") ON DELETE CASCADE`).
//		ForeignKey(`("license_id") REFERENCES "licenses" ("id") ON DELETE CASCADE`).
//		Exec(context.Background())
//	assert.NoError(t, err)
//}
//
//func TestGetInstance(t *testing.T) {
//	viper.Set(config.PostgresHost, "127.0.0.1")
//	viper.Set(config.PostgresPort, "5432")
//	viper.Set(config.PostgresDatabase, "licenses")
//	viper.Set(config.PostgresUsername, "postgres")
//	viper.Set(config.PostgresPassword, "123qweA#")
//
//	dbClient, err := NewPostgresClient(
//		viper.GetString(config.PostgresHost),
//		viper.GetString(config.PostgresPort),
//		viper.GetString(config.PostgresDatabase),
//		viper.GetString(config.PostgresUsername),
//		viper.GetString(config.PostgresPassword),
//	)
//	assert.NoError(t, err)
//	assert.NotNil(t, dbClient)
//
//	licenses := make([]entities.License, 0)
//	err = dbClient.NewSelect().Model(new(entities.License)).Relation("Policy").Relation("Product").Scan(context.Background(), &licenses)
//	assert.NoError(t, err)
//
//	fmt.Println(licenses[0].Product)
//}
