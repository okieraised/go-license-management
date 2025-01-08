package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-license-management/docs"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/middlewares"
	"go-license-management/server/api"
	"net/http"
	"os"
	"time"
)

// StartServer starts the API server
func StartServer(appService *api.AppService, quit chan os.Signal) {
	gin.SetMode(viper.GetString(config.ServerMode))
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{constants.AllowAllOrigins},
		AllowMethods: []string{http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodGet, http.MethodDelete},
		AllowHeaders: []string{constants.AccessControlAllowHeadersHeader, constants.OriginHeader, constants.AcceptHeader,
			constants.XRequestedWithHeader, constants.ContentTypeHeader, constants.AuthorizationHeader, constants.XAPIKeyHeader},
		ExposeHeaders:    []string{constants.ContentLengthHeader},
		AllowCredentials: true,
	}))

	router.Use(
		middlewares.RequestIDMW(), middlewares.TimeoutMW(), gzip.Gzip(gzip.DefaultCompression),
		middlewares.Recovery(), middlewares.LoggerMW(logging.GetInstance().GetLogger()), middlewares.HashHeaderMW(),
	)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rootRouter := api.New(appService)
	rootRouter.InitRouters(router)

	serverPort := "8888"
	if viper.GetString(config.ServerHttpPort) != "" {
		serverPort = viper.GetString(config.ServerHttpPort)
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", serverPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		var err error

		if viper.GetBool(config.ServerEnableTLS) {
			logging.GetInstance().GetLogger().Info("tls enabled")
			err = srv.ListenAndServeTLS(viper.GetString(config.ServerCertFile), viper.GetString(config.ServerKeyFile))
		} else {
			logging.GetInstance().GetLogger().Info("tls disabled")
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.GetInstance().GetLogger().Error(err.Error())
		}
	}()
	logging.GetInstance().GetLogger().Info(fmt.Sprintf("startup completed at: %s", serverAddr))

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.GetInstance().GetLogger().Error(fmt.Sprintf("error shutting down server: %s", err.Error()))
	}

	select {
	case <-ctx.Done():
		logging.GetInstance().GetLogger().Info("server shutdown completed")
	}
}
