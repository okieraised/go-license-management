package policies

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type PolicyRouter struct {
	tracer trace.Tracer
}

func NewPolicyRouter() *PolicyRouter {

	return &PolicyRouter{}
}

func (r *PolicyRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/groups")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:policy_id", r.retrieve)
		routes.PATCH("/:policy_id", r.update)
		routes.DELETE("/:policy_id", r.delete)
		routes.DELETE("/:policy_id", r.delete)
		routes.POST("/:policy_id/entitlements", r.attach)
		routes.DELETE("/:policy_id/entitlements", r.detach)
		routes.GET("/:policy_id/entitlements", r.listEntitlement)
	}
}

// create creates a new policy resource.
func (r *PolicyRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing policy.
func (r *PolicyRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified policy resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *PolicyRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a policy. It cannot be undone.
// This action also immediately deletes any licenses that the policy is associated with.
func (r *PolicyRouter) delete(ctx *gin.Context) {

}

// list returns a list of policies. The policies are returned sorted by creation date, with the most recent policies
// appearing first. Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a product, only policies of that specific product will be listed.
func (r *PolicyRouter) list(ctx *gin.Context) {

}

// attach attaches entitlements to a policy. This will immediately be taken into effect for all future license validations.
// Any license that implements the given policy will automatically possess all the policy's entitlements.
func (r *PolicyRouter) attach(ctx *gin.Context) {

}

// detach detaches entitlements from a policy. This will immediately be taken into effect for all future license validations.
func (r *PolicyRouter) detach(ctx *gin.Context) {

}

// listEntitlement returns a list of entitlements attached to the policy.
// The entitlements are returned sorted by creation date, with the most recent entitlements appearing first.
func (r *PolicyRouter) listEntitlement(ctx *gin.Context) {

}
