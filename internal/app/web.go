package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alex-cos/scvtemplate/config"
	"github.com/alex-cos/scvtemplate/pkg/logx"
	"github.com/gin-gonic/gin"
	"gopkg.in/tomb.v2"
)

func buildTLSConfig(certFile, keyFile, caFile, clientAuth string) (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Load client CA if required
	if clientAuth == "require" || clientAuth == "verify" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.ClientCAs = caCertPool

		if clientAuth == "require" {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
	}

	return tlsConfig, nil
}

// startGinServer starts a single HTTP server (one addr) with graceful shutdown handled via tomb.
func StartWebServer(
	t *tomb.Tomb,
	name string,
	cfg *config.WebConfig,
	router *gin.Engine,
) error {
	var (
		tlsConfig *tls.Config
		err       error
	)

	if cfg.EnableTLS {
		tlsConfig, err = buildTLSConfig(cfg.TLSCertFile, cfg.TLSKeyFile, cfg.TLSCAFile, cfg.TLSClientAuth)
		if err != nil {
			logx.L().Error("failed to load TLS config", "error", err)
			return err
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         tlsConfig,
		ReadTimeout:       0,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    0,
	}

	go func() {
		<-t.Dying()
		logx.L().Info("shutdown requested", slog.String("server", name), slog.String("addr", addr))
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logx.L().Error(
				"error during server shutdown",
				slog.String("server", name),
				slog.String("addr", addr),
				slog.Any("error", err),
			)
		} else {
			logx.L().Info(
				"server stopped gracefully",
				slog.String("server", name),
				slog.String("addr", addr),
			)
		}
	}()

	logx.L().Info("starting server", slog.String("server", name), slog.String("addr", addr))
	if cfg.EnableTLS {
		err = srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logx.L().Error(
				"server ListenAndServe error",
				slog.String("server", name),
				slog.String("addr", addr),
				slog.Any("error", err),
			)
			return err
		}
	} else {
		err = srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logx.L().Error(
				"server ListenAndServe error",
				slog.String("server", name),
				slog.String("addr", addr),
				slog.Any("error", err),
			)
			return err
		}
	}

	return nil
}
