package httpd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"rakewire.com/db"
	"rakewire.com/model"
	"testing"
	"time"
)

var (
	ws *Httpd
)

func TestMain(m *testing.M) {
	ws = &Httpd{
		Database: &FakeDb{},
	}
	go ws.Start(&model.HttpdConfiguration{
		Port:      4444,
		WebAppDir: "../test/public_html",
	})
	// TODO: how to catch error from Start
	// if err != nil {
	// 	fmt.Printf("Cannot start httpd: %s\n", err.Error())
	// 	os.Exit(1)
	// }
	status := m.Run()
	ws.Stop()
	os.Exit(status)
}

func TestStaticPaths(t *testing.T) {
	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("/")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 200, rsp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", rsp.Header.Get("Content-Type"))
	assert.Equal(t, "gzip", rsp.Header.Get("Content-Encoding"))
	assert.Equal(t, "23", rsp.Header.Get("Content-Length"))

	req = getRequest("/humans.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 200, rsp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", rsp.Header.Get("Content-Type"))
	assert.Equal(t, "gzip", rsp.Header.Get("Content-Encoding"))
	assert.Equal(t, "37", rsp.Header.Get("Content-Length"))

	req = getRequest("/hello/world.txt")
	rsp, err = c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 200, rsp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", rsp.Header.Get("Content-Type"))
	assert.Equal(t, "gzip", rsp.Header.Get("Content-Encoding"))
	assert.Equal(t, "41", rsp.Header.Get("Content-Length"))

}

func TestStaticRedirects(t *testing.T) {
	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("//")
	rsp, err := c.Do(req)
	assert.NotNil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 301, rsp.StatusCode)
	assert.Equal(t, "/", rsp.Header.Get("Location"))

	// TODO: static redirect cannot be to /./
	// req = getRequest("/index.html")
	// rsp, err = c.Do(req)
	// assert.NotNil(t, err)
	// assert.NotNil(t, rsp)
	// assert.Equal(t, 301, rsp.StatusCode)
	// assert.Equal(t, "/", rsp.Header.Get("Location"))

}

func TestAPIPath(t *testing.T) {
	require.NotNil(t, ws)

	c := getHTTPClient()

	req := getRequest("/api")
	rsp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 200, rsp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", rsp.Header.Get("Content-Type"))
	assert.Equal(t, "", rsp.Header.Get("Content-Encoding"))
	assert.Equal(t, "27", rsp.Header.Get("Content-Length"))

}

func getRequest(path string) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s%s", ws.listener.Addr(), path), nil)
	req.Header.Add("Accept-Encoding", "gzip")
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

type FakeDb struct {
}

func (z *FakeDb) GetFeeds() (map[string]*db.FeedInfo, error) {
	return nil, nil
}

func (z *FakeDb) SaveFeeds([]*db.FeedInfo) (int, error) {
	return 0, nil
}
