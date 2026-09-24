# Story 12 — Add Health Endpoints and Graceful Shutdown

**Time budget:** ~1.5 hours (Day 8)

## Goal

As an operator, I want to know whether the service is running and have it stop cleanly
without dropping in-flight requests.

## Endpoints

```
GET /health/live     liveness  (is the process up?)
GET /health/ready    readiness (can it serve?)
```

## Files to work in (`/src`)

- `internal/httpapi/httpapi.go` — `live`, `ready` (starter stubs), register them in `Routes`
- `cmd/api/main.go` — configured `http.Server`, signal handling, graceful shutdown

## Tasks

- [ ] Register `GET /health/live` → `live` and `GET /health/ready` → `ready` in `Routes`.
- [ ] Finish the health handlers to return JSON.
- [ ] In `cmd/api`, build an `http.Server{}` with **ReadTimeout, WriteTimeout,
      IdleTimeout, and ReadHeaderTimeout** set from config.
- [ ] Start the server in a goroutine so `main` can also wait for a shutdown signal.
- [ ] Listen for `SIGINT`/`SIGTERM` via `signal.NotifyContext`.
- [ ] On signal, call `server.Shutdown(ctx)` with a **deadline** context
      (`ShutdownTimeout`) so active requests get limited time to finish; then exit.
- [ ] Log (or print) shutdown errors.

## Acceptance criteria

- [ ] Health endpoints return JSON.
- [ ] The `http.Server` has read, write, idle, and header timeouts.
- [ ] Shutdown uses a context with a deadline.
- [ ] The process responds to `Ctrl+C` (SIGINT) and SIGTERM by shutting down.
- [ ] Shutdown failures are logged.

## Definition of done

```bash
go run ./cmd/api          # then Ctrl+C
# server logs "shutting down" and exits 0 without a panic
curl localhost:8080/health/live    # {"status":"ok"}
```

Commit: `story 12: health probes, server timeouts, graceful shutdown`.

## Hints (not solutions)

- `ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)`
  gives you a context cancelled on signal; `defer stop()`.
- `ListenAndServe` returns `http.ErrServerClosed` on a clean shutdown — treat that as
  success, not an error.
- Pattern: `go func(){ srv.ListenAndServe() }()`, then `<-ctx.Done()`, then
  `srv.Shutdown(shutdownCtx)`.
- `context.WithTimeout(context.Background(), cfg.ShutdownTimeout)` for the shutdown deadline;
  `defer cancel()`.
