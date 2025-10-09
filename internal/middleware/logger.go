package middleware

import (
	"log/slog"
	"time"

	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Check for incoming request ID or generate a new one
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", reqID)

		// Create a contextual logger
		logger := logx.L().With(
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_ip", c.ClientIP(),
		)

		// Inject logger into request context
		ctx := logx.WithLogger(c.Request.Context(), logger)
		ctx = logx.WithRequestID(ctx, reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logger.Info("request processed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", status),
			slog.Duration("duration", duration),
		)
	}
}
