package httpd

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"log"
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

// LogWriter is an io.Writer than writes to the log facility.
type LogWriter struct {
	name  string
	level string
}

func (z *LogWriter) Write(p []byte) (n int, err error) {
	log.Printf("%-7s %-7s %s", z.level, z.name, string(p))
	return len(p), nil
}

// LogAdapter log requests and responses
func LogAdapter(level string) Adapter {

	if level != "" {
		level = "[" + level + "]"
	} else {
		level = "[TRACE]"
	}

	logWriter := &LogWriter{
		name:  "[access]",
		level: level,
	}

	return func(h http.Handler) http.Handler {
		return gorillaHandlers.LoggingHandler(logWriter, h)
	}

}

// RedirectHandler permanently redirects to the given location
func RedirectHandler(location string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, location, http.StatusMovedPermanently)
	})
}

// NoCache adds cache-control headers so that the content is not cached
func NoCache() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache")
			h.ServeHTTP(w, r)
		})
	}
}
