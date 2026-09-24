# Background — Story 13: Add Structured Logging

Logging is how a running service explains itself. You'll learn Go's standard structured
logger `log/slog`, the middleware pattern for cross-cutting concerns, and the discipline of
logging *facts* (not secrets, not image bytes).

## Structured vs unstructured logging

Unstructured: `log.Printf("processed %s with %d variants", name, n)` — a human-readable
sentence that machines can't query.

Structured: key/value pairs a log system can filter and aggregate:

```json
{"time":"...","level":"INFO","msg":"processed asset","asset_id":"a1b2","variants_ok":3,"variants_failed":0,"duration_ms":142}
```

You can later ask "show me all assets where `variants_failed > 0`" without regex-parsing
prose. This is why every serious service logs structured. Go's built-in answer is `log/slog`
(standard library since 1.21 — no dependency needed).

## `log/slog` basics

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // JSON lines to stdout
logger.Info("processed asset",
    "asset_id", assetID,
    "filename", filename,
    "variants_requested", len(profiles),
    "variants_ok", len(result.Variants),
    "variants_failed", len(result.Errors),
    "duration_ms", time.Since(start).Milliseconds(),
)
```

- A **handler** decides the output format: `NewJSONHandler` (machine-friendly) or
  `NewTextHandler` (nicer for local dev). Pick per environment.
- `logger.Info(msg, key, value, key, value, ...)` takes a message then alternating
  key/value pairs. For type safety you can use attrs: `slog.String("asset_id", id)`,
  `slog.Int("variants_ok", n)`.
- Levels: `Debug`, `Info`, `Warn`, `Error`. Use `Error` for failures with an `"err",
  err.Error()` field.

## Inject the logger — no globals

`slog` has a package-global default (`slog.Info(...)`), and it's tempting to use it
everywhere. Don't — global mutable state is exactly what the architecture rules forbid.
Instead, create the logger in `main` and pass it into the things that need it:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
handlers := httpapi.NewHandlers(svc, cfg.MaxImageSizeBytes, logger)
mux := handlers.Routes()
srv.Handler = httpapi.LoggingMiddleware(logger, mux)
```

Injecting the logger makes dependencies explicit and lets tests capture logs (point the
handler at a `bytes.Buffer`). Same principle as config and storage: pass it in.

## Middleware: cross-cutting behavior

Some concerns — logging, auth, request IDs — apply to *every* request. You don't want to
copy that code into each handler. **Middleware** wraps a handler with extra behavior:

```go
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(sr, r) // run the wrapped handler, capturing its status
        logger.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", sr.status,
            "duration_ms", time.Since(start).Milliseconds(),
        )
    })
}
```

The shape is: take a `next http.Handler`, return a *new* handler that does something before
and/or after calling `next.ServeHTTP`. Wrapping is composable — you can nest middlewares.
Here we wrap the whole mux once, so every route (including health) gets logged.

## Capturing the status code

`http.ResponseWriter` doesn't expose the status that was written, so to log it you wrap the
writer and intercept `WriteHeader`:

```go
type statusRecorder struct {
    http.ResponseWriter // embedded: inherits Write, Header, etc.
    status int
}
func (sr *statusRecorder) WriteHeader(code int) {
    sr.status = code
    sr.ResponseWriter.WriteHeader(code) // still write it for real
}
```

This uses **embedding**: by embedding `http.ResponseWriter`, `statusRecorder` automatically
has all its methods, and you override just `WriteHeader`. Default `status` to 200 because a
handler that writes a body without calling `WriteHeader` implicitly returns 200 — you must
record that case too. Embedding-to-decorate is a very common Go pattern; recognize it.

## What NOT to log

Two firm rules for this project:
1. **Never log image bytes** or the raw upload — it's large, useless in logs, and could
   contain sensitive content. Log *metadata* (filename, size, dimensions, counts), never
   pixels.
2. **Don't leak internals to clients**, but *do* log them. The 500 response says "internal
   error"; the log line carries the real wrapped error with its stage. Your story-03/06
   error wrapping (`fmt.Errorf("decode: %w", err)`) already tags the stage — surface that in
   the `"err"` field so a failed upload is diagnosable.

## Tying it to processing

The most useful log line in this system is the per-asset summary after `svc.Process`:
asset id, filename, requested/ok/failed variant counts, and processing duration. You can log
it from the handler (it has the manifest and can time the call) or thread the logger into the
service. Either is fine; logging from the handler keeps the service transport-agnostic and is
simplest. Include the asset ID so a log line can be traced to files on disk.

## Try this

Implement `statusRecorder.WriteHeader` and `LoggingMiddleware`, wrap your mux, and hit the
server — watch the JSON request lines appear. Then add the processing summary line after
`Process`. Write one test that runs a request through the middleware with a
`bytes.Buffer`-backed logger and asserts the log line contains the status. Run
`go doc log/slog` to see the attr helpers.
