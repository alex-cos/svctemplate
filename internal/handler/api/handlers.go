package api

import (
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/alex-cos/scvtemplate/version"
	"github.com/gin-gonic/gin"
)

func (api *WebAPI) helloHandler(c *gin.Context) {
	log := logx.FromContext(c.Request.Context())
	reqID := logx.RequestIDFromContext(c.Request.Context())

	log.Info("handling hello endpoint", "request_id", reqID)
	c.JSON(200, gin.H{
		"message":    "hello world",
		"request_id": reqID,
	})
}

func (api *WebAPI) versionHandler(c *gin.Context) {
	log := logx.FromContext(c.Request.Context())
	reqID := logx.RequestIDFromContext(c.Request.Context())

	log.Info("handling version endpoint", "request_id", reqID)
	c.JSON(200, gin.H{
		"version":    version.GetVersion(),
		"build_date": version.GetBuildDate(),
	})
}
