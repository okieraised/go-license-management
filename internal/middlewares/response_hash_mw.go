package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
)

func HashHeaderMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writer := &ResponseWriterInterceptor{
			ResponseWriter: ctx.Writer,
			body:           make([]byte, 0),
		}
		ctx.Writer = writer

		ctx.Next()

	}
}

type ResponseWriterInterceptor struct {
	gin.ResponseWriter
	body []byte
}

func (w *ResponseWriterInterceptor) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	hash := sha256.Sum256(data)
	w.Header().Add(constants.ContentDigestHeader, fmt.Sprintf("sha256=%s", hex.EncodeToString(hash[:])))
	return w.ResponseWriter.Write(data)
}
