package httpd

//go:generate esc -o static.go -pkg httpd -prefix $PROJECT_HOME/web $PROJECT_HOME/web/public

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func (z *Service) mainRouter(useLocal bool) (*mux.Router, error) {

	// TODO: useLocal with go run, otherwise use embedded

	fs := Dir(useLocal, "/public")
	ofs := oneFS{name: "/index.html", root: fs}

	router := mux.NewRouter()

	// api router
	apiPrefix := "/api"
	router.PathPrefix(apiPrefix).Handler(
		Adapt(z.apiRouter(apiPrefix), NoCache()),
	)

	// HTML5 routes: any path without a dot (thus an extension)
	router.Path("/{route:[a-z0-9/-]+}").Handler(
		Adapt(http.FileServer(ofs), NoCache(), gorillaHandlers.CompressHandler),
	)

	// always redirect /index.html to /
	router.Path("/index.html").Handler(
		RedirectHandler("/"),
	)

	// static web site
	router.PathPrefix("/").Handler(
		Adapt(http.FileServer(fs), NoCache(), gorillaHandlers.CompressHandler),
	)

	return router, nil

}

func (z *Service) apiRouter(prefix string) *mux.Router {

	router := mux.NewRouter()

	var prefixAdmin = prefix + "/admin"
	var prefixAdminRepair = prefixAdmin + "/repairdb"
	router.Path(prefixAdminRepair).Methods(mPost).HandlerFunc(z.repairDatabase)
	router.Path(prefixAdminRepair).HandlerFunc(notSupported)
	router.Path(prefixAdmin).HandlerFunc(notFound)

	var prefixFeeds = prefix + "/feeds"
	router.Path(prefixFeeds).Methods(mGet).Queries("url", "{url:.+}").HandlerFunc(z.feedsGetFeedByURL)
	router.Path(prefixFeeds).Methods(mGet).HandlerFunc(z.feedsGet)
	router.Path(prefixFeeds).Methods(mPut).Headers(hContentType, mimeJSON).HandlerFunc(z.feedsSaveJSON)
	router.Path(prefixFeeds).Methods(mPut).Headers(hContentType, mimeText).HandlerFunc(z.feedsSaveText)
	router.Path(prefixFeeds).Methods(mPut).HandlerFunc(badMediaType)
	router.Path(prefixFeeds).HandlerFunc(notSupported)

	var prefixFeedsNext = prefixFeeds + "/next"
	router.Path(prefixFeedsNext).Methods(mGet).HandlerFunc(z.feedsGetFeedsNext)

	var prefixFeedsFeed = prefixFeeds + "/{feedID}"
	router.Path(prefixFeedsFeed).Methods(mGet).HandlerFunc(z.feedsGetFeedByID)
	router.Path(prefixFeedsFeed).HandlerFunc(notSupported)

	var prefixFeedsFeedLog = prefixFeeds + "/{feedID}/log"
	router.Path(prefixFeedsFeedLog).Methods(mGet).HandlerFunc(z.feedsGetFeedLogByID)
	router.Path(prefixFeedsFeedLog).HandlerFunc(notSupported)

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

type oneFS struct {
	name string
	root http.FileSystem
}

func (z oneFS) Open(name string) (http.File, error) {
	// ignore name and use z.name
	return z.root.Open(z.name)
}
