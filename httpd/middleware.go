package httpd

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"net/http"
	"rakewire/logger"
	"rakewire/middleware"
)

// LogWriter is an io.Writer than writes to the log facility.
type LogWriter struct {
	accessLogger *logger.Logger
}

func (z *LogWriter) Write(p []byte) (n int, err error) {
	log.Debugf("%s", string(p))
	return len(p), nil
}

// LogAdapter log requests and responses
func LogAdapter() middleware.Adapter {

	logWriter := &LogWriter{
		accessLogger: logger.New("access"),
	}

	return func(h http.Handler) http.Handler {
		return gorillaHandlers.LoggingHandler(logWriter, h)
	}

}
