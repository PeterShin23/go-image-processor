# Background — Story 01: Initialize the Go Project

You are about to meet Go for the first time. This story has almost no logic in it
on purpose: the goal is to get comfortable with how a Go project is *shaped* before
you write anything interesting.

## What a Go module is

A **module** is a project. It is defined by a single file, `go.mod`, at the root.
That file declares:

- the module's **import path** (a globally unique name, usually a URL-like string
  such as `github.com/PeterShin23/go-image-processor`), and
- the Go version, and
- the list of third-party dependencies.

When you write `import "github.com/PeterShin23/go-image-processor/internal/domain"`,
Go uses the module path in `go.mod` to know that this import lives *inside your own
project*, in the folder `internal/domain`.

You create a module once:

```bash
go mod init github.com/PeterShin23/go-image-processor
```

`/src` already contains a `go.mod`, so you don't have to run this — but you may
edit the module path on the first line if you prefer a different name. If you change
it, you must update every `import` line that references it.

`go mod tidy` scans your code and rewrites `go.mod`/`go.sum` to match the imports you
actually use. Run it whenever imports change.

## What a package is

A **package** is a folder of `.go` files that all start with the same
`package <name>` line. The package name is usually the folder name.

- `package main` is special: it produces an **executable** program and must contain a
  `func main()`. Both `cmd/cli` and `cmd/api` are `package main`.
- Every other folder is a **library** package that other code imports. `internal/domain`,
  `internal/config`, etc. are libraries.

### The `internal/` rule

Go gives one folder name magic powers: **`internal`**. Packages under `internal/` can
only be imported by code rooted at the parent of `internal/`. In our case that means
only *our own* project can import `internal/domain`; nobody who depends on our module
from outside could. This is exactly what we want — it keeps our guts private. You get
this for free just by putting packages under `internal/`.

## Exported vs unexported names

Go has no `public`/`private` keywords. Instead, **capitalization is visibility**:

- `Asset` (capital A) is **exported** — visible to other packages.
- `asset` (lowercase a) is **unexported** — visible only inside its own package.

This applies to types, functions, struct fields, constants, everything. Keep this in
mind constantly; it is the single most common surprise for people new to Go.

## The entry point

An executable's `main` is deliberately tiny in a well-structured Go program. Look at
`src/cmd/cli/main.go`: its whole job is to call into other packages and set
the process exit code. It must not contain image logic. A common idiom you'll see:

```go
func main() {
    if err := run(); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}
```

We keep `main` this thin so the real work lives in testable library functions, not in
`main` (which is awkward to test). You'll flesh out `run()` in later stories.

## Exit codes and the two output streams

Command-line programs communicate in three channels:

- **stdout** (`fmt.Println`, `os.Stdout`): the program's actual output — here, the JSON
  manifest. Machines read this.
- **stderr** (`fmt.Fprintln(os.Stderr, ...)`): diagnostics, errors, logs. Humans read this.
- **exit code**: `0` means success, non-zero means failure. `os.Exit(1)` sets it.

Keeping errors on stderr and results on stdout means someone can do
`go run ./cmd/cli ... > result.json` and get clean JSON without your error messages
mixed in. You'll rely on this in story 04.

## The typed placeholder error

Idiomatic Go represents failures as **values** of the `error` type, returned as the
last return value of a function. A **sentinel error** is a named error value you can
compare against:

```go
var ErrNotImplemented = errors.New("not implemented")
```

We put this in `internal/domain/errors.go`. Every stub function in the scaffold returns
it, so the code compiles and honestly reports "this isn't done yet" instead of pretending
to work. As you implement each story, you delete the `return ErrNotImplemented` lines.

## The tools you'll run constantly

Memorize these four. They are your feedback loop:

```bash
go build ./...   # does everything compile?  (./... means "this package and all below")
go test ./...    # do all tests pass?
go vet ./...     # does Go spot any suspicious code?
go fmt ./...     # auto-format every file to the one true Go style
```

Go has exactly one canonical formatting. Never argue about it; just run `go fmt`.

Also learn `go doc`:

```bash
go doc fmt.Println          # read the docs for a stdlib function
go doc os.Exit
```

## Try this

Before you touch anything, from `/src` run:

```bash
go build ./...
go test ./...
go run ./cmd/cli
go run ./cmd/api
```

Watch what each does. Then open `cmd/cli/main.go` and change the placeholder message,
rebuild, and rerun. That loop — edit, build, run — is the whole job.

## What "done" feels like

You haven't processed an image. You haven't parsed a flag. You have a tree of empty,
honest packages that compiles and tests green. That is exactly right for day one.
Resist the urge to jump ahead into `internal/media`; each later story will tell you when.
