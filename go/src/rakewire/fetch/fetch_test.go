package fetch

import (
	m "rakewire/model"
	"testing"
)

func TestFetch(t *testing.T) {

	t.SkipNow()

	feeds, err := m.ParseFeedsFromFile("../test/feedlist.txt")
	assertNoError(t, err)
	assertNotNil(t, feeds)

	logger.Debugf("feeds: %d\n", len(feeds))

	cfg := &Configuration{
		Workers: 20,
		Timeout: "20s",
	}

	requests := make(chan *m.Feed)
	responses := make(chan *m.Feed)

	ff := NewService(cfg, requests, responses)
	ff.Start()

	go func() {
		logger.Debugf("adding feeds: %d\n", len(feeds))
		for _, f := range feeds {
			logger.Debugf("adding feed: %s\n", f.URL)
			requests <- f
		}
		close(requests)
		logger.Debug("adding feeds done")
	}()

	go func() {
		logger.Debug("monitoring...")
		for rsp := range responses {
			logger.Debugf("%3d %s\n", rsp.Attempt.StatusCode, rsp.URL)
		}
		logger.Info("monitoring done")
	}()

	ff.Stop()

}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("Expected not nil value")
	}
}
