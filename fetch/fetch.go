package fetch

import (
	"io"
	"io/ioutil"
	"net/http"
	"rakewire.com/logging"
	"sync"
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
	Input        chan *Request
	Output       chan *Response
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
	if z != nil { // hack because on app close object is apparantly already garbage collected
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

func (z *Service) processFeed(req *Request, id int) {

	now := time.Now()
	result := &Response{
		FetcherID:   id,
		AttemptTime: &now,
		ID:          req.ID,
		URL:         req.URL,
	}

	rsp, err := z.client.Do(z.newRequest(req.URL))
	if rsp != nil && rsp.Body != nil {
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
	}

	if err == nil {
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

	z.Output <- result

}
