package server

import (
	"net/http"
	"path/filepath"

	"github.com/go-zoo/bone"
	"github.com/hamba/pkg/httpx"
	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/nrwiersma/ren"
)

// Application represents the main application.
type Application interface {
	log.Loggable
	stats.Statable

	// Render renders a template with the given data.
	Render(path string, data interface{}) ([]byte, error)
}

// NewMux creates a new server mux.
func NewMux(app Application) *bone.Mux {
	mux := httpx.NewMux()

	mux.GetFunc("/:group/:file", ImageHandler(app))

	return mux
}

// ImageHandler handles requests to render an image.
func ImageHandler(app Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := bone.GetValue(r, "group")
		file := bone.GetValue(r, "file")
		path := filepath.Join(group, file+".svg")

		data := map[string]string{}
		for k := range r.URL.Query() {
			data[k] = r.URL.Query().Get(k)
		}

		img, err := app.Render(path, data)
		if err != nil {
			switch err {
			case ren.ErrTemplateNotFound:
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			default:
				log.Error(app, "could not render template", "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = w.Write(img)
	}
}
