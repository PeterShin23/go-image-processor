# Go Media Processing Pipeline

Accepts an image and generates multiple optimized variants (thumbnail, square,
banner, ...) for different product surfaces. Built in three stages: a CLI, an
HTTP API that reuses the CLI's application service, and a React frontend.

> This README is a starting stub. You will complete it in **story 18**.

## Prerequisites

- Go 1.23+
- Node 18+ (only for the `web/` frontend)

## Layout

```
cmd/cli    CLI entry point
cmd/api    HTTP server entry point
internal/  application, config, domain, media, pipeline, storage packages
web/       React + TypeScript frontend (added in story 14)
testdata/  sample images for tests
data/      local storage for originals and generated variants
docs/      architecture and decision records
```

## Quick start

```bash
# CLI (placeholder until later stories)
go run ./cmd/cli --input ./testdata/images/cafe.jpg --output ./data/generated

# API (placeholder until later stories)
go run ./cmd/api

# Tests
go test ./...
```

## Make targets

`make run-cli run-api test test-race fmt vet web-install web-run web-build clean`

## Status

Scaffold only. See `specs/` for the story-by-story plan.
