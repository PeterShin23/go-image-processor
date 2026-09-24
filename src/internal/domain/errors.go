package domain

import "errors"

// ErrNotImplemented is returned by scaffold functions that have not been built
// yet. As you complete each story, replace the `return ErrNotImplemented`
// lines with real logic.
var ErrNotImplemented = errors.New("not implemented")
