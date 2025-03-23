package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logger.Debug("request started", "method", c.Request.Method, "path", c.Request.URL.Path, "remote_addr", c.ClientIP())

		c.Next()

		duration := time.Since(start)
		logger.Debug("request completed", "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "duration", duration.String())
	}
}
