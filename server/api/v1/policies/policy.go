package policies

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/permissions"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/policies/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type PolicyRouter struct {
	svc    *service.PolicyService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewPolicyRouter(svc *service.PolicyService) *PolicyRouter {
	tr := tracer.GetInstance().Tracer("policy_group")
	logger := logging.NewECSLogger()
	return &PolicyRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *PolicyRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/policies")
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyRead), r.list)
		routes.GET("/:policy_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyRead), r.retrieve)
		routes.PATCH("/:policy_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyUpdate), r.update)
		routes.DELETE("/:policy_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyDelete), r.delete)
		routes.POST("/:policy_id/entitlements", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyEntitlementsAttach), r.attach)
		routes.DELETE("/:policy_id/entitlements", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyEntitlementsDetach), r.detach)
		routes.GET("/:policy_id/entitlements", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.PolicyRead), r.listEntitlement)
	}
}

// create creates a new policy resource.
//
// @Summary 		API to register new policy resource
// @Description 	Creating new policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI 	true 	"path_param"
// @Param 			payload 			body 		policies.PolicyRegistrationRequest 	true 	"request"
// @Success 		201 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies [post]
func (r *PolicyRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy creation request")

	// serializer
	r.logger.GetLogger().Info("validating new policy creation request")
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq policy_attribute.PolicyCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyRegistrationRequest
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

	err = uriReq.Validate()
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
	result, err := r.svc.Create(ctx, bodyReq.ToPolicyRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrProductIDIsInvalid),
			errors.Is(err, cerrors.ErrPolicySchemeIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed creating new policy")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

// update updates the specified policy resource.
//
// @Summary 		API to update policy resource
// @Description 	Updating policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI     true 	"path_param"
// @Param 			payload 			body 		policies.PolicyUpdateRequest 	     true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id} [patch]
func (r *PolicyRouter) update(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received policy update request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	r.logger.GetLogger().Info("validating policy update request")
	var uriReq policy_attribute.PolicyCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyUpdateRequest
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

	err = uriReq.Validate()
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
	result, err := r.svc.Update(ctx, bodyReq.ToPolicyUpdateInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrProductIDIsInvalid),
			errors.Is(err, cerrors.ErrPolicySchemeIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed updating policy")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
	return
}

// retrieve retrieves the details of an existing policy.
//
// @Summary 		API to retrieve policy resource
// @Description 	Retrieving policy
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		policies.PolicyRetrievalRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id} [get]
func (r *PolicyRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy retrieval request")

	// serializer
	r.logger.GetLogger().Info("validating policy retrieval request")
	var req PolicyRetrievalRequest
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
	result, err := r.svc.Retrieve(ctx, req.ToPolicyRetrievalInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrProductIDIsInvalid):
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

// delete permanently deletes a policy. It cannot be undone.
// This action also immediately deletes any licenses that the policy is associated with.
//
// @Summary 		API to delete existing policy resource
// @Description 	Deleting existing policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		policies.PolicyDeletionRequest 	true 	"request"
// @Success 		204 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id} [delete]
func (r *PolicyRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy deletion request")

	// serializer
	r.logger.GetLogger().Info("validating policy deletion request")
	var req PolicyDeletionRequest
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
	result, err := r.svc.Delete(ctx, req.ToPolicyDeletionInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed deleting policy")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusNoContent, resp)
}

// list returns a list of policies. The policies are returned sorted by creation date, with the most recent policies appearing first.
//
// @Summary 		API to list policy resources
// @Description 	Listing policy resources
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI     true 	"path_param"
// @Param 			payload 			query 		policies.PolicyListRequest 	         true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies [get]
func (r *PolicyRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy history request")

	// serializer
	r.logger.GetLogger().Info("validating policy listing request")
	var uriReq policy_attribute.PolicyCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyListRequest
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
	result, err := r.svc.List(ctx, bodyReq.ToPolicyListInput(rootCtx, r.tracer, uriReq))
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

	r.logger.GetLogger().Info("completed listing policies")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

// attach attaches entitlements to a policy.
//
// @Summary 		API to attach entitlement to policy resource
// @Description 	Attaching entitlement to policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI     true 	"path_param"
// @Param 			payload 			body 		policies.PolicyAttachmentRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id}/entitlements [post]
func (r *PolicyRouter) attach(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy entitlement attachment request")

	// serializer
	r.logger.GetLogger().Info("validating policy attachment request")
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq policy_attribute.PolicyCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyAttachmentRequest
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
	result, err := r.svc.Attach(ctx, bodyReq.ToPolicyAttachmentInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		r.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrPolicyIDIsInvalid),
			errors.Is(err, cerrors.ErrEntitlementIDIsInvalid),
			errors.Is(err, cerrors.ErrPolicyEntitlementAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed attaching new entitlement to policy")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// detach detaches entitlements from a policy. This will immediately be taken into effect for all future license validations.
//
// @Summary 		API to detach entitlement from policy resource
// @Description 	Detaching entitlement from policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI     true 	"path_param"
// @Param 			payload 			body 		policies.PolicyDetachmentRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id}/entitlements [delete]
func (r *PolicyRouter) detach(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy entitlement detachment request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	r.logger.GetLogger().Info("validating policy detachment request")
	var uriReq policy_attribute.PolicyCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyDetachmentRequest
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
	result, err := r.svc.Detach(ctx, bodyReq.ToPolicyDetachmentInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		r.logger.GetLogger().Error(err.Error())
		cSpan.End()
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrPolicyIDIsInvalid),
			errors.Is(err, cerrors.ErrEntitlementIDIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed detaching entitlement from policy")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// listEntitlement returns a list of entitlements attached to the policy.
// The entitlements are returned sorted by creation date, with the most recent entitlements appearing first.
//
// @Summary 		API to list entitlements for policy resource
// @Description 	Listing entitlements for policy resource
// @Tags 			policy
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		policy_attribute.PolicyCommonURI     	true 	"request"
// @Param 			payload 			query 		policies.PolicyEntitlementListRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/policies/{policy_id}/entitlements [get]
func (r *PolicyRouter) listEntitlement(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new policy entitlement list request")

	// serializer
	r.logger.GetLogger().Info("validating policy entitlement listing request")
	var uriReq policy_attribute.PolicyCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq PolicyEntitlementListRequest
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
	result, err := r.svc.ListEntitlements(ctx, bodyReq.ToPolicyEntitlementListInput(rootCtx, r.tracer, uriReq))
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

	r.logger.GetLogger().Info("completed listing policy entitlements")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}
