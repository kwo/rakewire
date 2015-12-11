package httpd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (z *Service) feverRouter(prefix string) *mux.Router {

	router := mux.NewRouter()

	router.Queries("api", "").Methods(mPost).HandlerFunc(z.feverMux)
	router.Queries("api", "").HandlerFunc(notSupported)
	router.Path(prefix).HandlerFunc(notFound)

	return router

}

func (z *Service) feverMux(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, "cannot parse request\n", 400)
		return
	}

	// data includes api_key set to md5(email:password)
	// TODO: return json response with auth: 0
	if apiKey := req.PostFormValue("api_key"); apiKey == "" {
		http.Error(w, "no api_key given\n", 400)
		return
		// } else {
		// TODO: authenticate
	}

	// "api" key is guarenteed to be in form values
	nValues := len(req.URL.Query())
	if nValues == 1 {
		z.feverAuth(w, req)
		return
	} else if nValues > 2 {
		// TODO: ignore this error
		http.Error(w, "too many arguments, only one in addition to 'api' is allowed\n", 400)
		return
	}

	for k := range req.URL.Query() {

		switch k {
		case "api":
			continue
		case "feeds":
			z.feverFeeds(w, req)
			return
		default:
			http.Error(w, fmt.Sprintf("unrecognized argument: %s\n", k), 400)
			return
		}

	} // loop

}

// http://feedafever.com/api

func (z *Service) feverAuth(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "auth\n", 200)
}

func (z *Service) feverFeeds(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "feeds\n", 200)
}
