# Story 07 — Implement One Image Transformation

**Time budget:** ~1.5 hours (Day 5)

## Goal

As an intern, I want to transform one decoded source image into exactly one variant:
resize to fit, optionally center-crop to exact dimensions, then encode as JPEG or PNG.

## Files to work in (`/src`)

- `internal/media/transform.go` — `Transform` (starter provided)
- `internal/media/transform_test.go` — you create this
- `go.mod` — you'll likely add one dependency: `golang.org/x/image` (see background)

## Required behavior

- **Resize to fit**: scale so the image fits the target box *without distortion*
  (preserve aspect ratio).
- **Center crop**: for `ResizeCrop`, scale to *cover* the box, then crop the center to the
  exact `Width × Height`.
- **Exact output dimensions** for crop mode; fit mode may leave one side smaller.
- **JPEG and PNG encoding**, honoring `Quality` for JPEG.
- **Optional no-enlarge**: if `AllowEnlarge` is false and the source is smaller than the
  target, don't upscale.

## Tasks

- [ ] Implement `Transform(w, src, p)`: dispatch on `p.ResizeMode`, produce the resized
      (and possibly cropped) `image.Image`, encode to `w` per `p.Format`/`p.Quality`, and
      return the final dimensions.
- [ ] Break it into small helpers: `resizeFit`, `centerCrop`, `encode`.
- [ ] Add `golang.org/x/image/draw` for high-quality scaling (`go get golang.org/x/image`).
- [ ] Write tests that decode the output and assert its dimensions match expectations.

## Acceptance criteria

- [ ] Output dimensions match the profile (exact for crop mode).
- [ ] Resize-to-fit preserves aspect ratio (no distortion).
- [ ] Center crop produces exactly `Width × Height`.
- [ ] `Quality` is passed to the JPEG encoder.
- [ ] The encoded output can be decoded back successfully.
- [ ] Tests verify generated dimensions and pass: `go test ./internal/media/`.
- [ ] `Transform` writes to an `io.Writer` and never touches HTTP or files directly.

## Definition of done

Given a test source image, `Transform` yields a decodable JPEG at the target size.
Commit: `story 07: single-image resize/crop/encode`.

## Hints (not solutions)

- Compute the scale factor from the *ratio* of source to target; the smaller ratio fits,
  the larger ratio covers (for crop).
- `draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)` does a
  high-quality resize into a destination `*image.RGBA` you allocate.
- Center-crop with `src.(interface{ SubImage(image.Rectangle) image.Image })` or by
  drawing the centered region into a new image.
- `jpeg.Encode(w, img, &jpeg.Options{Quality: p.Quality})`, `png.Encode(w, img)`.
