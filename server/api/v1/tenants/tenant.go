package tenants

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/tenants/service"
	"go-license-management/server/models/v1/tenants"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"net/http"
)

const (
	tenantGroup = "tenant_group"
)

type TenantRouter struct {
	svc    *service.TenantService
	logger *slog.Logger
	tracer trace.Tracer
}

func NewTenantRouter(svc *service.TenantService) *TenantRouter {
	tr := tracer.GetInstance().Tracer(tenantGroup)
	logger := logging.GetInstance().With(slog.Group(tenantGroup))
	return &TenantRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *TenantRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/accounts")
		routes.POST("", r.create)
		//routes.GET("", r.list)
		//routes.GET("/:account_id", r.retrieve)
		//routes.PATCH("/:account_id", r.update)
		//routes.DELETE("/:account_id", r.delete)
	}
}

// create creates a new tenant resource.
func (r *TenantRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.InfoContext(ctx, "received new account creation request", slog.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)))

	// serializer
	var req tenants.TenantRegistrationRequest
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
