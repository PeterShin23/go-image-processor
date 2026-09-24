// Package pipeline coordinates generating many variants from one source image,
// collecting successes and failures. Story 08 does this sequentially; story 09
// adds bounded concurrency and context cancellation.
package pipeline

import (
	"context"
	"image"

	"github.com/PeterShin23/go-image-processor/internal/domain"
	"github.com/PeterShin23/go-image-processor/internal/storage"
)

// Pipeline turns one decoded source image into many stored variants.
type Pipeline struct {
	store         storage.Storage
	maxConcurrent int
}

// New constructs a Pipeline. Storage is injected (an interface), and the
// concurrency limit comes from config (used in story 09).
func New(store storage.Storage, maxConcurrent int) *Pipeline {
	return &Pipeline{store: store, maxConcurrent: maxConcurrent}
}

// Generate builds every requested variant for one asset and returns a result
// aggregating successes and per-profile failures. A single failing profile must
// NOT abort the others.
//
// STORY 08: implement this SEQUENTIALLY (a simple loop). For each profile:
// transform into a buffer (media.Transform), store it (p.store.SaveVariant),
// and record either a domain.Variant or a domain.ProcessingError.
//
// STORY 09: revisit to run profiles concurrently with a bound of
// p.maxConcurrent, respecting ctx cancellation.
func (p *Pipeline) Generate(ctx context.Context, assetID string, src image.Image, profiles []domain.MediaProfile) (domain.ProcessingResult, error) {
	return domain.ProcessingResult{}, domain.ErrNotImplemented
}
