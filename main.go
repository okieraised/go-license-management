package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-license-management/internal/config"
	"go-license-management/internal/infrastructure/database/postgres"
	"go-license-management/internal/infrastructure/logging"
	_ "go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	accountRepo "go-license-management/internal/repositories/v1/accounts"
	tenantRepo "go-license-management/internal/repositories/v1/tenants"
	accountSvc "go-license-management/internal/server/v1/accounts/service"
	tenantSvc "go-license-management/internal/server/v1/tenants/service"
	"go-license-management/server"
	"go-license-management/server/models"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Info(fmt.Sprintf("error reading config file, %s", err))
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
}

func newDataSource() (*models.DataSource, error) {
	dataSource := &models.DataSource{}

	// init logger
	logging.NewDefaultLogger()

	// tracer
	err := tracer.NewTracerProvider(
		viper.GetString(config.TracerURI),
		viper.GetString(""),
		"",
	)
	if err != nil {
		return dataSource, err
	}

	// database
	dbClient, err := postgres.NewPostgresClient(
		viper.GetString(config.PostgresHost),
		viper.GetString(config.PostgresDatabase),
		viper.GetString(config.PostgresUsername),
		viper.GetString(config.PostgresPassword),
	)
	if err != nil {
		return nil, err
	}

	dataSource.SetDatabase(dbClient)

	return dataSource, nil
}

func NewAppService(ds *models.DataSource) *models.AppService {
	appSvc := &models.AppService{}

	// register v1
	v1 := &models.V1AppService{}
	// tenant
	v1.SetTenant(tenantSvc.NewTenantService(tenantSvc.WithRepository(tenantRepo.NewTenantRepository(ds))))

	// account
	v1.SetAccount(accountSvc.NewAccountService(accountSvc.WithRepository(accountRepo.NewAccountRepository(ds))))

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
