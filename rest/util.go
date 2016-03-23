package rest

import (
	"net/http"
)

const (
	logName = "[fever]"
	logWarn = "[WARN]"
)

const (
	hContentType = "Content-Type"
	mGet         = "GET"
	//	mPost        = "POST"
	mPut     = "PUT"
	mimeJSON = "text/json; charset=utf-8"
)

func notFound(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusNotFound)
}

func notSupported(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusMethodNotAllowed)
}

func sendError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
