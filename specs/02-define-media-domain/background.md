# Background — Story 02: Define the Media Domain

This story is where you learn how Go models data. You'll build the nouns of the
system: a profile, a variant, a manifest. Get comfortable here — you'll touch these
types in every later story.

## Structs: grouping related data

A **struct** is a named bundle of fields. It is Go's equivalent of a plain object or
a record. No inheritance, no methods inside the braces — just data:

```go
type Variant struct {
    Name      string
    Width     int
    Height    int
    SizeBytes int64
}
```

You create one with a **composite literal**. Prefer naming the fields:

```go
v := Variant{Name: "thumbnail", Width: 200, Height: 200}
// unset fields get their zero value: Height would be 0 if omitted
```

### Zero values

Go has no `null` for basic types. Every type has a **zero value**: `0` for numbers,
`""` for strings, `false` for bools, `nil` for slices/maps/pointers. A freshly
declared `var v Variant` is fully usable with all fields at their zero values. This
matters: you rarely need constructors for plain data structs.

## Methods: functions attached to a type

A **method** is a function with a **receiver** — the value it's called on:

```go
func (p MediaProfile) Validate() error {
    // p is the receiver; call as profile.Validate()
}
```

`(p MediaProfile)` is a **value receiver**: the method gets a copy. Use value receivers
for small structs you only read. Use a **pointer receiver** `(p *MediaProfile)` when the
method must modify the struct or when the struct is large. For `Validate`, which only
reads, a value receiver is idiomatic.

## Named types over primitives

Notice `type OutputFormat string`. This is a **defined type** whose underlying type is
`string`. Why not just use `string`? Because the compiler now stops you from passing a
random string where a format is expected, and constants document the valid set:

```go
const (
    FormatJPEG OutputFormat = "jpeg"
    FormatPNG  OutputFormat = "png"
)
```

This pattern — a string (or int) type plus a block of typed constants — is Go's
lightweight enum. There is no built-in `enum` keyword. You are responsible for
validating that a value is one of the known constants (that's part of `Validate`).

### `iota` (optional)

For integer enums Go offers `iota`, which auto-increments in a `const` block:

```go
type Stage int
const (
    StageDecode Stage = iota // 0
    StageResize              // 1
    StageEncode              // 2
)
```

We used string constants instead because they serialize to readable JSON later. Know
`iota` exists; you don't need it here.

## Errors are values

Go has no exceptions. A function that can fail returns an `error` as its **last**
return value. `nil` means success:

```go
func (p MediaProfile) Validate() error {
    if p.Width <= 0 {
        return fmt.Errorf("profile %q: width must be > 0, got %d", p.Name, p.Width)
    }
    // ... more checks ...
    return nil // all good
}
```

`fmt.Errorf` builds an error with a formatted message. Include the offending value —
future-you debugging a failure will thank present-you. `%q` quotes a string, `%d`
formats an integer, `%v` formats anything.

Callers check errors immediately:

```go
if err := p.Validate(); err != nil {
    // handle it
}
```

You'll go deeper on error *wrapping* (`%w`) in story 03. For now, plain `fmt.Errorf`
messages are enough.

## Slices: Go's growable lists

A `[]MediaProfile` is a **slice** — a view over a backing array that can grow.

```go
profiles := []MediaProfile{p1, p2, p3}
for i, p := range profiles {
    _ = i // index
    _ = p // a COPY of the element
}
for _, p := range profiles { ... } // common: ignore the index with _
```

`range` gives you `(index, valueCopy)`. Because `p` is a copy, mutating it doesn't
change the slice. Append with `profiles = append(profiles, p4)`. A `nil` slice is a
valid empty slice — you can `range` over it and `append` to it.

## Maps: for the duplicate check

To reject duplicate profile names, walk the slice once and remember names you've seen:

```go
seen := make(map[string]bool)
for _, p := range profiles {
    if seen[p.Name] {
        return fmt.Errorf("duplicate profile name %q", p.Name)
    }
    seen[p.Name] = true
}
```

Reading a missing key returns the value type's zero value (`false` here), so you don't
need to check existence separately. `map[string]struct{}` is a common
memory-tiny variant when you only care about presence, but `map[string]bool` is clearer
while learning.

## Table-driven tests

This is *the* idiomatic Go test style. Define a slice of cases, loop, run each as a
subtest:

```go
func TestMediaProfileValidate(t *testing.T) {
    tests := []struct {
        name    string
        profile MediaProfile
        wantErr bool
    }{
        {name: "valid", profile: MediaProfile{Name: "t", Width: 200, Height: 200,
            ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: false},
        {name: "zero width", profile: MediaProfile{Name: "t", Width: 0, Height: 200,
            ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: true},
        // ... add the rest ...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.profile.Validate()
            if (err != nil) != tt.wantErr {
                t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

Key points:
- Test files end in `_test.go` and live in the same package.
- Test functions are `func TestXxx(t *testing.T)`.
- `t.Run(name, func)` creates a named subtest so failures point at the exact case.
- `(err != nil) != tt.wantErr` is the classic "did we get an error when we expected
  one?" check.
- Run with `go test ./internal/domain/` or `go test -run Validate ./...`.

Notice the test asserts *behavior* (error or not), not the exact message. That keeps
tests from breaking every time you reword an error.

## Why the domain package is "pure"

The acceptance criteria forbid importing HTTP or filesystem code here. This is
**dependency direction**: the domain is the center; everything depends on it, it
depends on nothing of ours. That keeps your core types reusable by the CLI, the API,
and the tests without dragging in a web server. If you ever `import "net/http"` inside
`domain`, stop — the type you're writing belongs in a different package.

## Try this

Run `go doc fmt.Errorf` and `go doc testing.T.Run`. Then implement `Validate` for the
width and height checks only, write the two matching test cases, and run
`go test ./internal/domain/`. Add the remaining checks one at a time, each with a test.
Small loops beat big-bang implementations.
