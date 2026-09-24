# Background — Story 10: Create the Shared Application Service

This story is the reason the whole codebase is shaped the way it is. Everything you built
— domain, config, storage, media, pipeline — gets composed here into one reusable flow.
The big idea is **dependency injection**: the service is handed its collaborators and
orchestrates them, without knowing whether it was triggered by a CLI or an HTTP request.

## What "application service" means

Think of layers:

```
transport (CLI flags / HTTP multipart)  <- knows the outside world
        |
application service (this story)         <- the use case: "process one image"
        |
domain + media + pipeline + storage      <- the building blocks
```

The **application service** is the seam. It exposes one method, `Process`, that expresses
the use case in transport-neutral terms: "here are some bytes, a filename, and options —
give me back a manifest." The CLI adapts flags→`ProcessInput`; the HTTP handler adapts an
upload→`ProcessInput`. Neither adapter contains business logic, and the service contains no
adapter logic. That's how one service serves both.

## Constructor injection

`NewService(store, pipe, limits, profiles, newID)` receives everything it depends on:

```go
type Service struct {
    store    storage.Storage      // an INTERFACE — swappable, fakeable
    pipe     *pipeline.Pipeline
    limits   media.Limits
    profiles []domain.MediaProfile
    newID    func() string
}
```

Two rules the acceptance criteria enforce:

1. **The service does not build its own storage.** It receives a `storage.Storage`. If it
   called `storage.NewLocalStorage(...)` internally, it would be welded to local disk and
   untestable without a filesystem. Instead, `main` decides the concrete type and injects
   it. This is the payoff of the story-05 interface.
2. **Dependencies come in through the constructor**, not through globals or package init.
   Every collaborator the service needs is visible in `NewService`'s signature — no hidden
   inputs.

Injecting `newID func() string` is a small but instructive touch: production passes a
random generator, tests pass `func() string { return "test-asset" }` so the output is
deterministic and asserting on paths is easy.

## The reader-consumed-once problem, again

`Process` must read the bytes *three* times: to validate, to store the original, and to
decode. But an `io.Reader` streams once (story 06). The clean fix: read the bytes into
memory once, up front, bounded by the size limit:

```go
data, err := io.ReadAll(io.LimitReader(in.Reader, s.limits.MaxBytes+1))
if err != nil { return domain.Manifest{}, fmt.Errorf("read input: %w", err) }
// then make fresh readers as needed:
src, err := media.Validate(bytes.NewReader(data), int64(len(data)), s.limits)
// ... store original from bytes.NewReader(data) ...
img, _, err := image.Decode(bytes.NewReader(data))
```

Each `bytes.NewReader(data)` is a brand-new reader over the same underlying bytes, so you
can "re-read" as many times as you like. Bounding with `LimitReader` keeps a malicious huge
upload from eating all your memory before validation runs.

## Orchestration, step by step

`Process` is mostly glue — call each collaborator, wrap errors with the stage name, stop on
fatal errors, and assemble the result:

```go
func (s *Service) Process(ctx context.Context, in ProcessInput) (domain.Manifest, error) {
    profiles := s.profiles
    if len(in.Profiles) > 0 { profiles = in.Profiles }
    if err := domain.ValidateProfiles(profiles); err != nil {
        return domain.Manifest{}, fmt.Errorf("profiles: %w", err)
    }

    data, err := io.ReadAll(io.LimitReader(in.Reader, s.limits.MaxBytes+1))
    // ... handle err ...

    src, err := media.Validate(bytes.NewReader(data), int64(len(data)), s.limits)
    if err != nil { return domain.Manifest{}, fmt.Errorf("validate: %w", err) }

    assetID := s.newID()
    if _, _, err := s.store.SaveOriginal(ctx, assetID, in.Filename, bytes.NewReader(data)); err != nil {
        return domain.Manifest{}, fmt.Errorf("store original: %w", err)
    }

    img, _, err := image.Decode(bytes.NewReader(data))
    if err != nil { return domain.Manifest{}, fmt.Errorf("decode: %w", err) }

    result, err := s.pipe.Generate(ctx, assetID, img, profiles)
    if err != nil { return domain.Manifest{}, fmt.Errorf("generate: %w", err) }

    return domain.Manifest{
        AssetID:          assetID,
        OriginalFilename: in.Filename,
        Status:           result.Status(len(profiles)),
        Variants:         result.Variants,
        Errors:           result.Errors,
    }, nil
}
```

Notice: validation and decoding fatal-error out (a bad file means no manifest), but
per-variant problems flow through as `result.Errors` with a `partial` status. Distinguish
"the whole request is invalid" (return an error) from "some variants failed" (return a
manifest describing it). Story 11 maps the former to a 4xx/5xx and the latter to a 200 with
a partial status.

## Generating an asset ID

The ID becomes a directory name and appears in URLs later, so keep it filesystem- and
URL-safe. A simple, collision-resistant option:

```go
import ("crypto/rand"; "encoding/hex")
func newAssetID() string {
    b := make([]byte, 16)
    _, _ = rand.Read(b) // crypto/rand basically never fails on normal systems
    return hex.EncodeToString(b)
}
```

Set this as the default when `newID == nil` in `NewService`. (A UUID library would also
work, but this needs no dependency.)

## Wiring it in `main`

Now `cmd/cli/main.go` composes the real thing — this is the **composition root**, the one
place concrete types are chosen:

```go
store := storage.NewLocalStorage(cfg.OriginalsDir, cfg.GeneratedDir)
pipe := pipeline.New(store, cfg.MaxConcurrent)
svc := app.NewService(store, pipe, limitsFromConfig(cfg), defaultProfiles(), nil)

f, err := os.Open(*input)
// ...
m, err := svc.Process(ctx, app.ProcessInput{Filename: filepath.Base(*input), Reader: f, Size: info.Size()})
```

`main` builds dependencies and hands them down; the service does the work. The API's `main`
(story 11) will do the *same* wiring and call the *same* `svc.Process` — that's the point of
the story.

## Emitting JSON

Add struct tags so the manifest serializes with clean field names, then encode:

```go
type Variant struct {
    Name      string `json:"name"`
    Width     int    `json:"width"`
    // ...
}
json.NewEncoder(os.Stdout).Encode(manifest)
```

`encoding/json` uses the tags for key names and reflects over exported fields. The same
encoding will be reused by the HTTP handler — define it once on the domain (or a small
response DTO) and both transports get identical output.

## Testing the service

Use a `t.TempDir()`-backed `LocalStorage` (or a fake), a deterministic `newID`, and a real
small image from `testdata/`. Assert: no error, status `completed`, `len(Variants) == 3`,
and that the files exist on disk. Because you injected `newID`, you know the exact asset
directory to check.

## Try this

Implement `Process` end to end, then write the service test. Only once
`go test ./internal/app/` is green, rewire `cmd/cli/main.go` and run the real command
against `testdata/images/cafe.jpg`. Open `data/generated/<assetID>/` and look at your
thumbnail, square, and banner — the first time you see them is the best moment of the week.
