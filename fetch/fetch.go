package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"sync"
	"time"
)

const (
	defaultTimeout   = time.Second * 20
	httpUserAgent    = "Rakewire " + m.VERSION
	hAcceptEncoding  = "Accept-Encoding"
	hContentEncoding = "Content-Encoding"
	hEtag            = "ETag"
	hIfModifiedSince = "If-Modified-Since"
	hIfNoneMatch     = "If-None-Match"
	hLastModified    = "Last-Modified"
	hUserAgent       = "User-Agent"
	mGET             = "GET"
	gzip             = "gzip"
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
		client:  newInternalClient(timeout),
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
	if z != nil { // #TODO:60 remove hack because on app close object is apparently already garbage collected
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
	req.Header.Set(hUserAgent, httpUserAgent)
	req.Header.Set(hAcceptEncoding, gzip)
	if feed.Last200 != nil {
		if feed.Last200.LastModified != nil {
			req.Header.Set(hIfModifiedSince, feed.Last200.LastModified.UTC().Format(http.TimeFormat))
		}
		if feed.Last200.ETag != "" {
			req.Header.Set(hIfNoneMatch, feed.Last200.ETag)
		}
	}
	return req
}

func (z *Service) processFeed(feed *m.Feed, id int) {

	startTime := time.Now().UTC().Truncate(time.Millisecond)
	now := startTime.Truncate(time.Second)
	feed.Attempt = &m.FeedLog{}

	feed.Attempt.StartTime = &now

	rsp, err := z.client.Do(z.newRequest(feed))
	if err != nil {
		feed.Attempt.Result = m.FetchResultClientError
		feed.Attempt.ResultMessage = err.Error()
	} else {

		// #TODO:0 test unexplicitly adding accept-encoding header for transparent decompression
		buf := &bytes.Buffer{}
		io.Copy(buf, rsp.Body)
		rsp.Body.Close()
		feed.Body = buf.Bytes()
		feed.Attempt.StatusCode = rsp.StatusCode

		if rsp.StatusCode == http.StatusMovedPermanently {
			feed.Attempt.Result = m.FetchResultRedirect
			feed.Attempt.ResultMessage = fmt.Sprintf("%s -> %s", feed.URL, rsp.Request.URL.String())
			feed.URL = rsp.Request.URL.String() // update feed
		} else if rsp.StatusCode == http.StatusOK || rsp.StatusCode == http.StatusNotModified {

			// #DOING:30 remove block
			feed.LastFetch = &now

			feed.Attempt.Result = m.FetchResultOK
			feed.Attempt.ETag = rsp.Header.Get(hEtag)
			feed.Attempt.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
			feed.Attempt.UsesGzip = usesGzip(rsp.Header.Get(hContentEncoding))

			if rsp.StatusCode == http.StatusOK {

				feed.Attempt.ContentLength = len(feed.Body)
				feed.Attempt.Checksum = checksum(feed.Body)
				feed.Attempt.UpdateCheck = m.UpdateCheckChecksum
				if feed.Last200 != nil && feed.Last200.Checksum != "" {
					if feed.Last200.Checksum != feed.Attempt.Checksum {
						// updated - reset back to minimum
						// #TODO:40 add UpdateCheckFeedEntries check
						feed.Attempt.IsUpdated = true
						feed.ResetInterval()
					} else {
						// not updated - use backoff policy to increase interval
						feed.Attempt.IsUpdated = false // not modified but site doesn't support conditional GETs
						feed.BackoffInterval()
					}
				} else {
					feed.Attempt.IsUpdated = true
				}

			} else if rsp.StatusCode == http.StatusNotModified { // 304 not modified
				// not updated - use backoff policy to increase interval
				feed.Attempt.IsUpdated = false
				feed.Attempt.UpdateCheck = m.UpdateCheck304
				feed.BackoffInterval()
			}

		} else if rsp.StatusCode >= 400 {
			// don't hammer site if error
			feed.BackoffIntervalError()
			feed.Attempt.Result = m.FetchResultServerError
		}

	} // err

	feed.Attempt.Duration = time.Now().Truncate(time.Millisecond).Sub(startTime)

	logger.Printf("fetch %2d: %3d %s %s\n", id, feed.Attempt.StatusCode, feed.URL, feed.Attempt.ResultMessage)
	z.output <- feed

}
