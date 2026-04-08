package rest

import (
	"net/http"
	"time"

	"github.com/parish/cmd/rest/handler"
	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/cache"
)

// Router configures and returns the HTTP router.
type Router struct {
	authHandler            *handler.AuthHandler
	scheduleHandler        *handler.ScheduleHandler
	parishGroupHandler     *handler.ParishGroupHandler
	eventHandler           *handler.EventHandler
	materialHandler        *handler.MaterialHandler
	roleHandler            *handler.RoleHandler
	authMiddleware         *middleware.AuthMiddleware
	rateLimiter            *middleware.RateLimiter
	resetPasswordRateLimit *middleware.RateLimiter
	corsOrigin             string
}

// NewRouter creates a new Router.
func NewRouter(
	authHandler *handler.AuthHandler,
	scheduleHandler *handler.ScheduleHandler,
	parishGroupHandler *handler.ParishGroupHandler,
	eventHandler *handler.EventHandler,
	materialHandler *handler.MaterialHandler,
	roleHandler *handler.RoleHandler,
	authMiddleware *middleware.AuthMiddleware,
	c cache.Cache,
	corsOrigin string,
) *Router {
	return &Router{
		authHandler:            authHandler,
		scheduleHandler:        scheduleHandler,
		parishGroupHandler:     parishGroupHandler,
		eventHandler:           eventHandler,
		materialHandler:        materialHandler,
		roleHandler:            roleHandler,
		authMiddleware:         authMiddleware,
		rateLimiter:            middleware.NewRateLimiter(c, 100, time.Minute),
		resetPasswordRateLimit: middleware.NewRateLimiterWithPrefix(c, 1, time.Minute, "rl:reset-password:"),
		corsOrigin:             corsOrigin,
	}
}

// Setup configures all routes and returns the root handler.
func (ref *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth (public)
	mux.HandleFunc("POST /api/auth/register", ref.authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", ref.authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", ref.authHandler.Logout)
	mux.Handle("POST /api/auth/reset-password", ref.resetPasswordRateLimit.Limit(http.HandlerFunc(ref.authHandler.ResetPassword)))
	mux.Handle("PUT /api/auth/password", ref.authMiddleware.Authenticate(http.HandlerFunc(ref.authHandler.ChangePassword)))

	// Schedule
	mux.HandleFunc("GET /api/schedule", ref.scheduleHandler.Get)
	mux.HandleFunc("PUT /api/schedule", ref.protected("schedule", true, ref.scheduleHandler.Update))

	// Parish groups
	mux.HandleFunc("GET /api/parish-groups", ref.parishGroupHandler.List)
	mux.HandleFunc("POST /api/parish-groups", ref.protected("parish_groups", true, ref.parishGroupHandler.Create))
	mux.HandleFunc("GET /api/parish-groups/{id}", ref.parishGroupHandler.Get)
	mux.HandleFunc("PUT /api/parish-groups/{id}", ref.protected("parish_groups", true, ref.parishGroupHandler.Update))
	mux.HandleFunc("DELETE /api/parish-groups/{id}", ref.protected("parish_groups", true, ref.parishGroupHandler.Delete))

	// Events
	mux.HandleFunc("GET /api/events", ref.eventHandler.List)
	mux.HandleFunc("POST /api/events", ref.protected("events", true, ref.eventHandler.Create))
	mux.HandleFunc("GET /api/events/{id}", ref.eventHandler.Get)
	mux.HandleFunc("PUT /api/events/{id}", ref.protected("events", true, ref.eventHandler.Update))
	mux.HandleFunc("DELETE /api/events/{id}", ref.protected("events", true, ref.eventHandler.Delete))

	// Materials
	mux.HandleFunc("GET /api/materials", ref.materialHandler.List)
	mux.HandleFunc("POST /api/materials", ref.protected("materials", true, ref.materialHandler.Create))
	mux.HandleFunc("GET /api/materials/{id}", ref.materialHandler.Get)
	mux.HandleFunc("PUT /api/materials/{id}", ref.protected("materials", true, ref.materialHandler.Update))
	mux.HandleFunc("DELETE /api/materials/{id}", ref.protected("materials", true, ref.materialHandler.Delete))

	// Roles
	mux.HandleFunc("GET /api/roles", ref.protected("roles", false, ref.roleHandler.List))
	mux.HandleFunc("POST /api/roles", ref.protected("roles", true, ref.roleHandler.Create))
	mux.HandleFunc("GET /api/roles/{id}", ref.protected("roles", false, ref.roleHandler.Get))
	mux.HandleFunc("PUT /api/roles/{id}", ref.protected("roles", true, ref.roleHandler.Update))
	mux.HandleFunc("DELETE /api/roles/{id}", ref.protected("roles", true, ref.roleHandler.Delete))

	// User role assignment
	mux.HandleFunc("PUT /api/users/{id}/roles", ref.protected("roles", true, ref.roleHandler.AssignRoles))

	// Global middleware chain (outermost applied first).
	h := middleware.RequestLogger(mux)
	h = middleware.MaxBodySize(1048576)(h)
	h = middleware.SecurityHeaders(h)
	h = middleware.CORS(ref.corsOrigin)(h)
	h = ref.rateLimiter.Limit(h)

	return h
}

// protected wraps a handler with authentication + permission checks.
func (ref *Router) protected(resource string, write bool, next http.HandlerFunc) http.HandlerFunc {
	wrapped := ref.authMiddleware.Authenticate(
		ref.authMiddleware.RequirePermission(resource, write)(
			http.HandlerFunc(next),
		),
	)
	return wrapped.ServeHTTP
}
