package fetch

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	r, err := os.Open("../test/feedlist.txt")
	require.Nil(t, err)
	require.NotNil(t, r)
	//feeds := []*Request{}
	feeds := URLListToRequestArray(r)
	r.Close()
	require.NotNil(t, feeds)

	logger.Printf("feeds: %d\n", len(feeds))

	cfg := &Configuration{
		Fetchers:           20,
		HTTPTimeoutSeconds: 20,
	}

	requests := make(chan *Request)
	responses := make(chan *Response)

	ff := NewService(cfg, requests, responses)
	ff.Start()

	go func() {
		logger.Printf("adding feeds: %d\n", len(feeds))
		for _, f := range feeds {
			//logger.Printf("adding feed: %s\n", f.URL)
			requests <- f
		}
		close(requests)
		logger.Println("adding feeds done")
	}()

	go func() {
		logger.Println("monitoring...")
		for rsp := range responses {
			logger.Printf("Worker: %2d, %4d %s %s\n", rsp.FetcherID, rsp.StatusCode, rsp.URL, rsp.Message)
		}
		logger.Println("monitoring done")
	}()

	ff.Stop()

}
