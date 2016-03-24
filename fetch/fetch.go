package fetch

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rakewire/feedparser"
	"rakewire/model"
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

const (
	fetchTimeout = "fetch.timeout"
	fetchWorkers = "fetch.workers"
)

var (
	fetchTimeoutDefault = time.Second * 20
	fetchWorkersDefault = 10
	httpUserAgent       = "Rakewire " + model.Version
)

// Service fetches feeds
type Service struct {
	sync.Mutex
	running bool
	input   chan *model.Feed
	output  chan *model.Harvest
	workers int
	latch   sync.WaitGroup
	client  *http.Client
}

// NewService create new fetcher service
func NewService(cfg *model.Configuration, input chan *model.Feed, output chan *model.Harvest) *Service {
	timeout, err := time.ParseDuration(cfg.GetStr(fetchTimeout, fetchTimeoutDefault.String()))
	if err != nil {
		timeout = fetchTimeoutDefault
	}
	return &Service{
		input:   input,
		output:  output,
		workers: cfg.GetInt(fetchWorkers, fetchWorkersDefault),
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

func (z *Service) processFeed(feed *model.Feed, id int) {

	harvest := &model.Harvest{
		Feed: feed,
	}

	startTime := time.Now().UTC().Truncate(time.Millisecond)
	now := startTime.Truncate(time.Second)

	harvest.Transmission = model.T.New(feed.ID)
	harvest.Transmission.URL = feed.URL
	harvest.Transmission.StartTime = now

	rsp, err := z.client.Do(newRequest(feed))
	if err != nil && (rsp == nil || rsp.StatusCode != http.StatusMovedPermanently) {
		processFeedClientError(harvest, err)
	} else {

		harvest.Transmission.StatusCode = rsp.StatusCode

		switch {

		case rsp.StatusCode == http.StatusMovedPermanently:
			processFeedMovedPermanently(harvest, rsp)

		case rsp.StatusCode == http.StatusOK:
			reader, _ := readBody(rsp)
			body := &ReadCounter{ReadCloser: reader}
			p := feedparser.NewParser()
			xmlFeed, err := p.Parse(body, harvest.Transmission.ContentType)

			processFeedOK(harvest, rsp)
			if err != nil || xmlFeed == nil {
				processFeedOKButCannotParse(harvest, err)
			} else {
				processFeedOKAndParse(harvest, body.Size, xmlFeed)
			}

		case rsp.StatusCode == http.StatusNotModified:
			processFeedNotModified(harvest, rsp)

		case rsp.StatusCode >= 400:
			processFeedServerError(harvest, rsp)

		case true:
			log.Printf("%-7s %-7s Uncaught Status Code: %d", logWarn, logName, rsp.StatusCode)

		} // switch

	} // not err

	harvest.Transmission.Duration = time.Now().Truncate(time.Millisecond).Sub(startTime)
	if feed.StatusSince.IsZero() || feed.Status != harvest.Transmission.Result {
		feed.StatusSince = time.Now()
	}
	feed.Status = harvest.Transmission.Result
	feed.StatusMessage = harvest.Transmission.ResultMessage

	z.output <- harvest

}

func processFeedOK(harvest *model.Harvest, rsp *http.Response) {
	harvest.Transmission.Result = model.FetchResultOK
	harvest.Transmission.ContentType = rsp.Header.Get(hContentType)
	harvest.Transmission.ETag = rsp.Header.Get(hEtag)
	harvest.Transmission.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
	harvest.Transmission.UsesGzip = usesGzip(rsp.Header.Get(hContentEncoding))
	harvest.Feed.ETag = harvest.Transmission.ETag
	harvest.Feed.LastModified = harvest.Transmission.LastModified
}

func processFeedOKButCannotParse(harvest *model.Harvest, err error) {
	harvest.Transmission.Result = model.FetchResultFeedError
	harvest.Transmission.ResultMessage = err.Error()
}

func processFeedOKAndParse(harvest *model.Harvest, size int, xmlFeed *feedparser.Feed) {
	harvest.Transmission.ContentLength = size
	harvest.Transmission.Flavor = xmlFeed.Flavor
	harvest.Transmission.Generator = xmlFeed.Generator
	harvest.Transmission.Title = xmlFeed.Title
	harvest.Feed.Title = xmlFeed.Title
	harvest.Feed.SiteURL = xmlFeed.LinkAlternate

	// convert Items to Items
	for _, xmlEntry := range xmlFeed.Entries {
		item := harvest.AddItem(xmlEntry.ID)
		item.Created = xmlEntry.Created
		item.Updated = xmlEntry.Updated
		item.Title = xmlEntry.Title
		item.URL = xmlEntry.LinkAlternate
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

func processFeedClientError(harvest *model.Harvest, err error) {
	harvest.Transmission.Result = model.FetchResultClientError
	harvest.Transmission.ResultMessage = err.Error()
}

func processFeedServerError(harvest *model.Harvest, rsp *http.Response) {
	harvest.Transmission.Result = model.FetchResultServerError
}

func processFeedMovedPermanently(harvest *model.Harvest, rsp *http.Response) {
	harvest.Transmission.Result = model.FetchResultRedirect
	newURL := rsp.Header.Get(hLocation)
	harvest.Transmission.ResultMessage = fmt.Sprintf("%s moved %s", harvest.Feed.URL, newURL)
	harvest.Feed.URL = newURL // update feed
}

func processFeedNotModified(harvest *model.Harvest, rsp *http.Response) {
	harvest.Transmission.Result = model.FetchResultOK
	harvest.Transmission.ETag = rsp.Header.Get(hEtag)
	harvest.Transmission.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
}

func newRequest(feed *model.Feed) *http.Request {
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
