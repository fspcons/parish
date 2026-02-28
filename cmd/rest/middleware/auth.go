package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/usecase"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	authUseCase usecase.Auth
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authUseCase usecase.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// Authenticate validates the authentication token
func (ref *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		user, err := ref.authUseCase.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission checks if the user has the required permission
func (ref *AuthMiddleware) RequirePermission(resource string, write bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hasPermission, err := ref.authUseCase.CheckPermission(r.Context(), user.ID, resource, write)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken reads the auth token from the auth_token cookie first,
// falling back to the Authorization: Bearer header.
func extractToken(r *http.Request) string {
	if c, err := r.Cookie("auth_token"); err == nil && c.Value != "" {
		return c.Value
	}

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
