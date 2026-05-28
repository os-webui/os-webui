package web

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/os-webui/os-webui/config"
	"github.com/os-webui/os-webui/internal/plugins"
	v1 "github.com/os-webui/os-webui/web/v1"
)

// Run bootstraps the Gin web engine with standard library native h2c support
func runWeb(cfg *config.WebConfig, dev bool, slog *slog.Logger) error {
	if !dev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	setupRoutes(r, cfg)

	// 1. Establish the native core network listener (supports tcp, tcp4, unix)
	listener, err := net.Listen(cfg.Network, cfg.Addr)
	if err != nil {
		slog.Error("failed to bind network listener", "network", cfg.Network, "address", cfg.Addr, "error", err)
		return err
	}
	defer listener.Close()

	// 2. Evaluate and inject TLS cryptography abstraction layer
	tlsConfig, err := cfg.TLS.MakeTLSConfig()
	if err != nil {
		slog.Error("failed to compile TLS configuration", "error", err)
		return err
	}

	if tlsConfig != nil {
		listener = tls.NewListener(listener, tlsConfig)
		slog.Info("TLS security layer activated successfully via ALPN transport")
	}

	// 3. 🌟 Initialize native protocol orchestration via Go standard library
	// This declares explicit support for both standard HTTP/1 and unencrypted HTTP/2 (Prior Knowledge h2c)
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true) // Native h2c engine engaged!

	// 4. Mount the handler and protocols straight into standard net/http Server architecture
	srv := &http.Server{
		Handler:   r,
		Protocols: protocols, // 🌟 Zero external dependencies, pure native compliance
	}

	// 5. Implement Graceful Shutdown mechanics via OS signals
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("os-webui engine interface online", "network", cfg.Network, "address", cfg.Addr, "native_h2c", true)
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("critical crash caught inside native server listener", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownCtx.Done()
	slog.Warn("shutdown signal intercepted, initializing graceful timeout sequences...")

	drainCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(drainCtx); err != nil {
		slog.Error("forced server termination triggered during socket drain window", "error", err)
		return err
	}

	if cfg.Network == "unix" {
		if err := os.Remove(cfg.Addr); err != nil && !os.IsNotExist(err) {
			slog.Warn("unable to purge dangling unix socket descriptor", "path", cfg.Addr, "error", err)
		} else {
			slog.Debug("unix socket descriptor successfully unlinked from host entry")
		}
	}

	plugins.DefaultPluginsManager.Cleanup(slog)

	slog.Info("os-webui runtime terminated gracefully. all systems clear.")
	return nil
}

func setupRoutes(r *gin.Engine, cfg *config.WebConfig) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"runtime": "go-native",
			"proto":   c.Request.Proto, // Prints "HTTP/2.0" instantly under native h2c client connection
		})
	})
	api := r.Group(`/api`)
	if len(cfg.Accounts) != 0 {
		api.Use(gin.BasicAuthForRealm(cfg.Accounts, `os-webui`))
	}

	v1.InitRouter(api)
}
