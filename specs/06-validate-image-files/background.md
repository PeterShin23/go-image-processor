# Background — Story 06: Validate Image Files

This story teaches Go's `image` package, the "sniff the bytes" security mindset, typed
sentinel errors, and the classic gotcha that **a reader can only be read once**.

## The `image` package and decoder registration

Go's standard library decodes images through `image` plus format sub-packages. The
crucial, surprising part: the sub-packages register themselves via **side-effect
imports**.

```go
import (
    "image"
    _ "image/jpeg" // registers the JPEG decoder
    _ "image/png"  // registers the PNG decoder
)
```

The `_` (blank import) means "import this package only for its `init()` side effects,
don't reference its names." Each format package's `init()` calls `image.RegisterFormat`.
If you forget these imports, `image.Decode` will fail on every file with "unknown format"
— a very common beginner trap. Import exactly the formats you support; importing only
JPEG and PNG is also how you *reject* GIF/BMP for free.

Two decode entry points:

```go
cfg, format, err := image.DecodeConfig(r) // reads only the HEADER: cheap
// cfg.Width, cfg.Height, format == "jpeg" | "png"

img, format, err := image.Decode(r)       // decodes ALL pixels: proves it's not corrupt
```

Use `DecodeConfig` to get dimensions and the format name without paying to decode every
pixel. Use `Decode` (or a successful `DecodeConfig` plus a full decode) to confirm the
file isn't truncated/corrupt. For this story, a successful `image.Decode` is the
strongest "it's a real image" signal.

## Don't trust the extension

A file named `cat.jpg` might be a text file, a GIF, or half a JPEG someone truncated.
Security rule: **detect format from content**. Two tools:

- `image.DecodeConfig` returns the real format string it recognized.
- `http.DetectContentType(first512Bytes)` returns a MIME type like `"image/jpeg"` by
  inspecting magic bytes. (Yes, it lives in `net/http`, but it doesn't require a
  request — it's a pure byte-sniffer. Using it here does not make your package "depend on
  HTTP" in the architectural sense; still, `image.DecodeConfig` alone is enough and keeps
  the import list clean.)

If the detected format isn't JPEG or PNG, return `ErrUnsupportedFormat`.

## The "a reader is consumed once" gotcha

An `io.Reader` is a stream. Once you read bytes out of it, they're **gone** — the next
read continues from where you stopped. So if you `image.DecodeConfig(r)` to sniff, then
try to `image.Decode(r)` again, the second call sees an empty/partial stream and fails.

Three ways to handle it:

1. **Read it all into memory first** (fine for images bounded by `MaxBytes`):

```go
data, err := io.ReadAll(io.LimitReader(r, limits.MaxBytes+1))
if err != nil { ... }
if int64(len(data)) > limits.MaxBytes { return domain.SourceImage{}, ErrTooLarge }
if len(data) == 0 { return domain.SourceImage{}, ErrEmptyFile }
// now decode from a fresh reader over the bytes, as many times as you like:
cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
```

`io.LimitReader(r, N)` caps how many bytes can be read — a cheap guard against someone
streaming a huge "image" to exhaust memory. Reading `MaxBytes+1` lets you detect
"too large" precisely.

2. **`io.TeeReader`**: reads from one reader while writing a copy into a buffer, so you can
   replay. More advanced; option 1 is simpler and fine here.

3. **Re-open the file**: if you have the path, just `os.Open` twice. Works for the CLI,
   not for an HTTP upload stream — so prefer option 1 for code the API will share.

Because the app service (story 10) feeds the *same* `Validate` from both a file (CLI) and
an upload (HTTP), option 1 is the portable choice.

## Typed sentinel errors and `errors.Is`

The starter defines package-level error values:

```go
var ErrEmptyFile = errors.New("image file is empty")
```

These are **sentinels**: named, comparable error values. A caller checks them with
`errors.Is`, which walks the wrap chain:

```go
_, err := media.Validate(r, size, limits)
if errors.Is(err, media.ErrUnsupportedFormat) {
    // map to a 415 in the HTTP layer (story 11)
}
```

You can wrap a sentinel with context and `errors.Is` still finds it:

```go
return domain.SourceImage{}, fmt.Errorf("validate cafe.jpg: %w", ErrCorruptImage)
```

The `%w` verb is what makes the wrapped error still "match" the sentinel. This is the
payoff of the error-wrapping habit from story 03: the HTTP layer can react to *categories*
of failure without string-matching messages.

> `errors.Is` compares identity (is this specific sentinel in the chain?).
> `errors.As` extracts a *typed* error into a variable (when you need its fields).
> You'll mostly use `errors.Is` here.

## Preparing test data

Put real files under `testdata/images/`:
- a small valid JPEG (`cafe.jpg`)
- a small valid PNG (`logo.png`)
- a deliberately corrupt file (e.g. `echo "not an image" > testdata/images/broken.jpg`)
- for the empty case, you can create a 0-byte file or pass `size: 0` with an empty reader
- for oversized, set a tiny `MaxBytes` in the test rather than needing a huge file

`testdata` is a special directory name: the Go toolchain ignores it when building, so
sample files there never get compiled or shipped.

## A test sketch

```go
func TestValidate(t *testing.T) {
    limits := media.Limits{MaxBytes: 5 << 20, MinWidth: 1, MinHeight: 1,
        MaxWidth: 10000, MaxHeight: 10000}
    tests := []struct {
        name    string
        path    string
        wantErr error // nil, or a sentinel to match with errors.Is
    }{
        {"valid jpeg", "testdata/images/cafe.jpg", nil},
        {"corrupt", "testdata/images/broken.jpg", media.ErrCorruptImage},
        // ...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, _ := os.Open(tt.path)
            defer f.Close()
            info, _ := f.Stat()
            _, err := media.Validate(f, info.Size(), limits)
            if tt.wantErr == nil && err != nil { t.Fatalf("unexpected: %v", err) }
            if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
                t.Fatalf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

## Try this

Wire the imports (`image`, `_ "image/jpeg"`, `_ "image/png"`), implement the empty +
oversized checks first (pure size logic, no decoding), test them, then add the
`DecodeConfig` format + dimensions checks. Run `go doc image.DecodeConfig` and
`go doc io.LimitReader` as you go.
