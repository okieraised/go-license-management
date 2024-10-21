package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.New().String()
		ctx.Request.Header.Set("", requestID)
		ctx.Set("", requestID)
		ctx.Next()
	}
}
