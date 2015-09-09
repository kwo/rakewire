package httpd

import (
	"net/http"
)

const (
	hCacheControl = "Cache-Control"
	vNoCache      = "no-cache"
)

const (
	optionNone = 0
)

// CacheControl struct contains the ServeHTTP method and the option to be used
type CacheControl struct {
	option int
}

// NoCache returns an HTTP handler which will modify the Cache-Control headers in ServeHTTP.
func NoCache() *CacheControl {
	return &CacheControl{
		option: optionNone,
	}
}

// ServeHTTP modifies Cache-Control headers in the http.ResponseWriter
func (h *CacheControl) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO add ETag and Last-Modified headers
	// TODO use negroni.NewResponseWriter?
	switch h.option {
	case optionNone:
		w.Header().Set(hCacheControl, vNoCache)
	}
	next(w, r)
}
