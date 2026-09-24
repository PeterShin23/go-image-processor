// Package config loads and validates application configuration from environment
// variables, providing sensible defaults for local development.
//
// STORY 03: implement Load. Keep config a plain value that gets passed into
// services — do NOT expose a global singleton.
package config

import (
	"time"

	"github.com/PeterShin23/go-image-processor/internal/domain"
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
//
// TODO (story 03): fill these with the defaults documented in .env.example.
func Default() Config {
	return Config{
		// TODO: Port: 8080, MaxImageSizeBytes: 10 << 20, etc.
	}
}

// Load builds a Config, starting from Default and overriding with any values
// present in the environment. getenv is injected so tests can supply a fake
// environment without touching real process state.
//
// TODO (story 03): for each setting, if getenv returns a non-empty string,
// parse it (strconv.Atoi / strconv.ParseInt / time.ParseDuration) and wrap
// parse failures with context. Return the finished Config or an error.
func Load(getenv func(string) string) (Config, error) {
	return Config{}, domain.ErrNotImplemented
}
