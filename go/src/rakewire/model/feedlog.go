package model

//go:generate gokv $GOFILE

import (
	"time"
)

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	ID            uint64        `json:"id" kv:"Time:2"`
	FeedID        uint64        `json:"feedId" kv:"+required,FeedTime:1"`
	Duration      time.Duration `json:"duration"`
	Result        string        `json:"result"`
	ResultMessage string        `json:"resultMessage"`
	StartTime     time.Time     `json:"startTime" kv:"Time:1,FeedTime:2"`
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
	ItemCount     int           `json:"itemCount"`
	NewItems      int           `json:"newItems"`
}

// FetchResults
const (
	FetchResultOK          = "OK"
	FetchResultRedirect    = "MV" // message contains old URL -> new URL
	FetchResultClientError = "EC" // message contains error text
	FetchResultServerError = "ES" // check http status code
	FetchResultFeedError   = "FP" // cannot parse feed
)

// NewFeedLog instantiates a new FeedLog with the required fields set.
func NewFeedLog(feedID uint64) *FeedLog {
	return &FeedLog{
		FeedID: feedID,
	}
}
