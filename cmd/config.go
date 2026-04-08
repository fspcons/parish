package main

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Port          string
	LogLevel      string
	CookieSecure  bool
	CorsOrigin    string
	TLSCert       string
	TLSKey        string
	AdminEmail    string
	AdminPassword string
	RedisURL      string

	// Local is true when running in development (log-only e-mail). When false, SendGrid is used for password reset.
	Local bool
	// SendGridAPIKey and EmailFrom are required when Local is false (set via Secret Manager or env in production).
	SendGridAPIKey string
	// SendGridAPIURL is the SendGrid v3 mail send URL; empty means use the default production URL.
	SendGridAPIURL string
	EmailFrom      string
}

func setupLogger() {
	// Setup structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}

func setupConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	cookieSecure := true
	if strings.EqualFold(os.Getenv("COOKIE_SECURE"), "false") {
		cookieSecure = false
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	local := strings.EqualFold(os.Getenv("LOCAL_DEV"), "true") ||
		strings.EqualFold(os.Getenv("ENV"), "development") ||
		strings.EqualFold(os.Getenv("ENV"), "local") ||
		os.Getenv("FIRESTORE_EMULATOR_HOST") != ""

	return &Config{
		Port:           port,
		LogLevel:       logLevel,
		CookieSecure:   cookieSecure,
		CorsOrigin:     corsOrigin,
		TLSCert:        os.Getenv("TLS_CERT"),
		TLSKey:         os.Getenv("TLS_KEY"),
		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
		RedisURL:       redisURL,
		Local:          local,
		SendGridAPIKey: os.Getenv("SENDGRID_API_KEY"),
		SendGridAPIURL: os.Getenv("SENDGRID_API_URL"),
		EmailFrom:      os.Getenv("EMAIL_FROM"),
	}
}
