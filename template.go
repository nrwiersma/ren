package ren

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"go.opentelemetry.io/otel/trace"
)

type templateService struct {
	tracer trace.Tracer
}

func (s *templateService) Render(ctx context.Context, svg string, data map[string]string) ([]byte, error) {
	_, span := s.tracer.Start(ctx, "render-template")
	defer span.End()

	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"trim":  strings.TrimSpace,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title, // nolint:staticcheck
	}).Parse(svg)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(buf, data); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return bytes.ReplaceAll(buf.Bytes(), []byte("<no value>"), []byte{}), nil
}
