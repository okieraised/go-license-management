package accounts

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/accounts/service"
	"go-license-management/server/models/v1/accounts"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"net/http"
)

const (
	accountGroup = "account_group"
)

type AccountRouter struct {
	svc    *service.AccountService
	logger *slog.Logger
	tracer trace.Tracer
}

func NewAccountRouter(svc *service.AccountService) *AccountRouter {
	tr := tracer.GetInstance().Tracer(accountGroup)
	logger := logging.GetInstance().With(slog.Group(accountGroup))
	return &AccountRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
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
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.InfoContext(ctx, "received new account creation request", slog.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	// serializer
	var req accounts.AccountCreateModelRequest
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

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	_, err = r.svc.Create(ctx, req.ToAccountRegistrationInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	ctx.JSON(http.StatusOK, resp)
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
