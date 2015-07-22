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
	httpUserAgent    = "Rakewire " + app.VERSION
	hEtag            = "Etag"
	hIfModifiedSince = "If-Modified-Since"
	hIfNoneMatch     = "If-None-Match"
	hLastModified    = "Last-Modified"
	hUserAgent       = "User-Agent"
	mGET             = "GET"
)

var (
	logger = logging.New("fetch")
)

// Service fetches feeds
type Service struct {
	input        chan *m.Feed
	output       chan *m.Feed
	fetcherCount int
	latch        sync.WaitGroup
	client       *http.Client
}

// NewService create new fetcher service
func NewService(cfg *Configuration, input chan *m.Feed, output chan *m.Feed) *Service {
	return &Service{
		input:        input,
		output:       output,
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
		z.input = nil
		z.output = nil
	}
	logger.Println("service stopped")
}

func (z *Service) run(id int) {

	logger.Printf("fetcher %2d starting...\n", id)

	for req := range z.input {
		z.processFeed(req, id)
	}

	logger.Printf("fetcher %2d exited.\n", id)
	z.latch.Done()

}

func (z *Service) newRequest(feed *m.Feed) *http.Request {
	req, _ := http.NewRequest(mGET, feed.URL, nil)
	if feed.LastModified != nil {
		req.Header.Set(hIfModifiedSince, feed.LastModified.UTC().Format(http.TimeFormat))
	}
	if feed.ETag != "" {
		req.Header.Set(hIfNoneMatch, feed.ETag)
	}
	req.Header.Set(hUserAgent, httpUserAgent)
	return req
}

func (z *Service) processFeed(feed *m.Feed, id int) {

	now := time.Now().Truncate(time.Second)

	feed.LastAttempt = &now

	status := 0
	message := ""
	rsp, err := z.client.Do(z.newRequest(feed))

	if rsp != nil && rsp.Body != nil {
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		status = rsp.StatusCode
	}

	if err == nil {
		feed.Failed = false
		if feed.URL != rsp.Request.URL.String() {
			feed.URL = rsp.Request.URL.String()
		} else {
			feed.LastFetch = &now
			feed.ETag = rsp.Header.Get(hEtag)
			m, err := http.ParseTime(rsp.Header.Get(hLastModified))
			if err == nil && !m.IsZero() {
				feed.LastModified = &m
			}
		}
	} else {
		feed.Failed = true
		message = err.Error()
	}

	logger.Printf("fetch %2d: %5t %3d %s - %s\n", id, feed.Failed, status, feed.URL, message)
	z.output <- feed

}
