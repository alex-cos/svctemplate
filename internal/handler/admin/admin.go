package admin

import (
	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/internal/middleware"
	"github.com/alex-cos/scvtemplate/pkg/dynamicLevel"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Admin router.
func Init(
	cmgr *config.ConfigMgmt,
	dynLevel *dynamicLevel.DynamicLevel,
) *gin.Engine {
	admin := gin.New()
	admin.Use(
		gin.Recovery(),
		cors.Default(),
		middleware.Logger(),
	)
	admin.GET("/healthz", healthHandler)
	admin.GET("/readyz", readyHandler)
	admin.GET("/admin/config", configHandler(cmgr))
	admin.POST("/admin/reload", reloadHandler(cmgr))
	admin.GET("/admin/loglevel", loglevelHandler(dynLevel))
	admin.POST("/admin/loglevel", loglevelHandler(dynLevel))

	return admin
}
