package fever

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	//m "rakewire/model"
	"strings"
	"time"
)

// http://feedafever.com/api

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

	useXML := req.URL.Query().Get("api") == "xml"

	// TODO: omit last_refreshed_on_time when unauthorized

	//var user *m.User
	rsp := &Response{
		Version: 3,
	}

	if apiKey := req.PostFormValue("api_key"); apiKey != "" {
		if u, err := z.db.UserGetByFeverHash(apiKey); err == nil && u != nil {
			rsp.Authorized = 1
			log.Printf("%-7s %-7s authorized: %s", logDebug, logName, u.Username)
			//user = u
		}
	}

	log.Printf("%-7s %-7s response: %v", logDebug, logName, rsp)

	if rsp.Authorized == 1 {

		for k := range req.URL.Query() {

			switch k {
			case "api":
				rsp.LastRefreshed = time.Now().Unix() // TODO: get last refreshed
			case "feeds":
				// add to response
			}

		} // loop

	}

	if useXML {
		w.Header().Set(hContentType, mimeXML)
		w.Write([]byte(strings.ToLower(strings.TrimSuffix(xml.Header, "\n")))) // lowercase without trailing newline
		if err := xml.NewEncoder(w).Encode(&rsp); err != nil {
			log.Printf("%-7s %-7s cannot serialize fever XML response: %s", logWarn, logName, err.Error())
		}
	} else {
		w.Header().Set(hContentType, mimeJSON)
		if err := json.NewEncoder(w).Encode(&rsp); err != nil {
			log.Printf("%-7s %-7s cannot serialize fever JSON response: %s", logWarn, logName, err.Error())
		}
	}

}
