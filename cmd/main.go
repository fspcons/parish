package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parish/cmd/rest"
	"github.com/parish/internal/repository"
)

func main() {
	cfg := setupConfig()
	setupLogger()

	container, cleanup := mustBuildDIC(cfg)
	defer cleanup()

	// Seed admin user/role on startup (idempotent, skipped when env vars are empty).
	if err := container.Invoke(func(
		userRepo repository.UserRepository,
		roleRepo repository.RoleRepository,
	) {
		go seedAdmin(context.Background(), cfg, userRepo, roleRepo)
	}); err != nil {
		slog.Error("failed to run seed", "error", err)
		os.Exit(1)
	}

	if err := container.Invoke(func(router *rest.Router) {
		srv := &http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           router.Setup(),
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}

		// Start server in a goroutine
		go func() {
			tlsEnabled := cfg.TLSCert != "" && cfg.TLSKey != ""
			slog.Info("server starting", "port", cfg.Port, "tls", tlsEnabled)

			var listenErr error
			if tlsEnabled {
				listenErr = srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
			} else {
				listenErr = srv.ListenAndServe()
			}
			if listenErr != nil && listenErr != http.ErrServerClosed {
				slog.Error("server failed to start", "error", listenErr)
				os.Exit(1)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server forced to shutdown", "error", err)
			os.Exit(1)
		}

		slog.Info("server stopped gracefully")
	}); err != nil {
		slog.Error("failed to invoke server", "error", err)
		os.Exit(1)
	}
}
