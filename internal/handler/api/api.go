package api

import (
	"github.com/alex-cos/scvtemplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

// API router.
func Init(mode string) *gin.Engine {
	api := gin.New()
	api.Use(middleware.All(mode)...)
	api.GET("/api/hello", helloHandler)
	api.GET("/api/version", versionHandler)

	return api
}
