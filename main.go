package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-license-management/internal/logging"
	_ "go-license-management/internal/logging"
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
	return dataSource, nil
}

func NewAppService(ds *models.DataSource) *models.AppService {
	appSvc := &models.AppService{}
	return appSvc
}

func main() {
	quit := make(chan os.Signal)
	serverQuit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	dataSources, err := newDataSource()
	if err != nil {
		logging.GetInstance().Error(err.Error())
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
		slog.Info("app shutdown completed")
	}
}
