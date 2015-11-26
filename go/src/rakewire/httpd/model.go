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

func serializeSaveFeedsResponse(count int) ([]byte, error) {
	return json.Marshal(SaveFeedsResponse{
		Count: count,
	})
}

func deserializeSaveFeedsResponse(r io.Reader) (int, error) {
	var result *SaveFeedsResponse
	err := json.NewDecoder(r).Decode(result)
	if result == nil {
		return 0, nil
	}
	return result.Count, err
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

func serializeLogs(logs []*m.FeedLog, w io.Writer) error {
	return json.NewEncoder(w).Encode(&logs)
}
