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

// Run bootstraps the Gin web engine with standard library native h2c support
func Run(cfg *config.Config) error {
	if !cfg.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

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
		slog.Info("os-webui engine interface online", "network", cfg.Web.Network, "address", cfg.Web.Addr, "native_h2c", true)
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

func setupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"runtime": "go-native",
			"proto":   c.Request.Proto, // Prints "HTTP/2.0" instantly under native h2c client connection
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		c.String(http.StatusOK, "websocket tunnel context endpoint ready")
	})
}
