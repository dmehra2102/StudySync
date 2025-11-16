package middleware

import (
	"time"

	"github.com/dmehra2102/StudySync/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		latency := time.Since(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()

		var errMsg string
		if lastErr := ctx.Errors.Last(); lastErr != nil {
			errMsg = lastErr.Error()
		}

		if query != "" {
			path = path + "?" + query
		}

		event := log.Info()
		if statusCode >= 400 {
			event = log.Error()
		}

		event = event.
			Str("client_ip", clientIP).
			Str("method", method).
			Int("status_code", statusCode).
			Str("path", path).
			Str("latency", latency.String())

		if errMsg != "" {
			event = event.Str("error", errMsg)
		}

		event.Msg("HTTP Request")
	}
}
