package ren

import (
	"bytes"
	"path/filepath"
	"text/template"
	"errors"
	"os"
)

var (
	ErrTemplate         = errors.New("error rendering template")
	ErrTemplateNotFound = errors.New("template not found")
	ErrTemplateInvalid  = errors.New("template invalid")
)

// Application represents the application.
type Application struct {
	templates string
}

// NewApplication creates an instance of Application.
func NewApplication(t string) *Application {
	return &Application{
		templates: t,
	}
}

// Render renders a template with the given data.
func (a *Application) Render(path string, data interface{}) ([]byte, error) {
	path = filepath.Join(a.templates, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrTemplateNotFound
	}

	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, ErrTemplate
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, data)
	if err != nil {
		return nil, ErrTemplateInvalid
	}

	return buf.Bytes(), nil
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
