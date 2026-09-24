# Background — Story 04: Build the Single-Image CLI Shell

Now you make the program actually *take input*. You'll meet the `flag` package, file
system checks, and the discipline of separating "read the command line" from "do the
work."

## The `flag` package

Go ships a command-line flag parser in the standard library. The pattern:

```go
import "flag"

func run() error {
    input := flag.String("input", "", "path to one source image")
    output := flag.String("output", "./data/generated", "output directory")
    flag.Parse()

    // IMPORTANT: these are POINTERS. Read the value with *input after Parse.
    if *input == "" {
        return fmt.Errorf("--input is required")
    }
    ...
}
```

- `flag.String(name, default, usage)` returns a `*string`. There are siblings:
  `flag.Int`, `flag.Bool`, `flag.Duration`.
- `flag.Parse()` reads `os.Args` and fills those pointers. Call it once, before you read
  any flag.
- Users pass flags as `--input x` or `--input=x` (single dash also works).
- `flag` auto-generates `-h`/`--help` from your usage strings — free documentation.

### Why pointers?

`flag.String` has to hand you something *before* parsing has happened, so it returns the
address of a variable it will fill in during `Parse()`. That's why you dereference with
`*input` afterward. (There's also `flag.StringVar(&myVar, ...)` if you'd rather bind to
your own variable — either is fine.)

## Checking the file exists

Two different jobs, two different calls:

```go
info, err := os.Stat(path) // ask about the file WITHOUT opening it
if err != nil {
    return fmt.Errorf("input %q: %w", path, err) // covers "does not exist" and more
}
if info.IsDir() {
    return fmt.Errorf("input %q is a directory, not a file", path)
}
```

- `os.Stat` returns file metadata (`os.FileInfo`): size, mode, is-it-a-directory.
- If the file is missing, `err` is non-nil and `errors.Is(err, os.ErrNotExist)` is true —
  handy if you want a friendlier message for that specific case.

Then actually open it to get a readable stream:

```go
f, err := os.Open(path)
if err != nil {
    return fmt.Errorf("open %q: %w", path, err)
}
defer f.Close()
```

## `defer`: cleanup that can't be forgotten

`defer f.Close()` schedules `f.Close()` to run when the surrounding function returns —
no matter how it returns (normal, early, or via an error path). This is Go's answer to
"always release the resource." You'll use `defer` for files, HTTP bodies, and mutexes
throughout the project. Put the `defer` immediately after the successful acquire, so the
cleanup is impossible to skip.

`*os.File` implements `io.Reader`, so once opened you can hand `f` to anything that reads
bytes — exactly what the application service will want in story 10.

## stdout vs stderr vs exit code (recap from story 01)

- Results → **stdout** (`fmt.Println`, or `fmt.Fprintln(os.Stdout, ...)`).
- Errors/diagnostics → **stderr** (`fmt.Fprintln(os.Stderr, ...)`).
- Failure → non-zero **exit code**. The starter `main` already does this: `run()`
  returns an `error`, and `main` prints it to stderr and calls `os.Exit(1)`.

So inside `run()` you never call `os.Exit` yourself — you just `return err`. That keeps
`run()` testable: a test can call `run`-like logic and inspect the returned error instead
of the process dying.

## The placeholder service

The CLI's real job is to hand the opened file to the shared application service. But that
service doesn't exist until story 10. So for now, "call the service" is a stand-in:
print the filename and a fake "processing..." line. This is intentional scaffolding —
you're proving the *plumbing* (flags → file → some callee → printed result) works before
the callee is real. When story 10 lands, you swap the placeholder for
`svc.Process(ctx, ...)`.

## Keep parsing separate from logic

Notice the shape: `main` handles process concerns (exit code), `run` handles the CLI
(flags, file checks, calling the service). Neither contains image code. This separation
is why the same application service can later be driven by an HTTP handler that does its
*own* input parsing (multipart instead of flags) but calls the identical service. Mixing
flag parsing into business logic would make that reuse impossible.

## Try this

Run `go doc flag` and skim it. Implement the `--input` required-check first and test the
three cases by hand (no flag, bad path, good path), checking `echo $?` each time. Then add
the `os.Stat` existence check. Keep `run()` returning errors rather than exiting.
