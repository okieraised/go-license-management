package entitlements

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type EntitlementRouter struct {
	tracer trace.Tracer
}

func NewEntitlementRouter() *EntitlementRouter {

	return &EntitlementRouter{}
}

func (r *EntitlementRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/entitlements")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:entitlement_id", r.retrieve)
		routes.PATCH("/:entitlement_id", r.update)
		routes.DELETE("/:entitlement_id", r.delete)
	}
}

// create creates a new entitlement resource.
func (r *EntitlementRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing entitlement.
func (r *EntitlementRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified entitlement resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *EntitlementRouter) update(ctx *gin.Context) {

}

// delete permanently deletes an entitlement.
// The entitlement will immediately be removed from all licenses and policies. It cannot be undone.
func (r *EntitlementRouter) delete(ctx *gin.Context) {

}

// list returns a list of entitlements. The entitlements are returned sorted by creation date,
// with the most recent entitlements appearing first. Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a license, only entitlements attached to that specific license will be listed.
func (r *EntitlementRouter) list(ctx *gin.Context) {

}
