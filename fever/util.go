package fever

import (
	"fmt"
	"net/http"
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

// encodeID takes a string ID from the Fever REST API and formats it as a model string ID
func encodeID(value string) string {
	return fmt.Sprintf("%010s", strings.TrimLeft(value, "0"))[:10]
}

// decodeID takes a string ID from model and converts it to a string for fever structs with string array.
func decodeID(value string) string {
	return strings.TrimLeft(value, "0")
}

// parseID takes a string ID from model and converts it to a uint64 for fever structs.
func parseID(value string) uint64 {
	if x, err := strconv.ParseUint(value, 10, 64); err == nil {
		return x
	}
	return 0
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
