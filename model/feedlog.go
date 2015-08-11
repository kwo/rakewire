package model

import (
	"encoding/json"
	"time"
)

// FetchResults
const (
	FetchResultOK          = 0
	FetchResultRedirect    = 10 // message contains old URL -> new URL
	FetchResultClientError = 20 // message contains error text
	FetchResultServerError = 30 // check http status code
)

// Check Levels for Update Status
const (
	UpdateCheck304         = 10 // HTTP Status Code 304
	UpdateCheckChecksum    = 20 // No 304 but checksums are the same
	UpdateCheckFeedEntries = 30 // Differing checksums but feed entries are all the same
)

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	Checksum      string        `json:"checksum,omitempty"`      // Checksum of HTTP payload (independent of etag)
	ContentLength int           `json:"contentLength,omitempty"` // Length of the HTTP Response
	Duration      time.Duration `json:"duration"`                // duration of http request plus feed processing
	ETag          string        `json:"etag,omitempty"`          // ETag from HTTP Request
	FeedID        string        `json:"feedId"`                  // UUID of feed
	IsUpdated     bool          `json:"updated"`                 // flag indicating updated content, see UpdateCheck
	LastModified  *time.Time    `json:"lastModified,omitempty"`  // Last-Modified time from HTTP Request
	Result        int           `json:"result"`                  // result code of fetch attempt
	ResultMessage string        `json:"resultMessage,omitempty"` // optional error message for result
	StartTime     *time.Time    `json:"startTime"`               // The time the fetch process started
	StatusCode    int           `json:"statusCode,omitempty"`    // HTTP status code
	UpdateCheck   int           `json:"updateCheck,omitempty"`   // 304, Checksum or Feed inspection
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
