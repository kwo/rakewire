package httpd

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	m "rakewire/model"
	"strconv"
	"time"
)

func (z *Service) feedsGet(w http.ResponseWriter, req *http.Request) {

	feeds, err := z.database.GetFeeds()
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeeds: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	log.Printf("%-7s %-7s Getting feeds: %d", logDebug, logName, len(feeds))

	w.Header().Set(hContentType, mimeJSON)
	err = serializeFeeds(feeds, w)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeeds: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Service) feedsGetFeedsNext(w http.ResponseWriter, req *http.Request) {

	maxTime := time.Now().Truncate(time.Second).Add(36 * time.Hour)
	feeds, err := z.database.GetFetchFeeds(maxTime)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFetchFeeds: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	log.Printf("%-7s %-7s Getting feeds: %d", logDebug, logName, len(feeds))

	w.Header().Set(hContentType, mimeJSON)
	err = serializeFeeds(feeds, w)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeedsNext: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Service) feedsGetFeedByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID, err := strconv.ParseUint(vars["feedID"], 10, 64)
	if err != nil {
		notFound(w, req)
		return
	}

	feed, err := z.database.GetFeedByID(feedID)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeedByID: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := serializeFeed(feed)
	if err != nil {
		log.Printf("%-7s %-7s Error in feed.Encode: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot serialize feed from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	w.Write(data)
	w.Write([]byte("\n"))

}

func (z *Service) feedsGetFeedByURL(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Query().Get("url")

	feed, err := z.database.GetFeedByURL(url)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeedByURL: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := serializeFeed(feed)
	if err != nil {
		log.Printf("%-7s %-7s Error in feed.Encode: %s", logWarn, logName, err.Error())
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
		log.Printf("%-7s %-7s Error deserializing feeds: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot deserialize feeds.", http.StatusInternalServerError)
		return
	}

	if len(feeds) == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	z.feedsSaveNative(w, feeds)

}

func (z *Service) feedsSaveNative(w http.ResponseWriter, feeds []*m.Feed) {

	for _, feed := range feeds {
		_, err := z.database.FeedSave(feed)
		if err != nil {
			log.Printf("%-7s %-7s Error in db.FeedSave: %s", logWarn, logName, err.Error())
			http.Error(w, "Cannot save feed to database.", http.StatusInternalServerError)
			return
		}
	}

	data, err := serializeFeedSavesResponse(len(feeds))
	if err != nil {
		log.Printf("%-7s %-7s Error serializing response: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot serialize response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	w.Write(data)

}

// feedsGetFeedLogByID get feed logs for the last 90 days
func (z *Service) feedsGetFeedLogByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID, err := strconv.ParseUint(vars["feedID"], 10, 64)
	if err != nil {
		badRequest(w, req)
		return
	}

	entries, err := z.database.GetFeedLog(feedID, 90*24*time.Hour)
	if err != nil {
		log.Printf("%-7s %-7s Error in db.GetFeedLog: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot retrieve feed logs from database.", http.StatusInternalServerError)
		return
	} else if entries == nil {
		notFound(w, req)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	err = serializeLogs(entries, w)
	if err != nil {
		log.Printf("%-7s %-7s Error encoding log entries: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot serialize feed logs from database.", http.StatusInternalServerError)
		return
	}

}
