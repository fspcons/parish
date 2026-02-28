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
	"github.com/parish/internal/repository"
	"github.com/parish/internal/repository/datastore"
	"github.com/parish/internal/usecase"
	"go.uber.org/dig"
)

func mustProvide(c *dig.Container, constructor any, opts ...dig.ProvideOption) {
	if err := c.Provide(constructor, opts...); err != nil {
		panic(fmt.Sprintf("failed to provide: %v", err))
	}
}

func mustBuildDIC(cfg *Config) (*dig.Container, func()) {
	c := dig.New()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		slog.Error("GCP_PROJECT_ID environment variable is required")
		os.Exit(1)
	}

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("failed to create Datastore client", "error", err)
		os.Exit(1)
	}

	redisCache, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to create Redis client", "error", err, "url", cfg.RedisURL)
		os.Exit(1)
	}

	cleanup := func() {
		if err := dsClient.Close(); err != nil {
			slog.Error("failed to close Datastore client", "error", err)
		}
		if err := redisCache.Close(); err != nil {
			slog.Error("failed to close Redis client", "error", err)
		}
	}

	// Datastore client
	mustProvide(c, func() *datastore.Client { return dsClient })

	// Cache
	mustProvide(c, func() cache.Cache { return redisCache })

	// Repositories (provided as interfaces)
	mustProvide(c, func(cl *datastore.Client) repository.ScheduleRepository { return datastore.NewScheduleRepository(cl) })
	mustProvide(c, func(cl *datastore.Client) repository.ParishGroupRepository {
		return datastore.NewParishGroupRepository(cl)
	})
	mustProvide(c, func(cl *datastore.Client) repository.EventRepository { return datastore.NewEventRepository(cl) })
	mustProvide(c, func(cl *datastore.Client) repository.MaterialRepository { return datastore.NewMaterialRepository(cl) })
	mustProvide(c, func(cl *datastore.Client) repository.UserRepository { return datastore.NewUserRepository(cl) })
	mustProvide(c, func(cl *datastore.Client) repository.RoleRepository { return datastore.NewRoleRepository(cl) })

	// Use cases
	mustProvide(c, func(r repository.ScheduleRepository) usecase.Schedule { return usecase.NewSchedule(r) })
	mustProvide(c, func(r repository.ParishGroupRepository) usecase.ParishGroup { return usecase.NewParishGroup(r) })
	mustProvide(c, func(r repository.EventRepository) usecase.Event { return usecase.NewEvent(r) })
	mustProvide(c, func(r repository.MaterialRepository) usecase.Material { return usecase.NewMaterial(r) })
	mustProvide(c, func(ur repository.UserRepository, rr repository.RoleRepository, cc cache.Cache) usecase.Auth {
		return usecase.NewAuth(ur, rr, cc)
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
