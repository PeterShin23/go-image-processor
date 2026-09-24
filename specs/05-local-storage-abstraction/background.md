# Background — Story 05: Create the Local Storage Abstraction

This story is your first real encounter with **interfaces** — the single most important
idea in Go's design — plus the `io` streaming model and safe path handling. Take your
time; the interface mindset here pays off in stories 10 and beyond.

## Interfaces: behavior, not data

A Go **interface** is a set of method signatures. Any type that has those methods
*satisfies* the interface automatically — there is **no `implements` keyword**. This is
called structural (or "duck") typing:

```go
type Storage interface {
    SaveOriginal(ctx context.Context, assetID, filename string, r io.Reader) (string, int64, error)
    // ...
}
```

`*LocalStorage` satisfies `Storage` simply by having methods with matching signatures.
Nothing declares the relationship. That's why we add this line:

```go
var _ Storage = (*LocalStorage)(nil)
```

It's a compile-time assertion: "a `*LocalStorage` must be usable as a `Storage`." If you
typo a method name or signature, the build fails *here* with a clear message instead of
somewhere confusing later. The `_` means "I don't need the variable, just check the type."

### Why an interface at all?

Because the app service (story 10) will depend on `Storage`, not on `LocalStorage`. That
means you can later write an `S3Storage` with the same methods and swap it in without
changing the service. It also means tests can pass a fake `Storage`. This is
**dependency injection**: high-level code depends on an abstraction, and the concrete
implementation is handed in from outside (via `NewService(store, ...)`).

Go idiom: **define interfaces where they're consumed, keep them small.** Our `Storage` has
three methods because the app needs three behaviors. Don't add methods "just in case" —
the acceptance criteria specifically warn against exposing unnecessary filesystem detail.

## The `io.Reader` / `io.Writer` model

Go models streaming bytes with two tiny interfaces:

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

Almost everything that produces or consumes bytes implements one of these: files
(`*os.File`), network connections, HTTP request bodies, in-memory buffers
(`bytes.Buffer`), even the multipart upload you'll handle in story 11. Because your
`SaveOriginal` takes an `io.Reader`, it doesn't care *where* the bytes come from — a file
on disk (CLI) or an HTTP upload (API) both work. That's the whole reason the CLI and API
can share this code.

The workhorse for moving bytes from a reader to a writer:

```go
n, err := io.Copy(dstFile, r) // returns bytes copied
```

You don't loop over `Read`/`Write` yourself; `io.Copy` does it efficiently. `n` is your
`size`.

### Readers vs ReadClosers

`io.Reader` just reads. `io.ReadCloser` also has `Close()`. Files are `ReadCloser`s — the
caller must close them. That's why `Open` returns `io.ReadCloser`: the caller `defer`s the
`Close`. Don't return a raw `*os.File`; returning the interface keeps your abstraction
honest (a future S3 reader isn't a file).

## Writing a file, step by step

```go
dir := filepath.Join(s.originalsDir, assetID)
if err := os.MkdirAll(dir, 0o755); err != nil {
    return "", 0, fmt.Errorf("create dir %q: %w", dir, err)
}
dst := filepath.Join(dir, safeName)
f, err := os.Create(dst)
if err != nil {
    return "", 0, fmt.Errorf("create file %q: %w", dst, err)
}
defer f.Close()
n, err := io.Copy(f, r)
if err != nil {
    return "", 0, fmt.Errorf("write %q: %w", dst, err)
}
return dst, n, nil
```

- `os.MkdirAll` creates every missing parent and is a no-op if the dir exists.
- `0o755` is an octal Unix permission (owner rwx, others r-x). The `0o` prefix is Go's
  octal literal syntax.
- `os.Create` truncates/creates the file for writing.

## Path traversal: the security bit

If a caller controls part of a filename, they might send `../../../../etc/passwd` or an
absolute path, trying to write outside your storage root. **Never** trust that input.

Two common defenses (use one or both):

1. **Strip directory components**: `filepath.Base(filename)` reduces
   `../../evil.txt` to `evil.txt`. Good when you only ever want a leaf name.
2. **Clean and verify containment**: build the full path, `filepath.Clean` it to resolve
   `..`, then confirm it still starts with your root:

```go
full := filepath.Clean(filepath.Join(root, assetID, name))
if !strings.HasPrefix(full, filepath.Clean(root)+string(os.PathSeparator)) {
    return "", 0, fmt.Errorf("unsafe path %q", name)
}
```

Your acceptance criteria require a test that a malicious filename cannot escape. Write it:
call `SaveOriginal` with `filename: "../../escape.txt"` and assert either an error or that
nothing was written outside the temp root.

## Testing with temporary directories

```go
func TestSaveOriginalRoundTrips(t *testing.T) {
    dir := t.TempDir() // auto-removed when the test ends
    s := NewLocalStorage(filepath.Join(dir, "orig"), filepath.Join(dir, "gen"))

    path, n, err := s.SaveOriginal(context.Background(), "asset1", "cafe.jpg",
        strings.NewReader("hello"))
    if err != nil { t.Fatal(err) }
    if n != 5 { t.Fatalf("size = %d, want 5", n) }

    rc, err := s.Open(context.Background(), path)
    if err != nil { t.Fatal(err) }
    defer rc.Close()
    got, _ := io.ReadAll(rc)
    if string(got) != "hello" { t.Fatalf("got %q", got) }
}
```

- `t.TempDir()` — isolated, auto-cleaned. Never write into the project's real `data/`.
- `strings.NewReader("hello")` is an easy `io.Reader` for tests — no real image needed to
  test storage mechanics.
- `context.Background()` is the empty root context; story 09 explains why the parameter
  exists. For now, pass it through.

## Try this

Implement `SaveOriginal` and `Open`, then write the round-trip test above and run
`go test ./internal/storage/`. Only after that works, add the path-traversal guard and its
test. Build in small, tested increments.
