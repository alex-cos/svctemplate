package app

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/pkg/dynamicLevel"
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/gin-gonic/gin"
	"gopkg.in/tomb.v2"
)

func waitForQuitSignal(t *tomb.Tomb) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		s := <-sig
		logx.L().Warn("os signal received, initiating shutdown", slog.String("signal", s.String()))
		t.Kill(nil)
	}()
}

func Execute(cmgr *config.ConfigMgmt, dynLevel *dynamicLevel.DynamicLevel) error {
	var t tomb.Tomb

	// parse mode
	mode := cmgr.GetMode()
	if strings.EqualFold(mode, "dev") {
		logx.L().Info("developppement mode detected", slog.Any("mode", mode))
		gin.SetMode(gin.DebugMode)
	} else {
		logx.L().Info("release mode detected")
		gin.SetMode(gin.ReleaseMode)
	}

	// watch configuration and reload it when change detceted
	cmgr.Watch(func(oldCfg, newCfg *config.Config) {
		logx.L().Warn("configuration change detected, restarting servers...")
		t.Kill(nil)
		_ = t.Wait()

		if oldCfg.LogLevel != newCfg.LogLevel {
			dynLevel.SetLevel(dynamicLevel.ParseLogLevel(newCfg.LogLevel))
		}
		// reset tomb and restart with new config
		t = tomb.Tomb{}
		StartAllServers(&t, cmgr, dynLevel)
	})

	// start servers
	StartAllServers(&t, cmgr, dynLevel)

	// signal handling
	waitForQuitSignal(&t)

	// wait for servers to finish
	if err := t.Wait(); err != nil {
		logx.L().Error("service stopped with error", slog.Any("error", err))
	} else {
		logx.L().Info("service stopped gracefully")
	}

	return nil
}
