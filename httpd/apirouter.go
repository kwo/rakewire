package httpd

import (
	"github.com/gorilla/mux"
	"net/http"
)

// APIRouter router for API
func APIRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{id}", APIHandler)
	return router
}

// APIHandler handles API requests
func APIHandler(w http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// key := vars["id"]
}
