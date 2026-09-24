# Story 06 — Validate Image Files

**Time budget:** ~1.5 hours (Day 4)

## Goal

As an intern, I want to validate an image before processing so corrupted and
unsupported files are rejected early with clear, typed errors.

## Files to work in (`/src`)

- `internal/media/validate.go` — `Validate`, `Limits`, typed errors (starter provided)
- `internal/media/validate_test.go` — you create this
- `testdata/images/` — add a valid JPEG, a valid PNG, a tiny corrupt file, and note an
  empty file case

## Required checks (in a sensible order)

1. File is not empty (`size == 0` → `ErrEmptyFile`).
2. File size is within `limits.MaxBytes` (→ `ErrTooLarge`).
3. Format is JPEG or PNG, detected from **content**, not the filename
   (→ `ErrUnsupportedFormat`).
4. The image actually decodes (→ `ErrCorruptImage`).
5. Width and height are within `limits` (→ `ErrDimensionsRange`).

On success, return a `domain.SourceImage{Width, Height, Format}`.

## Tasks

- [ ] Implement `Validate`. Detect format without decoding the whole image using
      `image.DecodeConfig` (returns dimensions + format name cheaply) and/or
      `http.DetectContentType` on the first 512 bytes.
- [ ] Ensure you can still read the image bytes afterward (a reader is consumed once —
      see background for the buffering trick).
- [ ] Return the **typed** errors so callers can `errors.Is` them.
- [ ] Register the image decoders you need (`import _ "image/jpeg"` and
      `import _ "image/png"`).
- [ ] Write table-driven tests using files under `testdata/images/`.

## Acceptance criteria

- [ ] Valid JPEG accepted; valid PNG accepted.
- [ ] Unsupported format (e.g. a `.txt` or GIF) rejected with `ErrUnsupportedFormat`.
- [ ] Corrupt file rejected with `ErrCorruptImage`.
- [ ] Empty file rejected with `ErrEmptyFile`.
- [ ] Oversized file rejected with `ErrTooLarge`.
- [ ] Validation does not import HTTP (it takes an `io.Reader`, not a request).
- [ ] Tests read from `testdata/images/` and pass: `go test ./internal/media/`.

## Definition of done

All five rejection paths and both accept paths are tested and green. Commit:
`story 06: content-based image validation with typed errors`.

## Hints (not solutions)

- Do not trust the extension: a file named `photo.jpg` can contain anything.
- `image.DecodeConfig(r)` reads only the header — cheap way to get width/height/format.
- To both sniff AND still decode, read into a buffer first (`bytes.Buffer` +
  `io.TeeReader`, or just `io.ReadAll` for small files) — see background.
- Test the typed errors with `errors.Is(err, media.ErrEmptyFile)`.
