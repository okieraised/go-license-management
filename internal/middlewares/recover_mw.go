package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"net/http"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				resp := response.NewResponse(ctx)
				logging.GetInstance().GetLogger().Error(string(debug.Stack()))
				resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer], cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer], err, nil, nil)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
				return
			}
		}()

		ctx.Next()
	}
}
