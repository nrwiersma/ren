// Package api implements a HTTP api.
package api

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-zoo/bone"
	"github.com/hamba/logger/v2"
	logCtx "github.com/hamba/logger/v2/ctx"
	httpx "github.com/hamba/pkg/v2/http"
	mdlw "github.com/hamba/pkg/v2/http/middleware"
	"github.com/hamba/statter/v2"
	"github.com/nrwiersma/ren"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Application represents the main application.
type Application interface {
	httpx.Health

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
	h := bone.New()
	h.NotFound(mdlw.WithStats("not-found", a.stats, http.NotFoundHandler()))
	h.Get(httpx.DefaultHealthPath, mdlw.WithStats("health", a.stats, httpx.NewHealthHandler(a.app)))

	h.Get("/:group/:file", mdlw.WithStats("image", a.stats, a.handleImage()))

	r := mdlw.WithRecovery(h, a.log)
	return otelhttp.NewHandler(r, "server")
}

// ServeHTTP serves and http request.
func (a *API) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.h.ServeHTTP(rw, req)
}

// handleImage handles requests to render an image.
func (a *API) handleImage() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, span := a.tracer.Start(req.Context(), "image")
		defer span.End()

		group := bone.GetValue(req, "group")
		file := bone.GetValue(req, "file")

		data := map[string]string{}
		for k := range req.URL.Query() {
			data[k] = req.URL.Query().Get(k)
		}

		img, err := a.app.Render(ctx, filepath.Join(group, file+".svg"), data)
		if err != nil {
			span.RecordError(err)

			if errors.Is(err, ren.ErrTemplateNotFound) {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			a.log.Error("Could not render template", logCtx.Error("error", err))
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "image/svg+xml")
		_, _ = rw.Write(img)
	}
}
