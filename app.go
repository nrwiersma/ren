// Package ren implements the business logic layer.
package ren

import (
	"context"
	"errors"

	"github.com/hamba/logger/v2"
	errorx "github.com/hamba/pkg/v2/errors"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/tags"
	"github.com/nrwiersma/ren/reader"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ErrTemplateNotFound occurs when a template cannot be found.
const ErrTemplateNotFound = errorx.Error("template not found")

// Reader represents a template Reader.
type Reader interface {
	Read(ctx context.Context, path string) (string, error)
}

// Application represents the application.
type Application struct {
	tmplSvc *templateService
	reader  Reader

	log    *logger.Logger
	stats  *statter.Statter
	tracer trace.Tracer
}

// NewApplication creates an instance of Application.
func NewApplication(r Reader, log *logger.Logger, stats *statter.Statter, tracer trace.TracerProvider) *Application {
	tmplSvc := &templateService{tracer: tracer.Tracer("template-service")}

	return &Application{
		tmplSvc: tmplSvc,
		reader:  r,
		log:     log,
		stats:   stats,
		tracer:  tracer.Tracer("app"),
	}
}

// Render renders a template with the given data.
func (a *Application) Render(ctx context.Context, path string, data map[string]string) ([]byte, error) {
	ctx, span := a.tracer.Start(ctx, "render", trace.WithAttributes(
		attribute.String("path", path),
	))
	defer span.End()

	svg, err := a.reader.Read(ctx, path)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, reader.ErrTemplateNotFound) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}

	a.stats.Counter("rendered", tags.Str("path", path)).Inc(1)

	return a.tmplSvc.Render(ctx, svg, data)
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
