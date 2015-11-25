package httpd

import (
	"encoding/json"
	"io"
	m "rakewire/model"
)

// SaveFeedsResponse response for SaveFeeds
type SaveFeedsResponse struct {
	Count int `json:"count"`
}

func serializeFeed(feed *m.Feed) ([]byte, error) {
	return nil, nil
}

func deserializeFeed(data []byte) (*m.Feed, error) {
	return nil, nil
}

func serializeFeeds(feeds []*m.Feed, w io.Writer) error {
	return json.NewEncoder(w).Encode(&feeds)
}

func deserializeFeeds(r io.Reader) ([]*m.Feed, error) {
	var feeds []*m.Feed
	err := json.NewDecoder(r).Decode(&feeds)
	return feeds, err
}
