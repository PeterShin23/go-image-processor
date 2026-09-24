# Go Media Processing Pipeline

## Internship Project Brief

### Project duration

1 week

### Project level

Beginner-to-intermediate Go project

### Project goal

Build a media-processing application that accepts image files and generates multiple optimized versions for different product use cases.

The project should simulate the kind of image pipeline used by:

* Social media platforms
* E-commerce applications
* Content management systems
* News websites
* Photography platforms
* Internal company asset systems

A single uploaded image may need to become several variants, such as:

* Profile thumbnail
* Mobile feed image
* Desktop feed image
* Square social media image
* Website banner
* Low-bandwidth preview
* High-quality archive version

The project will be built in three stages:

1. A Go CLI that processes one local image
2. A Go HTTP API that processes one uploaded image
3. A React interface that allows multiple files to be selected and processed

The application should teach practical Go development while remaining achievable within one week.

---

# Learning objectives

By completing this project, the developer should gain experience with:

1. Go modules and package organization
2. Structs and methods
3. Interfaces
4. Dependency injection
5. File input and output
6. Image decoding and encoding
7. Command-line applications
8. HTTP servers
9. Multipart file uploads
10. JSON responses
11. Error handling and error wrapping
12. Context cancellation
13. Goroutines and bounded concurrency
14. Unit testing
15. Integration testing
16. Table-driven tests
17. Configuration using environment variables
18. Structured logging
19. Graceful server shutdown
20. Building a small React client for a Go backend

---

# Product overview

The application receives an image and processes it according to a collection of predefined media profiles.

Each media profile describes the output that should be generated.

A media profile may include:

* Variant name
* Width
* Height
* Resize strategy
* Crop strategy
* Output format
* Compression quality
* Whether enlargement is allowed
* Whether original metadata should be preserved

Example media profiles:

| Variant   |    Dimensions | Format          | Purpose                   |
| --------- | ------------: | --------------- | ------------------------- |
| thumbnail |     200 × 200 | JPEG            | Avatar or search result   |
| square    |   1080 × 1080 | JPEG            | Social media image        |
| mobile    |     720 × 900 | JPEG            | Mobile content feed       |
| desktop   |    1200 × 800 | JPEG            | Desktop content feed      |
| banner    |    1600 × 600 | JPEG            | Website hero section      |
| preview   |     480 × 320 | JPEG            | Low-bandwidth preview     |
| archive   | Original size | Original format | High-quality preservation |

The exact output formats may depend on the image-processing library selected.

---

# Project stages

## Stage 1: Single-image CLI

The first version of the project should process one image from the local filesystem.

Example command:

```bash
go run ./cmd/cli \
  --input ./testdata/images/cafe.jpg \
  --output ./data/generated
```

The CLI should:

1. Read the input file.
2. Validate the image.
3. Decode the image.
4. Generate the configured variants.
5. Save the generated files.
6. Produce a JSON manifest.
7. Print the manifest to standard output.
8. Return an appropriate exit code.

The CLI should process exactly one source image per invocation.

The CLI must be completed before work begins on the HTTP API.

---

## Stage 2: Single-image HTTP API

After the CLI works, the same media-processing pipeline should be exposed through an HTTP server.

Endpoint:

```text
POST /v1/assets
```

The endpoint should accept one image using `multipart/form-data`.

The HTTP server should:

1. Receive the uploaded image.
2. Enforce upload-size limits.
3. Validate the image.
4. Save the original.
5. Call the same application service used by the CLI.
6. Generate the required variants.
7. Return a JSON manifest.

The HTTP layer should not contain image-transformation logic.

The CLI and HTTP API should share the same application service and pipeline.

---

## Stage 3: React multi-file interface

The React application should allow the user to:

1. Select multiple image files.
2. View the selected files.
3. Remove files before processing.
4. Start processing the list.
5. See the status of every file.
6. View successful results.
7. View errors.
8. Retry an individual failed file.

The React interface should send one HTTP request per image.

For example, selecting five files should result in five requests to:

```text
POST /v1/assets
```

The first version should not use a batch upload endpoint.

The browser owns the file queue, while the Go backend continues processing one image per request.

---

# System workflow

## CLI workflow

