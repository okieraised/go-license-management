package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/response"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				resp := response.NewResponse(ctx)
				resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer], comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer], nil, nil, nil)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
				return
			}
		}()

		ctx.Next()
	}
}
