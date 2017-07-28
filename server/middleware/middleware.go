package middleware

import (
	"context"
	"net/http"

	"github.com/msales/pkg/httpx/middleware"
)

// WithContext wraps pkg WithContext.
func WithContext(h http.Handler, ctx context.Context) http.Handler {
	return middleware.WithContext(h, ctx)
}

// Common wraps the handler with common middleware.
func Common(h http.Handler) http.Handler {
	h = middleware.WithResponseTime(h)
	h = middleware.WithRequestStats(h)

	return middleware.WithRecovery(h)
}
