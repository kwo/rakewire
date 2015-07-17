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
// Also a type of db.Feed
type Request struct {
	ID           string
	ETag         string
	LastModified *time.Time
	URL          string
}

// Response output from fetcher service.
// Also a type of Request.
// TODO: Also a type of db.Feed.
type Response struct {
	ID           string
	AttemptTime  *time.Time
	ETag         string
	Failed       bool
	FetcherID    int
	FetchTime    *time.Time
	LastModified *time.Time
	Message      string
	StatusCode   int
	URL          string
}
