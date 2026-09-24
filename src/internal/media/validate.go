package media

import (
	"errors"
	"io"

	"github.com/PeterShin23/go-image-processor/internal/domain"
)

// Typed validation errors. Callers (and the HTTP layer in story 11) can use
// errors.Is to map these to user-facing messages / status codes.
var (
	ErrEmptyFile         = errors.New("image file is empty")
	ErrTooLarge          = errors.New("image file exceeds the maximum allowed size")
	ErrUnsupportedFormat = errors.New("unsupported image format (only JPEG and PNG)")
	ErrCorruptImage      = errors.New("image could not be decoded")
	ErrDimensionsRange   = errors.New("image dimensions are outside the allowed range")
)

// Limits are the validation bounds, sourced from config.
type Limits struct {
	MaxBytes  int64
	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int
}

// Validate checks an image stream and returns its decoded facts on success.
// size is the known byte length (from the file info or upload header).
//
// TODO (story 06): reject empty and oversized inputs, sniff the real format
// (do NOT trust the extension), decode to confirm it isn't corrupt, and check
// the dimensions against limits. Return one of the typed errors above (wrapped
// with context) on failure.
func Validate(r io.Reader, size int64, limits Limits) (domain.SourceImage, error) {
	return domain.SourceImage{}, domain.ErrNotImplemented
}
