package model

//go:generate gokv $GOFILE

import (
	"time"
)

// Transmission represents an attempted HTTP request to a feed
type Transmission struct {
	ID            uint64 `kv:"Time:2"`
	FeedID        uint64 `kv:"+required,FeedTime:1"`
	Duration      time.Duration
	Result        string
	ResultMessage string
	StartTime     time.Time `kv:"Time:1,FeedTime:2"`
	URL           string
	ContentLength int
	ContentType   string
	ETag          string
	LastModified  time.Time
	StatusCode    int
	UsesGzip      bool
	Flavor        string
	Generator     string
	Title         string
	LastUpdated   time.Time
	ItemCount     int
	NewItems      int
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
		if id, _, err := fn(transmissionEntity, tx); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
