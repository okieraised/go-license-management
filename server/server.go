package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-license-management/internal/constants"
	"go-license-management/internal/middlewares"
	"go-license-management/server/models"
	"net/http"
	"os"
	"time"
)

func StartAPIServer(appService *models.AppService, quit chan os.Signal) {
	router := gin.New()
	router.Use(gin.Recovery())
	gin.SetMode("release")

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{constants.AllowAllOrigins},
		AllowMethods: []string{http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodGet, http.MethodDelete},
		AllowHeaders: []string{constants.AccessControlAllowHeadersHeader, constants.OriginHeader, constants.AcceptHeader,
			constants.XRequestedWithHeader, constants.ContentTypeHeader, constants.AuthorizationHeader, constants.XAPIKeyHeader},
		ExposeHeaders:    []string{constants.ContentLengthHeader},
		AllowCredentials: true,
	}))

	router.Use(middlewares.RequestIDMW(), middlewares.Recovery(), middlewares.TimeoutMW(), gzip.Gzip(gzip.DefaultCompression))

	serverAddr := "0.0.0.0:" + viper.GetString(comconstants.SERVER__HTTP_PORT)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.GetInstance().GetLogger().Fatal(err.Error())
		}
	}()
	logging.GetInstance().GetLogger().Info(fmt.Sprintf("startup completed at: %s", serverAddr))

	<-quit
	logging.GetInstance().GetLogger().Info("shutting down server in 5 seconds")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.GetInstance().GetLogger().Fatal(fmt.Sprintf("error shutting down server: %s", err.Error()))
	}

	select {
	case <-ctx.Done():
		logging.GetInstance().GetLogger().Info("server shutdown completed")
	}

}
