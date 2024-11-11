package policies

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/infrastructure/models/policy_attribute"
	"go-license-management/internal/infrastructure/tracer"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/policies/service"
	"go-license-management/server/models/v1/policies"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

const (
	policyGroup = "policy_group"
)

type PolicyRouter struct {
	svc    *service.PolicyService
	logger *logging.Logger
	tracer trace.Tracer
}

func NewPolicyRouter(svc *service.PolicyService) *PolicyRouter {
	tr := tracer.GetInstance().Tracer(policyGroup)
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
		routes.POST("", r.create)
		routes.GET("", r.list)
		routes.GET("/:policy_id", r.retrieve)
		routes.PATCH("/:policy_id", r.update)
		routes.DELETE("/:policy_id", r.delete)
		routes.POST("/:policy_id/entitlements", r.attach)
		routes.DELETE("/:policy_id/entitlements", r.detach)
		routes.GET("/:policy_id/entitlements", r.listEntitlement)
	}
}

// create creates a new policy resource.
func (r *PolicyRouter) create(ctx *gin.Context) {
	rootCtx, span := r.tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	resp := response.NewResponse(ctx)
	r.logger.WithCustomFields(zap.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField))).Info("received new policy creation request")

	// serializer
	_, cSpan := r.tracer.Start(rootCtx, "serializer")
	var uriReq policy_attribute.PolicyCommonURI
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericBadRequest], comerrors.ErrMessageMapper[comerrors.ErrGenericBadRequest], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var bodyReq policies.PolicyRegistrationRequest
	err = ctx.ShouldBind(&bodyReq)
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
	err = bodyReq.Validate()
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(comerrors.ErrCodeMapper[err], comerrors.ErrMessageMapper[err], nil, nil, nil)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	err = uriReq.Validate()
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
	result, err := r.svc.Create(ctx, bodyReq.ToPolicyRegistrationInput(rootCtx, r.tracer, uriReq))
	if err != nil {
		cSpan.End()
		r.logger.GetLogger().Error(err.Error())
		resp.ToResponse(result.Code, result.Message, result.Data, nil, nil)
		switch {
		case errors.Is(err, comerrors.ErrTenantNameIsInvalid),
			errors.Is(err, comerrors.ErrProductIDIsInvalid),
			errors.Is(err, comerrors.ErrPolicySchemeIsInvalid):
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

// retrieve retrieves the details of an existing policy.
func (r *PolicyRouter) retrieve(ctx *gin.Context) {

}

// update updates the specified policy resource by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged.
func (r *PolicyRouter) update(ctx *gin.Context) {

}

// delete permanently deletes a policy. It cannot be undone.
// This action also immediately deletes any licenses that the policy is associated with.
func (r *PolicyRouter) delete(ctx *gin.Context) {

}

// list returns a list of policies. The policies are returned sorted by creation date, with the most recent policies
// appearing first. Resources are automatically scoped to the authenticated bearer
// e.g. when authenticated as a product, only policies of that specific product will be listed.
func (r *PolicyRouter) list(ctx *gin.Context) {

}

// attach attaches entitlements to a policy. This will immediately be taken into effect for all future license validations.
// Any license that implements the given policy will automatically possess all the policy's entitlements.
func (r *PolicyRouter) attach(ctx *gin.Context) {

}

// detach detaches entitlements from a policy. This will immediately be taken into effect for all future license validations.
func (r *PolicyRouter) detach(ctx *gin.Context) {

}

// listEntitlement returns a list of entitlements attached to the policy.
// The entitlements are returned sorted by creation date, with the most recent entitlements appearing first.
func (r *PolicyRouter) listEntitlement(ctx *gin.Context) {

}
