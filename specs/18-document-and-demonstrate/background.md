# Background — Story 18: Document and Demonstrate the System

The code is done; now make it legible to the next person (often future-you). Good docs are
part of engineering, not an afterthought. This story also doubles as your presentation prep.

## What good project docs contain

Aim your README at a specific reader: **an engineer who just cloned the repo and wants it
running in five minutes.** That framing tells you what to include and in what order:

1. **One-paragraph purpose** — what it does and why it exists.
2. **Prerequisites** — Go version, Node version.
3. **Run each stage** — exact copy-pasteable commands for CLI, API, and web.
4. **Configuration** — the env vars (point at `.env.example`), with defaults.
5. **Example request + response** — copied from a *real* run, not invented.
6. **Test commands** — `go test ./...`, `-race`, `npm run build`.
7. **Known limitations / future work** — honesty builds trust.

Copy real output. Run the CLI and paste the actual JSON manifest; run the `curl` and paste
what came back. Invented examples drift from reality and mislead the next person.

## Architecture Decision Records (ADRs)

`docs/decisions.md` captures *why*, which code alone can't. The value isn't the decision you
made — it's the *reasoning*, so a future maintainer knows what to reconsider when constraints
change. A lightweight format per decision:

```
## Local disk storage (first implementation)

Context:   We need to persist originals and variants. Options: local disk, S3, a DB blob.
Decision:  Use local filesystem behind a storage.Storage interface.
Consequences: Simplest to run and test (t.TempDir); not durable or multi-node. The
              interface means swapping in S3 later touches one package, not the app service.
```

Write at least three. The four suggested ones aren't arbitrary — each maps to a question
you'll be asked in the final presentation:
- **Local disk storage** → "how could this later use S3?" (the interface is your answer).
- **Bounded variant concurrency** → "how is Go concurrency limited, and why?"
- **One request per image** → "why no batch endpoint?" (browser owns the queue; backend
  stays simple and stateless per request).
- **Shared application service** → "why do the CLI and API share one service?" (no duplicated
  business logic; both transports are thin adapters).

## The architecture document

`docs/architecture.md` should walk the three stages and the shared flow. You don't need
fancy diagrams — reuse the ASCII flows from `guideline.md`:

```
Process one source image
  ├── validate        (media.Validate)
  ├── create asset ID  (app.Service.newID)
  ├── store original   (storage.Storage.SaveOriginal)
  ├── generate variants(pipeline.Generate, bounded)
  ├── store variants   (storage.Storage.SaveVariant)
  └── return manifest  (domain.Manifest)
```

Then show how the CLI and the HTTP handler both funnel into `app.Service.Process`. The single
most important sentence in the whole document: *"the CLI and HTTP API are thin adapters over
one shared application service."* If a reader takes away only that, they understand the design.

## Final verification, honestly reported

Run every check and report results truthfully — including anything still failing:

```bash
cd src
go fmt ./...        # formatting (should produce no diff if you've been running it)
go vet ./...        # static checks
go test ./...       # unit + integration
go test -race ./... # the concurrency guarantee from story 09
cd web && npm run build   # the frontend type-checks and builds
```

If something fails, say so in the README's status/limitations section rather than hiding it.
"Known limitation: the `preview` profile occasionally off-by-one on odd dimensions" is far
more useful than silence. A reviewer trusts a doc that admits gaps.

## Map to the definition of done

The guideline lists 25 done-criteria. Skim them as a final checklist: single-image CLI, three
variants, JSON manifest, stored original, bounded concurrency, shared service, HTTP upload,
health endpoints, graceful shutdown, structured logs, multi-file React with per-file state,
one request per file, limited frontend concurrency, retry, passing tests, passing race
detector, building frontend, and a README covering all three workflows. Anything unchecked is
your remaining work.

## Preparing to present

The guideline's "Final intern presentation" lists 15 questions. Most are answered by your
decision records and by tracing one request through the code. Practice narrating: pick a real
uploaded file and talk through validate → store → generate (bounded) → manifest, then explain
the two independent concurrency controls (Go variants vs. React uploads). If you can do that
from memory, you understand the system.

## Try this

Do a real end-to-end demo first: start the API, start the web app, drop in 3–4 images
(including one bad file), and watch the results. Screenshot or copy the outputs. *Then* write
the README and ADRs from what you actually observed. Finish by running the full check suite and
recording the results.
