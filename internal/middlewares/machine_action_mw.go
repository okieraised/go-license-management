package middlewares

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/casbin_adapter"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/permissions"
	"go-license-management/internal/response"
	"net/http"
	"strings"
)

func MachineActionPermissionValidationMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actions := ctx.Param("machine_action")
		if actions == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
					"missing authorization header",
					nil,
					nil,
					nil,
				),
			)
			return
		}

		e, err := casbin.NewEnforcer(casbin_adapter.GetEnforcerModel(), casbin_adapter.GetAdapter())
		if err != nil {
			logging.GetInstance().GetLogger().Error(err.Error())
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
					cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
					nil,
					nil,
					nil,
				),
			)
			return
		}

		err = e.LoadPolicy()
		if err != nil {
			logging.GetInstance().GetLogger().Error(err.Error())
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
					cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
					nil,
					nil,
					nil,
				),
			)
			return
		}

		var permission string

		switch actions {
		case constants.MachineActionCheckout:
			permission = permissions.MachineCheckOut
		case constants.MachineActionPingHeartbeat:
			permission = permissions.MachineHeartbeatPing
		case constants.MachineActionResetHeartbeat:
			permission = permissions.MachineHeartbeatReset
		default:
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrAccountActionIsInvalid],
					cerrors.ErrMessageMapper[cerrors.ErrAccountActionIsInvalid],
					nil,
					nil,
					nil,
				),
			)
			return

		}
		permObjects := strings.Split(permission, ".")

		ok, err := e.Enforce(
			ctx.GetString(constants.ContextValueTenant),
			ctx.GetString(constants.ContextValueSubject),
			permObjects[0],
			permObjects[1],
		)
		if err != nil {
			logging.GetInstance().GetLogger().Error(err.Error())
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
					cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
					nil,
					nil,
					nil,
				),
			)
			return
		}

		if !ok {
			logging.GetInstance().GetLogger().Info(
				fmt.Sprintf("invalid permission: domain [%s] | subject [%s] | object [%s] | action [%s]",
					ctx.GetString(constants.ContextValueTenant),
					ctx.GetString(constants.ContextValueSubject),
					permObjects[0],
					permObjects[1]),
			)
			ctx.AbortWithStatusJSON(
				http.StatusForbidden,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericPermission],
					fmt.Sprintf("user [%s] does not have permission to perform the requested action", ctx.GetString(constants.ContextValueSubject)),
					nil,
					nil,
					nil,
				),
			)
			return
		}
		logging.GetInstance().GetLogger().Info(
			fmt.Sprintf("valid permission: domain [%s] | subject [%s] | object [%s] | action [%s]",
				ctx.GetString(constants.ContextValueTenant),
				ctx.GetString(constants.ContextValueSubject),
				permObjects[0],
				permObjects[1]),
		)

		ctx.Next()
	}
}
