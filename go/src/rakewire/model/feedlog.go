package model

import (
	"encoding/json"
	"time"
)

// FetchResults
const (
	FetchResultOK            = "OK"
	FetchResultRedirect      = "MV" // message contains old URL -> new URL
	FetchResultClientError   = "EC" // message contains error text
	FetchResultServerError   = "ES" // check http status code
	FetchResultFeedError     = "FP" // cannot parse feed
	FetchResultFeedTimeError = "FT" // cannot parse time
)

// Check Levels for Update Status
const (
	UpdateCheck304  = "NM" // HTTP Status Code 304
	UpdateCheckFeed = "LU" // No 304  but feed LastUpdated is the same
)

// NewFeedLog create new FeedLog
func NewFeedLog() *FeedLog {
	result := &FeedLog{}
	result.HTTP = HTTPInfo{}
	result.Feed = FeedInfo{}
	return result
}

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	Duration      time.Duration `json:"duration"`                // duration of http request plus feed processing
	IsUpdated     bool          `json:"updated"`                 // flag indicating updated content, see UpdateCheck
	Result        string        `json:"result"`                  // result code of fetch attempt
	ResultMessage string        `json:"resultMessage,omitempty"` // optional error message for result
	UpdateCheck   string        `json:"updateCheck,omitempty"`   // 304, Checksum or Feed inspection
	StartTime     time.Time     `json:"startTime"`               // The time the fetch process started
	URL           string        `json:"url"`
	Feed          FeedInfo      `json:"feed"`
	HTTP          HTTPInfo      `json:"http"`
}

// HTTPInfo log information about the HTTP request
type HTTPInfo struct {
	ContentLength int       `json:"contentLength,omitempty"` // Length of the HTTP Response
	ContentType   string    `json:"contentType,omitempty"`   // Content-Type header
	ETag          string    `json:"etag,omitempty"`          // ETag from HTTP Request
	LastModified  time.Time `json:"lastModified,omitempty"`  // Last-Modified time from HTTP Request
	StatusCode    int       `json:"statusCode,omitempty"`    // HTTP status code
	UsesGzip      bool      `json:"gzip,omitempty"`          // if the remote http server serves the feed compressed
}

// FeedInfo log information about the feed
type FeedInfo struct {
	Flavor    string    `json:"flavor"`
	Generator string    `json:"generator,omitempty"`
	Title     string    `json:"title"`
	Updated   time.Time `json:"updated"`
}

// Decode FeedLog object from bytes
func (z *FeedLog) Decode(data []byte) error {
	return json.Unmarshal(data, z)
}

// Encode FeedLog object to bytes
func (z *FeedLog) Encode() ([]byte, error) {
	return json.MarshalIndent(z, "", " ")
}
