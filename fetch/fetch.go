package fetch

// TODO: io timout errors at feed 235 - 286 always
// TODO: stop fetcher without sending special message?

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
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
	id int
}

const (
	fetcherCount = 20
)

var (
	httpClient = http.Client{
		Timeout: 60 * time.Second,
	}
)

// Fetch feeds in file
func Fetch(feedfile string) error {

	var feedlist, err = readFile(feedfile)
	if err != nil {
		return err
	}

	var total = len(feedlist)
	requests := make(chan status)
	responses := make(chan status, 5)
	signals := make(chan bool)

	initFetchers(requests, responses)
	go addFeeds(feedlist, requests)
	go processFeeds(total, responses, signals)
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

		var rsp, err = httpClient.Get(result.URL)

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
		f := &fetcher{id: i}
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

func addFeeds(feedlist []string, requests chan status) {

	for index, url := range feedlist {
		s := status{
			Index: index,
			URL:   url,
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

func readFile(feedfile string) ([]string, error) {

	var result []string

	f, err1 := os.Open(feedfile)
	if err1 != nil {
		return nil, err1
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			result = append(result, url)
		}
	}
	f.Close()

	return result, nil

}
