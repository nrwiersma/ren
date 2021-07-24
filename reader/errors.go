package reader

import "github.com/hamba/pkg/v2/errors"

// ErrTemplateNotFound is returned when the requested template does
// not exist.
const ErrTemplateNotFound = errors.Error("template not found")
