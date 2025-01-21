package authentications

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/authentication_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/authentications/service"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type AuthenticationRouter struct {
	svc    *service.AuthenticationService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewAuthenticationRouter(svc *service.AuthenticationService) *AuthenticationRouter {
	tr := tracer.GetInstance().Tracer("auth_group")
	logger := logging.NewECSLogger()
	return &AuthenticationRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *AuthenticationRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/auth")
		routes.POST("/login", r.login)
	}
}

// login validates existing account resource.
//
// @Summary 		API to validate existing account and return a corresponding jwt token
// @Description 	Validating account and generate a JWT token if valid, without tenant_name path parameter, one must provide the superadmin credentials
// @Tags 			authentication
// @Accept 			mpfd
// @Produce 		json
// @Param 			username 			formData 	string 					true 	"username"
// @Param 			password 			formData 	string 					true 	"password"
// @Param        	tenant_name    	    path     	string  				true  	"tenant_name"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/auth/login [post]
// @Router 			/auth/login [post]
func (r *AuthenticationRouter) login(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new account login request")

	// serializer
	var uriReq authentication_attribute.AuthenticationCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq AuthenticationLoginRequest
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

	// Super admin user can login with all paths
	if utils.DerefPointer(bodyReq.Username) != config.SuperAdminUsername {
		err = uriReq.Validate()
		if err != nil {
			cSpan.End()
			r.logger.GetLogger().Error(err.Error())
			resp.ToResponse(cerrors.ErrCodeMapper[err], cerrors.ErrMessageMapper[err], nil, nil, nil)
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	}
	cSpan.End()

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Login(ctx, bodyReq.ToAuthenticationLoginInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrGenericUnauthorized),
			errors.Is(err, cerrors.ErrAccountIsBanned),
			errors.Is(err, cerrors.ErrAccountIsInactive):
			ctx.JSON(http.StatusUnauthorized, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
	return
}
