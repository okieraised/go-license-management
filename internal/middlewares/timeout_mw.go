package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/response"
	"net/http"
	"time"
)

func TimeoutMW() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		timeoutDuration := 10 * time.Second

		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			ctx.Next()
			finish <- struct{}{}
		}()

		resp := response.NewResponse(ctx)
		select {
		case <-panicChan:
			resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericInternalServer], comerrors.ErrMessageMapper[comerrors.ErrGenericInternalServer], nil, nil, nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		case <-time.After(timeoutDuration):
			resp.ToResponse(comerrors.ErrCodeMapper[comerrors.ErrGenericRequestTimedOut], comerrors.ErrMessageMapper[comerrors.ErrGenericRequestTimedOut], nil, nil, nil)
			ctx.AbortWithStatusJSON(http.StatusGatewayTimeout, resp)
			return
		case <-finish:
		}
	}
}
