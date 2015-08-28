package httpd

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	m "rakewire/model"
	"testing"
)

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
	assert404NotFoundNoGzip(t, rsp, err)

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
	assert404NotFoundNoGzip(t, rsp, err)

}
