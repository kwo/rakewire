package fever

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	m "rakewire/model"
	"time"
)

// http://feedafever.com/api

const (
	// AuthParam must be sent with ever Fever request to authenticate to the service.
	AuthParam = "api_key"
)

// NewAPI creates a new Fever API instance
func NewAPI(prefix string, db Database) *API {
	return &API{
		prefix: prefix,
		db:     db,
	}
}

// Router returns the top-level Fever router
func (z *API) Router() *mux.Router {

	router := mux.NewRouter()

	router.Path(z.prefix).Queries("api", "").Methods(mPost).HandlerFunc(z.mux)
	router.Path(z.prefix).Queries("api", "").HandlerFunc(notSupported)
	router.Path(z.prefix).HandlerFunc(notFound)

	return router

}

func (z *API) mux(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, "cannot parse request\n", 400)
		return
	}

	if req.URL.Query().Get("api") == "xml" {
		http.Error(w, "xml api not supported\n", 400)
		return
	}

	rsp := &Response{
		Version: 3,
	}

	var user *m.User
	if apiKey := req.PostFormValue(AuthParam); apiKey != "" {
		if u, err := z.db.UserGetByFeverHash(apiKey); err == nil && u != nil {
			rsp.Authorized = 1
			log.Printf("%-7s %-7s authorized: %s", logDebug, logName, u.Username)
			user = u
		}
	}

	log.Printf("%-7s %-7s response: %v", logDebug, logName, rsp)

	if rsp.Authorized == 1 {

	Loop:
		for k := range req.URL.Query() {
			switch k {
			case "api":
				rsp.LastRefreshed = time.Now().Unix() // TODO: get last refreshed feed for user; need groups first
			case "groups":
				if groups, feedGroups, err := z.getGroupsAndFeedGroups(user.ID); err == nil {
					rsp.Groups = groups
					rsp.FeedGroups = feedGroups
				} else {
					log.Printf("%-7s %-7s error retrieving feeds and groups: %s", logWarn, logName, err.Error())
				}
				break Loop
			}
		} // loop

	} // authorized

	w.Header().Set(hContentType, mimeJSON)
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		log.Printf("%-7s %-7s cannot serialize fever JSON response: %s", logWarn, logName, err.Error())
	}

}
