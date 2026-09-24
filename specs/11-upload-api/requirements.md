# Story 11 — Create the Upload API

**Time budget:** ~1.5 hours (Day 8)

## Goal

As a client, I want to upload one image over HTTP so another application can use the
pipeline. Critically, the handler calls the **same** `app.Service` the CLI uses — no
image logic in the HTTP layer.

## Endpoint

```
POST /v1/assets     multipart/form-data, one image field (e.g. "image")
```

## Files to work in (`/src`)

- `internal/httpapi/httpapi.go` — `postAsset`, `Routes`, `writeJSON`, `writeJSONError`
  (starter provided)
- `internal/httpapi/httpapi_test.go` — you create this (use `httptest`)
- `cmd/api/main.go` — wire config → storage → pipeline → service → handlers → server

## Tasks

- [ ] In `Routes`, register `POST /v1/assets` → `postAsset`.
- [ ] In `postAsset`: wrap the body with `http.MaxBytesReader` to enforce `maxUpload`,
      parse the multipart form, read the image field with `r.FormFile`, build an
      `app.ProcessInput{Filename, Reader, Size}`, call `h.svc.Process`, and write the
      manifest as JSON with `writeJSON`.
- [ ] Implement `writeJSON` and `writeJSONError`.
- [ ] Map errors to status codes: bad/oversized/invalid input → 4xx (400/413/415),
      unexpected failures → 500. Use `errors.Is` against the `media` sentinels.
- [ ] Clean up: `defer file.Close()`; multipart temp files are released by the runtime,
      but close what you open.
- [ ] In `cmd/api/main.go`, construct everything and start `http.ListenAndServe` (timeouts
      and graceful shutdown come in story 12).
- [ ] Write tests using `net/http/httptest` that POST a real test image and assert 200 +
      a manifest, and that a non-image returns a 4xx.

## Acceptance criteria

- [ ] A valid image returns a success JSON response with a manifest.
- [ ] Invalid input returns a client-error (4xx) response.
- [ ] Internal failures return a server-error (5xx) response.
- [ ] All responses are JSON.
- [ ] Temporary multipart resources are cleaned up.
- [ ] Handlers stay thin (they call the service; they don't resize anything).
- [ ] Tests use `httptest` and pass: `go test ./internal/httpapi/`.

## Definition of done

```bash
go run ./cmd/api &
curl -X POST -F "image=@testdata/images/cafe.jpg" http://localhost:8080/v1/assets
```

returns a JSON manifest. Commit: `story 11: multipart upload API reusing the app service`.

## Hints (not solutions)

- Go 1.22+ supports method+path patterns: `mux.HandleFunc("POST /v1/assets", h.postAsset)`.
- `http.MaxBytesReader(w, r.Body, h.maxUpload)` makes oversized bodies fail as you read —
  detect that and return 413.
- `file, header, err := r.FormFile("image")` gives you the stream and `header.Filename`,
  `header.Size`.
- Reuse the JSON shape from story 10 so CLI and API emit identical manifests.
