package licenses

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/server/v1/licenses/service"
	"go.opentelemetry.io/otel/trace"
)

const (
	licenseGroup = "license_group"
)

type LicenseRouter struct {
	svc    *service.LicenseService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewLicenseRouter(svc *service.LicenseService) *LicenseRouter {
	tr := tracer.GetInstance().Tracer(licenseGroup)
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
		routes.POST("", r.generate)
		routes.GET("/:license_id", r.retrieve)
		routes.PATCH("/:license_id", r.update)
		routes.DELETE("/:license_id", r.delete)
		routes.POST("/delete", r.delete)
		routes.GET("", r.list)
		routes.POST("/:license_id/actions/:action", r.action)

	}
}

// generate creates a new license resource.
func (r *LicenseRouter) generate(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing license.
func (r *LicenseRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified license resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *LicenseRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a license. It cannot be undone.
// This action also immediately deletes any machines that the license is associated with.
func (r *LicenseRouter) delete(ctx *gin.Context) {

}

// list returns a list of licenses. The licenses are returned sorted by creation date,
// with the most recent licenses appearing first.
// Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a user, only licenses of that specific user will be listed.
func (r *LicenseRouter) list(ctx *gin.Context) {

}

// action Actions for the license resource.
func (r *LicenseRouter) action(ctx *gin.Context) {

}
