package fetch

import (
	"fmt"
	"net/http"
	"rakewire/feedparser"
	"rakewire/logging"
	m "rakewire/model"
	"sync"
	"time"
)

const (
	defaultTimeout   = time.Second * 20
	httpUserAgent    = "Rakewire " + m.VERSION
	hAcceptEncoding  = "Accept-Encoding"
	hContentEncoding = "Content-Encoding"
	hContentType     = "Content-Type"
	hEtag            = "ETag"
	hIfModifiedSince = "If-Modified-Since"
	hIfNoneMatch     = "If-None-Match"
	hLastModified    = "Last-Modified"
	hLocation        = "Location"
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
	latch   *sync.WaitGroup
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
		latch:   &sync.WaitGroup{},
		client:  newInternalClient(timeout),
	}
}

// Start service
func (z *Service) Start() {
	logger.Info("service starting...")
	// initialize fetchers
	for i := 0; i < z.workers; i++ {
		z.latch.Add(1)
		go z.run(i)
	} // for
	logger.Info("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Info("service stopping...")
	if z != nil { // TODO #RAKEWIRE-55: remove hack because on app close object is apparently already garbage collected
		z.latch.Wait()
		z.input = nil
		z.output = nil
	}
	logger.Info("service stopped")
}

func (z *Service) run(id int) {

	logger.Infof("fetcher %2d starting...\n", id)

	for req := range z.input {
		z.processFeed(req, id)
	}

	logger.Infof("fetcher %2d exited.\n", id)
	z.latch.Done()

}

func (z *Service) newRequest(feed *m.Feed) *http.Request {
	req, _ := http.NewRequest(mGET, feed.URL, nil)
	req.Header.Set(hUserAgent, httpUserAgent)
	req.Header.Set(hAcceptEncoding, "gzip")
	last200 := feed.GetLast200()
	if last200 != nil && last200.Result == m.FetchResultOK {
		if !last200.LastModified.IsZero() {
			req.Header.Set(hIfModifiedSince, last200.LastModified.UTC().Format(http.TimeFormat))
		}
		if last200.ETag != "" {
			req.Header.Set(hIfNoneMatch, last200.ETag)
		}
	}
	return req
}

func (z *Service) processFeed(feed *m.Feed, id int) {

	startTime := time.Now().UTC().Truncate(time.Millisecond)
	now := startTime.Truncate(time.Second)
	feed.Attempt = m.NewFeedLog(feed.ID)

	feed.Attempt.URL = feed.URL
	feed.Attempt.StartTime = now

	rsp, err := z.client.Do(z.newRequest(feed))
	if err != nil && (rsp == nil || rsp.StatusCode != http.StatusMovedPermanently) {
		feed.Attempt.Result = m.FetchResultClientError
		feed.Attempt.ResultMessage = err.Error()
	} else {

		feed.Attempt.StatusCode = rsp.StatusCode

		switch {

		case rsp.StatusCode == http.StatusMovedPermanently:
			feed.Attempt.Result = m.FetchResultRedirect
			newURL := rsp.Header.Get(hLocation)
			feed.Attempt.ResultMessage = fmt.Sprintf("%s moved %s", feed.URL, newURL)
			feed.URL = newURL // update feed

		case rsp.StatusCode == http.StatusOK:
			feed.Attempt.Result = m.FetchResultOK
			feed.Attempt.ContentType = rsp.Header.Get(hContentType)
			feed.Attempt.ETag = rsp.Header.Get(hEtag)
			feed.Attempt.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
			feed.Attempt.UsesGzip = usesGzip(rsp.Header.Get(hContentEncoding))

			reader, _ := readBody(rsp)
			body := &ReadCounter{ReadCloser: reader}
			p := feedparser.NewParser()
			xmlFeed, err := p.Parse(body, feed.Attempt.ContentType)

			if err != nil || xmlFeed == nil {
				// cannot parse feed
				feed.Attempt.Result = m.FetchResultFeedError
				feed.Attempt.ResultMessage = err.Error()
				feed.AdjustFetchTime(1 * time.Hour) // give us time to work on solution
			} else {
				feed.Attempt.ContentLength = body.Size
				feed.Attempt.Flavor = xmlFeed.Flavor
				feed.Attempt.Generator = xmlFeed.Generator
				feed.Attempt.Title = xmlFeed.Title
				feed.Attempt.LastUpdated = xmlFeed.Updated
				if feed.Title == "" {
					feed.Title = xmlFeed.Title
				}
				if xmlFeed.Updated.IsZero() {
					feed.Attempt.Result = m.FetchResultFeedTimeError
					feed.Attempt.IsUpdated = false
					feed.Attempt.UpdateCheck = m.UpdateCheckFeed
					feed.AdjustFetchTime(1 * time.Hour) // give us time to work on solution
				} else {
					feed.Feed = xmlFeed
					feed.Attempt.IsUpdated = isFeedUpdated(feed.Feed.Updated, feed.LastUpdated)
					feed.Attempt.UpdateCheck = m.UpdateCheckFeed
					feed.UpdateFetchTime(feed.Feed.Updated)
				}
			}

		case rsp.StatusCode == http.StatusNotModified:
			feed.Attempt.Result = m.FetchResultOK
			feed.Attempt.IsUpdated = false
			feed.Attempt.UpdateCheck = m.UpdateCheck304
			feed.Attempt.ETag = rsp.Header.Get(hEtag)
			feed.Attempt.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
			feed.UpdateFetchTime(time.Time{})

		case rsp.StatusCode >= 400:
			feed.Attempt.Result = m.FetchResultServerError
			feed.AdjustFetchTime(24 * time.Hour) // don't hammer site if error

		case true:
			logger.Warnf("Uncaught Status Code: %d", rsp.StatusCode)

		} // switch

	} // err

	feed.Attempt.Duration = time.Now().Truncate(time.Millisecond).Sub(startTime)

	logger.Infof("fetch %2d: %2s  %3d  %5t  %2s  %s  %s %s\n", id, feed.Attempt.Result, feed.Attempt.StatusCode, feed.Attempt.IsUpdated, feed.Attempt.UpdateCheck, feed.URL, feed.Attempt.ResultMessage, feed.Attempt.Flavor)
	z.output <- feed

}
