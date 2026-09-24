# Story 01 — Initialize the Go Project

**Time budget:** ~1.5 hours (Day 1)

## Goal

As an intern, I want a clean project structure so that I can implement the
application one package at a time. Nothing here processes an image yet. This story
is carpentry: understand the skeleton, make it compile, make the tests run.

> The scaffold already exists in `/src`. This story is about **reading and
> understanding** it, then confirming it builds. You'll write real code starting in
> story 02.

## What's already in `/src`

```
src/
├── go.mod                     module path (edit if you want a different one)
├── Makefile                   run-cli, run-api, test, test-race, fmt, vet, web-*, clean
├── .env.example               config settings you'll read in story 03
├── README.md                  project readme stub (you finish it in story 18)
├── cmd/cli/main.go            thin CLI entry point (calls run())
├── cmd/api/main.go            thin API entry point (calls run())
├── internal/
│   ├── app/app.go             shared application service (story 10)
│   ├── config/config.go       config loading (story 03)
│   ├── domain/domain.go       core types (story 02)
│   ├── domain/errors.go       ErrNotImplemented lives here
│   ├── media/media.go         decode/resize/crop/encode (stories 06, 07)
│   ├── pipeline/pipeline.go   variant generation (stories 08, 09)
│   └── storage/storage.go     storage interface + local impl (story 05)
├── testdata/images/           put sample .jpg/.png here for tests
├── data/originals/            saved originals (gitignored contents)
├── data/generated/            generated variants (gitignored contents)
└── docs/architecture.md, decisions.md   you write these in story 18
```

## Tasks

- [ ] Read `background.md` to understand modules, packages, `internal/`, exported vs
      unexported names, and the two output streams.
- [ ] From `/src`, run `go build ./...` and confirm it exits 0.
- [ ] Run `go run ./cmd/cli` and `go run ./cmd/api`; confirm each prints its
      "not implemented yet" line.
- [ ] Run `go test ./...` and confirm it passes (no test files yet is fine).
- [ ] Open `cmd/cli/main.go`, change the placeholder message, rebuild, rerun — feel
      the edit → build → run loop.
- [ ] Run `go mod tidy` (harmless now; habit-forming).
- [ ] Add one sample image (any small JPEG) into `testdata/images/` so later stories
      have something to work with. Name it `cafe.jpg` to match the examples.

## Acceptance criteria

- [ ] `go build ./...` exits 0.
- [ ] `go run ./cmd/cli` and `go run ./cmd/api` both run and print a placeholder line.
- [ ] `go test ./...` succeeds.
- [ ] Both `main.go` files contain **no** business logic — they only wire things up.
- [ ] No image-processing logic exists anywhere yet.

## Definition of done

`go test ./...` passes and the CLI prints something. Commit: `story 01: project scaffold`.

## Don't do yet

- No flag parsing (that's story 04).
- No real types beyond the placeholder packages (that's story 02).
- No image code at all.
