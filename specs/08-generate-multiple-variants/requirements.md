# Story 08 — Generate Multiple Variants

**Time budget:** ~1.5 hours (Day 6) — **this is the first big milestone**

## Goal

As a CLI user, I want one source image to produce several useful variants. You'll build
the sequential pipeline that loops over profiles, transforms + stores each, and
aggregates results — including partial failures.

## Files to work in (`/src`)

- `internal/pipeline/pipeline.go` — `Generate` (starter provided, do it sequentially now)
- `internal/pipeline/pipeline_test.go` — you create this
- `internal/domain/domain.go` — implement `ProcessingResult.Status`
- somewhere convenient (e.g. `internal/pipeline/profiles.go` or a config) — define the
  default profile set

## Required profiles

The full design includes `thumbnail, square, mobile, desktop, banner, preview`. The
first completed version only needs **three**:

```
thumbnail  200 x 200   crop   jpeg
square     1080 x 1080 crop   jpeg
banner     1600 x 600  crop   jpeg
```

Define these as a `[]domain.MediaProfile` you can pass into `Generate`.

## Tasks

- [ ] Implement `ProcessingResult.Status(requested int)`:
      all succeeded → `completed`; none → `failed`; some → `partial`.
- [ ] Implement `Pipeline.Generate` sequentially: for each profile, decode once
      (the caller passes an already-decoded `image.Image`), `media.Transform` into a
      `bytes.Buffer`, `store.SaveVariant`, and append a `Variant` or a `ProcessingError`.
- [ ] A failed profile records an error and continues — it must not stop the rest.
- [ ] Return the aggregated `ProcessingResult`.
- [ ] Write a test with an in-memory or temp-dir `Storage` and a fake source image,
      asserting three variants come back and their names/dimensions are right.

## Acceptance criteria

- [ ] One call generates at least three variants.
- [ ] All variants for one asset are grouped under one asset directory.
- [ ] Every requested profile appears in the result (as a success *or* an error).
- [ ] One failing profile does not hide the successful ones.
- [ ] `Status` distinguishes completed / partial / failed.
- [ ] `go test ./internal/pipeline/` passes.

## Definition of done

Wire this far enough that (with the CLI placeholder from story 04 pointed at it) one
command conceptually yields thumbnail + square + banner. Commit:
`story 08: sequential multi-variant pipeline with partial-failure handling`.

## Hints (not solutions)

- Decode the source ONCE before the loop; reuse the decoded `image.Image` for every
  profile (decoding per profile wastes work).
- `var buf bytes.Buffer; media.Transform(&buf, src, profile)` then
  `store.SaveVariant(ctx, assetID, profile.Name, profile.Format, &buf)`.
- To force a failure in a test, pass a profile with an impossible setting or a stub
  storage that errors on a chosen name.
- Keep `Generate` returning a top-level `error` only for *catastrophic* failures (bad
  input); per-profile problems belong in `result.Errors`.
