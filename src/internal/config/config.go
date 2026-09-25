// Package config loads and validates application configuration from environment
// variables, providing sensible defaults for local development.
//
package config

import (
	"fmt"
	"strconv"
	"time"
)

// Config holds all runtime settings. Zero image logic lives here.
type Config struct {
	Port              int
	MaxImageSizeBytes int64
	MinImageWidth     int
	MinImageHeight    int
	MaxImageWidth     int
	MaxImageHeight    int
	OriginalsDir      string
	GeneratedDir      string
	MaxConcurrent     int
	HTTPTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

// Default returns the config used when no environment variables are set.
func Default() Config {
	return Config{
		Port:              8080,
		MaxImageSizeBytes: 10 << 20,
		MinImageWidth:     1,
		MinImageHeight:    1,
		MaxImageWidth:     10000,
		MaxImageHeight:    10000,
		OriginalsDir:      "./data/originals",
		GeneratedDir:      "./data/generated",
		MaxConcurrent:     4,
		HTTPTimeout:       15 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}
}

func intEnv(getenv func(string) string, key string, fallback int) (int, error) {
	v := getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return n, nil
}

func int64Env(getenv func(string) string, key string, fallback int64) (int64, error) {
	v := getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return n, nil
}

func durationEnv(getenv func(string) string, key string, fallback time.Duration) (time.Duration, error) {
	v := getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return d, nil
}

// Load builds a Config, starting from Default and overriding with any values
// present in the environment. getenv is injected so tests can supply a fake
// environment without touching real process state.
func Load(getenv func(string) string) (Config, error) {
	cfg := Default()
	var err error

	if cfg.Port, err = intEnv(getenv, "PORT", cfg.Port); err != nil {
		return Config{}, err
	}
	if cfg.MaxImageSizeBytes, err = int64Env(getenv, "MAX_IMAGE_SIZE_BYTES", cfg.MaxImageSizeBytes); err != nil {
		return Config{}, err
	}
	if cfg.MinImageWidth, err = intEnv(getenv, "MIN_IMAGE_WIDTH", cfg.MinImageWidth); err != nil {
		return Config{}, err
	}
	if cfg.MinImageHeight, err = intEnv(getenv, "MIN_IMAGE_HEIGHT", cfg.MinImageHeight); err != nil {
		return Config{}, err
	}
	if cfg.MaxImageWidth, err = intEnv(getenv, "MAX_IMAGE_WIDTH", cfg.MaxImageWidth); err != nil {
		return Config{}, err
	}
	if cfg.MaxImageHeight, err = intEnv(getenv, "MAX_IMAGE_HEIGHT", cfg.MaxImageHeight); err != nil {
		return Config{}, err
	}
	if v := getenv("ORIGINALS_DIR"); v != "" {
		cfg.OriginalsDir = v
	}
	if v := getenv("GENERATED_DIR"); v != "" {
		cfg.GeneratedDir = v
	}
	if cfg.MaxConcurrent, err = intEnv(getenv, "MAX_CONCURRENT_TRANSFORMS", cfg.MaxConcurrent); err != nil {
		return Config{}, err
	}
	if cfg.HTTPTimeout, err = durationEnv(getenv, "HTTP_TIMEOUT", cfg.HTTPTimeout); err != nil {
		return Config{}, err
	}
	if cfg.ShutdownTimeout, err = durationEnv(getenv, "SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
