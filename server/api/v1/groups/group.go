package groups

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type GroupRouter struct {
	tracer trace.Tracer
}

func NewGroupRouter() *GroupRouter {

	return &GroupRouter{}
}

func (r *GroupRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/groups")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:group_id", r.retrieve)
		routes.PATCH("/:group_id", r.update)
		routes.DELETE("/:group_id", r.delete)
	}
}

// create creates a new group resource.
func (r *GroupRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing group.
func (r *GroupRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified group resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *GroupRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a group.
// The group will immediately be removed from all users, licenses and machines. It cannot be undone.
func (r *GroupRouter) delete(ctx *gin.Context) {

}

// list returns a list of groups. This will include all groups associated with the authenticated bearer,
// including groups they are an owner of, as well as groups they are a member of.
// The groups are returned sorted by creation date, with the most recent groups appearing first.
func (r *GroupRouter) list(ctx *gin.Context) {

}
