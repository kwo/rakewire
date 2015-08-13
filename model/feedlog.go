package model

import (
	"encoding/json"
	"time"
)

// FetchResults
const (
	FetchResultOK          = "OK"
	FetchResultRedirect    = "MV" // message contains old URL -> new URL
	FetchResultClientError = "EC" // message contains error text
	FetchResultServerError = "ES" // check http status code
	FetchResultFeedError   = "EF" // cannot parse feed
)

// Check Levels for Update Status
const (
	UpdateCheck304         = "NM" // HTTP Status Code 304
	UpdateCheckFeedEntries = "LU" // No 304  but feed LastUpdated is the same
)

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	ContentLength int           `json:"contentLength,omitempty"` // Length of the HTTP Response
	Duration      time.Duration `json:"duration"`                // duration of http request plus feed processing
	ETag          string        `json:"etag,omitempty"`          // ETag from HTTP Request
	IsUpdated     bool          `json:"updated"`                 // flag indicating updated content, see UpdateCheck
	LastModified  *time.Time    `json:"lastModified,omitempty"`  // Last-Modified time from HTTP Request
	Result        string        `json:"result"`                  // result code of fetch attempt
	ResultMessage string        `json:"resultMessage,omitempty"` // optional error message for result
	StartTime     *time.Time    `json:"startTime"`               // The time the fetch process started
	StatusCode    int           `json:"statusCode,omitempty"`    // HTTP status code
	UpdateCheck   string        `json:"updateCheck,omitempty"`   // 304, Checksum or Feed inspection
	UsesGzip      bool          `json:"gzip,omitempty"`          // if the remote http server serves the feed compressed
}

// Decode FeedLog object from bytes
func (z *FeedLog) Decode(data []byte) error {
	return json.Unmarshal(data, z)
}

// Encode FeedLog object to bytes
func (z *FeedLog) Encode() ([]byte, error) {
	return json.MarshalIndent(z, "", " ")
}
