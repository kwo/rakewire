package middleware

import (
	"net/http"
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
