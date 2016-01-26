package httpd

import (
	"encoding/json"
	"io"
	m "rakewire/model"
)

// FeedSavesResponse response for FeedSaves
type FeedSavesResponse struct {
	Count int `json:"count"`
}

func serializeFeedSavesResponse(count int) ([]byte, error) {
	return json.Marshal(&FeedSavesResponse{
		Count: count,
	})
}

func deserializeFeedSavesResponse(r io.Reader) (int, error) {
	result := &FeedSavesResponse{}
	err := json.NewDecoder(r).Decode(result)
	if result == nil {
		return 0, nil
	}
	return result.Count, err
}

func serializeFeed(feed *m.Feed) ([]byte, error) {
	return json.Marshal(feed)
}

func deserializeFeed(data []byte) (*m.Feed, error) {
	result := m.NewFeed("")
	err := json.Unmarshal(data, result)
	return result, err
}

func serializeFeeds(feeds []*m.Feed, w io.Writer) error {
	return json.NewEncoder(w).Encode(&feeds)
}

func deserializeFeeds(r io.Reader) ([]*m.Feed, error) {
	feeds := []*m.Feed{}
	err := json.NewDecoder(r).Decode(&feeds)
	return feeds, err
}

func serializeLogs(logs []*m.Transmission, w io.Writer) error {
	return json.NewEncoder(w).Encode(&logs)
}
