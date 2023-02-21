package reader

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"go.opentelemetry.io/otel/codes"
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
	if _, err := os.Stat(filepath.Clean(path)); err != nil {
		span.SetStatus(codes.Error, "Finding file")
		span.RecordError(err)

		if os.IsNotExist(err) {
			return "", ErrTemplateNotFound
		}

		return "", err
	}

	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		span.SetStatus(codes.Error, "Reading file")
		span.RecordError(err)

		return "", err
	}

	return string(b), nil
}
