package server

import (
	"net/http"
	"path/filepath"

	"github.com/go-zoo/bone"
	"github.com/msales/pkg/log"
	"github.com/nrwiersma/ren"
)

// Application represents the main application.
type Application interface {
	// Render renders a template with the given data.
	Render(path string, data interface{}) ([]byte, error)
	// IsHealthy checks the health of the Application.
	IsHealthy() error
}

// Server represents a http server handler.
type Server struct {
	app Application
	mux *bone.Mux
}

// New creates a new Server instance.
func New(app Application) *Server {
	s := &Server{
		app: app,
		mux: bone.New(),
	}

	s.mux.GetFunc("/:group/:file", s.ImageHandler)

	s.mux.GetFunc("/health", s.HealthHandler)
	s.mux.NotFound(NotFoundHandler())

	return s
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ImageHandler handles requests to render an image.
func (s *Server) ImageHandler(w http.ResponseWriter, r *http.Request) {
	group := bone.GetValue(r, "group")
	file := bone.GetValue(r, "file")
	path := filepath.Join(group, file+".svg")

	data := map[string]string{}
	for k := range r.URL.Query() {
		data[k] = r.URL.Query().Get(k)
	}

	img, err := s.app.Render(path, data)
	if err != nil {
		switch err {
		case ren.ErrTemplateNotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		default:
			log.Error(r.Context(), "could not render template", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(img)
}

// HealthHandler handles health requests.
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.app.IsHealthy(); err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// NotFoundHandler returns a 404.
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}
