package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
)

func RequestIDMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.New().String()
		ctx.Request.Header.Set(constants.RequestIDField, requestID)
		ctx.Set(constants.RequestIDField, requestID)
		ctx.Next()
	}
}
