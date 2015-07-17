package fetch

import (
	"time"
)

// Configuration configuration
type Configuration struct {
	Fetchers           int
	RequestBuffer      int
	HTTPTimeoutSeconds int
}

// Request input to fetcher service
type Request struct {
	FeedID       string
	ETag         string
	LastModified *time.Time
	URL          string
}

// Response output from fetcher service
type Response struct {
	AttemptTime  *time.Time
	ETag         string
	Failed       bool
	FetcherID    int
	FetchTime    *time.Time
	LastModified *time.Time
	Message      string
	Request      *Request
	StatusCode   int
	URL          string
}
