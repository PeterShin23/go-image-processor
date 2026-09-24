# Story 04 — Build the Single-Image CLI Shell

**Time budget:** ~1.5 hours (Day 3)

## Goal

As a user, I want to provide one image path on the command line so I can run the
application without starting a server. This story wires up flags and file checking —
the real processing is still a placeholder until story 10.

## Files to work in (`/src`)

- `cmd/cli/main.go` — implement `run()` (starter has `main` + empty `run`)
- `cmd/cli/main_test.go` — optional; you may test a small helper you extract

## Flags

```
--input     required   path to one source image
--output    optional   output directory (default ./data/generated)
--profiles  optional   ignored for now; wired up in story 08
```

Only `--input` needs to be required initially.

## Tasks

- [ ] In `run()`, define flags with the `flag` package.
- [ ] Parse them.
- [ ] If `--input` is empty, return a clear error (non-zero exit).
- [ ] Confirm the file exists and is a regular file (`os.Stat`); error clearly if not.
- [ ] Open the file with `os.Open` and `defer` closing it.
- [ ] Call a **placeholder** application service (for now, just print a fake manifest
      line or the filename; the real `app.Service` arrives in story 10).
- [ ] Print results to **stdout**, errors to **stderr**, and use non-zero exit codes on
      failure (the starter's `main` already maps a returned error to exit code 1).

## Acceptance criteria

- [ ] Missing `--input` produces a clear error and non-zero exit.
- [ ] A non-existent path produces a clear error and non-zero exit.
- [ ] The CLI processes exactly one image per invocation.
- [ ] Flag/path parsing is separate from any media logic (there is none yet).
- [ ] A successful run prints a placeholder result until later stories fill it in.

## Definition of done

```bash
go run ./cmd/cli                       # errors: --input required, exit 1
go run ./cmd/cli --input nope.jpg      # errors: file not found, exit 1
go run ./cmd/cli --input testdata/images/cafe.jpg   # prints placeholder, exit 0
```

Commit: `story 04: CLI flag parsing and input validation`.

## Hints (not solutions)

- Check exit codes with `echo $?` right after a run.
- `flag.String("input", "", "...")` returns a `*string`; read it *after* `flag.Parse()`.
- Keep `main` tiny; put logic in `run()` so it's testable later.
