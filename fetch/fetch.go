package fetch

import (
	"errors"
	"fmt"
	"github.com/kwo/rakewire/feedparser"
	"github.com/kwo/rakewire/logger"
	"github.com/kwo/rakewire/model"
	"net/http"
	"net/url"
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

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
	log        = logger.New("fetch")
)

// Configuration contains all parameters for the Fetch service
type Configuration struct {
	TimeoutSeconds int
	Workers        int
	UserAgent      string
}

// Service fetches feeds
type Service struct {
	sync.Mutex
	running   bool
	input     chan *model.Feed
	output    chan *model.Harvest
	workers   int
	latch     sync.WaitGroup
	client    *http.Client
	userAgent string
}

// NewService create new fetcher service
func NewService(cfg *Configuration, input chan *model.Feed, output chan *model.Harvest) *Service {
	return &Service{
		input:     input,
		output:    output,
		workers:   cfg.Workers,
		client:    newInternalClient(cfg.TimeoutSeconds),
		userAgent: cfg.UserAgent,
	}
}

// Start service
func (z *Service) Start() error {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Debugf("service already started, exiting...")
		return ErrRestart
	}

	log.Infof("starting...")
	log.Infof("timeout:    %s", z.client.Timeout.String())
	log.Infof("workers:    %d", z.workers)
	log.Infof("user agent: %s", z.userAgent)

	for i := 0; i < z.workers; i++ {
		z.latch.Add(1)
		go z.run(i)
	}
	z.running = true
	log.Infof("started")
	return nil

}

// Stop service
func (z *Service) Stop() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Debugf("service already stopped, exiting...")
		return
	}

	log.Debugf("stopping...")
	z.latch.Wait()
	z.input = nil
	z.output = nil
	z.running = false
	log.Infof("stopped")

}

// IsRunning indicated if the service is active or not.
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}

func (z *Service) run(id int) {

	log.Debugf("fetcher %2d starting...", id)

	for req := range z.input {
		z.processFeed(req, id)
	}

	log.Debugf("fetcher %2d exited", id)
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

	rsp, err := z.client.Do(z.newRequest(feed))
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
			xmlFeed, err := p.Parse(body)

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
			log.Debugf("Uncaught Status Code: %d", rsp.StatusCode)

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
	harvest.Feed.URL = resolveURL(harvest.Feed.URL, newURL)
}

func processFeedNotModified(harvest *model.Harvest, rsp *http.Response) {
	harvest.Transmission.Result = model.FetchResultOK
	harvest.Transmission.ETag = rsp.Header.Get(hEtag)
	harvest.Transmission.LastModified = parseDateHeader(rsp.Header.Get(hLastModified))
}

func (z *Service) newRequest(feed *model.Feed) *http.Request {
	req, _ := http.NewRequest(mGET, feed.URL, nil)
	req.Header.Set(hUserAgent, z.userAgent)
	req.Header.Set(hAcceptEncoding, "gzip")
	if !feed.LastModified.IsZero() {
		req.Header.Set(hIfModifiedSince, feed.LastModified.UTC().Format(http.TimeFormat))
	}
	if feed.ETag != "" {
		req.Header.Set(hIfNoneMatch, feed.ETag)
	}
	return req
}

func resolveURL(uOriginal, uNew string) string {

	urlOriginal, errParse1 := url.Parse(uOriginal)
	if errParse1 != nil {
		return uOriginal
	}

	urlNew, errParse2 := url.Parse(uNew)
	if errParse2 != nil {
		return uOriginal
	}

	return urlOriginal.ResolveReference(urlNew).String()

}
