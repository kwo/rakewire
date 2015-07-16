package httpd

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	m "rakewire.com/model"
	"strings"
)

// SaveFeedsResponse response for SaveFeeds
type SaveFeedsResponse struct {
	Count int `json:"count"`
}

func (z *Httpd) feedsGet(w http.ResponseWriter, req *http.Request) {

	feeds, err := z.Database.GetFeeds()
	if err != nil {
		logger.Printf("Error in m.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	err = feeds.Serialize(w)
	if err != nil {
		logger.Printf("Error in m.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Httpd) feedsGetFeedByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	feedID := vars["feedID"]

	feed, err := z.Database.GetFeedByID(feedID)
	if err != nil {
		logger.Printf("Error in m.GetFeedByID: %s\n", err.Error())
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

func (z *Httpd) feedsGetFeedByURL(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Query().Get("url")

	feed, err := z.Database.GetFeedByURL(url)
	if err != nil {
		logger.Printf("Error in m.GetFeedByURL: %s\n", err.Error())
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

func (z *Httpd) feedsSaveJSON(w http.ResponseWriter, req *http.Request) {

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

func (z *Httpd) feedsSaveText(w http.ResponseWriter, req *http.Request) {

	// curl -D - -X PUT -H "Content-Type: text/plain; charset=utf-8" --data-binary @feedlist.txt http://localhost:4444/api/feeds

	if req.ContentLength == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	feeds := m.NewFeeds()
	scanner := bufio.NewScanner(req.Body)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			feeds.Add(m.NewFeed(url))
		}
	}

	z.feedsSaveNative(w, feeds)

}

func (z *Httpd) feedsSaveNative(w http.ResponseWriter, feeds *m.Feeds) {

	n, err := z.Database.SaveFeeds(feeds)
	if err != nil {
		logger.Printf("Error in m.SaveFeeds: %s\n", err.Error())
		http.Error(w, "Cannot save feeds to database.", http.StatusInternalServerError)
		return
	}

	jsonRsp := SaveFeedsResponse{
		Count: n,
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
