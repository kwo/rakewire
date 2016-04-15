package fever

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	hContentType = "Content-Type"
	mimeJSON     = "text/json; charset=utf-8"
	mPost        = "POST"
)

const (
	itemRead      = "read"
	itemUnread    = "unread"
	itemStarred   = "saved"
	itemUnstarred = "unsaved"
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

// formatID takes a uint64 and formats it as a model string ID
func formatID(value uint64) string {
	return encodeID(strconv.FormatUint(value, 10))
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
