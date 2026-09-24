# Background — Story 12: Add Health Endpoints and Graceful Shutdown

This story is about the *lifecycle* of a server: how it starts, how it's monitored, and —
the tricky part — how it stops without hurting anyone. You'll learn OS signals, context
deadlines, and the goroutine dance that makes graceful shutdown work.

## Health endpoints (the easy part)

Two conventional probes:

- **Liveness** (`/health/live`): "is the process alive?" If this fails, a supervisor
  restarts you. Keep it trivial — return 200 unconditionally.
- **Readiness** (`/health/ready`): "can the process serve traffic right now?" A supervisor
  won't route requests to you until this passes. Keep it cheap; for this project a plain
  200 is fine (you have no external DB to check).

Both just `writeJSON(w, 200, {...})`. The value is operational: orchestrators and load
balancers poll these to decide whether to send you traffic or restart you.

## Configuring `http.Server` with timeouts

`http.ListenAndServe` (story 11) uses a default server with **no timeouts** — a slow or
malicious client can hold a connection open forever and exhaust your resources. In real
life you always construct the server explicitly:

```go
srv := &http.Server{
    Addr:              ":" + strconv.Itoa(cfg.Port),
    Handler:           handlers.Routes(),
    ReadTimeout:       cfg.HTTPTimeout,       // max time to read the whole request
    ReadHeaderTimeout: 5 * time.Second,       // max time to read just the headers
    WriteTimeout:      cfg.HTTPTimeout,       // max time to write the response
    IdleTimeout:       60 * time.Second,      // keep-alive idle limit
}
```

Each timeout closes a specific abuse vector. `ReadHeaderTimeout` in particular defends
against the "Slowloris" attack (dribbling headers to tie up connections). You don't need to
memorize the exact values — just know every production server sets all four.

## OS signals

When you press `Ctrl+C`, the OS sends your process `SIGINT`. Container platforms send
`SIGTERM` before killing a pod. By default these terminate the process *immediately*,
dropping any in-flight requests. Graceful shutdown means: catch the signal, stop accepting
new connections, let existing requests finish (up to a deadline), then exit.

Go's modern way to catch signals is `signal.NotifyContext`, which hands you a context that
cancels when a signal arrives:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

Now `ctx.Done()` fires on Ctrl+C or SIGTERM. This unifies "someone signaled us" with the
same context machinery you already use for cancellation — elegant.

## The shutdown dance

The problem: `srv.ListenAndServe()` **blocks** forever (it's serving). So you can't call it
and *also* wait for a signal on the same goroutine. The standard pattern runs the server on
a background goroutine and keeps `main` free to wait:

```go
srvErr := make(chan error, 1)
go func() {
    srvErr <- srv.ListenAndServe() // blocks here until the server stops
}()

select {
case err := <-srvErr:
    // server failed to start / crashed
    return err
case <-ctx.Done():
    // a signal arrived; begin graceful shutdown
}

shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
defer cancel()
if err := srv.Shutdown(shutdownCtx); err != nil {
    log.Printf("graceful shutdown failed: %v", err) // log, don't panic
    return err
}
```

What each piece does:
- The server runs in a goroutine; its return value (an error) comes back on `srvErr`.
- `select` waits for *either* the server dying *or* a shutdown signal — whichever first.
- `srv.Shutdown(ctx)` stops accepting new connections and waits for active requests to
  finish, but no longer than the context's deadline. If the deadline passes, it returns an
  error and any still-running requests are cut off. That deadline is your `ShutdownTimeout`.

## `ErrServerClosed` is not an error

When `Shutdown` completes, the blocked `ListenAndServe` returns `http.ErrServerClosed`.
That's the *expected* signal of a clean stop — treat it as success:

```go
if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
    // a real startup/serving error
}
```

Forgetting this makes clean shutdowns look like failures in your logs. It's a classic gotcha.

## Why context deadlines matter here

`context.WithTimeout` returns a context that auto-cancels after the duration, plus a
`cancel` function you must always `defer` to release resources. During shutdown it bounds
how long you'll wait for stragglers: "finish up, but I'm leaving in 10 seconds." This is the
same `Context` type flowing through your whole app — timeouts, cancellation, and signals all
speak it. That consistency is a big part of why Go services compose cleanly.

## Testing this

You can't easily unit-test `Ctrl+C`, and that's okay — verify it by hand
(`go run ./cmd/api`, then Ctrl+C, watch it log and exit 0). The health endpoints *are*
unit-testable with `httptest` exactly like story 11: request `/health/live`, assert 200 and
the JSON body. Do that much.

## Try this

Register and test the two health routes first (quick `httptest` wins). Then rebuild
`cmd/api` around an `http.Server` with timeouts and the signal/shutdown dance. Run it,
hit `/health/live` with curl, then Ctrl+C and confirm a clean exit. Run
`go doc signal.NotifyContext` and `go doc http.Server.Shutdown`.
