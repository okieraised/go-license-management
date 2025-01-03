package tenants

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/tenants/service"
	"go-license-management/internal/utils"
	"go-license-management/server/models/v1/tenants"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type TenantRouter struct {
	svc    *service.TenantService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewTenantRouter(svc *service.TenantService) *TenantRouter {
	tr := tracer.GetInstance().Tracer("tenant_group")
	logger := logging.NewECSLogger()
	return &TenantRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *TenantRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/tenants")
		routes.POST("", middlewares.JWTMasterValidationMW(), middlewares.PermissionValidationMW(constants.TenantCreate), r.create)
		routes.GET("", middlewares.JWTMasterValidationMW(), middlewares.PermissionValidationMW(constants.TenantRead), r.list)
		routes.GET("/:tenant_name", middlewares.JWTMasterValidationMW(), middlewares.PermissionValidationMW(constants.TenantRead), r.retrieve)
		routes.POST("/:tenant_name/regenerate", middlewares.JWTMasterValidationMW(), middlewares.PermissionValidationMW(constants.TenantUpdate), r.regenerate)
		routes.DELETE("/:tenant_name", middlewares.JWTMasterValidationMW(), middlewares.PermissionValidationMW(constants.TenantDelete), r.delete)
	}
}

// create creates a new tenant resource.
//
// @Summary 		API to register new tenant
// @Description 	Register new tenant
// @Tags 			tenant
// @Accept 			json
// @Produce 		json
// @Param 			Authorization 		header 		string 								true 	"authorization"
// @Param 			payload 			body 		tenants.TenantRegistrationRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants [post]
func (r *TenantRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new tenant creation request")

	// serializer
	var req tenants.TenantRegistrationRequest
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
	result, err := r.svc.Create(ctx, req.ToTenantRegistrationInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info(fmt.Sprintf("completed creating new tenant [%s]", utils.DerefPointer(req.Name)))
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

func (r *TenantRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new tenant list request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var req tenants.TenantListRequest
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
		r.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.ToResponse(comerrors.ErrCodeMapper[err], comerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.List(ctx, req.ToTenantListInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed retrieving tenants info")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

func (r *TenantRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new tenant retrieval request")

	// serializer
	var req tenants.TenantRetrievalRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&req)
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
	result, err := r.svc.Retrieve(ctx, req.ToTenantRetrievalInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed retrieving tenants info")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
	return
}

func (r *TenantRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new tenant deletion request")

	// serializer
	var req tenants.TenantDeletionRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&req)
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
	result, err := r.svc.Delete(ctx, req.ToTenantDeletionInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info(fmt.Sprintf("completed deleting tenant [%s]", utils.DerefPointer(req.TenantName)))
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusNoContent, resp)
	return
}

func (r *TenantRouter) regenerate(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new tenant regeneration request")

	// serializer
	var req tenants.TenantRegenerationRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&req)
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
	result, err := r.svc.Regenerate(ctx, req.ToTenantRegenerationInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info(fmt.Sprintf("completed updating tenant [%s]", utils.DerefPointer(req.TenantName)))
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
	return
}
