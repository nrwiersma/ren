// Package api implements a HTTP api.
package api

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hamba/logger/v2"
	lctx "github.com/hamba/logger/v2/ctx"
	mdlw "github.com/hamba/pkg/v2/http/middleware"
	"github.com/hamba/statter/v2"
	"github.com/nrwiersma/ren"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Application represents the main application.
type Application interface {
	// Render renders a template with the given data.
	Render(ctx context.Context, path string, data map[string]string) ([]byte, error)
}

// API is an http api handler.
type API struct {
	app Application

	h http.Handler

	log    *logger.Logger
	stats  *statter.Statter
	tracer trace.Tracer
}

// New returns an api handler.
func New(app Application, log *logger.Logger, stats *statter.Statter, tracer trace.Tracer) *API {
	api := &API{
		app:    app,
		log:    log,
		stats:  stats,
		tracer: tracer,
	}

	api.h = api.routes()

	return api
}

func (a *API) routes() http.Handler {
	mux := chi.NewRouter()
	mux.With(mdlw.Stats("not-found", a.stats)).NotFound(http.NotFound)
	mux.With(mdlw.Stats("image", a.stats)).Get("/{group}/{file}", a.handleRenderImage())

	r := mdlw.WithRecovery(mux, a.log)
	return otelhttp.NewHandler(r, "server")
}

// ServeHTTP serves and http request.
func (a *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.h.ServeHTTP(rw, req)
}

// handleRenderImage handles requests to render an image.
func (a *API) handleRenderImage() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, span := a.tracer.Start(req.Context(), "image")
		defer span.End()

		group := chi.URLParam(req, "group")
		file := chi.URLParam(req, "file")

		data := map[string]string{}
		for k := range req.URL.Query() {
			data[k] = req.URL.Query().Get(k)
		}

		img, err := a.app.Render(ctx, filepath.Join(group, file+".svg"), data)
		if err != nil {
			span.RecordError(err)

			switch {
			case errors.Is(err, ren.ErrTemplateNotFound):
				a.log.Debug("Could not find template")
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				a.log.Error("Could not render template", lctx.Error("error", err))
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		rw.Header().Set("Content-Type", "image/svg+xml")
		_, _ = rw.Write(img)
	}
}
