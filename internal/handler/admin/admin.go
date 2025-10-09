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
	admin.GET("/config", configHandler(cmgr))
	admin.POST("/reload", reloadHandler(cmgr))
	admin.GET("/loglevel", loglevelHandler(dynLevel))
	admin.POST("/loglevel", loglevelHandler(dynLevel))

	return admin
}
