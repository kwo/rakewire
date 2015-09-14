package httpd

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"net/http"
	"os"
)

const (
	hCacheControl = "Cache-Control"
	vNoCache      = "no-cache"
)

const (
	optionNone = 0
)

// Adapter creates middleware.
type Adapter func(http.Handler) http.Handler

// Adapt calls adapters for http handler
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

// LogAdapter log requests and responses
func LogAdapter() Adapter {
	return func(h http.Handler) http.Handler {
		return gorillaHandlers.CombinedLoggingHandler(os.Stdout, h)
	}
}

// NoCache adds cache-control headers so that the content is not cached
func NoCache() Adapter {
	return cacheControl(optionNone)
}

func cacheControl(option int) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch option {
			case optionNone:
				w.Header().Set(hCacheControl, vNoCache)
			}
			h.ServeHTTP(w, r)
		})
	}
}
