package fetch

import (
	"github.com/stretchr/testify/require"
	"os"
	m "rakewire.com/model"
	"testing"
)

func TestFetch(t *testing.T) {

	t.SkipNow()

	r, err := os.Open("../test/feedlist.txt")
	require.Nil(t, err)
	require.NotNil(t, r)
	//feeds := []*Request{}
	feeds := URLListToFeeds(r)
	r.Close()
	require.NotNil(t, feeds)

	logger.Printf("feeds: %d\n", feeds.Size())

	cfg := &Configuration{
		Fetchers:           20,
		HTTPTimeoutSeconds: 20,
	}

	requests := make(chan *m.Feed)
	responses := make(chan *m.Feed)

	ff := NewService(cfg, requests, responses)
	ff.Start()

	go func() {
		logger.Printf("adding feeds: %d\n", feeds.Size())
		for _, f := range feeds.Values {
			//logger.Printf("adding feed: %s\n", f.URL)
			requests <- f
		}
		close(requests)
		logger.Println("adding feeds done")
	}()

	go func() {
		logger.Println("monitoring...")
		for rsp := range responses {
			logger.Printf("%5t %s\n", rsp.Failed, rsp.URL)
		}
		logger.Println("monitoring done")
	}()

	ff.Stop()

}
