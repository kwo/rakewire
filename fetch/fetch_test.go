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
	//feeds := []*Request{}
	feeds := URLListToRequestArray(r)
	r.Close()
	require.NotNil(t, feeds)

	logger.Printf("feeds: %d\n", len(feeds))

	cfg := &Configuration{
		Fetchers:           20,
		HTTPTimeoutSeconds: 10,
	}

	chReq := make(chan *Request)
	chRsp := make(chan *Response)

	ff := NewService(cfg, chReq, chRsp)
	ff.Start()

	// add feeds
	go func() {
		logger.Printf("adding feeds: %d\n", len(feeds))
		for _, f := range feeds {
			logger.Printf("adding feed: %s\n", f.URL)
			chReq <- f
		}
		close(chReq)
	}()

	logger.Println("monitoring...")

monitor:
	for {
		select {
		case rsp := <-chRsp:
			logger.Printf("Worker: %2d, %4d %s %s\n", rsp.FetcherID, rsp.StatusCode, rsp.URL, rsp.Message)
		case <-time.After(time.Duration(cfg.HTTPTimeoutSeconds+1) * time.Second):
			break monitor
		}
	}

	ff.Stop()

}
