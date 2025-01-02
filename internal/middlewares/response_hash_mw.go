package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
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

		hash := sha256.Sum256(writer.body)
		hashString := hex.EncodeToString(hash[:])
		ctx.Header(constants.ContentDigestHeader, hashString)
	}
}

type ResponseWriterInterceptor struct {
	gin.ResponseWriter
	body []byte
}

func (w *ResponseWriterInterceptor) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)

	//hash := sha256.Sum256(data)
	//w.Header().Add(constants.ContentDigestHeader, fmt.Sprintf("sha256=%s", hex.EncodeToString(hash[:])))
	return w.ResponseWriter.Write(data)
}
