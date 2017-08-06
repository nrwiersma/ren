package ren

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"text/template"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
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
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return bytes.Replace(buf.Bytes(), []byte("<no value>"), []byte{}, -1), nil
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
