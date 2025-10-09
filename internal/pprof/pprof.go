package pprof

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/internal/middleware"
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"gopkg.in/tomb.v2"
)

// start starts pprof server on a single addr (used for dual binding too).
func Start(t *tomb.Tomb, cfg *config.PprofConfig) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	var handler http.Handler = mux
	addr := fmt.Sprintf(":%d", cfg.Port)

	if cfg.User != "" && cfg.Pass != "" {
		logx.L().Info("enabling basic auth for pprof", slog.String("addr", addr))
		handler = middleware.BasicAuthMiddleware(cfg.User, cfg.Pass, mux)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    0,
	}

	go func() {
		<-t.Dying()
		logx.L().Info("pprof shutdown requested", slog.String("addr", addr))
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logx.L().Error("error shutting down pprof", slog.String("addr", addr), slog.Any("error", err))
		}
	}()

	logx.L().Info("starting pprof server", slog.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logx.L().Error("pprof ListenAndServe error", slog.String("addr", addr), slog.Any("error", err))
		return err
	}

	return nil
}
