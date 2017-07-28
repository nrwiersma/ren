package ren

import (
	"bytes"
	"path/filepath"
	"text/template"
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

	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	t.Execute(buf, data)

	return buf.Bytes(), nil
}

// IsHealthy checks the health of the Application.
func (a *Application) IsHealthy() error {
	return nil
}
