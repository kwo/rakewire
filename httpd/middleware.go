package httpd

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"log"
	"net/http"
	"rakewire/middleware"
)

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
func LogAdapter(level string) middleware.Adapter {

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
