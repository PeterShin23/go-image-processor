# Background — Story 09: Add Bounded Variant Concurrency

This is the story people mean when they say "Go is good at concurrency." You'll learn
goroutines, channels, `WaitGroup`, semaphores, the race detector, and context
cancellation. Go slowly — concurrency bugs are subtle, and the race detector is your
friend.

## Goroutines

A **goroutine** is a function running concurrently, scheduled by the Go runtime onto OS
threads. You start one with `go`:

```go
go doWork(profile) // returns immediately; doWork runs "in the background"
```

Goroutines are cheap (a few KB of stack), so you can have thousands. But "cheap" is not
"free" — and image transforms are CPU- and memory-heavy. Launching one goroutine per
profile is fine for 6 profiles; launching unbounded goroutines for arbitrary work is how
you exhaust RAM. Hence **bounded** concurrency.

The catch: `main`/the calling function does not wait for goroutines automatically. If
`Generate` returns before its goroutines finish, their results are lost. You need to
**wait**.

## `sync.WaitGroup`: wait for N goroutines

```go
var wg sync.WaitGroup
for _, profile := range profiles {
    wg.Add(1)          // register one goroutine
    go func(pr domain.MediaProfile) {
        defer wg.Done() // signal completion, even on panic/early return
        // ... do work with pr ...
    }(profile)          // pass profile as an ARGUMENT (see loop-variable note)
}
wg.Wait() // blocks until every Done() has been called
```

Rules:
- `wg.Add(1)` **before** the `go`, never inside the goroutine (a race otherwise).
- `defer wg.Done()` as the goroutine's first line so it always runs.
- `wg.Wait()` after the loop; only then are all results in.

### The loop-variable gotcha

Historically, `for _, profile := range profiles { go func(){ use(profile) }() }` was a
classic bug: every goroutine captured the *same* `profile` variable, which the loop kept
overwriting. Go 1.22+ fixed the loop semantics so each iteration gets a fresh variable —
but passing `profile` as a function argument (as above) is the timeless, unambiguous fix.
Do that and you never have to think about which Go version you're on.

## Channels: safe communication

A **channel** is a typed pipe between goroutines. Sending and receiving are synchronized,
so channels are how goroutines share data **without** a lock:

```go
results := make(chan itemResult) // unbuffered
results <- r                     // send (blocks until someone receives)
r := <-results                   // receive
```

The Go mantra: *"Don't communicate by sharing memory; share memory by communicating."*
Instead of every goroutine writing into a shared slice (which races), have each **send**
its outcome on a channel, and a single collector **receive** them:

```go
type item struct {
    variant *domain.Variant
    err     *domain.ProcessingError
}
results := make(chan item, len(profiles)) // buffered so senders never block
// ... each goroutine sends one item ...
go func() { wg.Wait(); close(results) }() // close when all senders done
for it := range results {                 // range ends when channel is closed
    if it.err != nil { result.Errors = append(result.Errors, *it.err) }
    if it.variant != nil { result.Variants = append(result.Variants, *it.variant) }
}
```

Because only the collector touches the slices, there's no race and no mutex needed.
`close(results)` after `wg.Wait()` lets the `range` loop terminate. Buffer the channel to
`len(profiles)` so a goroutine can send and exit without waiting for the collector.

## Bounding concurrency: a semaphore channel

To cap *simultaneous* work at `maxConcurrent`, use a buffered channel as a counting
semaphore:

```go
sem := make(chan struct{}, p.maxConcurrent) // capacity = the limit
// inside each goroutine, before heavy work:
sem <- struct{}{}          // acquire a slot (blocks if all slots taken)
defer func() { <-sem }()   // release the slot when done
```

`struct{}{}` is a zero-size value — the channel is used purely for counting, carrying no
data. When the buffer is full, `sem <- struct{}{}` blocks, so at most `maxConcurrent`
goroutines are past that line at once. This is the whole trick.

> Alternative: `golang.org/x/sync/semaphore` gives a weighted semaphore with context
> support (`sem.Acquire(ctx, 1)`), which also handles cancellation cleanly. The channel
> version is more instructive the first time; either satisfies the story.

## `sync.Mutex`: the lock (know it, maybe don't need it)

A `sync.Mutex` guards shared state so only one goroutine touches it at a time:

```go
var mu sync.Mutex
mu.Lock()
result.Variants = append(result.Variants, v)
mu.Unlock()
```

This is the "share memory with a lock" style — valid, but if you use the results-channel
pattern above you won't need a mutex at all. Prefer the channel; reach for a mutex only if
it genuinely reads cleaner. Whichever you choose, the race detector will tell you if you
got it wrong.

## Context cancellation

`context.Context` is Go's standard way to signal "stop, the caller gave up" (a timeout, a
cancelled HTTP request, Ctrl-C). It flows *down* through call chains as the first argument.

```go
select {
case <-ctx.Done():
    return domain.ProcessingResult{}, ctx.Err() // caller cancelled; stop early
default:
}
```

`ctx.Done()` is a channel that closes when the context is cancelled; `ctx.Err()` says why
(`context.Canceled` or `context.DeadlineExceeded`). Check it before starting each new
profile so you don't launch work nobody's waiting for. In the HTTP server (story 12), the
request's context is cancelled when the client disconnects — this is how that
disconnection reaches all the way down to skip pointless image work.

`select` waits on multiple channel operations; `default` makes it non-blocking (a quick
"is it cancelled yet?" poll). Without `default`, `select` blocks until one case is ready.

## The race detector

A **data race** is two goroutines touching the same memory concurrently with at least one
write — undefined behavior, and the source of the worst "works on my machine" bugs. Go
ships a detector:

```bash
go test -race ./...
```

It instruments memory accesses and *reports races it actually observes at runtime*. So it
can't prove absence of races — it only catches ones your test exercises. Therefore: write
a test that runs `Generate` with several profiles (enough to force real overlap) and run
it under `-race`. Make `-race` part of your habit; a green `-race` run is an acceptance
criterion here.

## Putting it together (structure, not solution)

1. Make a buffered results channel and a semaphore channel.
2. `for` over profiles: check `ctx`; `wg.Add(1)`; `go` a function that acquires the
   semaphore, does transform+store, sends its `item`, releases the semaphore, `Done`s.
3. A separate goroutine: `wg.Wait()` then `close(results)`.
4. `range results` in the main goroutine to aggregate.
5. Return the aggregated result.

Keep the per-profile work identical to story 08; you're only changing *how* it's scheduled
and collected.

## Try this

First run your story-08 tests under `go test -race ./...` (sequential code should be
race-clean). Then refactor to the channel-collector pattern with the semaphore. Re-run
`-race`. Add a test that cancels the context and asserts `Generate` returns promptly with
`ctx.Err()`. Run `go doc sync.WaitGroup` and `go doc context.Context` alongside.
