# Background — Story 16: Add Limited Frontend Concurrency

This mirrors the Go concurrency story (09) but in JavaScript. The concept — a bounded
worker pool — is identical; only the mechanics differ. Understanding both sides is a
highlight of your final presentation.

## Why bound it

If you `fetch` all 10 files at once, you hammer the backend with 10 simultaneous image jobs,
each spawning its own bounded variant work in Go. The browser also caps connections per host
(~6). Two concurrent uploads is a gentle, predictable load that still feels fast. The number
lives in **one constant** so it's trivially tunable:

```ts
const UPLOAD_CONCURRENCY = 2;
```

## The two-concurrency mental model

Keep these straight — it's a common interview/presentation question:

| Control | Where | Limits | Set in |
| --- | --- | --- | --- |
| Upload concurrency | React (this story) | how many *images* upload at once | `UPLOAD_CONCURRENCY` |
| Variant concurrency | Go (story 09) | how many *variants of one image* transform at once | `MAX_CONCURRENT_TRANSFORMS` |

They're independent: 2 images uploading, each generating (say) up to 4 variants
concurrently in Go, is 2×4 work in flight — but each layer is bounded on its own axis.

## A worker-pool pattern

The cleanest approach: start `limit` "workers" that each pull the next unprocessed file
until the list is exhausted. Shared cursor, no fancy coordination:

```ts
async function start() {
  let cursor = 0;
  const next = (): QueuedFile | undefined => files[cursor++]; // hand out the next file

  async function worker() {
    let qf: QueuedFile | undefined;
    while ((qf = next())) {
      await processOne(qf); // sets uploading -> completed/failed, catches its own errors
    }
  }

  const workers = Array.from({ length: UPLOAD_CONCURRENCY }, () => worker());
  await Promise.all(workers); // wait for all workers to drain the queue
}
```

How it bounds concurrency: there are only `UPLOAD_CONCURRENCY` workers, and each processes
one file at a time (`await`), so at most that many uploads are ever in flight. When a worker's
`processOne` resolves, the `while` loop grabs the next file — "start new work when a slot
frees" falls out naturally.

`Promise.all(workers)` waits for every worker to finish (each finishes when `next()` returns
`undefined`). Because `processOne` catches its own errors and never throws, one failure can't
reject `Promise.all` and cancel the others — that satisfies "one failure doesn't cancel
unrelated uploads." (If `processOne` *could* throw, you'd use `Promise.allSettled` instead.)

## Alternative: `Promise.race` for a slot

Another idiom maintains a set of in-flight promises and races them:

```ts
const inFlight = new Set<Promise<void>>();
for (const qf of files) {
  const p = processOne(qf).finally(() => inFlight.delete(p));
  inFlight.add(p);
  if (inFlight.size >= UPLOAD_CONCURRENCY) {
    await Promise.race(inFlight); // wait until any one finishes, freeing a slot
  }
}
await Promise.all(inFlight); // drain the last batch
```

Both patterns are fine. The worker-pool version is usually easier to reason about the first
time; pick whichever you can explain out loud.

## Keeping order stable and results correct

The subtle bit: with concurrency, files finish **out of order** — file 3 might complete
before file 1. Two things keep the UI correct:

1. **Render order comes from the `files` array**, which you never reorder. So the *list* stays
   in selection order regardless of completion order.
2. **Update by id, not index.** Every state update matches on `f.id === id`
   (`prev.map(f => f.id === id ? {...f, ...} : f)`). Position-based updates would attach a
   result to the wrong row once completion order diverges from list order.

You already update by id from story 15, so you mostly get this for free — just don't be
tempted to sort by completion time.

## Try this

Refactor `start()` from the sequential loop into the worker-pool version with
`UPLOAD_CONCURRENCY = 2`. Select ~10 files and watch: at most two show "uploading" at any
moment, the list order never jumps, and each row lands on its own correct result. Temporarily
set the constant to 1 to confirm it reproduces the sequential behavior, then back to 2.
