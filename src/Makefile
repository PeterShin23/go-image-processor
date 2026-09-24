# Makefile for the Go Media Processing Pipeline.
# Run `make <target>`. Targets marked web-* operate on the ./web frontend.

# Default input/output used by run-cli; override on the command line, e.g.
#   make run-cli INPUT=./testdata/images/cafe.jpg
INPUT ?= ./testdata/images/cafe.jpg
OUTPUT ?= ./data/generated

.PHONY: run-cli run-api test test-race fmt vet web-install web-run web-build clean

run-cli:
	go run ./cmd/cli --input $(INPUT) --output $(OUTPUT)

run-api:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

web-install:
	cd web && npm install

web-run:
	cd web && npm run dev

web-build:
	cd web && npm run build

clean:
	rm -rf data/originals/* data/generated/*
	go clean
