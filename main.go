package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
	_ "go-license-management/internal/logging"
	"go-license-management/server/models"
	"log/slog"
	"strings"
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

	e, err := casbin.NewEnforcer("conf/rbac_model.conf", "conf/rbac_policy.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	subject := "superadmin" // the user who wants to access the resource
	domain := "domain1"     // the domain in which access is requested
	object := "data1"       // the resource to access
	action := "write"       // the action the user wants to perform

	// Check if the subject has permission
	ok, err := e.Enforce(subject, domain, object, action)
	if err != nil {
		fmt.Printf("Error checking permission: %v\n", err)
		return
	}
	if ok {
		fmt.Printf("Access granted for %s to %s %s in %s\n", subject, action, object, domain)
	} else {
		fmt.Printf("Access denied for %s to %s %s in %s\n", subject, action, object, domain)
	}

	//quit := make(chan os.Signal)
	//serverQuit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	//
	//dataSources, err := newDataSource()
	//if err != nil {
	//	logging.GetInstance().Error(err.Error())
	//	return
	//}
	//appSvc := NewAppService(dataSources)
	//
	//go func() {
	//	server.StartServer(appSvc, serverQuit)
	//}()
	//
	//<-quit
	//serverQuit <- syscall.SIGKILL
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	//
	//select {
	//case <-ctx.Done():
	//	slog.Info("app shutdown completed")
	//}
}
