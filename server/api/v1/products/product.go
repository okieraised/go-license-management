package products

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/products/service"
	"go-license-management/server/models/v1/products"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type ProductRouter struct {
	svc    *service.ProductService
	logger *logging.Logger
	tracer trace.Tracer
}

const (
	productGroup = "product_group"
)

func NewProductRouter(svc *service.ProductService) *ProductRouter {
	tr := tracer.GetInstance().Tracer(productGroup)
	logger := logging.NewECSLogger()
	return &ProductRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *ProductRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/products")
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:product_id", r.retrieve)
		routes.PATCH("/:product_id", r.update)
		routes.DELETE("/:product_id", r.delete)
		routes.POST("/:product_id/tokens", r.tokens)
	}
}

// create creates a new product resource.
func (r *ProductRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new product creation request")

	// serializer
	tenantName := ctx.Param("tenant_name")
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req products.ProductRegistrationRequest
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
	result, err := r.svc.Create(ctx, req.ToProductRegistrationInput(rootCtx, r.tracer, tenantName))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid),
			errors.Is(err, comerrors.ErrProductCodeAlreadyExist):
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

// retrieve retrieves the details of an existing product.
func (r *ProductRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new accounts retrieval request")

	// serializer
	var req products.ProductRetrievalRequest
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
	result, err := r.svc.Retrieve(ctx, req.ToProductRetrievalInput(rootCtx, r.tracer))
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

// update updates the specified product resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *ProductRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a product. It cannot be undone.
// This action also immediately deletes any policies, licenses and machines that the product is associated with.
func (r *ProductRouter) delete(ctx *gin.Context) {

}

// list returns a list of products. The products are returned sorted by creation date,
// with the most recent products appearing first.
func (r *ProductRouter) list(ctx *gin.Context) {

}

// tokens generates a new product token resource. Product tokens do not expire.
func (r *ProductRouter) tokens(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new product token creation request")

	// serializer
	tenantName := ctx.Param(constants.TenantNameField)
	if tenantName == "" {
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrTenantNameIsEmpty], comerrors.ErrMessageMapper[comerrors.ErrTenantNameIsEmpty], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	productID := ctx.Param("product_id")
	_, err := uuid.Parse(productID)
	if err != nil {
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrProductIDIsInvalid], comerrors.ErrMessageMapper[comerrors.ErrProductIDIsInvalid], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req products.ProductTokenRequest
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err = ctx.ShouldBind(&req)
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
	result, err := r.svc.Tokens(ctx, req.ToProductTokenInput(rootCtx, r.tracer, tenantName, productID))
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

	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return

}
