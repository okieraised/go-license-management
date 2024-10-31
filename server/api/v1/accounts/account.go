package accounts

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type AccountRouter struct {
	tracer trace.Tracer
}

func NewEntitlementRouter() *AccountRouter {

	return &AccountRouter{}
}

func (r *AccountRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/accounts")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:account_id", r.retrieve)
		routes.PATCH("/:account_id", r.update)
		routes.DELETE("/:account_id", r.delete)
	}
}

// create creates a new account resource.
func (r *AccountRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing account.
func (r *AccountRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified account resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *AccountRouter) update(ctx *gin.Context) {

}

// delete permanently deletes an account. It cannot be undone.
func (r *AccountRouter) delete(ctx *gin.Context) {

}

// list returns a list of accounts. The accounts are returned sorted by creation date,
// with the most recent accounts appearing first
func (r *AccountRouter) list(ctx *gin.Context) {

}
