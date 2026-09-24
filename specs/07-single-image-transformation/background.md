# Background — Story 07: Implement One Image Transformation

This is the algorithmic heart of the project. You'll learn Go's image model, how to
preserve aspect ratio, how center cropping works, and how to encode with quality
settings. It's also your first justified third-party dependency.

## Go's image model in 90 seconds

- `image.Image` is an interface: it can report its `Bounds()` (a `Rectangle`) and the
  color `At(x, y)`.
- `image.Rectangle` is defined by two points, `Min` and `Max`. Width is `Max.X - Min.X`.
  Bounds don't have to start at (0,0), but for images you create they usually do.
- Concrete images you'll allocate: `image.NewRGBA(image.Rect(0, 0, w, h))` — a writable
  image with 8-bit RGBA pixels.
- Decoding (`image.Decode`) gives you *some* concrete type (often `*image.YCbCr` for JPEG).
  You generally resize *into* a fresh `*image.RGBA` you control.

```go
b := src.Bounds()
srcW, srcH := b.Dx(), b.Dy() // Dx/Dy are width/height helpers
```

## Why you need a third-party scaler

The standard library can *draw* images (`image/draw`) but its built-in scaling is minimal.
High-quality resizing lives in the semi-official extended library
`golang.org/x/image/draw`. The `x/image` module is maintained by the Go team — it's the
right kind of dependency (meaningful image functionality, not a convenience wrapper), which
satisfies the project's "add third-party packages only when they earn their place" rule.

```bash
cd src
go get golang.org/x/image
go mod tidy
```

Then:

```go
import xdraw "golang.org/x/image/draw"

dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)
```

`CatmullRom` is a high-quality interpolator (smooth downscales). `Scale` resamples `src`
into `dst`. `draw.Over` is the compositing operator ("draw src over dst"). Alternatives:
`xdraw.ApproxBiLinear` (faster, softer), `xdraw.NearestNeighbor` (blocky, fast).

## Resize-to-fit: preserve aspect ratio

"Fit" means the whole image must fit inside the `W×H` box without stretching. Compute the
scale factor as the **smaller** of the two axis ratios, so neither side overflows:

```go
scale := math.Min(float64(targetW)/float64(srcW), float64(targetH)/float64(srcH))
```

Then the fitted size is `round(srcW*scale) × round(srcH*scale)`. One dimension will equal
the target; the other will be ≤ the target. Do **not** compute width and height with
independent scale factors — that stretches the image (distortion), which the acceptance
criteria forbid.

### No-enlarge

If `!p.AllowEnlarge` and `scale > 1` (the source is smaller than the box), clamp
`scale = 1`. That prevents blurry upscaling of tiny sources.

## Center crop: cover then trim

"Crop" means the output must be *exactly* `W×H`, filling the whole box even if that means
losing edges. Two steps:

1. **Cover**: scale using the **larger** ratio, so the image fully covers the box (one
   dimension will overflow):

```go
scale := math.Max(float64(targetW)/float64(srcW), float64(targetH)/float64(srcH))
```

2. **Center-crop** the overflow. After scaling to `coverW×coverH` (≥ target on both axes),
   take the centered `targetW×targetH` window:

```go
offX := (coverW - targetW) / 2
offY := (coverH - targetH) / 2
crop := image.Rect(offX, offY, offX+targetW, offY+targetH)
```

You can either scale into a cover-sized image then copy the centered rectangle into a
final `targetW×targetH` image with `xdraw.Draw(final, final.Bounds(), scaled, crop.Min, draw.Src)`,
or scale directly into the target while offsetting the source rectangle. Either works;
the copy-then-crop version is easier to reason about the first time.

## Encoding with quality

```go
import (
    "image/jpeg"
    "image/png"
)

// JPEG: quality 1..100, higher = better/larger
err := jpeg.Encode(w, img, &jpeg.Options{Quality: p.Quality})

// PNG: lossless; no quality knob, but you can pick a compression level
err := png.Encode(w, img)
```

`Encode` writes to any `io.Writer` — a file, a `bytes.Buffer`, or an HTTP response. That's
why `Transform` takes a `w io.Writer`: the pipeline (story 08) will pass a buffer, then
hand those bytes to storage. `Transform` itself never opens a file — separation of
concerns again.

Note these encoder imports are *direct* here (not blank `_`) because you call
`jpeg.Encode`/`png.Encode` by name, unlike the decoder registration in story 06.

## Testing dimensions

The cleanest test encodes to a buffer, decodes it back, and checks the size:

```go
func TestTransformCropExactSize(t *testing.T) {
    src := image.NewRGBA(image.Rect(0, 0, 800, 600)) // a fake source is fine
    var buf bytes.Buffer
    p := domain.MediaProfile{Name: "sq", Width: 200, Height: 200,
        ResizeMode: domain.ResizeCrop, Format: domain.FormatJPEG, Quality: 80}

    w, h, err := media.Transform(&buf, src, p)
    if err != nil { t.Fatal(err) }
    if w != 200 || h != 200 { t.Fatalf("got %dx%d, want 200x200", w, h) }

    cfg, _, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
    if err != nil { t.Fatal(err) }
    if cfg.Width != 200 || cfg.Height != 200 {
        t.Fatalf("decoded %dx%d", cfg.Width, cfg.Height)
    }
}
```

You don't need a real photo to test the math — a blank `image.NewRGBA` has real bounds.
For aspect-ratio (fit) tests, use a non-square source (e.g. 800×600 into a 300×300 box →
expect 300×225).

## Rounding gotcha

Float scaling plus integer pixels means off-by-one is easy. Use `int(math.Round(...))` and
guard against zero: a dimension must be at least 1 pixel. Decide crop offsets with integer
division (`/2`) and accept that an odd leftover pixel goes to one side — that's normal.

## Try this

Add `golang.org/x/image`, implement `resizeFit` first (fit mode), and test it against a
non-square source to confirm aspect ratio is preserved. Then implement `centerCrop` for
crop mode and test the exact-size guarantee. Encode last. Build up one helper at a time.
