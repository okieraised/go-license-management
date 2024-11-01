package api

import (
	"github.com/gin-gonic/gin"
	"go-license-management/server/api/v1/accounts"
	"go-license-management/server/models"
)

type RootRouter struct {
	AppService *models.AppService
}

func New(appService *models.AppService) *RootRouter {
	return &RootRouter{
		AppService: appService,
	}
}

func (rr *RootRouter) InitRouters(engine *gin.Engine) {
	// root
	rootRouter := engine.Group("/api")
	{
		v1Router := rootRouter.Group("/v1")

		// Account route
		accountRoute := accounts.NewAccountRouter(rr.AppService.GetV1Svc().GetAccount())
		accountRoute.Routes(v1Router, "")
	}

}
