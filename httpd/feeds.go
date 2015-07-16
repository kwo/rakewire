package httpd

import (
	"bufio"
	"encoding/json"
	"net/http"
	"rakewire.com/db"
	"strings"
)

// SaveFeedsResponse response for SaveFeeds
type SaveFeedsResponse struct {
	Count int `json:"count"`
}

func (z *Httpd) feedsGet(w http.ResponseWriter, req *http.Request) {

	feeds, err := z.Database.GetFeeds()
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot retrieve feeds from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, mimeJSON)
	err = feeds.Serialize(w)
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Httpd) feedsSave(w http.ResponseWriter, req *http.Request) {

	if req.ContentLength == 0 {
		sendError(w, http.StatusNoContent)
	} else {
		contentType := req.Header.Get(hContentType)
		if contentType == mimeJSON {
			z.feedsSaveJSON(w, req)
		} else if contentType == mimeText {
			z.feedsSaveText(w, req)
		} else {
			sendError(w, http.StatusUnsupportedMediaType)
		}
	}

}

func (z *Httpd) feedsSaveJSON(w http.ResponseWriter, req *http.Request) {

	feeds := db.NewFeeds()
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

	feeds := db.NewFeeds()
	scanner := bufio.NewScanner(req.Body)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			feeds.Add(db.NewFeed(url))
		}
	}

	z.feedsSaveNative(w, feeds)

}

func (z *Httpd) feedsSaveNative(w http.ResponseWriter, feeds *db.Feeds) {

	n, err := z.Database.SaveFeeds(feeds)
	if err != nil {
		logger.Printf("Error in db.SaveFeeds: %s\n", err.Error())
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