```text
CLI command
    |
    v
Read local file
    |
    v
Validate image
    |
    v
Store original
    |
    v
Generate variants
    |
    v
Store generated files
    |
    v
Return JSON manifest
```

---

## HTTP workflow

```text
HTTP upload
    |
    v
Parse multipart request
    |
    v
Validate image
    |
    v
Call application service
    |
    v
Generate variants
    |
    v
Return JSON manifest
```

---

## React workflow

```text
React file queue
    |
    ├── image-1.jpg ──> POST /v1/assets
    ├── image-2.png ──> POST /v1/assets
    ├── image-3.jpg ──> POST /v1/assets
    └── image-4.png ──> POST /v1/assets
                           |
                           v
                   Go media pipeline
```

Each file should have its own state:

```text
pending
uploading
processing
completed
failed
```

The first implementation should process files sequentially.

Once sequential processing works, the frontend may process two files concurrently.

---

# Example asset manifest

The completed pipeline should return an object conceptually similar to:

```json
{
  "asset_id": "generated-asset-id",
  "original_filename": "cafe.jpg",
  "status": "completed",
  "variants": [
    {
      "name": "thumbnail",
      "width": 200,
      "height": 200,
      "format": "jpeg",
      "path": "data/generated/generated-asset-id/thumbnail.jpg",
      "size_bytes": 18432
    },
    {
      "name": "banner",
      "width": 1600,
      "height": 600,
      "format": "jpeg",
      "path": "data/generated/generated-asset-id/banner.jpg",
      "size_bytes": 102400
    }
  ],
  "errors": []
}
```

The exact field names may be adjusted during implementation.

The coding agent should not implement the manifest-generation logic.

---

# Required capabilities

The completed project should support:

* JPEG input
* PNG input
* File validation
* Maximum file size
* Image dimension validation
* Local original-image storage
* Local generated-image storage
* Resize-to-fit behavior
* Center cropping
* Multiple media profiles
* JPEG output
* PNG output
* Configurable output quality where supported
* JSON manifests
* Single-image CLI processing
* Single-image HTTP upload
* Multi-file React upload queue
* Bounded Go concurrency
* Bounded frontend upload concurrency
* Structured logs
* Unit tests
* Integration tests
* Graceful server shutdown

---

# Optional stretch capabilities

These should only be attempted after the required functionality is complete:

* WebP output
* EXIF orientation correction
* EXIF metadata removal
* Blur placeholders
* Watermarking
* Content hashing
* Duplicate image detection
* Persistent job status
* Retry policies
* Object storage support
* S3-compatible storage
* Prometheus metrics
* Download links for generated variants
* Drag-and-drop upload
* Image previews
* Per-file progress percentages
* Configurable media profiles through JSON or YAML

---

# Non-goals

The first version should not include:

* Authentication
* User accounts
* A database
* Cloud deployment
* AWS infrastructure
* Kubernetes
* A distributed queue
* Kafka
* RabbitMQ
* Video processing
* Machine learning
* Content moderation
* Image generation
* Server-side rendering
* A complex frontend framework
* A frontend state-management library
* A batch image endpoint

The project should remain small enough to complete in one week.

---

# Recommended project structure

