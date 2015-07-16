package httpd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (z *Httpd) apiRouter(prefix string) *mux.Router {

	router := mux.NewRouter()

	var prefixAdmin = prefix + "/admin"
	var prefixAdminRepair = prefixAdmin + "/repairdb"
	router.Path(prefixAdminRepair).Methods(mPost).HandlerFunc(z.repairDatabase)
	router.Path(prefixAdminRepair).HandlerFunc(notSupported)
	router.Path(prefixAdmin).HandlerFunc(notFound)

	var prefixFeeds = prefix + "/feeds"
	router.Path(prefixFeeds).Methods(mGet).HandlerFunc(z.feedsGet)
	router.Path(prefixFeeds).Methods(mPut).Headers(hContentType, mimeJSON).HandlerFunc(z.feedsSaveJSON)
	router.Path(prefixFeeds).Methods(mPut).Headers(hContentType, mimeText).HandlerFunc(z.feedsSaveText)
	router.Path(prefixFeeds).Methods(mPut).HandlerFunc(badMediaType)
	router.Path(prefixFeeds).HandlerFunc(notSupported)

	var prefixFeedsFeed = prefixFeeds + "/{feedID}"
	router.Path(prefixFeedsFeed).Methods(mGet).HandlerFunc(z.feedsGetFeed)
	router.Path(prefixFeedsFeed).HandlerFunc(notSupported)

	return router

}

func sendOK(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusOK)
}

func badMediaType(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusUnsupportedMediaType)
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
