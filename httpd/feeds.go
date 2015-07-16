package httpd

import (
	"encoding/json"
	"net/http"
	"rakewire.com/db"
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

	w.Header().Set("Content-Type", mimeJSON)
	err = feeds.Serialize(w)
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Httpd) feedsSave(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")

	if contentType == mimeJSON {
		z.feedsSaveJSON(w, req)
	} else if contentType == mimeText {
		z.feedsSaveText(w, req)
	} else {
		sendError(w, http.StatusUnsupportedMediaType)
		return
	}

}

func (z *Httpd) feedsSaveJSON(w http.ResponseWriter, req *http.Request) {

	f := db.NewFeeds()
	err := f.Deserialize(req.Body)
	if err != nil {
		logger.Printf("Error deserializing feeds: %s\n", err.Error())
		http.Error(w, "Cannot deserialize feeds.", http.StatusInternalServerError)
		return
	}

	if f.Size() == 0 {
		sendError(w, http.StatusNoContent)
		return
	}

	l, err := z.Database.SaveFeeds(f)
	if err != nil {
		logger.Printf("Error in db.SaveFeeds: %s\n", err.Error())
		http.Error(w, "Cannot save feeds to database.", http.StatusInternalServerError)
		return
	}

	jsonRsp := SaveFeedsResponse{
		Count: l,
	}

	data, err := json.Marshal(jsonRsp)
	if err != nil {
		logger.Printf("Error serializing response: %s\n", err.Error())
		http.Error(w, "Cannot serialize response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimeJSON)
	w.Write(data)

}

func (z *Httpd) feedsSaveText(w http.ResponseWriter, req *http.Request) {

	// curl -D - -X PUT -H "Content-Type: text/plain; charset=utf-8" --data-binary @feedlist.txt http://localhost:4444/api/feeds

	sendError(w, http.StatusNotImplemented)

}