```text
go-media-pipeline/
├── cmd/
│   ├── cli/
│   │   └── main.go
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── domain/
│   ├── httpapi/
│   ├── media/
│   ├── pipeline/
│   └── storage/
├── web/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── types/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── testdata/
│   └── images/
├── data/
│   ├── originals/
│   └── generated/
├── docs/
│   ├── architecture.md
│   └── decisions.md
├── .env.example
├── .gitignore
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

---

# Package responsibilities

## `cmd/cli`

The command-line entry point.

Responsible for:

* Parsing command-line flags
* Loading configuration
* Constructing application dependencies
* Opening one source image
* Calling the shared application service
* Printing the result
* Returning the correct exit code

It must not contain image-processing business logic.

---

## `cmd/api`

The HTTP server entry point.

Responsible for:

* Loading configuration
* Constructing dependencies
* Starting the HTTP server
* Registering shutdown behavior
* Handling operating-system signals

It must not contain media-processing business logic.

---

## `internal/app`

The shared application layer.

Responsible for:

* Accepting one image-processing request
* Coordinating validation
* Saving the original image
* Calling the media pipeline
* Returning an asset manifest

The CLI and HTTP API should both call this package.

---

## `internal/config`

Configuration loading and validation.

Possible settings:

* HTTP port
* Maximum file size
* Minimum image width
* Minimum image height
* Maximum image width
* Maximum image height
* Original storage directory
* Generated storage directory
* Maximum concurrent transformations
* HTTP request timeout
* Shutdown timeout

---

## `internal/domain`

Core business types.

Possible concepts:

* Asset
* Source image
* Media profile
* Generated variant
* Asset manifest
* Processing result
* Processing error
* Processing status
* Resize mode
* Output format

This package should not depend on HTTP, React, or filesystem details.

---

## `internal/httpapi`

HTTP-specific behavior.

Responsible for:

* Routes
* Multipart form parsing
* Request-size limits
* JSON responses
* HTTP error mapping
* Health endpoints
* Request middleware

Handlers should remain thin.

---

## `internal/media`

Low-level image functionality.

Responsible for abstractions related to:

* Format detection
* Image decoding
* Resizing
* Cropping
* Encoding
* Output quality

This package should not know about HTTP requests or React.

---

## `internal/pipeline`

Coordinates variant generation.

Responsible for:

* Receiving a source image
* Processing all configured profiles
* Limiting concurrency
* Collecting successful results
* Collecting failures
* Respecting context cancellation
* Returning a processing result

---

## `internal/storage`

Storage abstraction and local filesystem implementation.

Responsible for:

* Saving original images
* Saving generated variants
* Reading stored assets when needed
* Creating asset directories
* Generating safe file paths
* Preventing path traversal

The first implementation should use local disk storage.

---

## `web`

The React and TypeScript frontend.

Responsible for:

* Selecting multiple files
* Managing the upload queue
* Sending one request per image
* Limiting upload concurrency
* Displaying per-file status
* Displaying generated variant metadata
* Retrying failed uploads

The React application should not perform image transformations.

---

# Architecture constraints

The project should follow these constraints:

1. Keep both `main.go` files small.
2. Do not place the application in one package.
3. Keep CLI logic separate from business logic.
4. Keep HTTP logic separate from media logic.
5. Keep storage behind an interface.
6. Pass `context.Context` through long-running operations.
7. Avoid global mutable state.
8. Use constructor functions for services with dependencies.
9. Wrap errors with useful context.
10. Do not use `panic` for expected failures.
11. Do not silently ignore failed transformations.
12. Do not launch unlimited goroutines.
13. Do not upload every selected browser file simultaneously.
14. Do not introduce a database.
15. Prefer the Go standard library where practical.
16. Add third-party packages only when they provide meaningful image-processing functionality.
17. Keep abstractions understandable to an intern.
18. Avoid creating interfaces before they provide a clear architectural benefit.
19. Ensure the CLI and HTTP API use the same application service.
20. Keep frontend state local and simple.

---

# User stories

## Story 1: Initialize the Go project

### Description

As an intern, I want a clean project structure so that I can implement the application one package at a time.

### Tasks

* Initialize the Go module.
* Create the recommended directories.
* Create `cmd/cli`.
* Create `cmd/api`.
* Add placeholder internal packages.
* Add `.gitignore`.
* Add `.env.example`.
* Add a Makefile.
* Add an initial README.
* Confirm that `go test ./...` succeeds.

### Acceptance criteria

* The project compiles.
* Both Go entry points exist.
* `go test ./...` succeeds.
* The `main.go` files contain no business logic.
* Image-processing logic has not been implemented.

### Go concepts

* Modules
* Packages
* Imports
* Entry points
* Exported and unexported names

---

## Story 2: Define the media domain

### Description

As an intern, I want domain models for assets, profiles, variants, and manifests so that the application has a consistent vocabulary.

### Tasks

Define the types needed to represent:

* Source image
* Asset
* Media profile
* Generated variant
* Asset manifest
* Processing result
* Processing status
* Processing error
* Resize strategy
* Output format

Add validation methods where appropriate.

### Acceptance criteria

* Domain types do not depend on HTTP.
* Domain types do not write files.
* Invalid width and height values are rejected.
* Invalid quality values are rejected.
* Duplicate profile names are rejected.
* Tests cover valid and invalid media profiles.

### Go concepts

* Structs
* Methods
* Custom types
* Constants
* Slices
* Validation
* Table-driven tests

---

## Story 3: Load application configuration

### Description

As an intern, I want configuration to come from environment variables so that application behavior can change without editing source code.

### Tasks

Support configuration for:

* Server port
* Maximum image size
* Original storage directory
* Generated storage directory
* Maximum concurrent transformations
* HTTP timeout
* Shutdown timeout

Provide defaults for local development.

### Acceptance criteria

* Optional settings have documented defaults.
* Invalid numeric values produce clear errors.
* Invalid durations produce clear errors.
* Configuration is passed into services.
* Configuration is not accessed through global variables.
* Tests isolate environment-variable changes.

### Go concepts

* Environment variables
* Integer parsing
* Duration parsing
* Error wrapping
* Constructor functions

---

## Story 4: Build the single-image CLI shell

### Description

As a user, I want to provide one image path through the command line so that I can run the application without starting a server.

### Example command

```bash
go run ./cmd/cli --input ./photo.jpg
```

### Initial flags

```text
--input
--output
--profiles
```

Only `--input` needs to be required initially.

### Tasks

* Parse command-line flags.
* Validate that an input path was provided.
* Validate that the file exists.
* Open the file.
* Pass it to a placeholder application service.
* Print output to standard output.
* Print errors to standard error.
* Use nonzero exit codes for failures.

### Acceptance criteria

* Missing input produces a clear error.
* An invalid path produces a clear error.
* The CLI processes exactly one image.
* CLI parsing is separate from media logic.
* Successful execution can print a placeholder result until later stories are complete.

### Go concepts

* `flag`
* `os.Args`
* Standard output
* Standard error
* Exit codes
* File paths

---

## Story 5: Create the local storage abstraction

### Description

As an intern, I want storage operations behind an interface so that local disk storage could later be replaced with object storage.

### Tasks

Design storage behavior for:

* Saving an original image
* Saving a generated variant
* Opening a stored asset
* Returning file metadata
* Creating asset directories

Implement local filesystem storage.

### Acceptance criteria

* Tests use temporary directories.
* Original and generated files are stored separately.
* Asset IDs isolate separate uploads.
* User-provided filenames cannot escape the storage root.
* Storage errors contain useful context.
* The interface does not expose unnecessary filesystem details.

### Go concepts

* Interfaces
* `io.Reader`
* `io.Writer`
* File operations
* Temporary directories
* Dependency injection
* Path handling

---

## Story 6: Validate image files

### Description

As an intern, I want to validate an image before processing so that corrupted and unsupported files are rejected.

### Required validation

* File exists
* File is not empty
* File size is within limits
* Format is JPEG or PNG
* Image can be decoded
* Width is within limits
* Height is within limits

Do not trust only the file extension.

### Acceptance criteria

* Valid JPEG files are accepted.
* Valid PNG files are accepted.
* Unsupported formats are rejected.
* Corrupted files are rejected.
* Empty files are rejected.
* Oversized files are rejected.
* Validation logic does not depend on HTTP.
* Tests use files under `testdata/images`.

### Go concepts

* Buffered readers
* MIME and format detection
* Image metadata
* Typed errors
* Error comparison

---

## Story 7: Implement one image transformation

### Description

As an intern, I want to transform one source image into one configured variant.

### Required behavior

* Resize to fit
* Center crop
* Exact output dimensions
* JPEG encoding
* PNG encoding
* Configurable quality where supported
* Optional prevention of enlargement

### Acceptance criteria

* Output dimensions match the media profile.
* Resize-to-fit does not distort the image.
* Center crop produces the exact target dimensions.
* Quality settings are passed to the encoder.
* Output can be decoded successfully.
* Tests verify generated image dimensions.
* Transformation logic does not write HTTP responses.

### Go concepts

* Image decoding
* Image encoding
* Buffers
* Interfaces
* `defer`
* Resource cleanup
* Context

---

## Story 8: Generate multiple variants

### Description

As a CLI user, I want one source image to produce several useful output variants.

### Initial profiles

The complete design may include:

```text
thumbnail
square
mobile
desktop
banner
preview
```

The first completed version only needs:

```text
thumbnail
square
banner
```

### Tasks

* Read the configured profiles.
* Process each profile.
* Save successful outputs.
* Record failures.
* Produce an asset manifest.

### Acceptance criteria

* One CLI command generates at least three variants.
* Generated files are grouped under one asset directory.
* Every configured profile appears in the final result.
* One failed profile does not hide successful profiles.
* The overall status distinguishes success, partial success, and failure.

### Go concepts

* Slices
* Loops
* Service structs
* Result aggregation
* Error handling

---

## Story 9: Add bounded variant concurrency

### Description

As an intern, I want image variants processed concurrently without exhausting system resources.

### Tasks

* Run independent transformations concurrently.
* Add a configurable concurrency limit.
* Collect results safely.
* Preserve profile-to-result associations.
* Stop unnecessary work when context is canceled.

### Acceptance criteria

* The concurrency limit is respected.
* Results do not contain data races.
* Every result belongs to the correct profile.
* Started goroutines exit normally.
* Context cancellation is respected.
* `go test -race ./...` succeeds.

### Go concepts

* Goroutines
* Channels
* Wait groups
* Worker pools
* Semaphores
* Mutexes
* Context cancellation
* Race detection

---

## Story 10: Create the shared application service

### Description

As an intern, I want one application service used by both the CLI and HTTP server so that business logic is not duplicated.

### Responsibilities

The application service should coordinate:

1. Image validation
2. Asset ID creation
3. Original file storage
4. Variant generation
5. Generated file storage
6. Manifest creation

### Expected conceptual input

* Image stream
* Original filename
* Media profiles
* Processing options

### Expected conceptual output

* Asset manifest
* Error

### Acceptance criteria

* The CLI calls the application service.
* The service contains no CLI flag parsing.
* The service contains no HTTP response-writing logic.
* The service receives its dependencies through construction.
* The service does not create concrete storage dependencies internally.

### Go concepts

* Dependency injection
* Interfaces
* Service structs
* Constructors
* Application boundaries

---

## Story 11: Create the upload API

### Description

As a client, I want to upload one image over HTTP so that another application can use the pipeline.

### Endpoint

```text
POST /v1/assets
```

### Request

Use `multipart/form-data` with one image field.

### Tasks

* Configure HTTP routes.
* Parse the multipart request.
* Enforce request-size limits.
* Read the uploaded filename.
* Call the shared application service.
* Return a JSON manifest.
* Map application errors to HTTP status codes.

### Acceptance criteria

* A valid image returns a success response.
* Invalid input returns a client-error response.
* Internal failures return a server-error response.
* Responses use JSON.
* Temporary multipart resources are cleaned up.
* Handlers remain thin.
* Tests use `httptest`.

### Go concepts

* `net/http`
* Handlers
* Multipart uploads
* JSON encoding
* Status codes
* `httptest`

---

## Story 12: Add health endpoints and graceful shutdown

### Description

As an operator, I want to know whether the service is running and have it stop cleanly.

### Endpoints

```text
GET /health/live
GET /health/ready
```

### Tasks

* Add liveness and readiness routes.
* Configure server timeouts.
* Listen for interrupt and termination signals.
* Stop accepting new requests during shutdown.
* Give active requests limited time to finish.

### Acceptance criteria

* Health endpoints return JSON.
* The HTTP server has read, write, idle, and header timeouts.
* Shutdown uses a context with a deadline.
* The process responds to operating-system signals.
* Shutdown failures are logged.

### Go concepts

* Signals
* Context deadlines
* Server lifecycle
* Graceful shutdown

---

## Story 13: Add structured logging

### Description

As an operator, I want useful processing logs so that I can understand what happened to each upload.

### Log fields

* Request method
* Request path
* Request duration
* Asset ID
* Original filename
* Number of requested variants
* Number of successful variants
* Number of failed variants
* Processing duration
* Error details

### Acceptance criteria

* Logs are structured.
* Every upload has an asset or request identifier.
* Errors identify the stage that failed.
* Image contents are never logged.
* Logging does not depend on global mutable state.

### Go concepts

* Structured logging
* Middleware
* Context values
* Durations
* Dependency injection

---

## Story 14: Initialize the React application

### Description

As a user, I want a simple browser interface where I can select multiple images.

### Technology requirements

Use:

* React
* TypeScript
* Vite
* Native `fetch`
* Basic CSS

Do not introduce:

* Redux
* Zustand
* React Router
* Server-side rendering
* A component framework
* Authentication
* A frontend database

### Tasks

* Initialize the React application.
* Add a multiple-file input.
* Store selected files in component state.
* Display file name and size.
* Allow individual files to be removed.
* Restrict selection to supported image formats.

### Acceptance criteria

* The frontend runs locally.
* Multiple files can be selected.
* Selected files appear in a list.
* Files can be removed before processing.
* The frontend does not transform images.

### TypeScript concepts

* Interfaces
* Component state
* File APIs
* Array updates
* Event handling

---

## Story 15: Build the frontend upload queue

### Description

As a user, I want selected files processed one at a time so that I can see the result of each upload.

### Initial queue behavior

For every file:

1. Mark it as pending.
2. Mark it as uploading.
3. Send it to the Go API.
4. Wait for the response.
5. Store the manifest.
6. Mark it as completed or failed.
7. Continue to the next file.

### File states

```text
pending
uploading
processing
completed
failed
```

The frontend may combine uploading and processing if the HTTP API does not expose separate server-side job status.

### Acceptance criteria

* Every selected file has independent state.
* Files are processed sequentially.
* A failed upload does not stop the remaining queue.
* Completed files store their manifest.
* Failed files display a useful error.
* The same queue cannot be started twice.
* The queue can be cleared after completion.

### TypeScript concepts

* Union types
* Async functions
* Error handling
* Immutable state updates
* API response types

---

## Story 16: Add limited frontend concurrency

### Description

As a user, I want a small number of files processed simultaneously without overwhelming the backend.

### Tasks

* Add a configurable upload concurrency limit.
* Default the limit to two.
* Start new work when an active request finishes.
* Preserve the display order of selected files.
* Keep results associated with the correct source files.

### Acceptance criteria

* No more than two uploads run concurrently by default.
* One failure does not cancel unrelated uploads.
* Every result remains attached to the correct file.
* Display order remains stable.
* The concurrency limit can be changed in one place.

### Important distinction

The project contains two separate concurrency controls:

1. React controls how many images are uploaded at once.
2. Go controls how many variants of one image are transformed at once.

These solve different problems and should remain independent.

---

## Story 17: Display results and retry failures

### Description

As a user, I want to inspect the generated outputs for each source image.

### Display for completed files

* Original filename
* Asset ID
* Overall status
* Generated variants
* Variant name
* Width
* Height
* Format
* File size
* Partial errors

### Tasks

* Render result data for each file.
* Display full success.
* Display partial success.
* Display total failure.
* Add retry behavior for individual failed files.

### Acceptance criteria

* Results are grouped by source file.
* Partial success is visually distinguishable.
* Failed files can be retried individually.
* Retrying does not resend successful files.
* The interface remains usable with at least ten files.

---

## Story 18: Document and demonstrate the system

### Description

As an intern, I want to explain the system so that another engineer can run, understand, and maintain it.

### README requirements

Document:

* Project purpose
* Architecture
* Prerequisites
* CLI setup
* CLI usage
* API setup
* API usage
* React setup
* Configuration
* Test commands
* Example request
* Example response
* Known limitations
* Future improvements

### Additional documents

Create:

```text
docs/architecture.md
docs/decisions.md
```

Record at least three decisions:

* Why local disk storage was selected
* Why variant concurrency is bounded
* Why the frontend sends one request per image
* Why the CLI and HTTP API share one application service

### Acceptance criteria

* A new developer can run the application from the README.
* The architecture document explains all three stages.
* Important tradeoffs are documented.
* The complete system can be demonstrated with multiple images.

---

# One-week schedule

## Day 1: Go foundation and project setup

Complete:

* Story 1
* Story 2
* Story 3
* Begin Story 4

Focus on:

* Modules
* Packages
* Structs
* Methods
* Validation
* Configuration
* Table-driven tests

End-of-day checkpoint:

```bash
go test ./...
```

must succeed.

The CLI should accept arguments, even if image processing is still a placeholder.

---

## Day 2: Storage and image validation

Complete:

* Story 4
* Story 5
* Story 6

Focus on:

* File I/O
* Interfaces
* Readers and writers
* Temporary test directories
* Image metadata
* Typed errors

End-of-day checkpoint:

* The CLI accepts one file.
* Valid JPEG and PNG files are recognized.
* Invalid files are rejected.
* Original files can be saved safely.

---

## Day 3: Image transformations and CLI completion

Complete:

* Story 7
* Story 8

Begin:

* Story 9

Focus on:

* Resizing
* Cropping
* Encoding
* Aspect ratios
* Manifests
* Result aggregation

End-of-day checkpoint:

One CLI command generates at least three variants:

```text
thumbnail
square
banner
```

This is the first major project milestone.

---

## Day 4: Concurrency and application service

Complete:

* Story 9
* Story 10

Focus on:

* Goroutines
* Channels
* Wait groups
* Concurrency limits
* Context cancellation
* Dependency injection

End-of-day checkpoint:

* Variant generation uses bounded concurrency.
* The CLI calls the shared application service.
* The race detector passes.

```bash
go test -race ./...
```

---

## Day 5: HTTP API

Complete:

* Story 11
* Story 12
* Story 13

Focus on:

* `net/http`
* Multipart uploads
* JSON responses
* Status codes
* Health routes
* Graceful shutdown
* Logging

End-of-day checkpoint:

An image can be uploaded with:

```bash
curl \
  -X POST \
  -F "image=@./testdata/images/cafe.jpg" \
  http://localhost:8080/v1/assets
