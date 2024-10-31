package users

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type UserRouter struct {
	tracer trace.Tracer
}

func NewAccountRouter() *UserRouter {

	return &UserRouter{}
}

func (r *UserRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/users")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:user_id", r.retrieve)
		routes.PATCH("/:user_id", r.update)
		routes.DELETE("/:user_id", r.delete)
		routes.POST("/:user_id/actions/:action_id", r.action)
	}
}

// create creates a new user resource. Users may be created with only an email address — no name or password is necessarily required.
// This can act as a way to associate an email address with a license, which can later be claimed and turned into a full user profile,
// if needed, using the password reset flow. This is particularly great for custom license recovery flows,
// where you may need to email a user their lost license keys.
func (r *UserRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing user.
func (r *UserRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified user resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *UserRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a user. It cannot be undone.
// This action also immediately deletes any licenses and machines that the user is associated with.
func (r *UserRouter) delete(ctx *gin.Context) {

}

// list returns a list of users. The users are returned sorted by creation date, with the most recent users appearing first.
// Resources are automatically scoped to the authenticated bearer e.g. when authenticated as a product,
// only users associated with the specific product, through a license, will be listed in the results.
func (r *UserRouter) list(ctx *gin.Context) {

}

// action actions for the user resource.
func (r *UserRouter) action(ctx *gin.Context) {

}
