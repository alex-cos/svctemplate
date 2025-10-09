package app

import (
	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/internal/handler/admin"
	"github.com/alex-cos/scvtemplate/internal/handler/api"
	"github.com/alex-cos/scvtemplate/internal/handler/metrics"
	"github.com/alex-cos/scvtemplate/internal/pprof"
	"github.com/alex-cos/scvtemplate/pkg/dynamicLevel"
	"gopkg.in/tomb.v2"
)

// startAllServers starts API/metrics/admin and optionally pprof.
func StartAllServers(
	t *tomb.Tomb,
	cmgr *config.ConfigMgmt,
	dynLevel *dynamicLevel.DynamicLevel,
) {
	cfg := cmgr.GetConfig()
	mode := cmgr.GetMode()

	// API router
	t.Go(func() error {
		return StartWebServer(t, "api", &cfg.APIConfig, api.Init(mode))
	})

	// Metrics router
	t.Go(func() error {
		return StartWebServer(t, "metrics", &cfg.MetricsConfig, metrics.Init())
	})

	// Admin router
	t.Go(func() error {
		return StartWebServer(t, "admin", &cfg.AdminConfig, admin.Init(cmgr, dynLevel))
	})

	// pprof: start one server per returned addr (so dual -> two pprof listeners)
	if cfg.PprofConfig.Enable {
		t.Go(func() error {
			return pprof.Start(t, &cfg.PprofConfig)
		})
	}
}
