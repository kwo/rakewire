package httpd

import (
	"github.com/gorilla/mux"
	"net/http"
	m "rakewire/model"
	"time"
)

func (z *Service) feedsGet(w http.ResponseWriter, req *http.Request) {

	feeds, err := z.Database.GetFeeds()
	if err != nil {
		logger.Warnf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	logger.Infof("Getting feeds: %d", len(feeds))

	w.Header().Set(hContentType, mimeJSON)
	err = serializeFeeds(feeds, w)
	if err != nil {
		logger.Warnf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Service) feedsGetFeedsNext(w http.ResponseWriter, req *http.Request) {

	maxTime := time.Now().Truncate(time.Second).Add(36 * time.Hour)
	feeds, err := z.Database.GetFetchFeeds(maxTime)
	if err != nil {
		logger.Warnf("Error in db.GetFetchFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	logger.Infof("Getting feeds: %d", len(feeds))

	w.Header().Set(hContentType, mimeJSON)
	err = serializeFeeds(feeds, w)
	if err != nil {
		logger.Warnf("Error in db.GetFeedsNext: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Service) feedsGetFeedByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID := vars["feedID"]

	feed, err := z.Database.GetFeedByID(feedID)
	if err != nil {
		logger.Warnf("Error in db.GetFeedByID: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := serializeFeed(feed)
	if err != nil {
		logger.Warnf("Error in feed.Encode: %s\n", err.Error())
		http.Error(w, "Cannot serialize feed from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	w.Write(data)
	w.Write([]byte("\n"))

}

func (z *Service) feedsGetFeedByURL(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Query().Get("url")

	feed, err := z.Database.GetFeedByURL(url)
	if err != nil {
		logger.Warnf("Error in db.GetFeedByURL: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := serializeFeed(feed)
	if err != nil {
		logger.Warnf("Error in feed.Encode: %s\n", err.Error())
		http.Error(w, "Cannot serialize feed from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	w.Write(data)
	w.Write([]byte("\n"))

}

func (z *Service) feedsSaveJSON(w http.ResponseWriter, req *http.Request) {

	if req.ContentLength == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	feeds, err := deserializeFeeds(req.Body)
	if err != nil {
		logger.Warnf("Error deserializing feeds: %s\n", err.Error())
		http.Error(w, "Cannot deserialize feeds.", http.StatusInternalServerError)
		return
	}

	if len(feeds) == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	z.feedsSaveNative(w, feeds)

}

func (z *Service) feedsSaveText(w http.ResponseWriter, req *http.Request) {

	// curl -D - -X PUT -H "Content-Type: text/plain; charset=utf-8" --data-binary @feedlist.txt http://localhost:4444/api/feeds

	if req.ContentLength == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	feeds := m.ParseFeedsFromReader(req.Body)
	z.feedsSaveNative(w, feeds)

}

func (z *Service) feedsSaveNative(w http.ResponseWriter, feeds []*m.Feed) {

	for _, feed := range feeds {
		err := z.Database.SaveFeed(feed)
		if err != nil {
			logger.Warnf("Error in db.SaveFeed: %s\n", err.Error())
			http.Error(w, "Cannot save feed to database.", http.StatusInternalServerError)
			return
		}
	}

	data, err := serializeSaveFeedsResponse(len(feeds))
	if err != nil {
		logger.Warnf("Error serializing response: %s\n", err.Error())
		http.Error(w, "Cannot serialize response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	w.Write(data)

}

func (z *Service) feedsGetFeedLogByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID := vars["feedID"]

	entries, err := z.Database.GetFeedLog(feedID, 7*24*time.Hour)
	if err != nil {
		logger.Warnf("Error in db.GetFeedLog: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed logs from database.", http.StatusInternalServerError)
		return
	} else if entries == nil {
		notFound(w, req)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	err = serializeLogs(entries, w)
	if err != nil {
		logger.Warnf("Error encoding log entries: %s\n", err.Error())
		http.Error(w, "Cannot serialize feed logs from database.", http.StatusInternalServerError)
		return
	}

}
