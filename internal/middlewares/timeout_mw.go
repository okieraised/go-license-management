package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/response"
	"net/http"
	"time"
)

func TimeoutMW() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		timeoutDuration := 10 * time.Second

		// Create a context with timeout to properly cancel downstream operations
		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeoutDuration)
		defer cancel()

		// Replace request context with timeout context
		ctx.Request = ctx.Request.WithContext(timeoutCtx)

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
			resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer], cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer], nil, nil, nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		case <-timeoutCtx.Done():
			// Check if it's a timeout or cancellation
			if timeoutCtx.Err() == context.DeadlineExceeded {
				resp.ToResponse(cerrors.ErrCodeMapper[cerrors.ErrGenericRequestTimedOut], cerrors.ErrMessageMapper[cerrors.ErrGenericRequestTimedOut], nil, nil, nil)
				ctx.AbortWithStatusJSON(http.StatusGatewayTimeout, resp)
			}
			return
		case <-finish:
			// Request completed successfully
		}
	}
}
