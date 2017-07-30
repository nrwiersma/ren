package server

import (
	"io"
	"net/http"

	"path/filepath"

	"github.com/go-zoo/bone"
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

	s.mux.Get("/:group/:file", ImageHandler(app))

	s.mux.Get("/health", HealthHandler(app))

	return s
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ImageHandler returns a image rendering handler using Render method from a Application instance.
func ImageHandler(app *ren.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := bone.GetValue(r, "group")
		file := bone.GetValue(r, "file")
		path := filepath.Join(group, file+".svg")

		data := map[string]string{}
		for k, _ := range r.URL.Query() {
			data[k] = r.URL.Query().Get(k)
		}

		img, err := app.Render(path, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(img)
	})
}

// HealthHandler returns a health handler using IsHealthy method from a Application instance.
func HealthHandler(app *ren.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := app.IsHealthy(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
