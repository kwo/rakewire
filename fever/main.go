package fever

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"rakewire/logger"
	"rakewire/model"
	"strings"
)

// http://feedafever.com/api

const (
	// AuthParam must be sent with ever Fever request to authenticate to the service.
	AuthParam = "api_key"
)

var log = logger.New("fever")

// API top level struct
type API struct {
	db model.Database
}

// New creates a new Fever API instance
func New(db model.Database) *API {
	return &API{
		db: db,
	}
}

// ServeHTTPC ist the context-aware http.Handler for Fever.
func (z *API) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "cannot parse request\n", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("api") == "xml" {
		http.Error(w, "xml api not supported\n", http.StatusBadRequest)
		return
	}

	rsp := &Response{
		Version: 3,
	}

	z.db.Update(func(tx model.Transaction) error {

		var user *model.User
		if apiKey := r.PostFormValue(AuthParam); apiKey != "" {
			z.db.Select(func(tx model.Transaction) error {
				u := model.U.GetByFeverhash(tx, apiKey)
				if u != nil {
					rsp.Authorized = 1
					log.Debugf("authorized: %s", u.Username)
					user = u
				}
				return nil
			})
		}

		log.Debugf("request query: %v", r.URL.Query())
		log.Debugf("request form:  %v", r.PostForm)

		if rsp.Authorized == 1 {

			for k := range r.URL.Query() {
				switch k {

				case "api":
					tr := model.T.GetLast(tx)
					if tr != nil {
						rsp.LastRefreshed = tr.StartTime.Unix()
					} else {
						log.Debugf("error retrieving last transmission fetch time: %s", "not found")
					}

					uMark := r.PostFormValue("mark")
					uAs := r.PostFormValue("as")
					uID := r.PostFormValue("id")
					uBefore := r.PostFormValue("before")
					if uMark != "" {
						if err := z.updateItems(user.ID, uMark, uAs, uID, uBefore, tx); err != nil {
							log.Debugf("error updating items: %s", err.Error())
						}
					}

				case "feeds":
					if feeds, feedGroups, err := z.getFeeds(user.ID, tx); err == nil {
						rsp.Feeds = feeds
						rsp.FeedGroups = feedGroups
					} else {
						log.Debugf("error retrieving feeds and feed_groups: %s", err.Error())
					}

				case "groups":
					if groups, feedGroups, err := z.getGroups(user.ID, tx); err == nil {
						rsp.Groups = groups
						rsp.FeedGroups = feedGroups
					} else {
						log.Debugf("error retrieving groups and feed_groups: %s", err.Error())
					}

				case "items":
					rsp.ItemCount = model.E.Query(tx, user.ID).Count()
					if id := r.URL.Query().Get("since_id"); len(id) > 0 {
						items, err := z.getItemsNext(user.ID, id, tx)
						if err == nil {
							rsp.Items = items
						} else {
							log.Debugf("error retrieving items: %s", err.Error())
						}
					} else if id := r.URL.Query().Get("max_id"); len(id) > 0 {
						items, err := z.getItemsPrev(user.ID, id, tx)
						if err == nil {
							rsp.Items = items
						} else {
							log.Debugf("error retrieving items: %s", err.Error())
						}
					} else if ids := r.URL.Query().Get("with_ids"); len(ids) > 0 {
						idArray := strings.Split(ids, ",")
						items, err := z.getItemsByIds(user.ID, idArray, tx)
						if err == nil {
							rsp.Items = items
						} else {
							log.Debugf("error retrieving items: %s", err.Error())
						}
					} else {
						items, err := z.getItemsAll(user.ID, tx)
						if err == nil {
							rsp.Items = items
						} else {
							log.Debugf("error retrieving items all: %s", err.Error())
						}
					}

				case "unread_item_ids":
					if itemIDs, err := z.getUnreadItemIDs(user.ID, tx); err == nil {
						rsp.UnreadItemIDs = itemIDs
					} else {
						log.Debugf("error retrieving unread item IDs: %s", err.Error())
					}

				case "saved_item_ids":
					if itemIDs, err := z.getSavedItemIDs(user.ID, tx); err == nil {
						rsp.SavedItemIDs = itemIDs
					} else {
						log.Debugf("error retrieving saved item IDs: %s", err.Error())
					}

				}
			}

		} // authorized

		return nil

	})

	w.Header().Set(hContentType, mimeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		log.Debugf("cannot serialize fever JSON response: %s", err.Error())
	}

}
