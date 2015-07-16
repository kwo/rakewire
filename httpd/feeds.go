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

	w.Header().Set("Content-Type", "application/json")
	err = feeds.Serialize(w)
	if err != nil {
		logger.Printf("Error in db.GetFeeds: %s\n", err.Error())
		http.Error(w, "Cannot serialize feeds from database.", http.StatusInternalServerError)
		return
	}

}

func (z *Httpd) feedsSave(w http.ResponseWriter, req *http.Request) {

	f := db.NewFeeds()
	f.Deserialize(req.Body)

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
