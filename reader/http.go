package reader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HTTPReader is an http file reader.
type HTTPReader struct {
	base   *url.URL
	client *http.Client

	tracer trace.Tracer
}

// NewHTTPReader returns an http file reader.
func NewHTTPReader(uri string, tracer trace.Tracer) (*HTTPReader, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid base url %q: %w", uri, err)
	}

	return &HTTPReader{
		base: base,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		tracer: tracer,
	}, nil
}

// Read reads the file at the given path.
func (r *HTTPReader) Read(ctx context.Context, path string) (string, error) {
	_, span := r.tracer.Start(ctx, "read", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	u, err := r.base.Parse(strings.TrimLeft(path, "/"))
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	switch resp.StatusCode {
	case http.StatusOK:
		return string(b), nil
	case http.StatusNotFound:
		span.RecordError(ErrTemplateNotFound)
		return "", ErrTemplateNotFound
	default:
		err = fmt.Errorf("unexpected status code %d", resp.StatusCode)
		span.RecordError(err)
		return "", err
	}
}
