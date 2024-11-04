package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/response"
	"net/http"
)

func SuperAdminAuthMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := response.NewResponse(ctx)
		user, ok := ctx.Get(gin.AuthUserKey)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		fmt.Println(user)

		ctx.Next()
	}
}
