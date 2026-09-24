# Background — Story 08: Generate Multiple Variants

This story is about **aggregation** and **partial failure** — patterns you'll reuse
constantly in real services. The Go concepts are modest (slices, loops, structs); the
design thinking is the point. Do it sequentially first; story 09 makes it concurrent.

## The shape of the pipeline

`Generate` receives an already-decoded `image.Image` and a list of profiles, and returns
a `ProcessingResult`. The core loop:

```go
func (p *Pipeline) Generate(ctx context.Context, assetID string, src image.Image,
    profiles []domain.MediaProfile) (domain.ProcessingResult, error) {

    var result domain.ProcessingResult
    for _, profile := range profiles {
        var buf bytes.Buffer
        w, h, err := media.Transform(&buf, src, profile)
        if err != nil {
            result.Errors = append(result.Errors, domain.ProcessingError{
                Profile: profile.Name, Stage: "transform", Message: err.Error(),
            })
            continue // keep going — one failure must not stop the rest
        }
        path, size, err := p.store.SaveVariant(ctx, assetID, profile.Name, profile.Format, &buf)
        if err != nil {
            result.Errors = append(result.Errors, domain.ProcessingError{
                Profile: profile.Name, Stage: "store", Message: err.Error(),
            })
            continue
        }
        result.Variants = append(result.Variants, domain.Variant{
            Name: profile.Name, Width: w, Height: h, Format: profile.Format,
            Path: path, SizeBytes: size,
        })
    }
    return result, nil
}
```

Study the structure:
- **Decode once** (the caller already did) — reuse `src` for every profile.
- Each profile gets a fresh `bytes.Buffer`. `Transform` writes encoded bytes into it;
  `SaveVariant` reads them out. A buffer is both an `io.Writer` and an `io.Reader`.
- On error, record a `ProcessingError` and `continue`. Do **not** `return`.
- `Stage` labels *where* it broke ("transform" vs "store"), which is gold when debugging.

## Partial failure is a feature, not a bug

A naive implementation returns on the first error. That's wrong here: if `banner` fails
but `thumbnail` and `square` succeed, the user should still get their two good variants
plus a clear note about the one that failed. This "collect errors, keep going" pattern
appears everywhere: batch jobs, bulk APIs, import tools. The `ProcessingResult` holding
*both* `Variants` and `Errors` is what makes it expressible.

The top-level `error` return is reserved for catastrophes that make the whole call
meaningless (e.g. a nil source). Per-profile issues live in `result.Errors`.

## Deriving overall status

`ProcessingResult.Status(requested int)` collapses the detail into one word:

```go
func (r ProcessingResult) Status(requested int) ProcessingStatus {
    switch {
    case len(r.Variants) == requested:
        return StatusCompleted
    case len(r.Variants) == 0:
        return StatusFailed
    default:
        return StatusPartial
    }
}
```

`requested` is `len(profiles)`. Passing it in (rather than storing it on the result) keeps
`ProcessingResult` a plain data bag. This three-way status feeds both the CLI exit code and
the HTTP response later.

## `bytes.Buffer`: an in-memory stream

`bytes.Buffer` is the swiss-army pipe between producers and consumers of bytes:

```go
var buf bytes.Buffer
media.Transform(&buf, src, profile) // writes into buf  (buf as io.Writer)
p.store.SaveVariant(..., &buf)      // reads out of buf  (buf as io.Reader)
```

Pass `&buf` (a pointer) so both sides share the same buffer. After a full read the buffer
is drained — fine here since each profile uses its own. This is why `Transform` took an
`io.Writer` and `SaveVariant` took an `io.Reader`: `bytes.Buffer` bridges them with no file
in between.

## Appending to slices

`result.Variants = append(result.Variants, v)` — remember `append` returns a (possibly
new) slice, so you must assign it back. Starting from a `nil` slice is fine; `append` to
`nil` allocates. You never pre-size these; let `append` grow them.

## Where do the profiles come from?

For now, hardcode the three-profile set as a function returning `[]domain.MediaProfile`
(e.g. in `internal/pipeline/profiles.go` or passed from the CLI). Keep it validated:
call `domain.ValidateProfiles(profiles)` before running so a typo (duplicate name, bad
dimensions) is caught up front rather than mid-loop. Making profiles configurable via JSON
is a stretch goal — don't do it yet.

## Testing with a fake or temp storage

You have two easy options:
1. **Real `LocalStorage` in a `t.TempDir()`** — exercises the actual disk path.
2. **A tiny fake `Storage`** implemented in the test — lets you *force* a `SaveVariant`
   error for one profile name to prove partial failure works:

```go
type fakeStore struct{ failOn string }
func (f fakeStore) SaveVariant(_ context.Context, _, name string, _ domain.OutputFormat, r io.Reader) (string, int64, error) {
    if name == f.failOn { return "", 0, errors.New("boom") }
    n, _ := io.Copy(io.Discard, r) // drain it
    return "/fake/" + name, n, nil
}
// ...implement the other Storage methods as no-ops to satisfy the interface
```

The fake demonstrates the power of the story-05 interface: because the pipeline depends on
`storage.Storage`, not `*LocalStorage`, you can substitute a controllable test double.
`io.Discard` is a writer that throws bytes away — handy when you only care about the count.

## Try this

Implement `Status` and test its three cases (pure logic, no images). Then implement the
sequential `Generate` and test the happy path with a temp-dir store. Finally add the
fake-store partial-failure test. Three small commits beat one big one.
