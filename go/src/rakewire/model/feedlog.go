package model

import (
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

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	ID     string `db:"primary-key"`
	FeedID string `db:"indexFeedTime:1"`

	Duration      time.Duration // duration of http request plus feed processing
	IsUpdated     bool          // flag indicating updated content, see UpdateCheck
	Result        string        // result code of fetch attempt, see FetchResults
	ResultMessage string        // optional error message for result
	UpdateCheck   string        // 304, Checksum or Feed inspection
	StartTime     time.Time     `db:"indexFeedTime:2"` // The time the fetch process started
	URL           string        // URL of the feed

	ContentLength int       // Length of the HTTP Response
	ContentType   string    // HTTP Content-Type header
	ETag          string    // ETag from HTTP Request
	LastModified  time.Time // Last-Modified time from HTTP Request
	StatusCode    int       // HTTP status code
	UsesGzip      bool      // if the remote http server serves the feed compressed

	Flavor    string    // Feed type: RSSx, Atom or RDF
	Generator string    // Feed generator (for example, Wordpress)
	Title     string    // Title as specified by the feed
	Updated   time.Time // most recent entry contained in the feed, NOT the feed updated time
}
