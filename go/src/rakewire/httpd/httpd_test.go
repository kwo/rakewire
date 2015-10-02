package httpd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	"strconv"
	"testing"
	"time"
)

const (
	feedURL          = "http://localhost:5555/feed.xml"
	hLastModified    = "Last-Modified"
	hIfModifiedSince = "If-Modified-Since"
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
		Port: 4444,
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
	assertHTML(t, rsp, err)

	req = newRequest(mGet, "/robots.txt")
	rsp, err = c.Do(req)
	assertText(t, rsp, err)

	req = newRequest(mGet, "/fonts/MaterialIcons-Regular.woff2")
	rsp, err = c.Do(req)
	assert200OK(t, rsp, err, "application/octet-stream")

}

func TestHTML5Paths(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/")
	rsp, err := c.Do(req)
	assertHTML(t, rsp, err)
	abody, err := ioutil.ReadAll(rsp.Body)
	require.Nil(t, err)
	require.NotNil(t, abody)
	assert.Equal(t, strconv.Itoa(len(abody)), rsp.Header.Get(hContentLength))

	req = newRequest(mGet, "/route")
	rsp, err = c.Do(req)
	assertHTML(t, rsp, err)
	body, err := ioutil.ReadAll(rsp.Body)
	require.Nil(t, err)
	require.NotNil(t, body)
	assert.Equal(t, strconv.Itoa(len(abody)), rsp.Header.Get(hContentLength))
	assert.Equal(t, abody, body)

	// also multi-level paths
	req = newRequest(mGet, "/route/route")
	rsp, err = c.Do(req)
	assertHTML(t, rsp, err)
	body, err = ioutil.ReadAll(rsp.Body)
	require.Nil(t, err)
	require.NotNil(t, body)
	assert.Equal(t, strconv.Itoa(len(abody)), rsp.Header.Get(hContentLength))
	assert.Equal(t, abody, body)

	// include paths with uuids
	req = newRequest(mGet, "/feed/bf24f476-5899-11e5-af27-5cf938992b62/log")
	rsp, err = c.Do(req)
	assertHTML(t, rsp, err)
	body, err = ioutil.ReadAll(rsp.Body)
	require.Nil(t, err)
	require.NotNil(t, body)
	assert.Equal(t, strconv.Itoa(len(abody)), rsp.Header.Get(hContentLength))
	assert.Equal(t, abody, body)

	req = newRequest(mGet, "/home")
	rsp, err = c.Do(req)
	assertHTML(t, rsp, err)
	body, err = ioutil.ReadAll(rsp.Body)
	require.Nil(t, err)
	require.NotNil(t, body)
	assert.Equal(t, strconv.Itoa(len(abody)), rsp.Header.Get(hContentLength))
	assert.Equal(t, abody, body)

	// only all lowercase
	req = newRequest(mGet, "/Route")
	rsp, err = c.Do(req)
	assert404NotFound(t, rsp, err)

	// only a-z, not dot or slashes
	req = newRequest(mGet, "/route.html")
	rsp, err = c.Do(req)
	assert404NotFound(t, rsp, err)

	// no paths with extensions
	req = newRequest(mGet, "/route/route.txt")
	rsp, err = c.Do(req)
	assert404NotFound(t, rsp, err)

}

func TestStatic404(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/favicon.ico")
	rsp, err := c.Do(req)
	assert404NotFound(t, rsp, err)

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

	// static redirect cannot be to /./
	req = newRequest(mGet, "/index.html")
	rsp, err = c.Do(req)
	assert.NotNil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 301, rsp.StatusCode)
	assert.Equal(t, "/", rsp.Header.Get("Location"))

}

func TestNotModified(t *testing.T) {

	require.NotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/")
	rsp, err := c.Do(req)
	assertHTML(t, rsp, err)
	lastModified := rsp.Header.Get(hLastModified)

	req = newRequest(mGet, "/")
	req.Header.Add(hIfModifiedSince, lastModified)
	rsp, err = c.Do(req)
	assert304NotModified(t, rsp, err)

	req = newRequest(mGet, "/robots.txt")
	rsp, err = c.Do(req)
	assertText(t, rsp, err)
	lastModified = rsp.Header.Get(hLastModified)

	req = newRequest(mGet, "/robots.txt")
	req.Header.Add(hIfModifiedSince, lastModified)
	rsp, err = c.Do(req)
	assert304NotModified(t, rsp, err)

	req = newRequest(mGet, "/fonts/MaterialIcons-Regular.woff2")
	rsp, err = c.Do(req)
	assert200OK(t, rsp, err, "application/octet-stream")
	lastModified = rsp.Header.Get(hLastModified)

	req = newRequest(mGet, "/fonts/MaterialIcons-Regular.woff2")
	req.Header.Add(hIfModifiedSince, lastModified)
	rsp, err = c.Do(req)
	assert304NotModified(t, rsp, err)

}

func assertHTML(t *testing.T, rsp *http.Response, err error) {
	assert200OK(t, rsp, err, mimeHTML)
}

func assertText(t *testing.T, rsp *http.Response, err error) {
	assert200OK(t, rsp, err, mimeText)
}

func assertJSONAPI(t *testing.T, rsp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
}

func assert200OK(t *testing.T, rsp *http.Response, err error, mimeType string) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeType, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, hAcceptEncoding, rsp.Header.Get(hVary))
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
	assert.NotEqual(t, "", rsp.Header.Get(hLastModified))
}

func assert200OKAPI(t *testing.T, rsp *http.Response, err error, mimeType string) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeType, rsp.Header.Get(hContentType))
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
	assert.NotEqual(t, "", rsp.Header.Get(hLastModified))
}

func assert304NotModified(t *testing.T, rsp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotModified, rsp.StatusCode)
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
}

func assert404NotFound(t *testing.T, rsp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, hAcceptEncoding, rsp.Header.Get(hVary))
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
}

func assert404NotFoundAPI(t *testing.T, rsp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, vNoCache, rsp.Header.Get(hCacheControl))
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

func unzipReader(data io.Reader) ([]byte, error) {

	r, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}

	var uncompressedData bytes.Buffer
	if _, err = io.Copy(&uncompressedData, r); err != nil {
		return nil, err
	}

	return uncompressedData.Bytes(), nil

}
