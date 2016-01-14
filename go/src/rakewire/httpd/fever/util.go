package fever

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	logName  = "[fever]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

const (
	hAcceptEncoding  = "Accept-Encoding"
	hContentEncoding = "Content-Encoding"
	hContentType     = "Content-Type"
	mPost            = "POST"
	mimeJSON         = "text/json; charset=utf-8"
)

func boolToUint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func parseID(values url.Values, key string) uint64 {
	if value := values.Get(key); value != "" {
		x, _ := strconv.ParseUint(value, 10, 64)
		return x
	}
	return 0
}

func parseIDArray(values url.Values, key string) []uint64 {
	if value := values.Get(key); value != "" {
		valueElements := strings.Split(value, ",")
		idElements := []uint64{}
		for _, valueElement := range valueElements {
			idElement, _ := strconv.ParseUint(valueElement, 10, 64)
			if idElement > 0 {
				idElements = append(idElements, idElement)
			}
		}
		return idElements
	}
	return nil
}

func notFound(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusNotFound)
}

func notSupported(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusMethodNotAllowed)
}

func sendError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
