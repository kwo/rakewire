package httpd

import (
	"net/http"
)

// NoCache adds cache-control headers so that the content is not cached
func NoCache() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache")
			h.ServeHTTP(w, r)
		})
	}
}
