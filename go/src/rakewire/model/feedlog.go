package model

import (
	"github.com/pborman/uuid"
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

// NewFeedLog instantiates a new FeedLog with the required fields set.
func NewFeedLog(feedID string) *FeedLog {
	return &FeedLog{
		ID:     uuid.NewUUID().String(),
		FeedID: feedID,
	}
}

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	ID            string        `json:"id" kv:"key"`
	FeedID        string        `json:"feedId" kv:"FeedTime:1"`
	Duration      time.Duration `json:"duration"`
	IsUpdated     bool          `json:"updated"`
	Result        string        `json:"result"`
	ResultMessage string        `json:"resultMessage"`
	UpdateCheck   string        `json:"updateCheck"`
	StartTime     time.Time     `json:"startTime" kv:"FeedTime:2"`
	URL           string        `json:"url"`
	ContentLength int           `json:"contentLength"`
	ContentType   string        `json:"contentType"`
	ETag          string        `json:"etag"`
	LastModified  time.Time     `json:"lastModified"`
	StatusCode    int           `json:"statusCode"`
	UsesGzip      bool          `json:"gzip"`
	Flavor        string        `json:"flavor"`
	Generator     string        `json:"generator"`
	Title         string        `json:"title"`
	LastUpdated   time.Time     `json:"lastUpdated"`
}
