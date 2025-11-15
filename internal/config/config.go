package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AppConfig holds runtime configuration derived from env vars.
type AppConfig struct {
	DatabaseURL      string
	HTTPAddr         string
	MetricsAddr      string
	SeedData         bool
	Environment      string
	GracefulShutdown time.Duration
}

// Load loads configuration with sane defaults for local dev.
func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable"),
		HTTPAddr:         getEnv("HTTP_ADDR", ":8080"),
		MetricsAddr:      getEnv("METRICS_ADDR", ":9090"),
		Environment:      getEnv("ENVIRONMENT", "local"),
		GracefulShutdown: 10 * time.Second,
	}

	seed := getEnv("SEED_DATA", "true")
	if v, err := strconv.ParseBool(seed); err == nil {
		cfg.SeedData = v
	} else {
		return nil, fmt.Errorf("invalid SEED_DATA value: %w", err)
	}

	if v := os.Getenv("GRACEFUL_SHUTDOWN_SECONDS"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil {
			cfg.GracefulShutdown = time.Duration(sec) * time.Second
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
