package model

//go:generate gokv $GOFILE

import (
	"time"
)

// Transmission represents an attempted HTTP request to a feed
type Transmission struct {
	ID            string        `json:"id" kv:"Time:2"`
	FeedID        string        `json:"feedID" kv:"+required,FeedTime:1"`
	Duration      time.Duration `json:"duration"`
	Result        string        `json:"result,omitempty"`
	ResultMessage string        `json:"resultMessage,omitempty"`
	StartTime     time.Time     `json:"startTime" kv:"Time:1,FeedTime:2"`
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

// FetchResults
const (
	FetchResultOK          = "OK"
	FetchResultRedirect    = "MV" // message contains old URL -> new URL
	FetchResultClientError = "EC" // message contains error text
	FetchResultServerError = "ES" // check http status code
	FetchResultFeedError   = "FP" // cannot parse feed
)

// NewTransmission instantiates a new Transmission with the required fields set.
func NewTransmission(feedID string) *Transmission {
	return &Transmission{
		FeedID: feedID,
	}
}

func (z *Transmission) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
