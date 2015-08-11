package fetch

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
		Workers: 20,
		Timeout: "20s",
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
			logger.Printf("%3d %s\n", rsp.Attempt.StatusCode, rsp.URL)
		}
		logger.Println("monitoring done")
	}()

	ff.Stop()

}

func TestHash(t *testing.T) {

	f, err := os.Open("../test/feed.xml")
	require.Nil(t, err)
	require.NotNil(t, f)

	hash := sha256.New()
	_, err = io.Copy(hash, f)
	assert.Nil(t, err)
	f.Close()

	d := hash.Sum(nil)
	assert.NotNil(t, d)
	assert.Equal(t, 32, len(d))
	cs := hex.EncodeToString(d)
	assert.NotNil(t, cs)
	assert.Equal(t, 64, len(cs))

	//logger.Printf("file hash: %s", cs)

}