```

The response should contain an asset manifest.

---

## Day 6: React interface

Complete:

* Story 14
* Story 15
* Begin Story 17

Focus on:

* Multiple file selection
* File queue state
* Sequential uploads
* API integration
* Success and failure display

End-of-day checkpoint:

* Several files can be selected.
* Each file is uploaded independently.
* Files are processed sequentially.
* Results are displayed per file.

---

## Day 7: Frontend concurrency, testing, and documentation

Complete:

* Story 16
* Story 17
* Story 18
* Any unfinished required work

Run Go checks:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Run frontend checks:

```bash
cd web
npm run build
```

End-of-day checkpoint:

* Multiple files can be selected.
* Two files can be uploaded concurrently.
* Each uploaded image generates multiple variants.
* Every file displays its own result.
* Failed files can be retried.
* The project is documented.

---

# Definition of done

The project is complete when:

1. A local JPEG or PNG can be processed through the CLI.
2. The CLI processes exactly one image.
3. The CLI generates at least three variants.
4. The CLI returns a JSON manifest.
5. The original file is stored.
6. Generated files are stored under an asset directory.
7. Image validation rejects unsupported and corrupted files.
8. Variant processing uses bounded concurrency.
9. The CLI and HTTP API use the same application service.
10. One image can be uploaded through the HTTP API.
11. The HTTP API returns a JSON manifest.
12. Health endpoints are available.
13. The server shuts down gracefully.
14. Logs contain useful processing information.
15. The React application accepts multiple selected files.
16. The frontend sends one HTTP request per file.
17. The frontend tracks state independently for every file.
18. Failed files do not stop the remaining queue.
19. Frontend upload concurrency is limited.
20. Results are shown for each source image.
21. Failed files can be retried.
22. Unit and integration tests pass.
23. The Go race detector passes.
24. The React application builds successfully.
25. The README explains the CLI, API, and frontend workflows.

---

# Final intern presentation

At the end of the project, the developer should be prepared to explain:

1. How a CLI request moves through the application
2. How an HTTP upload moves through the application
3. Why the CLI and HTTP API share one application service
4. Why media logic is separate from HTTP logic
5. Why storage is represented by an interface
6. How aspect ratios are preserved
7. How center cropping works
8. How partial failures are represented
9. How Go concurrency is limited
10. How frontend concurrency is limited
11. Why the API processes one image per request
12. How context cancellation propagates
13. How the system could later use S3
14. How the system could later use a queue
15. What would be required to run the system in production

---

# Coding agent instructions

## Role

Act as a senior Go engineer preparing a one-week internship project for a developer learning Go.

Create the project scaffold, development tooling, package boundaries, placeholder types, tests, and documentation.

Do not implement the completed application.

The project should teach the developer how to build the system.

---

## Required project stages

The scaffold must support three progressive stages:

1. A Go CLI that processes one local image
2. A Go HTTP API that processes one uploaded image
3. A React and TypeScript interface that processes a list of selected files

The CLI must be designed first.

The HTTP API must reuse the CLI's application service.

The React application must send one request per image.

Do not create a batch upload endpoint.

---

## Required scaffold

Create:

* Go module
* `cmd/cli`
* `cmd/api`
* Shared internal packages
* React TypeScript application under `web`
* Local data directories
* Test-data directories
* Documentation directories
* Makefile
* `.env.example`
* `.gitignore`
* README
* Sample media profile configuration

---

## Agent responsibilities

The coding agent should:

1. Initialize the Go module.
2. Create the complete directory structure.
3. Add package declarations.
4. Add package documentation.
5. Add domain structs where useful.
6. Add interfaces where they clarify architecture.
7. Add constructor signatures.
8. Add typed placeholder errors.
9. Add function signatures.
10. Add TODO comments.
11. Add skipped or intentionally incomplete tests.
12. Add empty React components.
13. Add frontend response types.
14. Add placeholder API client functions.
15. Add Makefile commands.
16. Add setup and run instructions.
17. Ensure the scaffold compiles whenever practical.
18. Map project stories to the files where they should be implemented.

---

## Agent restrictions

Do not implement:

* Image resizing
* Image cropping
* Image encoding
* Media profile execution
* The completed processing pipeline
* The completed application service
* The local storage implementation
* The HTTP upload handler
* The frontend upload queue
* Frontend concurrency management
* Retry behavior
* Completed test answers

Do not:

* Add a database
* Add cloud infrastructure
* Add authentication
* Add a frontend framework beyond React
* Add complex state management
* Add a batch upload endpoint
* Put all business logic in `main.go`
* Hide implementation inside generated helpers
* Provide line-by-line solutions inside comments
* Over-engineer the architecture

---

## Placeholder behavior

Function bodies may:

* Return `ErrNotImplemented`
* Return placeholder values
* Contain a short TODO
* Validate basic constructor inputs
* Define expected inputs and outputs

Function bodies must not contain completed business logic.

The scaffold should remain compilable where possible.

---

## Required architecture

The CLI and HTTP API should call the same application service.

Conceptual application flow:

```text
Process one source image
    |
    ├── validate
    ├── create asset ID
    ├── store original
    ├── generate variants
    ├── store generated files
    └── return manifest
