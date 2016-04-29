package model

import (
	"encoding/json"
	"time"
)

const (
	entityTransmission        = "Transmission"
	indexTransmissionTime     = "Time"
	indexTransmissionFeedTime = "FeedTime"
)

var (
	indexesTransmission = []string{
		indexTransmissionTime, indexTransmissionFeedTime,
	}
)

// FetchResults
const (
	FetchResultOK          = "OK"
	FetchResultRedirect    = "MV" // message contains old URL -> new URL
	FetchResultClientError = "EC" // message contains error text
	FetchResultServerError = "ES" // check http status code
	FetchResultFeedError   = "FP" // cannot parse feed
)

// Transmission represents an attempted HTTP request to a feed
type Transmission struct {
	ID            string        `json:"id"`
	FeedID        string        `json:"feedId"`
	Duration      time.Duration `json:"duration"`
	Result        string        `json:"result,omitempty"`
	ResultMessage string        `json:"resultMessage,omitempty"`
	StartTime     time.Time     `json:"startTime"`
	URL           string        `json:"url"`
	ContentLength int           `json:"contentLength,omitempty"`
	ContentType   string        `json:"contentType,omitempty"`
	ETag          string        `json:"etag,omitempty"`
	LastModified  time.Time     `json:"lastModified,omitempty"`
	StatusCode    int           `json:"statusCode,omitempty"`
	UsesGzip      bool          `json:"gzip,omitempty"`
	Flavor        string        `json:"flavor,omitempty"`
	Generator     string        `json:"generator,omitempty"`
	Title         string        `json:"title,omitempty"`
	LastUpdated   time.Time     `json:"lastUpdated,omitempty"`
	ItemCount     int           `json:"itemCount,omitempty"`
	NewItems      int           `json:"newItems,omitempty"`
}

// GetID returns the unique ID for the object
func (z *Transmission) GetID() string {
	return z.ID
}

func (z *Transmission) clear() {
	z.ID = empty
	z.FeedID = empty
	z.Duration = 0
	z.Result = empty
	z.ResultMessage = empty
	z.StartTime = time.Time{}
	z.URL = empty
	z.ContentLength = 0
	z.ContentType = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.StatusCode = 0
	z.UsesGzip = false
	z.Flavor = empty
	z.Generator = empty
	z.Title = empty
	z.LastUpdated = time.Time{}
	z.ItemCount = 0
	z.NewItems = 0
}

func (z *Transmission) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Transmission) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Transmission) hasIncrementingID() bool {
	return true
}

func (z *Transmission) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexTransmissionTime] = []string{keyEncodeTime(z.StartTime), z.ID}
	result[indexTransmissionFeedTime] = []string{z.FeedID, keyEncodeTime(z.StartTime)}
	return result
}

func (z *Transmission) setID(tx Transaction) error {
	id, err := tx.NextID(entityTransmission)
	if err != nil {
		return err
	}
	z.ID = keyEncodeUint(id)
	return nil
}

// Transmissions is a collection Transmission objects
type Transmissions []*Transmission

// Reverse reverses the order of the collection
func (z Transmissions) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

func (z *Transmissions) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Transmissions) encode() ([]byte, error) {
	return json.Marshal(z)
}
