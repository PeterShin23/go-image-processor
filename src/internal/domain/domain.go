// Package domain holds the core business types for the media pipeline:
// assets, media profiles, generated variants, manifests, and their statuses.
//
// This package must not import HTTP, filesystem, or React concerns. It is the
// shared vocabulary every other package speaks.
//
// STORY 02: flesh out the types and their Validate methods below. The fields
// shown are a starting shape — add, rename, or remove fields as you see fit,
// but keep the package free of transport/storage details.
package domain

import "fmt"

// OutputFormat identifies an encoded image format.
//
// TODO (story 02): decide the set of supported formats and add a Valid method.
type OutputFormat string

const (
	FormatJPEG OutputFormat = "jpeg"
	FormatPNG  OutputFormat = "png"
)

// ResizeMode describes how a source image is fit into a variant's dimensions.
type ResizeMode string

const (
	// ResizeFit scales to fit within the box without cropping (may leave one
	// dimension smaller than the target).
	ResizeFit ResizeMode = "fit"
	// ResizeCrop scales to cover the box, then center-crops to exact size.
	ResizeCrop ResizeMode = "crop"
)

// ProcessingStatus is the overall outcome of processing one asset.
type ProcessingStatus string

const (
	StatusCompleted ProcessingStatus = "completed"
	StatusPartial   ProcessingStatus = "partial"
	StatusFailed    ProcessingStatus = "failed"
)

// MediaProfile describes one desired output variant.
type MediaProfile struct {
	Name         string
	Width        int
	Height       int
	ResizeMode   ResizeMode
	Format       OutputFormat
	Quality      int  // 1..100, meaningful for JPEG
	// TODO (story 02): add fields you need and remove ones you don't.
}

func (v ResizeMode) Valid() bool {
	switch v {
		case ResizeCrop, ResizeFit:
			return true
		default:
			return false
	}
}

func (v OutputFormat) Valid() bool {
	switch v {
		case FormatJPEG, FormatPNG:
			return true
		default:
			return false
	}
}

// Validate reports whether the profile is well-formed.
//
// TODO (story 02): reject non-positive width/height, out-of-range quality,
// unknown resize modes, and unknown formats. Return a descriptive error.
func (p MediaProfile) Validate() error {
	if p.Width <= 0 || p.Height <= 0 {
		return fmt.Errorf("width and height must be >= 0: %d width, %d height", p.Width, p.Height)
	}
	if p.Quality < 1 || p.Quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100: %d quality", p.Quality)
	}
	
	if !p.ResizeMode.Valid() {
		return fmt.Errorf("unknown resizeMode: %q", p.ResizeMode)
	}

	if !p.Format.Valid() {
		return fmt.Errorf("unknown format: %q", p.Format)
	}

	return nil
}

// ValidateProfiles checks a set of profiles, including that names are unique.
//
// TODO (story 02): reject duplicate profile names and any invalid profile.
func ValidateProfiles(profiles []MediaProfile) error {
	seen := make(map[string]bool)

	for _, p := range profiles {
		if err := p.Validate(); err != nil {
			return err
		}
		if seen[p.Name] {
			return fmt.Errorf("duplicate profile name: %q", p.Name)
		}
		seen[p.Name] = true
	}

	return nil
}

// SourceImage describes a decoded input image (not its pixels — just the facts
// validation and the pipeline need to know).
type SourceImage struct {
	Width  int
	Height int
	Format OutputFormat
}

// Variant is one successfully generated output file's metadata.
type Variant struct {
	Name      string
	Width     int
	Height    int
	Format    OutputFormat
	Path      string
	SizeBytes int64
}

// ProcessingError records why a single profile failed, without stopping others.
type ProcessingError struct {
	Profile string
	Stage   string // e.g. "decode", "resize", "encode", "store"
	Message string
}

// ProcessingResult is what the pipeline returns: the variants it managed to
// build and the failures it collected. The app service turns this (plus asset
// identity) into a Manifest.
type ProcessingResult struct {
	Variants []Variant
	Errors   []ProcessingError
}

// Status derives the overall status from the counts of successes and failures.
//
// TODO (story 08): completed when all requested profiles succeeded, failed when
// none did, partial otherwise.
func (r ProcessingResult) Status(requested int) ProcessingStatus {
	return StatusFailed // TODO: replace with real logic
}

// Manifest is the result returned for one processed asset.
type Manifest struct {
	AssetID          string
	OriginalFilename string
	Status           ProcessingStatus
	Variants         []Variant
	Errors           []ProcessingError
}
