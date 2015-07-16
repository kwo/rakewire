package httpd

import (
	"github.com/gorilla/mux"
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

	return router

}
