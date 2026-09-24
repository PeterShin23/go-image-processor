# Story 10 — Create the Shared Application Service

**Time budget:** ~1.5 hours (Day 7) — **the architectural keystone**

## Goal

As an intern, I want ONE application service used by both the CLI and the HTTP server so
business logic is never duplicated. After this story, the CLI stops using its placeholder
and calls the real service. The API (story 11) will call the *same* service.

## Files to work in (`/src`)

- `internal/app/app.go` — `Service`, `NewService`, `ProcessInput`, `Process` (starter provided)
- `internal/app/app_test.go` — you create this
- `cmd/cli/main.go` — replace the story-04 placeholder with a real `Service` call

## The flow `Process` coordinates

1. Read the input bytes (bounded by the max size).
2. Validate the image (`media.Validate`).
3. Create an asset ID (`s.newID()`).
4. Store the original (`s.store.SaveOriginal`).
5. Decode the image once (`image.Decode`).
6. Generate variants (`s.pipe.Generate`).
7. Assemble a `domain.Manifest` (asset id, filename, variants, errors, derived status).

## Tasks

- [ ] Implement `NewService` including a default `newID` when the caller passes `nil`.
- [ ] Implement `Process` following the flow above; choose the profile set
      (`in.Profiles` if non-empty, else `s.profiles`).
- [ ] In `cmd/cli/main.go`, construct the dependencies (`storage.NewLocalStorage`,
      `pipeline.New`, `app.NewService`) and call `svc.Process`, then print the manifest
      as JSON.
- [ ] Write `app_test.go` with a fake or temp-dir storage and a real image, asserting a
      completed manifest with the expected number of variants.

## Acceptance criteria

- [ ] The CLI calls the application service (no image logic in `main`).
- [ ] `Process` contains no CLI flag parsing.
- [ ] `Process` contains no HTTP response writing.
- [ ] The service receives its dependencies through `NewService` (constructor injection).
- [ ] The service does **not** construct concrete storage itself — it's handed a
      `storage.Storage`.
- [ ] `go test ./...` passes.

## Definition of done

```bash
go run ./cmd/cli --input testdata/images/cafe.jpg --output ./data/generated
```

prints a JSON manifest and writes real variant files under `data/generated/<assetID>/`.
Commit: `story 10: shared application service; CLI now uses it`.

## Hints (not solutions)

- Read bytes once into memory (bounded), then build fresh readers from them for validate,
  store-original, and decode — the reader-consumed-once rule from story 06 applies.
- A default ID generator: `crypto/rand` + hex, or a timestamp+random string. Keep it
  URL/path-safe (it becomes a directory name).
- Encode the manifest with `encoding/json`: add `json:"..."` tags to the domain types (or
  a small response struct) and `json.NewEncoder(os.Stdout).Encode(manifest)`.
- The CLI's exit code should reflect failure: a `failed` status or a returned error → non-zero.
