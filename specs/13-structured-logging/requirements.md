# Story 13 — Add Structured Logging

**Time budget:** ~1.5 hours (Day 9)

## Goal

As an operator, I want useful, structured processing logs so I can understand what
happened to each upload — without ever logging image contents.

## Files to work in (`/src`)

- `internal/httpapi/middleware.go` — `LoggingMiddleware`, `statusRecorder` (starter provided)
- `internal/httpapi/httpapi.go` — accept a `*slog.Logger` and log processing outcomes
- `internal/app/app.go` — optionally accept a logger to log per-asset processing stats
- `cmd/api/main.go` and `cmd/cli/main.go` — create the `slog.Logger` and inject it

## Log fields to cover (across request + processing)

request method, request path, request duration, asset ID, original filename,
number of requested variants, number of successful variants, number of failed variants,
processing duration, error details.

## Tasks

- [ ] Implement `statusRecorder` (capture the status code) and `LoggingMiddleware`
      (time the request, log method/path/status/duration).
- [ ] Wrap your `Routes()` with `LoggingMiddleware` in `cmd/api`.
- [ ] Create a `slog.Logger` (JSON handler) in each `main` and inject it — no global logger.
- [ ] Log a processing summary after `svc.Process` (asset id, filename, requested/
      successful/failed variant counts, processing duration).
- [ ] Ensure errors identify the stage that failed (you already wrap with stage context —
      surface that).

## Acceptance criteria

- [ ] Logs are structured (key/value, e.g. JSON), not free-form `fmt.Println`.
- [ ] Every upload log carries an asset ID or request identifier.
- [ ] Errors identify the failing stage.
- [ ] Image bytes are never logged.
- [ ] Logging uses an injected logger, not global mutable state.

## Definition of done

Hitting `POST /v1/assets` prints a structured request line plus a processing summary line
with variant counts. Commit: `story 13: structured logging via slog, injected`.

## Hints (not solutions)

- `slog.New(slog.NewJSONHandler(os.Stdout, nil))` gives a JSON logger; use
  `slog.NewTextHandler` for human-friendly local logs.
- `logger.Info("request", "method", r.Method, "status", sr.status, "duration", d)` —
  alternating key/value pairs, or `slog.String(...)`, `slog.Int(...)` attrs.
- The status-capturing wrapper only needs to override `WriteHeader`; default it to 200 in
  case the handler never calls `WriteHeader`.
- Prefer passing the logger into constructors (`NewHandlers(svc, max, logger)`) over a
  package global.
