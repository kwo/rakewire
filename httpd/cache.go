package httpd

import (
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"net/http"
)

// NoCache adds cache-control headers so that the content is not cached
func NoCache(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTPC(ctx, w, r)
	})
}
