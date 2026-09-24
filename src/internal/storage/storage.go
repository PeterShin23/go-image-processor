// Package storage abstracts where images live. The first implementation writes
// to the local filesystem, but the interface should allow swapping in object
// storage later without touching the rest of the app.
//
// STORY 05: implement LocalStorage's methods. The Storage interface below is a
// starting shape — adjust the method set if your design needs it, but keep
// callers depending on the interface, not the concrete type.
package storage

import (
	"context"
	"io"

	"github.com/PeterShin23/go-image-processor/internal/domain"
)

// Storage is the behavior the app service depends on. It hides where and how
// bytes are persisted.
type Storage interface {
	// SaveOriginal stores the uploaded/source image under the asset's directory
	// and returns the stored path and byte count.
	SaveOriginal(ctx context.Context, assetID, filename string, r io.Reader) (path string, size int64, err error)

	// SaveVariant stores one generated variant under the asset's directory.
	SaveVariant(ctx context.Context, assetID, variantName string, format domain.OutputFormat, r io.Reader) (path string, size int64, err error)

	// Open returns a reader for a previously stored path.
	Open(ctx context.Context, path string) (io.ReadCloser, error)
}

// LocalStorage writes to the local filesystem under two root directories.
type LocalStorage struct {
	originalsDir string
	generatedDir string
}

// NewLocalStorage constructs a LocalStorage. Dependencies (the two roots) are
// injected — the type never reaches for globals.
func NewLocalStorage(originalsDir, generatedDir string) *LocalStorage {
	return &LocalStorage{originalsDir: originalsDir, generatedDir: generatedDir}
}

// Compile-time proof that *LocalStorage satisfies Storage. If a method is
// missing or has the wrong signature, the build fails here with a clear message.
var _ Storage = (*LocalStorage)(nil)

func (s *LocalStorage) SaveOriginal(ctx context.Context, assetID, filename string, r io.Reader) (string, int64, error) {
	// TODO (story 05): build a safe path under originalsDir/<assetID>/, create
	// the directory, copy r to the file, and return the path and size.
	return "", 0, domain.ErrNotImplemented
}

func (s *LocalStorage) SaveVariant(ctx context.Context, assetID, variantName string, format domain.OutputFormat, r io.Reader) (string, int64, error) {
	// TODO (story 05): build a safe path under generatedDir/<assetID>/, create
	// the directory, copy r to the file, and return the path and size.
	return "", 0, domain.ErrNotImplemented
}

func (s *LocalStorage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	// TODO (story 05): open the file for reading. Guard against paths that
	// escape the storage roots.
	return nil, domain.ErrNotImplemented
}
