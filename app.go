package ren

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
)

// Application errors
var (
	ErrTemplateNotFound = errors.New("template not found")
)

// Reader represents a template Reader.
type Reader interface {
	Read(path string) (string, error)
}

// Application represents the application.
type Application struct {
	logger  log.Logger
	statter stats.Statter

	Reader Reader
}

// NewApplication creates an instance of Application.
func NewApplication(l log.Logger, s stats.Statter) *Application {
	return &Application{
		logger:  l,
		statter: s,
	}
}

// Render renders a template with the given data.
func (a *Application) Render(path string, data interface{}) ([]byte, error) {
	svg, err := a.Reader.Read(path)
	if err != nil {
		return nil, ErrTemplateNotFound
	}

	tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(svg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(buf, data); err != nil {
		return nil, err
	}

	return bytes.Replace(buf.Bytes(), []byte("<no value>"), []byte{}, -1), nil
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}

// Logger returns the Logger attached to the Application.
func (a *Application) Logger() log.Logger {
	return a.logger
}

// Statter returns the Statter attached to the Application.
func (a *Application) Statter() stats.Statter {
	return a.statter
}
