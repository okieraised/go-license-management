package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"log/slog"
	"time"
)

func LoggerMW(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()

		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors.Errors() {
				logger.Error(e)
			}
			return
		}

		latency := time.Since(start).Milliseconds()
		logger.LogAttrs(context.Background(), slog.LevelInfo, "request handled",
			slog.String(constants.RequestIDField, ctx.GetString(constants.RequestIDField)),
			slog.Int("status", ctx.Writer.Status()),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("full-path", ctx.FullPath()),
			slog.String("query", query),
			slog.String("ip", ctx.ClientIP()),
			slog.String("user-agent", ctx.Request.UserAgent()),
			slog.Int64("latency", latency),
		)
	}
}
