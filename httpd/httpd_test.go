package httpd

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"rakewire.com/db/bolt"
	"rakewire.com/model"
	"testing"
	"time"
)

var (
	ws *Httpd
)

func TestMain(m *testing.M) {

	testDatabaseFile := "../test/test.db"

	cfg := model.DatabaseConfiguration{
		Location: testDatabaseFile,
	}
	testDatabase := bolt.Database{}
	err := testDatabase.Open(&cfg)
	if err != nil {
		fmt.Printf("Cannot open database: %s\n", err.Error())
		os.Exit(1)
	}

	chErrors := make(chan error)
	ws = &Httpd{
		Database: &testDatabase,
	}
	go ws.Start(&model.HttpdConfiguration{
		Port:      4444,
		WebAppDir: "../test/public_html",
	}, chErrors)

	// TODO: probably need to wait before jumping to default case

	select {
	case err := <-chErrors:
		fmt.Printf("Cannot start httpd: %s\n", err.Error())
		testDatabase.Close()
		os.Remove(testDatabaseFile)
		os.Exit(1)
	default:
		status := m.Run()
		ws.Stop()
		testDatabase.Close()
		os.Remove(testDatabaseFile)
		os.Exit(status)
	}

}

func TestStaticPaths(t *testing.T) {

	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("/")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeHTML, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "23", rsp.Header.Get(hContentLength))

	req = getRequest("/humans.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "37", rsp.Header.Get(hContentLength))

	req = getRequest("/hello/world.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeText, rsp.Header.Get(hContentType))
	assert.Equal(t, "gzip", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, "41", rsp.Header.Get(hContentLength))

}

func Test404(t *testing.T) {

	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("/favicon.ico")
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

	c := getHTTPClient()

	req := getRequest("//")
	rsp, err := c.Do(req)
	assert.NotNil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusMovedPermanently, rsp.StatusCode)
	assert.Equal(t, "/", rsp.Header.Get("Location"))
	assert.Equal(t, 0, int(rsp.ContentLength))

	// TODO: static redirect cannot be to /./
	// req = getRequest("/index.html")
	// rsp, err = c.Do(req)
	// assert.NotNil(t, err)
	// assert.NotNil(t, rsp)
	// assert.Equal(t, 301, rsp.StatusCode)
	// assert.Equal(t, "/", rsp.Header.Get("Location"))

}

func TestFeedGet(t *testing.T) {

	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("/api/feeds")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	assert.Equal(t, mimeJSON, rsp.Header.Get(hContentType))
	assert.Equal(t, "", rsp.Header.Get(hContentEncoding))
	assert.Equal(t, 5, int(rsp.ContentLength))

}

func TestFeedPutNoContent(t *testing.T) {

	require.NotNil(t, ws)

	c := getHTTPClient()

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

	// jsonRsp := SaveFeedsResponse{}
	// err = json.NewDecoder(rsp.Body).Decode(&jsonRsp)
	// assert.Nil(t, err)
	// assert.NotNil(t, jsonRsp)
	// assert.Equal(t, 0, jsonRsp.Count)

}

func TestFeedPost(t *testing.T) {

	require.NotNil(t, ws)

	c := getHTTPClient()

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

func getBodyAsString(r io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return string(buf.Bytes()[:]), err
}

func getZBodyAsString(r io.Reader) (string, error) {
	data, err := unzipReader(r)
	return string(data[:]), err
}

func getRequest(path string) *http.Request {
	return newRequest(mGet, path)
}

func newRequest(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s%s", ws.listener.Addr(), path), nil)
	req.Header.Add(hAcceptEncoding, "gzip")
	return req
}

func getHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrNotSupported
		},
		Timeout: 1 * time.Second,
	}
}
