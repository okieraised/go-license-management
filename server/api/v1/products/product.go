package products

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/server/v1/products/service"
	"go.opentelemetry.io/otel/trace"
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

}

// retrieve retrieves the details of an existing product.
func (r *ProductRouter) retrieve(ctx *gin.Context) {

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

}
