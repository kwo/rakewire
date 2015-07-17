package fetch

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {

	t.SkipNow()

	r, err := os.Open("../test/feedlist.txt")
	require.Nil(t, err)
	require.NotNil(t, r)
	feeds := URLListToRequestArray(r)
	r.Close()
	require.NotNil(t, feeds)

	cfg := &Configuration{
		Fetchers:           20,
		RequestBuffer:      500,
		HTTPTimeoutSeconds: 10,
	}
	ff := NewService(cfg)
	go ff.Start()
	ff.Add(feeds)

monitor:
	for {
		select {
		case rsp := <-ff.Harvest():
			logger.Printf("Worker: %2d, %4d %s %s\n", rsp.FetcherID, rsp.StatusCode, rsp.URL, rsp.Message)
		case <-time.After(5 * time.Second):
			break monitor
		}
	}

	ff.Stop()

}
