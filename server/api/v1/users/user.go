package users

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/server/models/v1/users"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"net/http"
)

const (
	userGroup = "user_group"
)

type UserRouter struct {
	logger *slog.Logger
	tracer trace.Tracer
}

func NewAccountRouter() *UserRouter {
	tr := tracer.GetInstance().Tracer(userGroup)
	logger := logging.GetInstance().With(slog.Group(userGroup))
	return &UserRouter{
		logger: logger,
		tracer: tr,
	}
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
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.InfoContext(ctx, "received new user creation request", slog.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	// serializer
	var req users.UserRegistrationRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.BindJSON(&req)
	if err != nil {
		cSpan.End()
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// validation
	_, cSpan = r.tracer.Start(rootCtx, "validation")
	err = req.Validate()
	if err != nil {
		cSpan.End()
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	ctx.JSON(http.StatusOK, resp)
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
