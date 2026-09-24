# Specs — Go Media Processing Pipeline

This folder breaks the project brief (`../guideline.md`) into 18 stories you can
implement one at a time. It is written for an intern learning Go **while** building
this project, working roughly **1.5 hours per day**.

## Where the code lives: `/src`

The actual project is one Go module in the root **`/src`** folder. It ships as a
**scaffold**: every package already exists with a doc comment, and as you progress the
stories fill in type stubs, function signatures returning `ErrNotImplemented`, and
`TODO`s. You are never staring at a blank file — but the real logic is yours to write.
That is the whole point: you learn Go by replacing the TODOs.

Run all Go commands from inside `/src`:

```bash
cd src
go build ./...
go test ./...
```

## How each story folder works

Every `NN-<story>/` folder contains two files:

| File | What it is |
| ---- | ---------- |
| `requirements.md` | The "ticket". What to build, which `/src` files to touch, and the acceptance criteria you must hit. Read this first. |
| `background.md` | The Go (or TypeScript) concepts you need, taught from scratch, plus hints and gotchas. Read this second. |

## How to use a story

1. Read `requirements.md` — note which files under `/src` it names.
2. Read `background.md` — do the small "try this" exercises if any.
3. Open those `/src` files and replace the `TODO`s until the acceptance criteria pass.
4. Run the checks from `/src`: `go test ./...`, then `go vet ./...`, then `go fmt ./...`.
5. Commit. Move to the next story tomorrow.

## The mentor contract

Ask a mentor (or an AI assistant) to **review** your code and give **hints**, not to
write the story for you. You will remember what you struggled with; you will forget
what was handed to you.

## Re-paced schedule (~10 days at ~1.5 hrs/day)

The original brief assumes 7 full days. At 1.5 hours/day the same work spreads across
about 10 sessions. Suggested grouping (each line ≈ one 1.5-hour session):

| Day | Stories | Milestone |
| --- | ------- | --------- |
| 1 | 01 Initialize project · 02 Define domain (start) | Project compiles, `go test ./...` passes |
| 2 | 02 Define domain (finish) · 03 Load config | Domain + config with table-driven tests |
| 3 | 04 CLI shell · 05 Storage abstraction (start) | CLI accepts `--input` and prints a placeholder |
| 4 | 05 Storage (finish) · 06 Validate images | Valid JPEG/PNG accepted, bad files rejected |
| 5 | 07 One transformation | One variant generated at exact dimensions |
| 6 | 08 Generate multiple variants | **Milestone: CLI makes thumbnail + square + banner** |
| 7 | 09 Bounded concurrency · 10 Shared app service | `go test -race ./...` passes; CLI uses the service |
| 8 | 11 Upload API · 12 Health + shutdown | `curl` upload returns a manifest; server shuts down cleanly |
| 9 | 13 Structured logging · 14 React init · 15 Upload queue | Frontend selects files and uploads sequentially |
| 10 | 16 Frontend concurrency · 17 Results + retry · 18 Docs | Full demo with multiple files; README written |

If a session runs long, stop at a compiling checkpoint and carry the rest into the
next day. Ending on green (`go test ./...` passing) every day matters more than
finishing a story.

## Dependency order

Stories build on each other. Do them in number order. The hard dependencies are:

- 02 (domain) is used by almost everything after it.
- 05 (storage) + 06 (validation) + 07 (transform) feed 08 (variants).
- 08 → 09 (concurrency) → 10 (app service).
- 10 is reused by both 04 (CLI) and 11 (API) — that shared service is the
  architectural heart of the project.
- 14–17 are the React frontend and only need a running API (stories 11–13).

## The frontend shortcut

You plan to generate ~80% of the React frontend with AI. That is fine — stories 14–17
still tell you the required behavior and state model so you can prompt for it precisely
and review what you get. Focus your learning hours on the Go backend.
