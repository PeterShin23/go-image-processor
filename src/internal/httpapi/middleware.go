package httpapi

import (
	"log/slog"
	"net/http"
)

// LoggingMiddleware wraps a handler to emit one structured log line per request
// (method, path, status, duration). The logger is injected — no global logger.
//
// STORY 13: implement this. Capture the status code (wrap ResponseWriter),
// time the request, and log with logger.Info(...). Never log image bytes.
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO (story 13): record start time; wrap w to capture the status;
		// call next.ServeHTTP; then logger.Info("request", "method", r.Method,
		// "path", r.URL.Path, "status", sw.status, "duration", time.Since(start)).
		next.ServeHTTP(w, r)
	})
}

// statusRecorder wraps http.ResponseWriter to remember the status code written.
//
// TODO (story 13): implement WriteHeader to store the code, then delegate.
type statusRecorder struct {
	http.ResponseWriter
	status int
}
