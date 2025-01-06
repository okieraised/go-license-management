package api

import (
	"github.com/gin-gonic/gin"
	"go-license-management/server/api/v1/accounts"
	"go-license-management/server/api/v1/authentications"
	"go-license-management/server/api/v1/entitlements"
	"go-license-management/server/api/v1/licenses"
	"go-license-management/server/api/v1/machines"
	"go-license-management/server/api/v1/policies"
	"go-license-management/server/api/v1/products"
	"go-license-management/server/api/v1/tenants"
)

type RootRouter struct {
	AppService *AppService
}

func New(appService *AppService) *RootRouter {
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

		superAdminRoute := authentications.NewAuthenticationRouter(rr.AppService.GetV1Svc().GetAuth())
		superAdminRoute.Routes(v1Router, "")

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

		// Entitlement routes
		entitlementRoute := entitlements.NewEntitlementRouter(rr.AppService.GetV1Svc().GetEntitlement())
		entitlementRoute.Routes(v1Router, prefix)

		// Policy routes
		policyRoute := policies.NewPolicyRouter(rr.AppService.GetV1Svc().GetPolicy())
		policyRoute.Routes(v1Router, prefix)

		// License routes
		licenseRoute := licenses.NewLicenseRouter(rr.AppService.GetV1Svc().GetLicense())
		licenseRoute.Routes(v1Router, prefix)

		// Machine routes
		machineRoute := machines.NewMachineRouter(rr.AppService.GetV1Svc().GetMachine())
		machineRoute.Routes(v1Router, prefix)
	}

}
