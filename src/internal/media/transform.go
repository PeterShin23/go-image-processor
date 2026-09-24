package media

import (
	"image"
	"io"

	"github.com/PeterShin23/go-image-processor/internal/domain"
)

// Transform turns one already-decoded source image into one variant according
// to the profile, encodes it, and writes the bytes to w. It returns the final
// output dimensions.
//
// STORY 07: implement the resize + optional center-crop + encode. Keep the
// steps as small helpers (e.g. resizeFit, centerCrop, encode) so each is
// testable. Honor p.ResizeMode, p.Format, p.Quality, and p.AllowEnlarge.
//
// TODO: replace the placeholder return.
func Transform(w io.Writer, src image.Image, p domain.MediaProfile) (width int, height int, err error) {
	return 0, 0, domain.ErrNotImplemented
}
