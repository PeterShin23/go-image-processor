# Story 09 — Add Bounded Variant Concurrency

**Time budget:** ~1.5 hours (Day 7)

## Goal

As an intern, I want variants processed concurrently without exhausting system
resources. You'll turn the sequential loop from story 08 into a bounded worker pattern
that respects a concurrency limit and context cancellation — and prove it's race-free.

## Files to work in (`/src`)

- `internal/pipeline/pipeline.go` — rewrite `Generate` to run profiles concurrently,
  bounded by `p.maxConcurrent`
- `internal/pipeline/pipeline_test.go` — add a race/concurrency test

## Tasks

- [ ] Run independent profile transformations concurrently.
- [ ] Cap in-flight work at `p.maxConcurrent` (a counting semaphore via a buffered
      channel, or `golang.org/x/sync/semaphore`).
- [ ] Collect results safely — no data races on the shared `ProcessingResult`.
- [ ] Keep each result associated with its correct profile.
- [ ] Stop launching / short-circuit work when `ctx` is cancelled.
- [ ] Ensure every started goroutine exits (no leaks) — use a `sync.WaitGroup`.

## Acceptance criteria

- [ ] The concurrency limit is respected (never more than `maxConcurrent` transforms at
      once).
- [ ] No data races: `go test -race ./...` passes.
- [ ] Every result belongs to the correct profile.
- [ ] All goroutines exit normally.
- [ ] Context cancellation is respected.

## Definition of done

`go test -race ./...` is green and the pipeline still produces the same variants as
story 08. Commit: `story 09: bounded concurrent variant generation`.

## Hints (not solutions)

- A buffered channel `sem := make(chan struct{}, p.maxConcurrent)` is a semaphore:
  send before starting work, receive (in a `defer`) when done.
- Prefer collecting each goroutine's output over a **results channel**, then draining it
  after `wg.Wait()` — that sidesteps needing a mutex around the slices.
- If you *do* share the slices directly, guard them with a `sync.Mutex`.
- Check `ctx.Err() != nil` (or `select { case <-ctx.Done(): ... }`) to bail early.
- Run `go test -race` often; the race detector only reports races it actually observes,
  so exercise the concurrency with several profiles.
