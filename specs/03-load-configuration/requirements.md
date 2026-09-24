# Story 03 — Load Application Configuration

**Time budget:** ~1.5 hours (Day 2)

## Goal

As an intern, I want configuration to come from environment variables so that behavior
can change without editing source code. Defaults must make local development work with
zero setup.

## Files to work in (`/src`)

- `internal/config/config.go` — `Config`, `Default`, and `Load` (starter provided)
- `internal/config/config_test.go` — you create this
- `.env.example` — already lists every setting; keep it in sync if you add fields

## Tasks

- [ ] Fill in `Default()` with the values from `.env.example` (port 8080, 10 MiB max
      size, dirs `./data/originals` and `./data/generated`, 4 concurrent transforms,
      `15s` HTTP timeout, `10s` shutdown timeout, sensible dimension bounds).
- [ ] Implement `Load(getenv func(string) string) (Config, error)`:
  - start from `Default()`
  - for each setting, if `getenv(KEY)` is non-empty, parse and override
  - integers via `strconv.Atoi`, 64-bit via `strconv.ParseInt`, durations via
    `time.ParseDuration`
  - wrap parse errors so the message names the offending variable
- [ ] Write `config_test.go` that:
  - checks defaults are returned when the fake env is empty
  - checks an override is applied (e.g. `PORT=9090`)
  - checks an invalid value (e.g. `PORT=abc`, `HTTP_TIMEOUT=nope`) returns a clear error

## Acceptance criteria

- [ ] Optional settings have documented defaults.
- [ ] Invalid numeric values produce clear errors (naming the variable).
- [ ] Invalid durations produce clear errors.
- [ ] `Config` is passed into services, not read from a global.
- [ ] Config is not accessed through package-level mutable variables.
- [ ] Tests isolate environment changes by passing a fake `getenv` — no reliance on the
      real process environment.
- [ ] `go test ./internal/config/` passes.

## Definition of done

Config loads with defaults and overrides, tests pass, `go vet ./...` clean. Commit:
`story 03: env-driven configuration`.

## Hints (not solutions)

- A fake env in tests is just a map: `func(k string) string { return m[k] }`.
- Wrap with `%w`: `fmt.Errorf("PORT: %w", err)` — story 03's background explains why.
- Keep `Load` boring and linear; a small helper like
  `parseIntEnv(getenv, "PORT", def)` reduces repetition but isn't required.
