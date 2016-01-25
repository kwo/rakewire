package fetch

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rakewire/feedparser"
	m "rakewire/model"
	"sync"
	"time"
)

const (
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

const (
	logName  = "[fetch]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
)

var (
	defaultTimeout = time.Second * 20
	httpUserAgent  = "Rakewire " + m.Version
)

// Configuration configuration
type Configuration struct {
	Workers int
	Timeout string
}

// Service fetches feeds
type Service struct {
	sync.Mutex
	running bool
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
func (z *Service) Start() error {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Printf("%-7s %-7s service already started, exiting...", logWarn, logName)
		return ErrRestart
	}

	log.Printf("%-7s %-7s service starting...", logDebug, logName)
	for i := 0; i < z.workers; i++ {
		z.latch.Add(1)
		go z.run(i)
	}
	z.running = true
	log.Printf("%-7s %-7s service started", logInfo, logName)
	return nil

}

// Stop service
func (z *Service) Stop() {

	// TODO #RAKEWIRE-55: remove hack because on app close object is apparently already garbage collected
	if z == nil {
		log.Printf("%-7s %-7s service is nil, exiting...", logError, logName)
		return
	}

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Printf("%-7s %-7s service already stopped, exiting...", logWarn, logName)
		return
	}

	log.Printf("%-7s %-7s service stopping...", logDebug, logName)
	z.latch.Wait()
	z.input = nil
	z.output = nil
	z.running = false
	log.Printf("%-7s %-7s service stopped", logInfo, logName)

}

// IsRunning indicated if the service is active or not.
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}

func (z *Service) run(id int) {

	log.Printf("%-7s %-7s fetcher %2d starting...", logDebug, logName, id)

	for req := range z.input {
		z.processFeed(req, id)
	}

	log.Printf("%-7s %-7s fetcher %2d exited", logDebug, logName, id)
	z.latch.Done()

}

func (z *Service) processFeed(feed *m.Feed, id int) {

	startTime := time.Now().UTC().Truncate(time.Millisecond)
	now := startTime.Truncate(time.Second)
	feed.Attempt = m.NewFeedLog(feed.ID)

	feed.Attempt.URL = feed.URL
	feed.Attempt.StartTime = now

	rsp, err := z.client.Do(newRequest(feed))
	if err != nil && (rsp == nil || rsp.StatusCode != http.StatusMovedPermanently) {
		processFeedClientError(feed, err)
	} else {

		feed.Attempt.StatusCode = rsp.StatusCode

		switch {

		case rsp.StatusCode == http.StatusMovedPermanently:
			processFeedMovedPermanently(feed, rsp)

		case rsp.StatusCode == http.StatusOK:
			reader, _ := readBody(rsp)
			body := &ReadCounter{ReadCloser: reader}
			p := feedparser.NewParser()
			xmlFeed, err := p.Parse(body, feed.Attempt.ContentType)

			processFeedOK(feed, rsp)
			if err != nil || xmlFeed == nil {
				processFeedOKButCannotParse(feed, err)
			} else {
				processFeedOKAndParse(feed, body.Size, xmlFeed)
			}

		case rsp.StatusCode == http.StatusNotModified:
			processFeedNotModified(feed, rsp)

		case rsp.StatusCode >= 400:
			processFeedServerError(feed, rsp)

		case true:
			log.Printf("%-7s %-7s Uncaught Status Code: %d", logWarn, logName, rsp.StatusCode)

		} // switch

	} // not err

	feed.Attempt.Duration = time.Now().Truncate(time.Millisecond).Sub(startTime)
	if feed.StatusSince.IsZero() || feed.Status != feed.Attempt.Result {
		feed.StatusSince = time.Now()
	}
	feed.Status = feed.Attempt.Result
	feed.StatusMessage = feed.Attempt.ResultMessage

	z.output <- feed

}

func processFeedOK(feed *m.Feed, rsp *http.Response) {
	feed.Attempt.Result = m.FetchResultOK
	feed.Attempt.ContentType = rsp.Header.Get(hContentType)
	feed.Attempt.ETag = rsp.Header.Get(hEtag)
	feed.Attempt.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
	feed.Attempt.UsesGzip = usesGzip(rsp.Header.Get(hContentEncoding))
	feed.ETag = feed.Attempt.ETag
	feed.LastModified = feed.Attempt.LastModified
}

func processFeedOKButCannotParse(feed *m.Feed, err error) {
	feed.Attempt.Result = m.FetchResultFeedError
	feed.Attempt.ResultMessage = err.Error()
}

func processFeedOKAndParse(feed *m.Feed, size int, xmlFeed *feedparser.Feed) {
	feed.Attempt.ContentLength = size
	feed.Attempt.Flavor = xmlFeed.Flavor
	feed.Attempt.Generator = xmlFeed.Generator
	feed.Attempt.Title = xmlFeed.Title
	feed.Title = xmlFeed.Title
	feed.SiteURL = xmlFeed.GetLinkAlternate()

	// convert Items to Items
	for _, xmlEntry := range xmlFeed.Entries {
		item := feed.AddItem(xmlEntry.ID)
		item.Created = xmlEntry.Created
		item.Updated = xmlEntry.Updated
		item.Title = xmlEntry.Title
		item.URL = xmlEntry.GetLinkAlternate()
		if len(xmlEntry.Authors) > 0 {
			item.Author = xmlEntry.Authors[0]
		}
		if xmlEntry.Content != "" {
			item.Content = xmlEntry.Content
		} else {
			item.Content = xmlEntry.Summary
		}
	}

}

func processFeedClientError(feed *m.Feed, err error) {
	feed.Attempt.Result = m.FetchResultClientError
	feed.Attempt.ResultMessage = err.Error()
}

func processFeedServerError(feed *m.Feed, rsp *http.Response) {
	feed.Attempt.Result = m.FetchResultServerError
}

func processFeedMovedPermanently(feed *m.Feed, rsp *http.Response) {
	feed.Attempt.Result = m.FetchResultRedirect
	newURL := rsp.Header.Get(hLocation)
	feed.Attempt.ResultMessage = fmt.Sprintf("%s moved %s", feed.URL, newURL)
	feed.URL = newURL // update feed
}

func processFeedNotModified(feed *m.Feed, rsp *http.Response) {
	feed.Attempt.Result = m.FetchResultOK
	feed.Attempt.ETag = rsp.Header.Get(hEtag)
	feed.Attempt.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
}

func newRequest(feed *m.Feed) *http.Request {
	req, _ := http.NewRequest(mGET, feed.URL, nil)
	req.Header.Set(hUserAgent, httpUserAgent)
	req.Header.Set(hAcceptEncoding, "gzip")
	if !feed.LastModified.IsZero() {
		req.Header.Set(hIfModifiedSince, feed.LastModified.UTC().Format(http.TimeFormat))
	}
	if feed.ETag != "" {
		req.Header.Set(hIfNoneMatch, feed.ETag)
	}
	return req
}
