package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		got, err := Load(func(string) string { return "" })
		if err != nil {
			t.Fatal(err)
		}
		if want := Default(); got != want {
			t.Fatalf("Load() = %+v, want %+v", got, want)
		}
	})

	t.Run("overrides", func(t *testing.T) {
		env := map[string]string{
			"PORT":                      "9090",
			"MAX_IMAGE_SIZE_BYTES":      "2048",
			"MIN_IMAGE_WIDTH":           "2",
			"MIN_IMAGE_HEIGHT":          "3",
			"MAX_IMAGE_WIDTH":           "4000",
			"MAX_IMAGE_HEIGHT":          "5000",
			"ORIGINALS_DIR":             "/originals",
			"GENERATED_DIR":             "/generated",
			"MAX_CONCURRENT_TRANSFORMS": "8",
			"HTTP_TIMEOUT":              "30s",
			"SHUTDOWN_TIMEOUT":          "20s",
		}
		got, err := Load(func(key string) string { return env[key] })
		if err != nil {
			t.Fatal(err)
		}
		want := Config{
			Port:              9090,
			MaxImageSizeBytes: 2048,
			MinImageWidth:     2,
			MinImageHeight:    3,
			MaxImageWidth:     4000,
			MaxImageHeight:    5000,
			OriginalsDir:      "/originals",
			GeneratedDir:      "/generated",
			MaxConcurrent:     8,
			HTTPTimeout:       30 * time.Second,
			ShutdownTimeout:   20 * time.Second,
		}
		if got != want {
			t.Fatalf("Load() = %+v, want %+v", got, want)
		}
	})

	for _, tt := range []struct {
		key   string
		value string
	}{
		{key: "PORT", value: "abc"},
		{key: "MAX_IMAGE_SIZE_BYTES", value: "abc"},
		{key: "HTTP_TIMEOUT", value: "nope"},
	} {
		t.Run("invalid "+tt.key, func(t *testing.T) {
			_, err := Load(func(key string) string {
				if key == tt.key {
					return tt.value
				}
				return ""
			})
			if err == nil || !strings.Contains(err.Error(), tt.key) {
				t.Fatalf("Load() error = %v, want error containing %q", err, tt.key)
			}
		})
	}
}
