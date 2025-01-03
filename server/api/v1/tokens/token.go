package tokens

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"go-license-management/internal/middlewares"
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
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(constants.UserCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(constants.UserCreate), r.list)
		routes.GET("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(constants.UserCreate), r.retrieve)
		routes.DELETE("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(constants.UserCreate), r.revoke)
		routes.PUT("/:token_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(constants.UserCreate), r.regenerate)
	}
}

// create Generate a new token resource for a user, using the user's email and password.
// Server does not store your tokens for security reasons. After a token is generated, it cannot be recovered if lost.
// The token will need to be revoked if lost, and a new token should be generated.
// Alternatively, the existing token can be regenerated (rolled).
func (r *TokenRouter) create(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing token.
func (r *TokenRouter) retrieve(ctx *gin.Context) {

}

// revoke permanently revokes a token. It cannot be undone.
func (r *TokenRouter) revoke(ctx *gin.Context) {

}

// list returns a list of tokens. The tokens are returned sorted by creation date,
// with the most recent tokens appearing first.
func (r *TokenRouter) list(ctx *gin.Context) {

}

// regenerate regenerates an existing token resource. This will replace the token attribute with a new secure token,
// and extend the token's expiry by 2 weeks from the current time.
func (r *TokenRouter) regenerate(ctx *gin.Context) {

}
