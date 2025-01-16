package middlewares

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/casbin_adapter"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"net/http"
	"strings"
)

func PermissionValidationMW(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		e, err := casbin.NewEnforcer(casbin_adapter.GetEnforcerModel(), casbin_adapter.GetAdapter())
		if err != nil {
			logging.GetInstance().GetLogger().Error(err.Error())
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.NewResponse(ctx).ToResponse(
					comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer],
					comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer],
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
					comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer],
					comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer],
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
					comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer],
					comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer],
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
					comerrors.ErrCodeMapper[comerrors.ErrGenericPermission],
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
