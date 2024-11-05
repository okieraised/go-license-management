package machines

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/server/v1/machines/service"
	"go.opentelemetry.io/otel/trace"
)

const (
	machineGroup = "machine_group"
)

type MachineRouter struct {
	svc    *service.MachineService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewMachineRouter(svc *service.MachineService) *MachineRouter {
	tr := tracer.GetInstance().Tracer(machineGroup)
	logger := logging.NewECSLogger()
	return &MachineRouter{
		svc:    svc,
		logger: logger,
		tracer: tr,
	}
}

func (r *MachineRouter) Routes(engine *gin.RouterGroup, path string) {
	routes := engine.Group(path)
	{
		routes = routes.Group("/machines")
		routes.POST("/activate", r.activate)
		routes.GET("", r.list)
		routes.GET("/:machine_id", r.retrieve)
		routes.PATCH("/:machine_id", r.update)
		routes.DELETE("/:machine_id", r.deactivate)
		routes.POST("/:machine_id/actions/:action", r.action)
		routes.PUT("/:machine_id/owner", r.changeOwner)
		routes.PUT("/:machine_id/group", r.changeGroup)
	}
}

// activate creates, or activates, a new machine resource for a license.
func (r *MachineRouter) activate(ctx *gin.Context) {

}

// retrieve retrieves the details of an existing machine.
func (r *MachineRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified machine resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *MachineRouter) update(ctx *gin.Context) {

}

// delete permanently deletes, or deactivates, a machine. It cannot be undone.
// This will immediately delete all processes and components associated with the machine.
func (r *MachineRouter) deactivate(ctx *gin.Context) {

}

// list returns a list of machines. The machines are returned sorted by creation date, with the most recent machines appearing first.
// Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a user, only machines for that specific user will be listed.
func (r *MachineRouter) list(ctx *gin.Context) {

}

// action actions to check out a machine. This will generate a snapshot of the machine at time of checkout,
// encoded into a machine file certificate that can be decoded and used for licensing offline and air-gapped environments.
// The algorithm will depend on the license policy's scheme.
func (r *MachineRouter) action(ctx *gin.Context) {

}

// changeOwner changes a machine's owner relationship. This will immediately transfer the machine resource to the new owner.
func (r *MachineRouter) changeOwner(ctx *gin.Context) {

}

// Change a machine's group relationship. This will immediately transfer the machine resource to the new group.
// Changing the machine's group will not retroactively change the group of its user or license.
func (r *MachineRouter) changeGroup(ctx *gin.Context) {

}
