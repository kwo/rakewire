package httpd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (z *Httpd) apiRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api", feedsHandler)
	return router
}

// func (z *Httpd) feedsRouter() *mux.Router {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/", feedsHandler)
// 	return router
// }

// APIHandler handles API requests
func feedsHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Welcome to the feedsHandler"))
	// vars := mux.Vars(req)
	// key := vars["id"]
}
