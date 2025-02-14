package accounts

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/account_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/permissions"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/accounts/service"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type AccountRouter struct {
	svc    *service.AccountService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewAccountRouter(svc *service.AccountService) *AccountRouter {
	tr := tracer.GetInstance().Tracer("account_group")
	logger := logging.NewECSLogger()
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
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserRead), r.list)
		routes = routes.Group("/:username")
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserRead), r.retrieve)
		routes.PATCH("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserUpdate), r.update)
		routes.DELETE("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.UserDelete), r.delete)
		routes.POST("/actions/:action", middlewares.JWTValidationMW(), middlewares.AccountActionPermissionValidationMW(), r.actions)

	}
}

// create creates a new account resource.
//
// @Summary 		API to register new account resource
// @Description 	Register new account resource
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param        	tenant_name    	    path     	string  				true  	"tenant_name"
// @Param 			payload 			body 		accounts.AccountRegistrationRequest 	true 	"request"
// @Success 		201 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts [post]
func (r *AccountRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new account creation request")

	// serializer
	var uriReq account_attribute.AccountCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq AccountRegistrationRequest
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
	result, err := r.svc.Create(ctx, bodyReq.ToAccountRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrAccountUsernameAlreadyExist),
			errors.Is(err, cerrors.ErrAccountEmailAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info(fmt.Sprintf("completed creating new account [%s]", utils.DerefPointer(bodyReq.Username)))
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

// retrieve retrieves the details of an existing account.
//
// @Summary 		API to retrieve existing account
// @Description 	Retrieving account
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		accounts.AccountRetrievalRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts/{username} [get]
func (r *AccountRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new accounts retrieval request")

	// serializer
	var req AccountRetrievalRequest
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

	if req.Username == nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountUsernameIsEmpty], cerrors.ErrMessageMapper[cerrors.ErrAccountUsernameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Retrieve(ctx, req.ToAccountRetrievalInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed retrieving account info")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// update updates the specified account resource.
//
// @Summary 		API to retrieve existing account
// @Description 	Retrieving account
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		account_attribute.AccountCommonURI	 	true 	"path_param"
// @Param 			payload 			body 		accounts.AccountUpdateRequest	        true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts/{username} [patch]
func (r *AccountRouter) update(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new accounts update request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq account_attribute.AccountCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq AccountUpdateRequest
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

	if uriReq.Username == nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountUsernameIsEmpty], cerrors.ErrMessageMapper[cerrors.ErrAccountUsernameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	cSpan.End()

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
	result, err := r.svc.Update(ctx, bodyReq.ToAccountUpdateInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrAccountUsernameIsInvalid), errors.Is(err, cerrors.ErrTenantNameIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed updating account info")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// delete permanently deletes an account.
//
// @Summary 		API to delete existing account
// @Description 	Deleting account
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		accounts.AccountDeletionRequest	 	true 	"request"
// @Success 		204 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts/{username} [delete]
func (r *AccountRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new accounts deletion request")

	// serializer
	var req AccountDeletionRequest
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
	result, err := r.svc.Delete(ctx, req.ToAccountDeletionInput(rootCtx, r.tracer))
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

// list returns a list of accounts. The accounts are returned sorted by creation date,
// with the most recent accounts appearing first.
//
// @Summary 		API to list existing account resources
// @Description 	Listing account resources
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param        	tenant_name    	    path     	string  				        true  	"tenant_name"
// @Param 			payload 			query 		accounts.AccountListRequest	 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts [get]
func (r *AccountRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new accounts list request")

	// serializer
	var uriReq account_attribute.AccountCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq AccountListRequest
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
	result, err := r.svc.List(ctx, bodyReq.ToAccountListInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed retrieval request")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

// actions performs account action
//
// @Summary 		API to perform action on account resource
// @Description 	Performing actions on account resource
// @Tags 			account
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param        	param       	    path     	account_attribute.AccountCommonURI  true  	"path_param"
// @Param 			payload 			body 		accounts.AccountActionRequest	 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/accounts/{username}/actions/{action} [post]
func (r *AccountRouter) actions(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new accounts action request")

	// serializer
	var uriReq account_attribute.AccountCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if uriReq.Action == nil {
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountActionIsEmpty], cerrors.ErrMessageMapper[cerrors.ErrAccountActionIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq AccountActionRequest
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

	if utils.DerefPointer(uriReq.Action) == constants.AccountActionUpdatePassword {
		if bodyReq.CurrentPassword == nil {
			resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountCurrentPasswordIsEmpty],
				cerrors.ErrMessageMapper[cerrors.ErrAccountCurrentPasswordIsEmpty], nil, nil, nil)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}

		if bodyReq.NewPassword == nil {
			resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountNewPasswordIsEmpty],
				cerrors.ErrMessageMapper[cerrors.ErrAccountNewPasswordIsEmpty], nil, nil, nil)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	} else if utils.DerefPointer(uriReq.Action) == constants.AccountActionResetPassword {
		if bodyReq.NewPassword == nil {
			resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountNewPasswordIsEmpty],
				cerrors.ErrMessageMapper[cerrors.ErrAccountNewPasswordIsEmpty], nil, nil, nil)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}

		if bodyReq.ResetToken == nil {
			resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrAccountResetTokenIsEmpty],
				cerrors.ErrMessageMapper[cerrors.ErrAccountResetTokenIsEmpty], nil, nil, nil)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Action(ctx, bodyReq.ToAccountActionInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrAccountUsernameIsInvalid),
			errors.Is(err, cerrors.ErrAccountPasswordNotMatch),
			errors.Is(err, cerrors.ErrAccountResetTokenIsInvalid):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed account action request")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}
