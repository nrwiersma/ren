package server

import (
	"io"
	"net/http"

	"path/filepath"

	"github.com/go-zoo/bone"
	"github.com/msales/pkg/log"
	"github.com/nrwiersma/ren"
)

// Server represents a http server handler.
type Server struct {
	mux *bone.Mux
}

// New creates a new Server instance.
func New(app *ren.Application) *Server {
	s := &Server{
		mux: bone.New(),
	}

	s.mux.Get("/:group/:file", NewImageHandler(app))

	s.mux.Get("/health", NewHealthHandler(app))
	s.mux.NotFound(NotFoundHandler())

	return s
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type ImageHandler struct {
	Render func(string, interface{}) ([]byte, error)
}

// ImageHandler returns a image rendering handler using Render method from a Application instance.
func NewImageHandler(a *ren.Application) *ImageHandler {
	return &ImageHandler{
		Render: a.Render,
	}
}

func (h ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	group := bone.GetValue(r, "group")
	file := bone.GetValue(r, "file")
	path := filepath.Join(group, file+".svg")

	data := map[string]string{}
	for k := range r.URL.Query() {
		data[k] = r.URL.Query().Get(k)
	}

	img, err := h.Render(path, data)
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

type HealthHandler struct {
	IsHealthy func() error
}

// NewHealthHandler returns a health handler using IsHealthy method from a Application instance.
func NewHealthHandler(a *ren.Application) *HealthHandler {
	return &HealthHandler{
		IsHealthy: a.IsHealthy,
	}
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.IsHealthy(); err != nil {
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
