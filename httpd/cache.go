package httpd

import (
	"net/http"

	"golang.org/x/net/context"
)

// NoCache adds cache-control headers so that the content is not cached
func NoCache() Middleware {
	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache")
			next.ServeHTTPC(ctx, w, r)
		})
	}
}
