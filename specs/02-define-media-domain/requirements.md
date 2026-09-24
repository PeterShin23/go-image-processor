# Story 02 — Define the Media Domain

**Time budget:** ~1.5–3 hours (Days 1–2)

## Goal

As an intern, I want domain models for assets, profiles, variants, and manifests so
that the whole application shares one consistent vocabulary. This is the most
important early story: every later package imports `internal/domain`.

## Files to work in (`/src`)

- `internal/domain/domain.go` — the types and validation (starter shapes provided)
- `internal/domain/errors.go` — `ErrNotImplemented` already here; add typed errors if useful
- `internal/domain/domain_test.go` — table-driven test skeletons (remove the `t.Skip`s)

## Tasks

- [ ] Review the starter types: `OutputFormat`, `ResizeMode`, `ProcessingStatus`,
      `MediaProfile`, `Variant`, `ProcessingError`, `Manifest`. Adjust fields to what
      you actually need.
- [ ] Add a `SourceImage` concept (e.g. decoded dimensions + format) — you'll use it
      in stories 06/07. It can be a small struct now.
- [ ] Implement `MediaProfile.Validate() error`:
  - reject width ≤ 0 and height ≤ 0
  - reject quality outside 1..100
  - reject unknown `ResizeMode` and unknown `OutputFormat`
- [ ] Implement `ValidateProfiles([]MediaProfile) error`:
  - reject duplicate profile names
  - reject if any individual profile is invalid
- [ ] Optionally add small helpers like `OutputFormat.Valid()` / `ResizeMode.Valid()`.
- [ ] Fill in the table-driven tests and delete the `t.Skip` lines.

## Acceptance criteria

- [ ] Domain types do not import HTTP, filesystem, or any transport package.
- [ ] Domain types do not read or write files.
- [ ] Invalid width/height values are rejected by `Validate`.
- [ ] Invalid quality values are rejected.
- [ ] Duplicate profile names are rejected by `ValidateProfiles`.
- [ ] Tests cover both valid and invalid profiles and pass: `go test ./internal/domain/`.

## Definition of done

`go test ./internal/domain/` passes with real (non-skipped) tests, and `go vet ./...`
is clean. Commit: `story 02: media domain types + validation`.

## Hints (not solutions)

- Return errors with context, e.g. `fmt.Errorf("profile %q: width must be > 0", p.Name)`.
- A `map[string]bool` (or `map[string]struct{}`) is the idiomatic way to detect
  duplicate names in one pass.
- Keep JSON tags off for now unless you want them; the manifest is printed in later
  stories where you'll decide on the wire format.
