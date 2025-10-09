package middleware

import (
	"strconv"
	"time"

	"github.com/alex-cos/scvtemplate/observability"
	"github.com/gin-gonic/gin"
)

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method

		observability.HttpActiveRequests.WithLabelValues(method, path).Inc()
		c.Next()
		duration := time.Since(start).Seconds()

		status := strconv.Itoa(c.Writer.Status())
		observability.HttpRequests.WithLabelValues(method, path, status).Inc()
		observability.HttpDuration.WithLabelValues(method, path).Observe(duration)
		observability.HttpActiveRequests.WithLabelValues(method, path).Dec()
	}
}
