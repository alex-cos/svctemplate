package api

import (
	"strings"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WebAPI struct {
	mode      string
	ratelimit *config.RateLimiterConfig
}

func New(mode string, ratelimit *config.RateLimiterConfig) *WebAPI {
	return &WebAPI{
		mode:      mode,
		ratelimit: ratelimit,
	}
}

// API router.
func (api *WebAPI) Init() *gin.Engine {
	engine := gin.New()

	rateLimiter := middleware.NewRateLimiter(api.ratelimit.Rate, api.ratelimit.Burst)

	handlers := []gin.HandlerFunc{
		gin.Recovery(),
		middleware.Security(),
		middleware.Prometheus(),
		middleware.Logger(),
	}
	if strings.EqualFold(api.mode, "dev") {
		handlers = append(handlers, cors.Default())
	} else {
		handlers = append(handlers, rateLimiter.Middleware())
	}
	engine.Use(handlers...)

	engine.GET("/api/hello", api.helloHandler)
	engine.GET("/api/version", api.versionHandler)

	return engine
}
