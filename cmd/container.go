package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/parish/cmd/rest"
	"github.com/parish/cmd/rest/handler"
	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/cache"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/email"
	"github.com/parish/internal/repository"
	repfs "github.com/parish/internal/repository/firestore"
	"github.com/parish/internal/usecase"
	"go.uber.org/dig"
)

func mustProvide(c *dig.Container, constructor any, opts ...dig.ProvideOption) {
	if err := c.Provide(constructor, opts...); err != nil {
		panic(fmt.Sprintf("failed to provide: %v", err))
	}
}

func mustBuildDIC(cfg *Config) (*dig.Container, func()) {
	if !cfg.Local {
		if domain.IsEmpty(cfg.SendGridAPIKey) || domain.IsEmpty(cfg.EmailFrom) {
			slog.Error("SENDGRID_API_KEY and EMAIL_FROM are required when not running locally (set LOCAL_DEV=true or use the Firestore emulator for log-only e-mail)")
			os.Exit(1)
		}
	}

	c := dig.New()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		slog.Error("GCP_PROJECT_ID environment variable is required")
		os.Exit(1)
	}

	ctx := context.Background()
	fsStore, err := repfs.NewStore(ctx, projectID)
	if err != nil {
		slog.Error("failed to create Firestore client", "error", err)
		os.Exit(1)
	}

	redisCache, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to create Redis client", "error", err, "url", cfg.RedisURL)
		os.Exit(1)
	}

	cleanup := func() {
		if err := fsStore.Close(); err != nil {
			slog.Error("failed to close Firestore client", "error", err)
		}
		if err := redisCache.Close(); err != nil {
			slog.Error("failed to close Redis client", "error", err)
		}
	}

	// Firestore client
	mustProvide(c, func() *repfs.Store { return fsStore })

	// Cache
	mustProvide(c, func() cache.Cache { return redisCache })

	// E-mail: log-only locally; SendGrid in deployed environments (typical on GCP with API key in Secret Manager).
	mustProvide(c, func() email.Sender {
		if cfg.Local {
			return email.LogSender{}
		}
		return email.NewSendGridSender(cfg.SendGridAPIKey, cfg.EmailFrom, cfg.SendGridAPIURL)
	})

	// Repositories (provided as interfaces)
	mustProvide(c, func(st *repfs.Store) repository.ScheduleRepository { return repfs.NewScheduleRepository(st) })
	mustProvide(c, func(st *repfs.Store) repository.ParishGroupRepository {
		return repfs.NewParishGroupRepository(st)
	})
	mustProvide(c, func(st *repfs.Store) repository.EventRepository { return repfs.NewEventRepository(st) })
	mustProvide(c, func(st *repfs.Store) repository.MaterialRepository { return repfs.NewMaterialRepository(st) })
	mustProvide(c, func(st *repfs.Store) repository.UserRepository { return repfs.NewUserRepository(st) })
	mustProvide(c, func(st *repfs.Store) repository.RoleRepository { return repfs.NewRoleRepository(st) })

	// Use cases
	mustProvide(c, func(r repository.ScheduleRepository) usecase.Schedule { return usecase.NewSchedule(r) })
	mustProvide(c, func(r repository.ParishGroupRepository) usecase.ParishGroup { return usecase.NewParishGroup(r) })
	mustProvide(c, func(r repository.EventRepository) usecase.Event { return usecase.NewEvent(r) })
	mustProvide(c, func(r repository.MaterialRepository) usecase.Material { return usecase.NewMaterial(r) })
	mustProvide(c, func(ur repository.UserRepository, rr repository.RoleRepository, cc cache.Cache, m email.Sender) usecase.Auth {
		return usecase.NewAuth(ur, rr, cc, m)
	})
	mustProvide(c, func(r repository.RoleRepository) usecase.Role { return usecase.NewRole(r) })

	// Handlers
	mustProvide(c, func(uc usecase.Auth) *handler.AuthHandler { return handler.NewAuthHandler(uc, cfg.CookieSecure) })
	mustProvide(c, func(uc usecase.Schedule) *handler.ScheduleHandler { return handler.NewScheduleHandler(uc) })
	mustProvide(c, func(uc usecase.ParishGroup) *handler.ParishGroupHandler { return handler.NewParishGroupHandler(uc) })
	mustProvide(c, func(uc usecase.Event) *handler.EventHandler { return handler.NewEventHandler(uc) })
	mustProvide(c, func(uc usecase.Material) *handler.MaterialHandler { return handler.NewMaterialHandler(uc) })
	mustProvide(c, func(rc usecase.Role, ac usecase.Auth) *handler.RoleHandler {
		return handler.NewRoleHandler(rc, ac)
	})

	// Auth middleware
	mustProvide(c, func(uc usecase.Auth) *middleware.AuthMiddleware { return middleware.NewAuthMiddleware(uc) })

	// Router
	mustProvide(c, func(
		ah *handler.AuthHandler,
		sh *handler.ScheduleHandler,
		pgh *handler.ParishGroupHandler,
		eh *handler.EventHandler,
		mh *handler.MaterialHandler,
		rh *handler.RoleHandler,
		am *middleware.AuthMiddleware,
		cc cache.Cache,
	) *rest.Router {
		return rest.NewRouter(ah, sh, pgh, eh, mh, rh, am, cc, cfg.CorsOrigin)
	})

	return c, cleanup
}
