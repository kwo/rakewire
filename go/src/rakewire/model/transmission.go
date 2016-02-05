package model

//go:generate gokv $GOFILE

import (
	"time"
)

// Transmission represents an attempted HTTP request to a feed
type Transmission struct {
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

// NewTransmission instantiates a new Transmission with the required fields set.
func NewTransmission(feedID uint64) *Transmission {
	return &Transmission{
		FeedID: feedID,
	}
}

func (z *Transmission) setIDIfNecessary(fn fnNextID, tx Transaction) error {
	if z.ID == 0 {
		if id, err := fn(transmissionEntity, tx); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
