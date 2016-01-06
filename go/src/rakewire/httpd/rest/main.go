package rest

import (
	"github.com/gorilla/mux"
)

// NewAPI creates a new REST API instance
func NewAPI(prefix string, db Database) *API {
	return &API{
		prefix: prefix,
		db:     db,
	}
}

// Router returns the top-level Fever router
func (z *API) Router() *mux.Router {

	router := mux.NewRouter()

	prefixUsers := "/users"
	router.Path(z.prefix + prefixUsers + "/{username}").Methods(mGet).HandlerFunc(z.usersGet)
	router.Path(z.prefix + prefixUsers + "/{username}").HandlerFunc(notSupported)
	router.Path(z.prefix + prefixUsers).HandlerFunc(notFound)

	router.Path(z.prefix).HandlerFunc(notFound)

	return router

}
