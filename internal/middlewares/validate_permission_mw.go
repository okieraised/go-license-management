package middlewares

import (
	"github.com/gin-gonic/gin"
)

func PermissionValidationMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()
	}
}
