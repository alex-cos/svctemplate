package admin

import (
	"net/http"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/observability"
	"github.com/alex-cos/scvtemplate/pkg/dynamicLevel"
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/gin-gonic/gin"
)

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}

func readyHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ready"})
}

func configHandler(cmgr *config.ConfigMgmt) func(c *gin.Context) {
	return func(c *gin.Context) {
		log := logx.FromContext(c.Request.Context())
		log.Info("config requested via admin endpoint")

		c.JSON(http.StatusOK, cmgr.GetConfig())
	}
}

func reloadHandler(cmgr *config.ConfigMgmt) func(c *gin.Context) {
	return func(c *gin.Context) {
		log := logx.FromContext(c.Request.Context())
		log.Info("reload requested via admin endpoint")
		if cmgr.NotifyChange() {
			observability.RecordConfigReload(true)
			c.JSON(200, gin.H{"status": "reload triggered"})
		} else {
			c.JSON(429, gin.H{"status": "reload already in progress"})
		}
	}
}

func loglevelHandler(dynLevel *dynamicLevel.DynamicLevel) func(c *gin.Context) {
	return func(c *gin.Context) {
		log := logx.FromContext(c.Request.Context())
		log.Info("log requested via admin endpoint")
		if c.Request.Method == http.MethodPost {
			var req LoglevelRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			dynLevel.SetLevel(dynamicLevel.ParseLogLevel(req.Level))
		}
		c.JSON(http.StatusOK, gin.H{"level": dynLevel.Level()})
	}
}