```

The React application should manage:

```text
Selected files
    |
    ├── pending
    ├── uploading
    ├── processing
    ├── completed
    └── failed
```

The first frontend queue should be sequential.

A later story should add a concurrency limit of two uploads.

---

## Required Makefile commands

Include commands conceptually similar to:

```text
make run-cli
make run-api
make test
make test-race
make fmt
make vet
make web-install
make web-run
make web-build
make clean
```

Exact command names may vary if documented clearly.

---

## Required agent output

After generating the scaffold, provide:

1. The complete file tree
2. A description of every Go package
3. A description of the React structure
4. The first story the intern should begin
5. The exact files associated with that story
6. Commands for running the CLI scaffold
7. Commands for running the API scaffold
8. Commands for running the React scaffold
9. Commands for running tests
10. A list of architectural assumptions
11. Confirmation that the implementation remains intentionally unfinished

---

# Mentor rules for future coding assistance

The coding assistant should follow these rules while helping the intern:

* Explain the relevant Go concept before proposing code.
* Ask the intern to attempt the implementation first.
* Review the intern's code rather than replacing it.
* Point out incorrect behavior directly.
* Provide small hints before complete solutions.
* Do not rewrite an entire package unless explicitly requested.
* Prefer idiomatic Go.
* Explain why an approach is idiomatic.
* Require tests for every story.
* Encourage `go doc`.
* Encourage `go test`.
* Encourage `go vet`.
* Encourage the race detector.
* Keep the scope achievable within one week.
* Treat mistakes as learning opportunities.
* Avoid unnecessary abstractions.
* Favor clear code over clever code.
