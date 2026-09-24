# Story 18 — Document and Demonstrate the System

**Time budget:** ~1 hour (Day 10)

## Goal

As an intern, I want to explain the system so another engineer can run, understand, and
maintain it — and so you can present it confidently.

## Files to work in (`/src`)

- `README.md` — flesh out the stub
- `docs/architecture.md` — the three stages and the request flow
- `docs/decisions.md` — at least three decision records

## README must document

project purpose · architecture overview · prerequisites · CLI setup · CLI usage · API
setup · API usage · React setup · configuration (env vars) · test commands · an example
request · an example response · known limitations · future improvements.

## `docs/decisions.md` — record at least three of

- Why local disk storage was selected first (vs. object storage).
- Why variant concurrency is bounded.
- Why the frontend sends one request per image (no batch endpoint).
- Why the CLI and HTTP API share one application service.

Use a short format per decision: **Context → Decision → Consequences**.

## Final checks (run all, paste results into the README's status section)

```bash
cd src
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
cd web && npm run build
```

## Acceptance criteria

- [ ] A new developer can run all three stages from the README alone.
- [ ] `docs/architecture.md` explains the CLI, API, and React stages and the shared flow.
- [ ] At least three decisions are recorded with tradeoffs.
- [ ] The full system can be demonstrated end-to-end with multiple images.
- [ ] All the final checks above pass.

## Definition of done

The whole checklist in the guideline's "Definition of done" (items 1–25) is satisfied and
demonstrable. Commit: `story 18: documentation and final polish`.

## Hints (not solutions)

- Write the example request/response from a real run (`curl ... | jq`), don't invent it.
- Keep the architecture doc diagram-light; the ASCII flows from `guideline.md` are enough.
- For your presentation, rehearse the 15 questions in the guideline's "Final intern
  presentation" section — the decision records answer several of them.
