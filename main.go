package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/casbin_adapter"
	"go-license-management/internal/infrastructure/database/postgres"
	"go-license-management/internal/infrastructure/logging"
	_ "go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	accountRepo "go-license-management/internal/repositories/v1/accounts"
	authRepo "go-license-management/internal/repositories/v1/authentications"
	entitlementRepo "go-license-management/internal/repositories/v1/entitlements"
	licenseRepo "go-license-management/internal/repositories/v1/licenses"
	machineRepo "go-license-management/internal/repositories/v1/machines"
	policyRepo "go-license-management/internal/repositories/v1/policies"
	productRepo "go-license-management/internal/repositories/v1/products"
	tenantRepo "go-license-management/internal/repositories/v1/tenants"
	accountSvc "go-license-management/internal/services/v1/accounts/service"
	authSvc "go-license-management/internal/services/v1/authentications/service"
	entitlementSvc "go-license-management/internal/services/v1/entitlements/service"
	licenseSvc "go-license-management/internal/services/v1/licenses/service"
	machineSvc "go-license-management/internal/services/v1/machines/service"
	policySvc "go-license-management/internal/services/v1/policies/service"
	productSvc "go-license-management/internal/services/v1/products/service"
	tenantSvc "go-license-management/internal/services/v1/tenants/service"
	"go-license-management/server"
	"go-license-management/server/api"
	"go-license-management/server/api/v1"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	// init logger
	logging.NewDefaultLogger()

	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("error reading config file: %v", err))
		os.Exit(1)
	}
	logging.GetInstance().GetLogger().Info("successfully read config file")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	// Seeding database
	_, err = postgres.NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresPort),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("failed to initialize license database connection: %v", err))
		os.Exit(1)
	}

	err = postgres.CreateSchemaIfNotExists()
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("failed to initialize license database schemas: %v", err))
		os.Exit(1)
	}

	// Seeding roles and superadmin user
	err = postgres.SeedingDatabase()
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("failed to initialize license initial data: %v", err))
		os.Exit(1)
	}

	// Seeding Casbin permissions
	_, err = casbin_adapter.NewCasbinAdapter(viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresPort),
	)
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("failed to initialize casbin: %v", err))
		os.Exit(1)
	}

	err = casbin_adapter.SeedingCasbinPermissions()
	if err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
		os.Exit(1)
	}
}

func newDataSource() (*api.DataSource, error) {
	dataSource := &api.DataSource{}

	// initialize tracer
	err := tracer.NewTracerProvider(
		viper.GetString(config.TracerURI),
		constants.AppName,
		"",
	)
	if err != nil {
		return dataSource, err
	}

	// initialize database
	dataSource.SetDatabase(postgres.GetInstance())

	// initialize casbin adapter
	dataSource.SetCasbin(casbin_adapter.GetAdapter())

	return dataSource, nil
}

func NewAppService(ds *api.DataSource) *api.AppService {
	appSvc := &api.AppService{}

	// register v1
	v1Svc := &v1.V1AppService{}

	// tenant
	v1Svc.SetTenant(tenantSvc.NewTenantService(tenantSvc.WithRepository(tenantRepo.NewTenantRepository(ds))))

	// auth
	v1Svc.SetAuth(authSvc.NewAuthenticationService(authSvc.WithRepository(authRepo.NewAuthenticationRepository(ds))))

	// account
	v1Svc.SetAccount(accountSvc.NewAccountService(
		accountSvc.WithRepository(accountRepo.NewAccountRepository(ds)),
		accountSvc.WithCasbinAdapter(ds.GetCasbin())),
	)

	// product
	v1Svc.SetProduct(productSvc.NewProductService(productSvc.WithRepository(productRepo.NewProductRepository(ds))))

	// policy
	v1Svc.SetPolicy(policySvc.NewPolicyService(policySvc.WithRepository(policyRepo.NewPolicyRepository(ds))))

	// entitlements
	v1Svc.SetEntitlement(entitlementSvc.NewEntitlementService(entitlementSvc.WithRepository(entitlementRepo.NewEntitlementRepository(ds))))

	// licenses
	v1Svc.SetLicense(licenseSvc.NewLicenseService(licenseSvc.WithRepository(licenseRepo.NewLicenseRepository(ds))))

	// machines
	v1Svc.SetMachine(machineSvc.NewMachineService(machineSvc.WithRepository(machineRepo.NewMachineRepository(ds))))

	appSvc.SetV1Svc(v1Svc)
	return appSvc
}

// @title           Go License Management API
// @version         0.1.0
// @description     Go License Management Server.
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8888
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	quit := make(chan os.Signal)
	serverQuit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	dataSources, err := newDataSource()
	if err != nil {
		logging.GetInstance().GetLogger().Error(err.Error())
		return
	}
	appSvc := NewAppService(dataSources)

	dfs := func() {
		server.StartServer(appSvc, serverQuit)
	}

	go dfs()

	//monigoInstance := &monigo.Monigo{
	//	ServiceName:             "data-api", // Mandatory field
	//	DashboardPort:           8080,       // Default is 8080
	//	DataPointsSyncFrequency: "5s",       // Default is 5 Minutes
	//	DataRetentionPeriod:     "4d",       // Default is 7 days. Supported values: "1h", "1d", "1w", "1m"
	//	TimeZone:                "Local",    // Default is Local timezone. Supported values: "Local", "UTC", "Asia/Kolkata", "America/New_York" etc. (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
	//}
	//
	//go monigo.TraceFunction(dfs)
	//
	//go monigoInstance.Start()

	<-quit
	serverQuit <- syscall.SIGKILL

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		logging.GetInstance().GetLogger().Info("app shutdown completed")
	}
}
