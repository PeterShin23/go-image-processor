# Background — Story 03: Load Application Configuration

Configuration is where two big Go ideas show up for real: **error wrapping** and
**avoiding global state through dependency injection**. Both are habits you'll use for
the rest of the project.

## Reading the environment

The standard library reads env vars with `os.Getenv`:

```go
os.Getenv("PORT") // returns "" if unset
```

But notice our `Load` takes a parameter `getenv func(string) string` instead of calling
`os.Getenv` directly. That's deliberate. A function value can be passed around like any
other value. In production, `main` passes the real one:

```go
cfg, err := config.Load(os.Getenv)
```

In tests, you pass a fake backed by a map, so tests never mutate the real process
environment (which would make them flaky and order-dependent):

```go
env := map[string]string{"PORT": "9090"}
cfg, err := config.Load(func(k string) string { return env[k] })
```

This is **dependency injection** in its simplest form: the function depends on an
input you hand it, not on a hidden global. The acceptance criterion "tests isolate
environment changes" falls out for free.

> Go also offers `t.Setenv("PORT", "9090")` to set a real env var that auto-restores
> after the test. It works, but injecting `getenv` is cleaner and needs no cleanup.

## Parsing strings into typed values

Env vars are always strings. Convert them:

```go
import "strconv"

n, err := strconv.Atoi("8080")          // string -> int
b, err := strconv.ParseInt("10485760", 10, 64) // -> int64 (base 10, 64-bit)

import "time"
d, err := time.ParseDuration("15s")     // -> time.Duration
```

`time.Duration` is a named integer type counting nanoseconds. Duration strings look
like `"15s"`, `"1m30s"`, `"500ms"`. You'll compare and pass these around later for
HTTP timeouts.

## The defaults-then-override pattern

The clean shape for config loading:

1. Start from a fully-populated `Default()`.
2. For each key, only override if the env var is present and non-empty.

```go
cfg := Default()
if v := getenv("PORT"); v != "" {
    n, err := strconv.Atoi(v)
    if err != nil {
        return Config{}, fmt.Errorf("PORT: %w", err)
    }
    cfg.Port = n
}
// ... repeat per setting ...
return cfg, nil
```

Because you start from `Default()`, any variable the user didn't set keeps its default.
No setting is ever accidentally zero.

## Error wrapping with `%w`

In story 02 you built errors with `%v`. Now meet `%w`. When you wrap an underlying
error with `fmt.Errorf("PORT: %w", err)`, you produce a new error whose message is
`"PORT: invalid syntax"` **and** which still remembers the original error inside it.

Why care? Because callers can later ask "is this failure ultimately a specific known
error?" using `errors.Is`, or extract a typed error using `errors.As`:

```go
if errors.Is(err, strconv.ErrSyntax) { ... }
```

Rule of thumb:
- Use `%w` when you want the caller to be able to inspect the underlying cause.
- Use `%v` when you're just formatting a message and the cause doesn't need to travel.

For config, `%w` is the right call: it preserves the parse error while adding which
variable failed. You'll rely on wrapping heavily in the pipeline and app service, where
knowing *which stage* failed matters.

## Why no globals

You could stash config in a package-level `var Current Config` and read it anywhere.
Resist. Global mutable state makes code hard to test (tests fight over the global) and
hides dependencies (a function's real inputs become invisible). Instead, `Load` returns
a `Config` **value**, and `main` passes it into each constructor:

```go
cfg, _ := config.Load(os.Getenv)
store := storage.NewLocalStorage(cfg.OriginalsDir, cfg.GeneratedDir)
svc := app.NewService(store, ... , cfg.MaxConcurrent)
```

Every service states, in its constructor, exactly what config it needs. This is the
"pass dependencies in" discipline the whole project is built on.

## A tidy helper (optional)

If the repetition bugs you, extract a helper — but keep it readable:

```go
func intEnv(getenv func(string) string, key string, def int) (int, error) {
    v := getenv(key)
    if v == "" {
        return def, nil
    }
    n, err := strconv.Atoi(v)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", key, err)
    }
    return n, nil
}
```

Don't over-abstract on day one, though. Straight-line code that's obviously correct
beats a clever helper you have to decode.

## Try this

Run `go doc strconv.Atoi` and `go doc time.ParseDuration`. Implement `Default()` and the
`PORT` branch of `Load`, then write two tests: one asserting the default when the env is
empty, one asserting `PORT=9090` overrides it. Add a test for `PORT=abc` returning an
error. Then fill in the remaining settings the same way.
