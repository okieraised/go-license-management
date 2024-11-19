package accounts

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/accounts/models"
	"go-license-management/internal/server/v1/accounts/service"
	"go-license-management/internal/utils"
	"go-license-management/server/models/v1/accounts"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

const (
	accountGroup = "account_group"
)

type AccountRouter struct {
	svc    *service.AccountService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewAccountRouter(svc *service.AccountService) *AccountRouter {
	tr := tracer.GetInstance().Tracer(accountGroup)
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
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes = routes.Group("/:account_name")
		routes.GET("", r.retrieve)
		routes.PATCH("", r.update)
		routes.DELETE("", r.delete)
		routes.POST("/actions/:action_name", r.actions)

	}
}

// create creates a new account resource.
func (r *AccountRouter) create(ctx *gin.Context) {

	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new account creation request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req accounts.AccountRegistrationRequest
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
	result, err := r.svc.Create(ctx, req.ToAccountRegistrationInput(rootCtx, r.tracer, tenantName))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid), errors.Is(err, comerrors.ErrAccountUsernameAlreadyExist):
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

// retrieve retrieves the details of an existing account.
func (r *AccountRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new accounts retrieval request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req accounts.AccountRetrievalRequest
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

	// handler
	_, cSpan = r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.Retrieve(ctx, req.ToAccountRetrievalInput(rootCtx, r.tracer, tenantName))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// update updates the specified account resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *AccountRouter) update(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new accounts deletion request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	username := ctx.Param("account_name")
	if username == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrAccountUsernameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrAccountUsernameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req accounts.AccountUpdateRequest
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
	result, err := r.svc.Update(ctx, req.ToAccountUpdateInput(rootCtx, r.tracer, tenantName, username))
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

// delete permanently deletes an account. It cannot be undone.
func (r *AccountRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new accounts deletion request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req accounts.AccountDeletionRequest
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
	result, err := r.svc.Delete(ctx, req.ToAccountDeletionInput(rootCtx, r.tracer, tenantName))
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
// with the most recent accounts appearing first
func (r *AccountRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new accounts list request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// handler
	_, cSpan := r.tracer.Start(rootCtx, "handler")
	result, err := r.svc.List(ctx, &models.AccountListInput{
		TracerCtx:  rootCtx,
		Tracer:     r.tracer,
		TenantName: utils.RefPointer(tenantName),
	})
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

func (r *AccountRouter) actions(ctx *gin.Context) {}
