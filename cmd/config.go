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

	return &Config{
		Port:          port,
		LogLevel:      logLevel,
		CookieSecure:  cookieSecure,
		CorsOrigin:    corsOrigin,
		TLSCert:       os.Getenv("TLS_CERT"),
		TLSKey:        os.Getenv("TLS_KEY"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		RedisURL:      redisURL,
	}
}
