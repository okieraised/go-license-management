package licenses

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type LicenseRouter struct {
	tracer trace.Tracer
}

func NewLicenseRouter() *LicenseRouter {

	return &LicenseRouter{}
}

func (r *LicenseRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/license")
		routes.POST("/generate", r.generate)
		routes.POST("/renew", r.renew)
		routes.POST("/lookup", r.lookup)
		routes.POST("/revoke", r.revoke)
		routes.POST("/delete", r.delete)
	}
}

func (r *LicenseRouter) generate(ctx *gin.Context) {

}

func (r *LicenseRouter) renew(ctx *gin.Context) {

}

func (r *LicenseRouter) lookup(ctx *gin.Context) {

}

func (r *LicenseRouter) revoke(ctx *gin.Context) {

}

func (r *LicenseRouter) delete(ctx *gin.Context) {

}
