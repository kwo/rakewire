package httpd

import (
	"golang.org/x/net/context"
	"net/http"
)

// HandlerC is a context-aware http handler
type HandlerC interface {
	ServeHTTPC(context.Context, http.ResponseWriter, *http.Request)
}

// HandlerFuncC is a function to create HandlerC
type HandlerFuncC func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTPC calls f(ctx, w, r).
func (f HandlerFuncC) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

// Middleware creates middleware.
type Middleware func(next HandlerC) HandlerC

// Chain calls middleware for HandlerC
func Chain(h HandlerC, middlewares ...Middleware) HandlerC {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// Adapt adapts a context-aware handler to a normal http handler
func Adapt(ctx context.Context, h HandlerC) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTPC(ctx, w, r)
	})
}
