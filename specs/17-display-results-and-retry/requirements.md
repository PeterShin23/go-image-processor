# Story 17 — Display Results and Retry Failures

**Time budget:** ~1 hour (Day 10)

## Goal

As a user, I want to inspect the generated outputs for each source image, tell full
success from partial from failure, and retry a single failed file.

## Files to work in (`/src/web`)

- `src/components/FileList.tsx` — render results per file
- `src/hooks/useUploadQueue.ts` — add `retry(id)`
- optionally a `VariantList` sub-component and some CSS

## Display for completed files

original filename · asset id · overall status · each variant's (name, width, height,
format, size) · any partial errors.

## Tasks

- [ ] Render each file's status: `completed`, `partial` (has variants AND errors),
      `failed`.
- [ ] Make `partial` visually distinct (e.g. a yellow badge) from full success/failure.
- [ ] For completed/partial files, list the variants from the manifest.
- [ ] For failed files, show the error message.
- [ ] Add a retry button on failed files that re-runs ONLY that file.
- [ ] Retrying must not resend already-successful files.

## Acceptance criteria

- [ ] Results are grouped by source file.
- [ ] Partial success is visually distinguishable.
- [ ] Failed files can be retried individually.
- [ ] Retrying does not resend successful files.
- [ ] The interface stays usable with at least 10 files.
- [ ] `npm run build` succeeds.

## Definition of done

A run with a mix of good and bad files shows per-file variants, a distinct partial state,
and working per-file retry. Commit: `story 17: results display + per-file retry`.

## Hints (not solutions)

- Derive the badge from `manifest.status` for completed uploads; a client `failed`
  (network/HTTP) is separate from a server `partial`.
- `retry(id)` = reset that file to `pending`/`uploading`, call `uploadImage` again, update
  by id — reuse `processOne`.
- Keep it simple: a `<ul>` of variant rows is enough; image thumbnails/preview are a stretch
  goal.
