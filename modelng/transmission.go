package modelng

import (
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

// Transmissions is a collection Transmission objects
type Transmissions []*Transmission

// Reverse reverses the order of the collection
func (z Transmissions) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

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
