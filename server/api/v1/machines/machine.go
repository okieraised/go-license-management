package machines

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/machines/service"
	"go-license-management/server/models/v1/machines"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

const (
	machineGroup = "machine_group"
)

type MachineRouter struct {
	svc    *service.MachineService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewMachineRouter(svc *service.MachineService) *MachineRouter {
	tr := tracer.GetInstance().Tracer(machineGroup)
	logger := logging.NewECSLogger()
	return &MachineRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *MachineRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/machines")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:machine_id", r.retrieve)
		routes.PATCH("/:machine_id", r.update)
		routes.DELETE("/:machine_id", r.deactivate)
		routes.POST("/:machine_id/actions/:action", r.action)
	}
}

// create creates, or activates, a new machine resource for a license.
func (r *MachineRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new machine activation request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq machine_attribute.MachineCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericBadRequest], comerrors.ErrMessageMapper[comerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq machines.MachineRegistrationRequest
	err = ctx.ShouldBind(&bodyReq)
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
	err = bodyReq.Validate()
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
	result, err := r.svc.Create(ctx, bodyReq.ToMachineRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid),
			errors.Is(err, comerrors.ErrMachineLicenseIsInvalid),
			errors.Is(err, comerrors.ErrMachineFingerprintAssociatedWithLicense),
			errors.Is(err, comerrors.ErrLicenseIsSuspended),
			errors.Is(err, comerrors.ErrLicenseIsBanned),
			errors.Is(err, comerrors.ErrLicenseIsExpired):
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

// retrieve retrieves the details of an existing machine.
func (r *MachineRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified machine resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *MachineRouter) update(ctx *gin.Context) {

}

// delete permanently deletes, or deactivates, a machine. It cannot be undone.
// This will immediately delete all processes and components associated with the machine.
func (r *MachineRouter) deactivate(ctx *gin.Context) {

}

// list returns a list of machines. The machines are returned sorted by creation date, with the most recent machines appearing first.
// Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a user, only machines for that specific user will be listed.
func (r *MachineRouter) list(ctx *gin.Context) {

}

// action actions to check out a machine. This will generate a snapshot of the machine at time of checkout,
// encoded into a machine file certificate that can be decoded and used for licensing offline and air-gapped environments.
// The algorithm will depend on the license policy's scheme.
func (r *MachineRouter) action(ctx *gin.Context) {

}
