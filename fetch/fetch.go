package fetch

// TODO: stop fetcher without sending special message?

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	m "rakewire.com/model"
	"time"
)

type status struct {
	halt       bool
	Index      int
	URL        string
	WorkerID   int
	StatusCode int
	Message    string
}

type fetcher struct {
	id     int
	client *http.Client
}

const (
	fetcherCount = 10
)

// Fetch feeds in file
func Fetch(feeds *m.Feeds) error {

	requests := make(chan status)
	responses := make(chan status, 5)
	signals := make(chan bool)

	initFetchers(requests, responses)
	go addFeeds(feeds, requests)
	go processFeeds(feeds.Size(), responses, signals)
	<-signals
	go destroyFetchers(requests, signals)
	<-signals

	log.Println("exiting...")

	return nil

}

func (f *fetcher) getFeed(chReq chan status, chRsp chan status) {

	log.Printf("Fetcher %2d started", f.id)

	for {

		result := <-chReq
		if result.halt {
			break
		}

		result.WorkerID = f.id

		var rsp, err = f.client.Get(result.URL)
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()

		if err != nil {
			result.StatusCode = 5000
			result.Message = err.Error()
		} else if result.URL != rsp.Request.URL.String() {
			result.StatusCode = 3000
			result.URL = rsp.Request.URL.String()
		} else {
			result.StatusCode = rsp.StatusCode
		}

		chRsp <- result

	}

	log.Printf("Fetcher %2d exited", f.id)

}

func initFetchers(requests chan status, responses chan status) {

	for i := 0; i < fetcherCount; i++ {
		f := &fetcher{
			id: i,
			client: &http.Client{
				Timeout: 60 * time.Second,
			},
		}
		go f.getFeed(requests, responses)
	}

}

func destroyFetchers(requests chan status, signals chan bool) {

	for i := 0; i < fetcherCount; i++ {
		s := status{
			halt: true,
		}
		requests <- s
	}

	signals <- true

}

func addFeeds(feeds *m.Feeds, requests chan status) {

	for index, feed := range feeds.Values {
		s := status{
			Index: index,
			URL:   feed.URL,
		}
		requests <- s
	}

}

func processFeeds(total int, responses chan status, signals chan bool) {

	var counter int
	for {
		s := <-responses
		log.Printf("Worker: %2d, Feed %4d: %4d %s %s\n", s.WorkerID, s.Index, s.StatusCode, s.URL, s.Message)
		counter++
		if counter >= total {
			break
		}
	}

	signals <- true

}
