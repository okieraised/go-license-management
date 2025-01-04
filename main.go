package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/casbin_adapter"
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

	// tracer
	err := tracer.NewTracerProvider(
		viper.GetString(config.TracerURI),
		constants.AppName,
		"",
	)
	if err != nil {
		return dataSource, err
	}

	// database
	dataSource.SetDatabase(postgres.GetInstance())

	// casbin adapter
	dataSource.SetCasbin(casbin_adapter.GetAdapter())

	return dataSource, nil
}

func NewAppService(ds *api.DataSource) *api.AppService {
	appSvc := &api.AppService{}

	// register v1
	v1 := &api.V1AppService{}

	// tenant
	v1.SetTenant(tenantSvc.NewTenantService(tenantSvc.WithRepository(tenantRepo.NewTenantRepository(ds))))

	// auth
	v1.SetAuth(authSvc.NewAuthenticationService(authSvc.WithRepository(authRepo.NewAuthenticationRepository(ds))))

	// account
	v1.SetAccount(accountSvc.NewAccountService(
		accountSvc.WithRepository(accountRepo.NewAccountRepository(ds)),
		accountSvc.WithCasbinAdapter(ds.GetCasbin())),
	)

	// product
	v1.SetProduct(productSvc.NewProductService(productSvc.WithRepository(productRepo.NewProductRepository(ds))))

	// policy
	v1.SetPolicy(policySvc.NewPolicyService(policySvc.WithRepository(policyRepo.NewPolicyRepository(ds))))

	// entitlements
	v1.SetEntitlement(entitlementSvc.NewEntitlementService(entitlementSvc.WithRepository(entitlementRepo.NewEntitlementRepository(ds))))

	// licenses
	v1.SetLicense(licenseSvc.NewLicenseService(licenseSvc.WithRepository(licenseRepo.NewLicenseRepository(ds))))

	// machines
	v1.SetMachine(machineSvc.NewMachineService(machineSvc.WithRepository(machineRepo.NewMachineRepository(ds))))

	appSvc.SetV1Svc(v1)
	return appSvc
}

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

	go func() {
		server.StartServer(appSvc, serverQuit)
	}()

	<-quit
	serverQuit <- syscall.SIGKILL

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		logging.GetInstance().GetLogger().Info("app shutdown completed")
	}
}
