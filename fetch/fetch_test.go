package fetch

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	r, err := os.Open("../test/feedlist.txt")
	require.Nil(t, err)
	require.NotNil(t, r)
	feeds := URLListToRequestArray(r)
	r.Close()
	require.NotNil(t, feeds)

	logger.Printf("feeds: %d\n", len(feeds))

	cfg := &Configuration{
		Fetchers:           20,
		RequestBuffer:      500,
		HTTPTimeoutSeconds: 10,
	}

	ch := &Channels{
		Requests:  make(chan *Request),
		Responses: make(chan *Response),
	}

	chErrors := make(chan error)

	ff := NewService(cfg, ch)
	go ff.Start(chErrors)

	// add feeds
	go func() {
		logger.Printf("adding feeds: %d\n", len(feeds))
		for _, f := range feeds {
			logger.Printf("adding feed: %s\n", f.URL)
			ch.Requests <- f
		}
	}()

	logger.Println("monitoring...")

monitor:
	for {
		select {
		case rsp := <-ch.Responses:
			logger.Printf("Worker: %2d, %4d %s %s\n", rsp.FetcherID, rsp.StatusCode, rsp.URL, rsp.Message)
		case <-time.After(5 * time.Second):
			break monitor
		}
	}

	ff.Stop()

}
