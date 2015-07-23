package fetch

import (
	"bytes"
	"io"
	"net/http"
	"rakewire.com/app"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"sync"
	"time"
)

const (
	defaultTimeout   = time.Second * 20
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

// Configuration configuration
type Configuration struct {
	Workers int
	Timeout string
}

// Service fetches feeds
type Service struct {
	input   chan *m.Feed
	output  chan *m.Feed
	workers int
	latch   sync.WaitGroup
	client  *http.Client
}

// NewService create new fetcher service
func NewService(cfg *Configuration, input chan *m.Feed, output chan *m.Feed) *Service {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}
	return &Service{
		input:   input,
		output:  output,
		workers: cfg.Workers,
		client: &http.Client{
			// CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 	return http.ErrNotSupported
			// },
			Timeout: timeout,
		},
	}
}

// Start service
func (z *Service) Start() {
	logger.Println("service starting...")
	// initialize fetchers
	for i := 0; i < z.workers; i++ {
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

	message := ""
	rsp, err := z.client.Do(z.newRequest(feed))
	if err != nil {
		feed.StatusCode = 999
		message = err.Error()
	} else {

		buf := &bytes.Buffer{}
		io.Copy(buf, rsp.Body)
		rsp.Body.Close()
		feed.Body = buf.Bytes()
		feed.StatusCode = rsp.StatusCode

		if feed.URL != rsp.Request.URL.String() {
			feed.URL = rsp.Request.URL.String()
			feed.StatusCode = 300
		} else if rsp.StatusCode == 200 || rsp.StatusCode == 304 {

			feed.LastFetch = &now
			feed.ETag = rsp.Header.Get(hEtag)
			feed.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))

			if rsp.StatusCode == 200 {

				cs := checksum(feed.Body)
				if feed.Checksum != "" {
					if feed.Checksum != cs {
						// updated - reset back to minimum
						feed.Interval = m.FeedIntervalMin
					} else {
						// not updated - use backoff policy to increase interval
						feed.Interval *= 2
						if feed.Interval > m.FeedIntervalMax {
							feed.Interval = m.FeedIntervalMax
						}
					}
				}
				feed.Checksum = cs

			} else { // 304
				// not updated - use backoff policy to increase interval
				feed.Interval *= 2
				if feed.Interval > m.FeedIntervalMax {
					feed.Interval = m.FeedIntervalMax
				}
			}

		} // 200 or 304

	} // err

	logger.Printf("fetch %2d: %3d %s %s\n", id, feed.StatusCode, feed.URL, message)
	z.output <- feed

}
