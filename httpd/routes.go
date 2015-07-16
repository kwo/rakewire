package httpd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	mGet = "GET"
	mPut = "PUT"
)

func (z *Httpd) apiRouter(prefix string) *mux.Router {
	router := mux.NewRouter()

	var prefixFeeds = prefix + "/feeds"
	router.Methods(mGet).Path(prefixFeeds).HandlerFunc(z.feedsGet)
	router.Methods(mPut).Path(prefixFeeds).HandlerFunc(z.feedsSave)
	router.Path(prefixFeeds).HandlerFunc(notSupported)

	return router

}

func notSupported(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusMethodNotAllowed)
}

func sendError(w http.ResponseWriter, code int) {
	text := fmt.Sprintf("%d %s", code, http.StatusText(code))
	http.Error(w, text, code)
}
