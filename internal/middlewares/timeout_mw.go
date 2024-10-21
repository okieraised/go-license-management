package middlewares

import (
	"github.com/gin-gonic/gin"
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

		select {
		case <-panicChan:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		case <-time.After(timeoutDuration):
			ctx.AbortWithStatusJSON(http.StatusGatewayTimeout, nil)
			return
		case <-finish:
		}
	}
}
