package httpd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	m "rakewire/model"
	"testing"
	"time"
)

const (
	feedURL = "http://localhost:5555/feed.xml"
)

var (
	ws     *Service
	feedID string
)

func TestMain(m *testing.M) {

	testDatabaseFile := "../../../test/httpd.db"

	cfg := db.Configuration{
		Location: testDatabaseFile,
	}
	testDatabase := &bolt.Database{}
	err := testDatabase.Open(&cfg)
	if err != nil {
		fmt.Printf("Cannot open database: %s\n", err.Error())
		os.Exit(1)
	}

	chErrors := make(chan error)
	ws = &Service{
		Database: testDatabase,
	}
	go ws.Start(&Configuration{
		Port:      4444,
		WebAppDir: "../../../test/public_html",
	}, chErrors)

	select {
	case err := <-chErrors:
		fmt.Printf("Error: %s\n", err.Error())
		testDatabase.Close()
		os.Remove(testDatabaseFile)
		os.Exit(1)
	case <-time.After(1 * time.Second):
		status := m.Run()
		ws.Stop()
		testDatabase.Close()
		os.Remove(testDatabaseFile)
		os.Exit(status)
	}

}

func TestStaticPaths(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeHTML, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "23", rsp.Header.Get(hContentLength))

	req = newRequest(mGet, "/humans.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "37", rsp.Header.Get(hContentLength))

	req = newRequest(mGet, "/hello/world.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "41", rsp.Header.Get(hContentLength))

}

func TestStatic404(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/favicon.ico")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))

	expectedText := "404 page not found\n"
	assert.Equal(t, 43 /* len(expectedText) */, int(rsp.ContentLength)) // gzip expands from 19 to 43
	bodyText, err := getZBodyAsString(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedText, bodyText)

}

func TestStaticRedirects(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "//")
	rsp, err := c.Do(req)
	assert.NotNil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusMovedPermanently, rsp.StatusCode)
	assert.Equal(t, "/", rsp.Header.Get("Location"))
	assert.Equal(t, 0, int(rsp.ContentLength))

	// #TODO:90 static redirect cannot be to /./
	// req = newRequest(mGet, "/index.html")
	// rsp, err = c.Do(req)
	// assert.NotNil(t, err)
	// assert.NotNil(t, rsp)
	// assert.Equal(t, 301, rsp.StatusCode)
	// assert.Equal(t, "/", rsp.Header.Get("Location"))

}

func TestFeedsPut(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPut, "/api/feeds")
	req.Header.Add(hContentType, mimeJSON)

	buf := bytes.Buffer{}
	feeds := m.NewFeeds()
	feed := m.NewFeed(feedURL)
	feedID = feed.ID
	feeds.Add(feed)
	feeds.Serialize(&buf)
	req.Body = ioutil.NopCloser(&buf)

	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))

	jsonRsp := SaveFeedsResponse{}
	err = json.NewDecoder(rsp.Body).Decode(&jsonRsp)
	assert.Nil(t, err)
	assert.NotNil(t, jsonRsp)
	assert.Equal(t, 1, jsonRsp.Count)

}

func TestFeedsPutNoContent(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPut, "/api/feeds")
	req.Header.Add(hContentType, mimeJSON)
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNoContent, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))

	// expectedText := "204 No Content\n"
	// assert.Equal(t, len(expectedText), int(rsp.ContentLength))
	// bodyText, err := getBodyAsString(rsp.Body)
	// assert.Nil(t, err)
	// assert.Equal(t, expectedText, bodyText)

}

func TestFeedsMethodNotAllowed(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPost, "/api/feeds")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusMethodNotAllowed, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))

	expectedText := "Method Not Allowed\n"
	assert.Equal(t, len(expectedText), int(rsp.ContentLength))
	bodyText, err := getBodyAsString(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedText, bodyText)

}

func TestFeedsGet(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))
	//assert.Equal(t, 98, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, rsp.ContentLength, n)
	feeds := m.NewFeeds()
	err = feeds.Deserialize(&buf)
	assert.Nil(t, err)
	assert.Equal(t, 1, feeds.Size())
	feed := feeds.Values[0]
	assert.Equal(t, feedURL, feed.URL)

}

func TestFeedGetByURL(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds?url=http%3A%2F%2Flocalhost%3A5555%2Ffeed.xml")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))
	//assert.Equal(t, 106, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, rsp.ContentLength, n)
	feed := m.Feed{}
	err = feed.Decode(buf.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost:5555/feed.xml", feed.URL)

}

func TestFeedGetByURL404(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds?url=http%3A%2F%2Flocalhost%3A5555%2Ffeed.XML")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))

	expectedText := "Not Found\n"
	assert.Equal(t, len(expectedText), int(rsp.ContentLength))
	bodyText, err := getBodyAsString(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedText, bodyText)

}

func TestFeedGetByID(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds/"+feedID)
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))
	//assert.Equal(t, 106, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, rsp.ContentLength, n)
	feed := m.Feed{}
	err = feed.Decode(buf.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost:5555/feed.xml", feed.URL)

}

func TestFeedGetByID404(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds/helloworld")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))

	expectedText := "Not Found\n"
	assert.Equal(t, len(expectedText), int(rsp.ContentLength))
	bodyText, err := getBodyAsString(rsp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedText, bodyText)

}

func getBodyAsString(r io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return string(buf.Bytes()[:]), err
}

func getZBodyAsString(r io.Reader) (string, error) {
	data, err := unzipReader(r)
	return string(data[:]), err
}

func newRequest(method string, path string) *http.Request {
	fmt.Printf("ws: %t %t\n", ws == nil, ws.listener == nil)
	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s%s", ws.listener.Addr(), path), nil)
	req.Header.Add(hAcceptEncoding, "gzip")
	return req
}

func newHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrNotSupported
		},
		Timeout: 1 * time.Second,
	}
}
