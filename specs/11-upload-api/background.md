# Background — Story 11: Create the Upload API

Now the same pipeline gets a second front door: HTTP. The learning here is Go's
`net/http`, multipart uploads, JSON responses, mapping errors to status codes, and
testing servers with `httptest`. The architectural lesson is that the handler is *thin* —
it's an adapter, not a worker.

## `net/http` basics

An HTTP handler is anything satisfying:

```go
type Handler interface { ServeHTTP(w http.ResponseWriter, r *http.Request) }
```

You rarely implement that directly; you write handler *functions* and register them on a
**mux** (router):

```go
mux := http.NewServeMux()
mux.HandleFunc("POST /v1/assets", h.postAsset) // Go 1.22+ method+path patterns
```

- `http.ResponseWriter` is how you write the response: set headers, `WriteHeader(status)`,
  then write the body. **Order matters**: set headers and status *before* writing the body;
  the first body write locks in a 200 if you haven't called `WriteHeader`.
- `*http.Request` holds everything about the request: method, URL, headers, and `Body`
  (an `io.ReadCloser`).

Note `h.postAsset` is a **method value** — it carries `h` (with its `svc`) along, which is
how the handler reaches the injected application service without any globals.

## Multipart uploads

Browsers and `curl -F` send files as `multipart/form-data`: a body split into parts, each
with its own headers. Go parses it for you:

```go
func (h *Handlers) postAsset(w http.ResponseWriter, r *http.Request) {
    r.Body = http.MaxBytesReader(w, r.Body, h.maxUpload) // cap the whole body

    file, header, err := r.FormFile("image") // "image" = the form field name
    if err != nil {
        writeJSONError(w, http.StatusBadRequest, "missing or invalid 'image' field")
        return
    }
    defer file.Close()

    in := app.ProcessInput{
        Filename: header.Filename, // e.g. "cafe.jpg"
        Reader:   file,            // a multipart.File is an io.Reader
        Size:     header.Size,
    }
    manifest, err := h.svc.Process(r.Context(), in)
    // ...map err to a status, else writeJSON(w, 200, manifest)...
}
```

Key points:
- `r.FormFile(name)` returns the file stream, a header (filename + size), and an error. It
  triggers multipart parsing under the hood; large uploads may spill to temp files, which
  the runtime cleans up (calling `file.Close()` is still correct hygiene).
- `multipart.File` implements `io.Reader`, so it drops straight into `ProcessInput.Reader` —
  the exact same `Process` the CLI calls, fed from a different source. *This is the payoff of
  story 10.*
- `r.Context()` is the request's context. It's cancelled if the client disconnects, and it
  flows into `Process` → `pipeline.Generate` → the `ctx` checks you wrote in story 09.

## Enforcing the size limit

`http.MaxBytesReader(w, r.Body, n)` wraps the body so reads past `n` bytes fail with a
specific error and it signals the client. Enforce it *before* you buffer the whole thing, so
a 1 GB "image" can't exhaust memory. When the limit trips, respond `413 Request Entity Too
Large`.

## Mapping errors to status codes

The handler's job on failure is to translate a Go error into the right HTTP status. Reuse
the typed sentinels from story 06 with `errors.Is`:

```go
switch {
case errors.Is(err, media.ErrUnsupportedFormat):
    writeJSONError(w, http.StatusUnsupportedMediaType, err.Error()) // 415
case errors.Is(err, media.ErrTooLarge):
    writeJSONError(w, http.StatusRequestEntityTooLarge, err.Error()) // 413
case errors.Is(err, media.ErrEmptyFile), errors.Is(err, media.ErrCorruptImage),
     errors.Is(err, media.ErrDimensionsRange):
    writeJSONError(w, http.StatusBadRequest, err.Error()) // 400
case err != nil:
    writeJSONError(w, http.StatusInternalServerError, "internal error") // 500 — don't leak internals
}
```

Two habits worth forming: (1) client mistakes are 4xx, server mistakes are 5xx; (2) for 5xx,
send a generic message to the client and keep the detailed error in your logs (story 13) —
don't leak stack traces or internal paths to callers.

A *partial* success (some variants failed) is **not** an HTTP error — return 200 with a
manifest whose `status` is `partial`. The HTTP status describes the request's handling; the
manifest describes the media outcome.

## Writing JSON responses

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}
```

Set the header first, then the status, then stream the body. `any` is Go's alias for
`interface{}` (Go 1.18+) — "any value." `json.NewEncoder(w).Encode(v)` writes directly to
the response without building an intermediate string.

## Testing with `httptest` (no real network)

`net/http/httptest` lets you exercise handlers in-process — fast and deterministic:

```go
func TestPostAsset(t *testing.T) {
    // build a multipart body containing a real test image
    var body bytes.Buffer
    mw := multipart.NewWriter(&body)
    fw, _ := mw.CreateFormFile("image", "cafe.jpg")
    f, _ := os.Open("testdata/images/cafe.jpg"); io.Copy(fw, f); f.Close()
    mw.Close()

    req := httptest.NewRequest("POST", "/v1/assets", &body)
    req.Header.Set("Content-Type", mw.FormDataContentType()) // sets the multipart boundary
    rec := httptest.NewRecorder()

    handlers.Routes().ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
    }
    // decode rec.Body into a manifest and assert on it
}
```

- `httptest.NewRequest` builds a request without a socket.
- `httptest.NewRecorder` captures the response (`.Code`, `.Body`, `.Header()`).
- `mw.FormDataContentType()` gives the exact `Content-Type` with the boundary — forgetting
  this header is the #1 reason multipart tests fail.
- Wire the handler through your real `Routes()` so you also test routing.

## Wiring `cmd/api/main.go`

The API's composition root mirrors the CLI's:

```go
cfg, _ := config.Load(os.Getenv)
store := storage.NewLocalStorage(cfg.OriginalsDir, cfg.GeneratedDir)
pipe := pipeline.New(store, cfg.MaxConcurrent)
svc := app.NewService(store, pipe, limitsFromConfig(cfg), defaultProfiles(), nil)
handlers := httpapi.NewHandlers(svc, cfg.MaxImageSizeBytes)

http.ListenAndServe(":"+strconv.Itoa(cfg.Port), handlers.Routes())
```

Same `svc` construction as the CLI — ideally factor that shared wiring into a small helper
both `main`s call, so there's truly one way to build the service. `ListenAndServe` blocks,
serving until the process dies. Story 12 replaces it with a configured `http.Server` plus
graceful shutdown.

## Try this

Implement `writeJSON`/`writeJSONError`, then `postAsset` for the happy path only, then the
`httptest` success test. Once green, add the error mapping and a test posting a `.txt` file
to confirm you get a 4xx. Wire `cmd/api` last and try the real `curl`.
