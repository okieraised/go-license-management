package licenses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/license_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/permissions"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/licenses/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type LicenseRouter struct {
	svc    *service.LicenseService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewLicenseRouter(svc *service.LicenseService) *LicenseRouter {
	tr := tracer.GetInstance().Tracer("license_group")
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
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.LicenseCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.LicenseRead), r.list)
		routes.GET("/:license_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.LicenseRead), r.retrieve)
		routes.PATCH("/:license_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.LicenseUpdate), r.update)
		routes.DELETE("/:license_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.LicenseDelete), r.delete)
		routes.POST("/actions/:action", middlewares.JWTValidationMW(), middlewares.LicenseActionPermissionValidationMW(), r.action)
	}
}

// create creates a new license resource.
//
// @Summary 		API to register new license resource
// @Description 	Register new license
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param 			    path 		license_attribute.LicenseCommonURI   	true 	"path_param"
// @Param 			payload 			body 		licenses.LicenseRegistrationRequest 	true 	"request"
// @Success 		201 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses [post]
func (r *LicenseRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new license creation request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	r.logger.GetLogger().Info("validating license creation request")
	var uriReq license_attribute.LicenseCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq LicenseRegistrationRequest
	err = ctx.ShouldBind(&bodyReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
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
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Create(ctx, bodyReq.ToLicenseRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrPolicyIDIsInvalid),
			errors.Is(err, cerrors.ErrProductIDIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed generating new license")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

// retrieve retrieves the details of an existing license.
//
// @Summary 		API to retrieve license resource
// @Description 	Retrieve license
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		licenses.LicenseRetrievalRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses/{license_id} [get]
func (r *LicenseRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new license retrieval request")

	// serializer
	r.logger.GetLogger().Info("validating license retrieval request")
	var req LicenseRetrievalRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
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
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Retrieve(ctx, req.ToLicenseRetrievalInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed retrieval request")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// update updates the specified license resource.
//
// @Summary 		API to update license resource
// @Description 	Updating license
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param 			    path 		license_attribute.LicenseCommonURI   	true 	"path_param"
// @Param 			payload 			body 		licenses.LicenseUpdateRequest 	        true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses/{license_id} [patch]
func (r *LicenseRouter) update(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received license update request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	r.logger.GetLogger().Info("validating license update request")
	var uriReq license_attribute.LicenseCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq LicenseUpdateRequest
	err = ctx.ShouldBind(&bodyReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
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
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Update(ctx, bodyReq.ToLicenseUpdateInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrPolicyIDIsInvalid),
			errors.Is(err, cerrors.ErrProductIDIsInvalid),
			errors.Is(err, cerrors.ErrLicenseIDIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed updating license")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return

}

// delete permanently deletes a license. It cannot be undone.
// This action also immediately deletes any machines that the license is associated with.
//
// @Summary 		API to delete license resource
// @Description 	Deleting license
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		licenses.LicenseDeletionRequest 	true 	"request"
// @Success 		204 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses/{license_id} [delete]
func (r *LicenseRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new license deletion request")

	// serializer
	var req LicenseDeletionRequest
	r.logger.GetLogger().Info("validating license deletion request")
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
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
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Delete(ctx, req.ToLicenseDeletionInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusNoContent, resp)
}

// list returns a list of licenses. The licenses are returned sorted by creation date,
// with the most recent licenses appearing first.
//
// @Summary 		API to list existing license resources
// @Description 	Listing existing license resources
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param 			    path 		license_attribute.LicenseCommonURI   	true 	"path_param"
// @Param 			payload 			query 		licenses.LicenseListRequest 	        true 	"request"
// @Success 		204 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses [get]
func (r *LicenseRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new license history request")

	// serializer
	r.logger.GetLogger().Info("validating license listing request")
	var uriReq license_attribute.LicenseCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq LicenseListRequest
	err = ctx.ShouldBind(&bodyReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
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
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.List(ctx, bodyReq.ToLicenseListInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed listing licenses")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

// action Actions for the license resource.
//
//   - Validate: Action to validate a license key. This will look up the license by its key and check the following:
//     if the license is suspended, if the license is expired, if the license is overdue for check-in, and
//     if the license meets its machine requirements (if strict).
//   - suspend: Action to temporarily suspend (ban) a license. This will cause the license to fail validation until reinstated.
//   - reinstate: Action to reinstate a suspended license.
//   - renew: Action to renew a license. Extends license expiry by the policy's duration, according to the policy's renewal basis.
//   - checkout: Action to check out a license. This will generate a snapshot of the license at time of checkout,
//     encoded into a license file certificate that can be decoded and used for licensing offline and air-gapped
//     environments. The algorithm will depend on the policy's scheme.
//   - checkin: Action to check in a license. Sets the license's lastCheckIn to the current time
//   - increment-usage: Action to increment a license's uses attribute in accordance with its policy's maxUses attribute.
//     When the policy's maxUses limit is exceeded, the increment attempt will fail. When the policy's maxUses is
//     set to null, there is no limit on usage.
//   - decrement-usage: Action to decrement a license's uses attribute in accordance with its policy's maxUses attribute.
//   - reset-usage: Action to reset a license's uses attribute to 0.
//
// @Summary 		API to perform action on license resource
// @Description 	Performing action on license resource
// @Tags 			license
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param 			    path 		license_attribute.LicenseCommonURI   	true 	"path_param"
// @Param 			payload 			body 		licenses.LicenseActionsRequest 	        true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/licenses/actions/{action} [post]
func (r *LicenseRouter) action(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new license action request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	r.logger.GetLogger().Info("validating license action request")
	var uriReq license_attribute.LicenseCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq LicenseActionsRequest
	err = ctx.ShouldBind(&bodyReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	cSpan.End()

	// validation
	_, cSpan = r.tracer.Start(rootCtx, "validation")
	err = uriReq.Validate()
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	err = bodyReq.Validate()
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Actions(ctx, bodyReq.ToLicenseActionsInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrGenericInternalServer):
			ctx.JSON(http.StatusInternalServerError, resp)
		default:
			ctx.JSON(http.StatusBadRequest, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed handling license actions")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}
