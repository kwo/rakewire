package fetch

import (
	"io"
	"io/ioutil"
	"net/http"
	"rakewire.com/app"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"sync"
	"time"
)

const (
	httpUserAgent = "Rakewire " + app.VERSION
)

var (
	logger = logging.New("fetch")
)

// Service fetches feeds
type Service struct {
	Input        chan *m.Feed
	Output       chan *m.Feed
	fetcherCount int
	latch        sync.WaitGroup
	client       *http.Client
}

// NewService create new fetcher service
func NewService(cfg *Configuration) *Service {
	return &Service{
		fetcherCount: cfg.Fetchers,
		client: &http.Client{
			// CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 	return http.ErrNotSupported
			// },
			Timeout: time.Duration(cfg.HTTPTimeoutSeconds) * time.Second,
		},
	}
}

// Start service
func (z *Service) Start() {
	logger.Println("service starting...")
	// initialize fetchers
	for i := 0; i < z.fetcherCount; i++ {
		z.latch.Add(1)
		go z.run(i)
	} // for
	logger.Println("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Println("service stopping...")
	if z != nil { // hack because on app close object is apparently already garbage collected
		z.latch.Wait()
		z.Input = nil
		z.Output = nil
	}
	logger.Println("service stopped")
}

func (z *Service) run(id int) {

	logger.Printf("fetcher %2d starting...\n", id)

	for req := range z.Input {
		z.processFeed(req, id)
	}

	logger.Printf("fetcher %2d exited.\n", id)
	z.latch.Done()

}

func (z *Service) newRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", httpUserAgent)
	return req
}

func (z *Service) processFeed(feed *m.Feed, id int) {

	now := time.Now()

	feed.LastAttempt = &now

	status := 0
	rsp, err := z.client.Do(z.newRequest(feed.URL))
	if rsp != nil && rsp.Body != nil {
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		status = rsp.StatusCode
	}

	if err == nil {
		if feed.URL != rsp.Request.URL.String() {
			feed.URL = rsp.Request.URL.String()
		} else {
			feed.LastFetch = &now
			feed.ETag = rsp.Header.Get("etag")
			m, err := http.ParseTime(rsp.Header.Get("Last-Modified"))
			if err != nil && !m.IsZero() {
				feed.LastModified = &m
			}
		}
	} else {
		feed.Failed = true
	}

	logger.Printf("fetch %2d: %5t %3d %s\n", id, feed.Failed, status, feed.URL)
	z.Output <- feed

}
