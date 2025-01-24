package tokens

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/permissions"
	"go.opentelemetry.io/otel/trace"
)

type TokenRouter struct {
	tracer trace.Tracer
}

func NewTokenRouter() *TokenRouter {

	return &TokenRouter{}
}

func (r *TokenRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/tokens")
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.list)
		routes.GET("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.retrieve)
		routes.DELETE("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.revoke)
		routes.PUT("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.regenerate)
	}
}

// create Generate a new token resource for a user
func (r *TokenRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing token.
func (r *TokenRouter) retrieve(ctx *gin.Context) {

}

// revoke permanently revokes a token. It cannot be undone.
func (r *TokenRouter) revoke(ctx *gin.Context) {

}

// list returns a list of tokens.
// The tokens are returned sorted by creation date, with the most recent tokens appearing first.
func (r *TokenRouter) list(ctx *gin.Context) {

}

// regenerate regenerates an existing token resource.
func (r *TokenRouter) regenerate(ctx *gin.Context) {

}
