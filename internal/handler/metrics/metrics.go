package metrics

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics router.
func Init() *gin.Engine {
	metrics := gin.New()
	metrics.Use(
		gin.Recovery(),
		cors.Default(),
	)
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return metrics
}
