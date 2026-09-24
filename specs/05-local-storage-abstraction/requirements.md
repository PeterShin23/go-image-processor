# Story 05 — Create the Local Storage Abstraction

**Time budget:** ~1.5–3 hours (Days 3–4)

## Goal

As an intern, I want storage behind an interface so local disk could later be swapped
for object storage without touching the rest of the app. You implement disk storage
now; the interface keeps that choice replaceable.

## Files to work in (`/src`)

- `internal/storage/storage.go` — `Storage` interface + `LocalStorage` (starter provided)
- `internal/storage/storage_test.go` — you create this (use `t.TempDir()`)

## Tasks

- [ ] Implement `NewLocalStorage` (starter provided) and the three methods:
  - `SaveOriginal(ctx, assetID, filename, r)` → write under `originalsDir/<assetID>/`
  - `SaveVariant(ctx, assetID, variantName, format, r)` → write under
    `generatedDir/<assetID>/<variantName>.<ext>`
  - `Open(ctx, path)` → return a reader for a stored file
- [ ] Create the asset directory if missing (`os.MkdirAll`).
- [ ] Copy the reader to the destination file (`io.Copy`) and return the path and byte
      count.
- [ ] Build safe paths: a caller-supplied filename must **not** be able to escape the
      storage root (no `../../etc/passwd`). Sanitize with `filepath.Base` and/or verify
      the cleaned absolute path stays within the root.

## Acceptance criteria

- [ ] Tests use temporary directories (`t.TempDir()`), never the real `data/` dir.
- [ ] Original and generated files land in separate roots.
- [ ] Different asset IDs produce isolated directories.
- [ ] A malicious filename (e.g. `../../escape.txt`) cannot write outside the root — add
      a test proving this.
- [ ] Storage errors are wrapped with useful context (which path/op failed).
- [ ] The interface exposes behavior, not filesystem internals (no `*os.File` leaking).
- [ ] `go test ./internal/storage/` passes.

## Definition of done

Round-trip test passes: save bytes, `Open` them back, compare. Path-traversal test
passes. Commit: `story 05: local filesystem storage behind an interface`.

## Hints (not solutions)

- `t.TempDir()` gives you a fresh directory auto-deleted after the test.
- `filepath.Join(root, assetID, name)` builds paths portably; then `filepath.Clean`
  and check the result still has `root` as a prefix.
- `io.Copy(dst, src)` returns the number of bytes copied — that's your `size`.
- Return `(path, size, err)`; wrap errors like `fmt.Errorf("save original %q: %w", ...)`.
