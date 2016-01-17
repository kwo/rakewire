package httpd

//go:generate esc -o static.go -pkg httpd -prefix $PROJECT_HOME/web $PROJECT_HOME/web/public

import (
	"github.com/gorilla/mux"
	"net/http"
	"rakewire/httpd/fever"
	"rakewire/httpd/rest"
	"rakewire/middleware"
	"rakewire/model"
)

func (z *Service) mainRouter(useLocal, useLegacy bool) (*mux.Router, error) {

	router := mux.NewRouter()

	if useLegacy {

		// api router
		apiPrefix := "/api0"
		router.PathPrefix(apiPrefix).Handler(
			z.apiRouter(apiPrefix),
		)

	}

	// rest api router
	restPrefix := "/api"
	restAPI := rest.NewAPI(restPrefix, z.database)
	router.PathPrefix(restPrefix).Handler(
		middleware.Adapt(restAPI.Router(), middleware.BasicAuth(&middleware.BasicAuthOptions{
			Database: model.NewBoltDatabase(z.database.BoltDB()), Realm: "Rakewire",
		})),
	)

	// fever api router
	feverPrefix := "/fever/"
	feverAPI := fever.NewAPI(feverPrefix, z.database)
	router.PathPrefix(feverPrefix).Handler(
		feverAPI.Router(),
	)

	if useLegacy {

		fs := Dir(useLocal, "/public")
		ofs := oneFS{name: "/index.html", root: fs}

		// HTML5 routes: any path without a dot (thus an extension)
		router.Path("/{route:[a-z0-9/-]+}").Handler(
			http.FileServer(ofs),
		)

		// always redirect /index.html to /
		router.Path("/index.html").Handler(
			http.RedirectHandler("/", http.StatusMovedPermanently),
		)

		// static web site
		router.PathPrefix("/").Handler(
			http.FileServer(fs),
		)

	}

	return router, nil

}

func (z *Service) apiRouter(prefix string) *mux.Router {

	router := mux.NewRouter()

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

func badRequest(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusBadRequest)
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
