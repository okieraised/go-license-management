package api

import (
	"github.com/gin-gonic/gin"
	"go-license-management/server/api/v1/accounts"
	"go-license-management/server/api/v1/authentications"
	"go-license-management/server/api/v1/policies"
	"go-license-management/server/api/v1/products"
	"go-license-management/server/api/v1/tenants"
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

		// tenant route
		tenantRoute := tenants.NewTenantRouter(rr.AppService.GetV1Svc().GetTenant())
		tenantRoute.Routes(v1Router, "")

		// common path prefix
		prefix := "tenants/:tenant_name"

		// Authentication routes
		authenRoute := authentications.NewAuthenticationRouter(rr.AppService.GetV1Svc().GetAuth())
		authenRoute.Routes(v1Router, prefix)

		// Account routes
		accountRoute := accounts.NewAccountRouter(rr.AppService.GetV1Svc().GetAccount())
		accountRoute.Routes(v1Router, prefix)

		// Product routes
		productRoute := products.NewProductRouter(rr.AppService.GetV1Svc().GetProduct())
		productRoute.Routes(v1Router, prefix)

		// Policy routes
		policyRoute := policies.NewPolicyRouter(rr.AppService.GetV1Svc().GetPolicy())
		policyRoute.Routes(v1Router, prefix)
	}

}
