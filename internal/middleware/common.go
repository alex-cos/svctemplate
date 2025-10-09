package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// all middlewares.
func All(mode string) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{
		gin.Recovery(),
		Security(),
		Prometheus(),
		Logger(),
	}

	if strings.EqualFold(mode, "dev") {
		handlers = append(handlers, cors.Default())
	}

	return handlers
}

// all middlewares.
func Necessary(mode string) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{
		gin.Recovery(),
		Logger(),
	}

	if strings.EqualFold(mode, "dev") {
		handlers = append(handlers, cors.Default())
	}

	return handlers
}
