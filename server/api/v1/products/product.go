package products

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/product_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/middlewares"
	"go-license-management/internal/permissions"
	"go-license-management/internal/response"
	"go-license-management/internal/services/v1/products/service"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/attribute"
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
		routes.POST("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductCreate), r.create)
		routes.GET("", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductRead), r.list)
		routes.GET("/:product_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductRead), r.retrieve)
		routes.PATCH("/:product_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductUpdate), r.update)
		routes.DELETE("/:product_id", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductDelete), r.delete)
		routes.POST("/:product_id/tokens", middlewares.JWTValidationMW(), middlewares.PermissionValidationMW(permissions.ProductTokensGenerate), r.tokens)
	}
}

// create creates a new product resource.
//
// @Summary 		API to register new product resource
// @Description 	Register new product
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		product_attribute.ProductCommonURI   	true 	"path_param"
// @Param 			payload 			body 		products.ProductRegistrationRequest 	true 	"request"
// @Success 		201 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products [post]
func (r *ProductRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product creation request")

	// serializer
	r.logger.GetLogger().Info("validating product registration request")
	var uriReq product_attribute.ProductCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq ProductRegistrationRequest
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
	result, err := r.svc.Create(ctx, bodyReq.ToProductRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrProductCodeAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed creating new product")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}

// retrieve retrieves the details of an existing product.
//
// @Summary 		API to retrieve existing product resource
// @Description 	Retrieving existing product resource
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		products.ProductRetrievalRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products/{product_id} [get]
func (r *ProductRouter) retrieve(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product retrieval request")

	// serializer
	r.logger.GetLogger().Info("validating product retrieval request")
	var req ProductRetrievalRequest
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
	result, err := r.svc.Retrieve(ctx, req.ToProductRetrievalInput(rootCtx, r.tracer))
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

	r.logger.GetLogger().Info("completed retrieving product info")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
}

// update updates the specified product resource.
//
// @Summary 		API to update product resource
// @Description 	Updating product
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		product_attribute.ProductCommonURI 	true 	"path_param"
// @Param 			payload 			body 		products.ProductUpdateRequest 	    true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products/{product_id} [patch]
func (r *ProductRouter) update(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product creation request")

	// serializer
	r.logger.GetLogger().Info("validating product update request")
	var uriReq product_attribute.ProductCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq ProductUpdateRequest
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
	result, err := r.svc.Update(ctx, bodyReq.ToProductUpdateInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, cerrors.ErrTenantNameIsInvalid),
			errors.Is(err, cerrors.ErrProductCodeAlreadyExist):
			ctx.JSON(http.StatusBadRequest, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info(fmt.Sprintf("completed updating product [%s]", utils.DerefPointer(uriReq.ProductID)))
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusOK, resp)
	return
}

// delete permanently deletes a product. It cannot be undone.
// This action also immediately deletes any policies, licenses and machines that the product is associated with.
//
// @Summary 		API to delete product resource
// @Description 	Delete product
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		products.ProductDeletionRequest 	true 	"request"
// @Success 		204 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products/{product_id} [delete]
func (r *ProductRouter) delete(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product deletion request")

	// serializer
	r.logger.GetLogger().Info("validating product deletion request")
	var req ProductDeletionRequest
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
	result, err := r.svc.Delete(ctx, req.ToProductDeletionInput(rootCtx, r.tracer))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	cSpan.End()

	r.logger.GetLogger().Info("completed deleting product")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusNoContent, resp)
}

// list returns a list of products. The products are returned sorted by creation date,
// with the most recent products appearing first.
//
// @Summary 		API to list product resources
// @Description 	Listing products
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			payload 			path 		products.ProductListRequest 	true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products [get]
func (r *ProductRouter) list(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product listing request")

	// serializer
	r.logger.GetLogger().Info("validating product listing request")
	var uriReq product_attribute.ProductCommonURI
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq ProductListRequest
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
	result, err := r.svc.List(ctx, bodyReq.ToProductListInput(rootCtx, r.tracer, uriReq))
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

	r.logger.GetLogger().Info("completed listing products")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, result.Count)
	ctx.JSON(http.StatusOK, resp)
	return
}

// tokens generates a new product token resource. Product tokens do not expire.
//
// @Summary 		API to generate product token resource
// @Description 	Generate product token
// @Tags 			product
// @Accept 			json
// @Produce 		json
// @Security        BearerAuth
// @Param 			param    			path 		product_attribute.ProductCommonURI 	true 	"path_param"
// @Param 			payload 			body 		products.ProductTokenRequest 	    true 	"request"
// @Success 		200 				{object} 	response.Response
// @Failure 		400 				{object} 	response.Response
// @Failure 		500 				{object} 	response.Response
// @Router 			/tenants/{tenant_name}/products/{product_id}/tokens [post]
func (r *ProductRouter) tokens(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path, trace.WithAttributes(attribute.KeyValue{
		Key:   constants.RequestIDField,
		Value: attribute.StringValue(ctx.GetString(constants.RequestIDField)),
	}))
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(
		zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
		zap.String(constants.ContextValueSubject, ctx.GetString(constants.ContextValueSubject)),
	).Info("received new product token creation request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq product_attribute.ProductCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericBadRequest], cerrors.ErrMessageMapper[cerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq ProductTokenRequest
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
	result, err := r.svc.Tokens(ctx, bodyReq.ToProductTokenInput(rootCtx, r.tracer, uriReq))
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

	r.logger.GetLogger().Info("completed creating product tokens")
	resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
	ctx.JSON(http.StatusCreated, resp)
	return
}
