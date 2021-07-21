package reader

import "errors"

// ErrTemplateNotFound is returned when the requested template does
// not exist.
var ErrTemplateNotFound = errors.New("template not found")
