package reader

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// FileReader is a local file reader.
type FileReader struct {
	base string

	tracer trace.Tracer
}

// NewFileReader returns a file reader.
func NewFileReader(path string, tracer trace.Tracer) *FileReader {
	if strings.HasPrefix(path, "/.") {
		path = path[1:]
	}

	return &FileReader{
		base:   path,
		tracer: tracer,
	}
}

// Read reads the file at the given path.
func (r *FileReader) Read(ctx context.Context, path string) (string, error) {
	_, span := r.tracer.Start(ctx, "read")
	defer span.End()

	path = filepath.Join(r.base, path)
	if _, err := os.Stat(path); err != nil {
		span.RecordError(err)

		if os.IsNotExist(err) {
			return "", ErrTemplateNotFound
		}

		return "", err
	}
	path = filepath.Clean(path)

	b, err := os.ReadFile(path)
	if err != nil {
		span.RecordError(err)

		return "", err
	}

	return string(b), nil
}
