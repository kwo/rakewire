package fever

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rakewire/db"
	m "rakewire/model"
)

// http://feedafever.com/api

const (
	// AuthParam must be sent with ever Fever request to authenticate to the service.
	AuthParam = "api_key"
)

// NewAPI creates a new Fever API instance
func NewAPI(prefix string, db db.Database) *API {
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

	log.Printf("%-7s %-7s request query: %v", logDebug, logName, req.URL.Query())
	log.Printf("%-7s %-7s request form:  %v", logDebug, logName, req.PostForm)

	if rsp.Authorized == 1 {

		for k := range req.URL.Query() {
			switch k {

			case "api":
				startTime, err := z.db.FeedLogGetLastFetchTime()
				if err == nil {
					rsp.LastRefreshed = startTime.Unix()
				} else {
					log.Printf("%-7s %-7s error retrieving last feedlog fetch time: %s", logWarn, logName, err.Error())
				}

				uMark := req.PostFormValue("mark")
				uAs := req.PostFormValue("as")
				uID := req.PostFormValue("id")
				uBefore := req.PostFormValue("before")
				if uMark != "" {
					if err := z.updateItems(user.ID, uMark, uAs, uID, uBefore); err != nil {
						log.Printf("%-7s %-7s error updating items: %s", logWarn, logName, err.Error())
					}
				}

			case "feeds":
				if feeds, feedGroups, err := z.getFeeds(user.ID); err == nil {
					rsp.Feeds = feeds
					rsp.FeedGroups = feedGroups
				} else {
					log.Printf("%-7s %-7s error retrieving feeds and feed_groups: %s", logWarn, logName, err.Error())
				}

			case "groups":
				if groups, feedGroups, err := z.getGroups(user.ID); err == nil {
					rsp.Groups = groups
					rsp.FeedGroups = feedGroups
				} else {
					log.Printf("%-7s %-7s error retrieving groups and feed_groups: %s", logWarn, logName, err.Error())
				}

			case "items":
				if count, err := z.db.UserEntryGetTotalForUser(user.ID); err == nil {
					rsp.ItemCount = count
				} else {
					log.Printf("%-7s %-7s error retrieving item count: %s", logWarn, logName, err.Error())
				}
				if id := parseID(req.URL.Query(), "since_id"); id > 0 {
					items, err := z.getItemsNext(user.ID, id)
					if err == nil {
						rsp.Items = items
					} else {
						log.Printf("%-7s %-7s error retrieving items: %s", logWarn, logName, err.Error())
					}
				} else if id := parseID(req.URL.Query(), "max_id"); id > 0 {
					items, err := z.getItemsPrev(user.ID, id)
					if err == nil {
						rsp.Items = items
					} else {
						log.Printf("%-7s %-7s error retrieving items: %s", logWarn, logName, err.Error())
					}
				} else if ids := parseIDArray(req.URL.Query(), "with_ids"); ids != nil {
					items, err := z.getItemsByIds(user.ID, ids)
					if err == nil {
						rsp.Items = items
					} else {
						log.Printf("%-7s %-7s error retrieving items: %s", logWarn, logName, err.Error())
					}
				} else {
					items, err := z.getItemsAll(user.ID)
					if err == nil {
						rsp.Items = items
					} else {
						log.Printf("%-7s %-7s error retrieving items all: %s", logWarn, logName, err.Error())
					}
				}

			case "unread_item_ids":
				if itemIDs, err := z.getUnreadItemIDs(user.ID); err == nil {
					rsp.UnreadItemIDs = itemIDs
				} else {
					log.Printf("%-7s %-7s error retrieving unread item IDs: %s", logWarn, logName, err.Error())
				}

			case "saved_item_ids":
				if itemIDs, err := z.getSavedItemIDs(user.ID); err == nil {
					rsp.SavedItemIDs = itemIDs
				} else {
					log.Printf("%-7s %-7s error retrieving saved item IDs: %s", logWarn, logName, err.Error())
				}

			}
		}

	} // authorized

	w.WriteHeader(http.StatusOK)
	w.Header().Set(hContentType, mimeJSON)
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		log.Printf("%-7s %-7s cannot serialize fever JSON response: %s", logWarn, logName, err.Error())
	}

}
