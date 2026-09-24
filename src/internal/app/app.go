// Package app is the shared application service used by BOTH the CLI and the
// HTTP API. It coordinates validation, storing the original, running the
// pipeline, and building a manifest. It owns no transport (CLI/HTTP) details.
//
// STORY 10: implement NewService and Process. The service depends on
// abstractions (a storage.Storage, a *pipeline.Pipeline) passed into its
// constructor — it must not build those concrete dependencies itself.
package app

import (
	"context"
	"io"

	"github.com/PeterShin23/go-image-processor/internal/domain"
	"github.com/PeterShin23/go-image-processor/internal/media"
	"github.com/PeterShin23/go-image-processor/internal/pipeline"
	"github.com/PeterShin23/go-image-processor/internal/storage"
)

// Service is the one place the end-to-end "process an image" flow lives.
type Service struct {
	store    storage.Storage
	pipe     *pipeline.Pipeline
	limits   media.Limits
	profiles []domain.MediaProfile
	newID    func() string
}

// NewService wires the service's dependencies. newID generates a unique asset
// ID (inject one so tests can make it deterministic); pass nil to use a default.
//
// TODO (story 10): if newID is nil, set a sensible default generator.
func NewService(store storage.Storage, pipe *pipeline.Pipeline, limits media.Limits, profiles []domain.MediaProfile, newID func() string) *Service {
	return &Service{store: store, pipe: pipe, limits: limits, profiles: profiles, newID: newID}
}

// ProcessInput is the transport-neutral request. The CLI builds it from a file;
// the HTTP handler builds it from a multipart upload. Neither leaks into here.
type ProcessInput struct {
	Filename string
	Reader   io.Reader
	Size     int64
	// Profiles overrides the service default when non-empty (optional).
	Profiles []domain.MediaProfile
}

// Process runs the full flow and returns a manifest.
//
// TODO (story 10): read the input bytes (bounded), validate (media.Validate),
// create an asset ID, store the original, decode the image, run the pipeline
// (pipe.Generate), and assemble a domain.Manifest with the derived status.
func (s *Service) Process(ctx context.Context, in ProcessInput) (domain.Manifest, error) {
	return domain.Manifest{}, domain.ErrNotImplemented
}
