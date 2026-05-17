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
)

// Run bootstraps the Gin web engine and manages the server lifecycle
func Run(cfg *config.Config) error {
	// Set Gin runtime mode to Release to strip debug noise and optimize allocations
	gin.SetMode(gin.ReleaseMode)

	// Initialize raw Gin engine
	r := gin.New()

	// Inject recovery middleware to gracefully catch unhandled panics within routines
	r.Use(gin.Recovery())

	// Bind core application routing metrics
	setupRoutes(r)

	// 1. Establish the native core network listener (supports tcp, tcp4, unix)
	listener, err := net.Listen(cfg.Web.Network, cfg.Web.Addr)
	if err != nil {
		slog.Error("failed to bind network listener", "network", cfg.Web.Network, "address", cfg.Web.Addr, "error", err)
		return err
	}
	defer listener.Close()

	// 2. Evaluate and inject TLS cryptography abstraction layer
	tlsConfig, err := cfg.Web.TLS.MakeTLSConfig()
	if err != nil {
		slog.Error("failed to compile TLS configuration", "error", err)
		return err
	}

	if tlsConfig != nil {
		listener = tls.NewListener(listener, tlsConfig)
		slog.Info("TLS security layer activated successfully via ALPN transport")
	}

	// 3. Mount handler into standard net/http Server architecture
	srv := &http.Server{
		Handler: r,
	}

	// 4. Implement Graceful Shutdown mechanics via OS signals
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Spawn the network server worker thread asynchronously
	go func() {
		slog.Info("os-webui engine interface online", "network", cfg.Web.Network, "address", cfg.Web.Addr)
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("critical crash caught inside native server listener", "error", err)
			os.Exit(1)
		}
	}()

	// Block main thread context until OS signal is intercepted
	<-shutdownCtx.Done()
	slog.Warn("shutdown signal intercepted, initializing graceful timeout sequences...")

	// Enforce a strict 5-second maximum drain window for pending streaming sockets
	drainCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(drainCtx); err != nil {
		slog.Error("forced server termination triggered during socket drain window", "error", err)
		return err
	}

	// If using Unix Domain Sockets, sweep the dead .sock file off the host filesystem cleanly
	if cfg.Web.Network == "unix" {
		if err := os.Remove(cfg.Web.Addr); err != nil && !os.IsNotExist(err) {
			slog.Warn("unable to purge dangling unix socket descriptor", "path", cfg.Web.Addr, "error", err)
		} else {
			slog.Debug("unix socket descriptor successfully unlinked from host entry")
		}
	}

	slog.Info("os-webui runtime terminated gracefully. all systems clear.")
	return nil
}

// setupRoutes handles operational HTTP routing parameters
func setupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"runtime": "go-native",
		})
	})

	// Placeholder route: Where your gorilla/websocket connection handler will attach to Yaegi stream execution contexts
	r.GET("/ws", func(c *gin.Context) {
		c.String(http.StatusOK, "websocket tunnel context endpoint ready")
	})
}
