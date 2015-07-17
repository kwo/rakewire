package fetch

import (
	"io"
	"io/ioutil"
	"net/http"
	"rakewire.com/logging"
	"time"
)

const (
	httpUserAgent = "Rakewire Bot 0.0.1"
)

var (
	logger = logging.New("fetch")
)

// Service fetches feeds
type Service struct {
	fetchers     []*fetcher
	requests     chan *Request
	responses    chan *Response
	fetcherCount int
	httpTimeout  time.Duration
}

type fetcher struct {
	id         int
	client     *http.Client
	requests   chan *Request
	responses  chan *Response
	killsignal chan bool
}

// NewService create new fetcher service
func NewService(cfg *Configuration) *Service {
	return &Service{
		requests:     make(chan *Request, cfg.RequestBuffer),
		responses:    make(chan *Response),
		fetcherCount: cfg.Fetchers,
		httpTimeout:  time.Duration(cfg.HTTPTimeoutSeconds) * time.Second,
	}
}

// Start service
func (z *Service) Start() {

	logger.Println("starting service")

	for i := 0; i < z.fetcherCount; i++ {

		f := &fetcher{
			id: i,
			client: &http.Client{
				// CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 	return http.ErrNotSupported
				// },
				Timeout: z.httpTimeout,
			},
			requests:   z.requests,
			responses:  z.responses,
			killsignal: make(chan bool),
		}

		z.fetchers = append(z.fetchers, f)

		go f.start()

	}

}

// Stop service
func (z *Service) Stop() {

	logger.Println("stopping service")
	for _, f := range z.fetchers {
		f.stop()
	}
	z.fetchers = nil

}

// Add feeds to pool
func (z *Service) Add(requests []*Request) {
	for _, req := range requests {
		z.requests <- req
	}
}

// Harvest responses from service
func (z *Service) Harvest() chan *Response {
	return z.responses
}

func (z *fetcher) start() {

Loop:
	for {
		select {
		case <-z.killsignal:
			break Loop
		case req := <-z.requests:
			z.processFeed(req)
		}
	}

}

func (z *fetcher) stop() {
	logger.Printf("stopping fetcher: %2d\n", z.id)
	z.killsignal <- true
	close(z.killsignal)
	z.client = nil
	//logger.Printf("exiting fetcher %2d", z.id)
}

func (z *fetcher) newRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", httpUserAgent)
	return req
}

func (z *fetcher) processFeed(req *Request) {

	now := time.Now()
	result := &Response{
		FetcherID:   z.id,
		AttemptTime: &now,
		ID:          req.ID,
		URL:         req.URL,
	}

	rsp, err := z.client.Do(z.newRequest(req.URL))

	if err == nil {
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		result.StatusCode = rsp.StatusCode
		if req.URL != rsp.Request.URL.String() {
			result.URL = rsp.Request.URL.String()
			result.StatusCode = 3000
		} else {
			result.FetchTime = &now
			result.ETag = rsp.Header.Get("etag")
			m, err := http.ParseTime(rsp.Header.Get("lastmodified"))
			if err != nil {
				result.LastModified = &m
			}
		}
	} else {
		result.Failed = true
		result.Message = err.Error()
		result.StatusCode = 5000
	}

	z.responses <- result

}
