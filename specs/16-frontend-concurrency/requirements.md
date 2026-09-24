# Story 16 — Add Limited Frontend Concurrency

**Time budget:** ~1 hour (Day 10)

## Goal

As a user, I want a small number of files processed simultaneously without overwhelming
the backend. Default: **2 at once**.

## Files to work in (`/src/web`)

- `src/hooks/useUploadQueue.ts` — replace the sequential loop with a bounded worker pool

## Tasks

- [ ] Add a concurrency limit constant, default `2`, defined in ONE place.
- [ ] Process up to `limit` files at once; when one finishes, start the next pending file.
- [ ] Preserve the display order of selected files (order of *results in the list* stays
      stable even though completion order may vary).
- [ ] Keep each result attached to its correct file (update by id, never by position).
- [ ] One failure must not cancel unrelated in-flight or pending uploads.

## Acceptance criteria

- [ ] No more than 2 uploads run concurrently by default.
- [ ] One failure does not cancel other uploads.
- [ ] Every result stays attached to the correct file.
- [ ] Display order remains stable.
- [ ] The limit can be changed in one place.
- [ ] `npm run build` succeeds.

## Important distinction (know this for your presentation)

There are **two independent** concurrency controls in this project:
1. **React** limits how many *images* upload at once (this story).
2. **Go** limits how many *variants of one image* transform at once (story 09).

They solve different problems and must stay independent.

## Definition of done

Selecting 10 files processes them 2-at-a-time; the list order never jumps around. Commit:
`story 16: bounded frontend upload concurrency (2)`.

## Hints (not solutions)

- A simple pool: keep an index into the file list; spawn `limit` async "workers" that each
  pull the next unprocessed file until none remain.
- Or `Promise`-based: maintain a set of in-flight promises; `await Promise.race(...)` to
  learn when a slot frees, then launch the next.
- Because you update state by `id`, out-of-order completion is fine — the list render order
  comes from the `files` array, which you never reorder.
