package httpd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	m "rakewire.com/model"
	"time"
)

// SaveFeedsResponse response for SaveFeeds
type SaveFeedsResponse struct {
	Count int `json:"count"`
}

func (z *Service) feedsGet(w http.ResponseWriter, req *http.Request) {

	feeds, err := z.Database.GetFeeds()
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	logger.Printf("Getting feeds: %d", feeds.Size())

	w.Header().Set(hContentType, mimeJSON)
	err = feeds.Serialize(w)
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Service) feedsGetFeedsNext(w http.ResponseWriter, req *http.Request) {

	maxTime := time.Now().Truncate(time.Second).Add(48 * time.Hour)
	feeds, err := z.Database.GetFetchFeeds(&maxTime)
	if err != nil {
		logger.Printf("Error in db.GetFetchFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	logger.Printf("Getting feeds: %d", feeds.Size())

	w.Header().Set(hContentType, mimeText)
	line := fmt.Sprintf("%-8s  %-8s %6s   %s\n", "Next", "Last", "Status", "URL")
	w.Write([]byte(line))

	for _, f := range feeds.Values {
		dtNext := f.NextFetch.Local().Format("15:04:05")
		dtLast := ""
		if f.Last != nil {
			dtLast = f.Last.StartTime.Local().Format("15:04:05")
		}
		line := fmt.Sprintf("%-8s  %-8s %6d   %s\n", dtNext, dtLast, f.Last.StatusCode, f.URL)
		w.Write([]byte(line))
	}

}

func (z *Service) feedsGetFeedByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID := vars["feedID"]

	feed, err := z.Database.GetFeedByID(feedID)
	if err != nil {
		logger.Printf("Error in db.GetFeedByID: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := feed.Encode()
	if err != nil {
		logger.Printf("Error in feed.Encode: %s\n", err.Error())
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
		logger.Printf("Error in db.GetFeedByURL: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed from database.", http.StatusInternalServerError)
		return
	} else if feed == nil {
		notFound(w, req)
		return
	}

	data, err := feed.Encode()
	if err != nil {
		logger.Printf("Error in feed.Encode: %s\n", err.Error())
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

	feeds := m.NewFeeds()
	err := feeds.Deserialize(req.Body)
	if err != nil {
		logger.Printf("Error deserializing feeds: %s\n", err.Error())
		http.Error(w, "Cannot deserialize feeds.", http.StatusInternalServerError)
		return
	}

	if feeds.Size() == 0 {
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

	feeds := m.ParseListToFeeds(req.Body)
	z.feedsSaveNative(w, feeds)

}

func (z *Service) feedsSaveNative(w http.ResponseWriter, feeds *m.Feeds) {

	err := z.Database.SaveFeeds(feeds)
	if err != nil {
		logger.Printf("Error in db.SaveFeeds: %s\n", err.Error())
		http.Error(w, "Cannot save feeds to database.", http.StatusInternalServerError)
		return
	}

	jsonRsp := SaveFeedsResponse{
		Count: len(feeds.Values),
	}

	data, err := json.Marshal(jsonRsp)
	if err != nil {
		logger.Printf("Error serializing response: %s\n", err.Error())
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
		logger.Printf("Error in db.GetFeedLog: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feed logs from database.", http.StatusInternalServerError)
		return
	} else if entries == nil {
		notFound(w, req)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		logger.Printf("Error encoding log entries: %s\n", err.Error())
		http.Error(w, "Cannot serialize feed logs from database.", http.StatusInternalServerError)
		return
	}

}
