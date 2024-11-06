package licenses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/licenses/service"
	"go-license-management/server/models/v1/license"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

const (
	licenseGroup = "license_group"
)

type LicenseRouter struct {
	svc    *service.LicenseService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewLicenseRouter(svc *service.LicenseService) *LicenseRouter {
	tr := tracer.GetInstance().Tracer(licenseGroup)
	logger := logging.NewECSLogger()
	return &LicenseRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *LicenseRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/licenses")
		routes.POST("", r.generate)
		routes.GET("/:license_id", r.retrieve)
		routes.PATCH("/:license_id", r.update)
		routes.DELETE("/:license_id", r.delete)
		routes.POST("/delete", r.delete)
		routes.GET("", r.list)
		routes.POST("/:license_id/actions/:action", r.action)

	}
}

// generate creates a new license resource.
func (r *LicenseRouter) generate(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new license creation request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req license.LicenseRegistrationRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBind(&req)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericBadRequest], comerrors.ErrMessageMapper[comerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// validation
	_, cSpan = r.tracer.Start(rootCtx, "validation")
	err = req.Validate()
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[err], comerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Create(ctx, req.ToLicenseRegistrationInput(rootCtx, r.tracer, tenantName))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid), errors.Is(err, comerrors.ErrEntitlementCodeAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

// retrieve retrieves the details of an existing license.
func (r *LicenseRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified license resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *LicenseRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a license. It cannot be undone.
// This action also immediately deletes any machines that the license is associated with.
func (r *LicenseRouter) delete(ctx *gin.Context) {

}

// list returns a list of licenses. The licenses are returned sorted by creation date,
// with the most recent licenses appearing first.
// Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a user, only licenses of that specific user will be listed.
func (r *LicenseRouter) list(ctx *gin.Context) {

}

// action Actions for the license resource.
func (r *LicenseRouter) action(ctx *gin.Context) {

}
