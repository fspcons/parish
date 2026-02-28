package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/parish/internal/cache"
)

// CORS returns middleware that adds CORS headers for the given allowed origin.
// A specific origin (not "*") is required so that credentials (cookies) are accepted.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Strict transport security (HTTPS only)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content security policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

// RateLimiter implements a sliding-window rate limiter backed by Redis sorted sets.
type RateLimiter struct {
	cache  cache.Cache
	limit  int
	window time.Duration
}

// NewRateLimiter creates a new Redis-backed rate limiter.
func NewRateLimiter(c cache.Cache, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		cache:  c,
		limit:  limit,
		window: window,
	}
}

// Limit rate limits requests based on IP address using a Redis sorted set.
func (ref *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		key := "rl:" + ip
		now := time.Now()
		windowStart := fmt.Sprintf("%d", now.Add(-ref.window).UnixNano())

		ctx := r.Context()

		// Remove entries outside the sliding window.
		_ = ref.cache.ZRemRangeByScore(ctx, key, "-inf", windowStart)

		// Count remaining entries within the window.
		count, err := ref.cache.ZCard(ctx, key)
		if err != nil {
			slog.Error("rate limiter: failed to count requests", "error", err, "ip", ip)
			// On cache failure, allow the request through.
			next.ServeHTTP(w, r)
			return
		}

		if int(count) >= ref.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Add the current request with the current timestamp as score.
		member := fmt.Sprintf("%d", now.UnixNano())
		_ = ref.cache.ZAdd(ctx, key, float64(now.UnixNano()), member)

		// Set key expiry so Redis auto-cleans after the window elapses.
		_ = ref.cache.Expire(ctx, key, ref.window)

		next.ServeHTTP(w, r)
	})
}

// getIP extracts the IP address from the request
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// RequestLogger logs incoming requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
