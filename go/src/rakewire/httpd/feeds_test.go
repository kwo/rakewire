package httpd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	m "rakewire/model"
	"testing"
)

func TestFeedsPut(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPut, "/api/feeds")
	req.Header.Add(hContentType, mimeJSON)

	buf := bytes.Buffer{}
	var feeds []*m.Feed
	feed := m.NewFeed(feedURL)
	feedID = feed.ID
	feeds = append(feeds, feed)
	err := serializeFeeds(feeds, &buf)
	assertNoError(t, err)
	assertEqual(t, 1, len(feeds))
	req.Body = ioutil.NopCloser(&buf)

	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertJSONAPI(t, rsp, err)

	count, err := deserializeSaveFeedsResponse(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, 1, count)

}

func TestFeedsPutNoContent(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPut, "/api/feeds")
	req.Header.Add(hContentType, mimeJSON)
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusNoContent, rsp.StatusCode)
	assertEqual(t, mimeText, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))

	// expectedText := "204 No Content\n"
	// assertEqual(t, len(expectedText), int(rsp.ContentLength))
	// bodyText, err := getBodyAsString(rsp.Body)
	// assertNoError(t, err)
	// assertEqual(t, expectedText, bodyText)

}

func TestFeedsMethodNotAllowed(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mPost, "/api/feeds")
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusMethodNotAllowed, rsp.StatusCode)
	assertEqual(t, mimeText, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))

	expectedText := "Method Not Allowed\n"
	assertEqual(t, len(expectedText), int(rsp.ContentLength))
	bodyText, err := getBodyAsString(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, expectedText, bodyText)

}

func TestFeedsGet(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds")
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusOK, rsp.StatusCode)
	assertEqual(t, mimeJSON, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))
	//assertEqual(t, 98, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, rsp.ContentLength, n)
	feeds, err := deserializeFeeds(&buf)
	assertNoError(t, err)
	assertEqual(t, 1, len(feeds))
	feed := feeds[0]
	assertEqual(t, feedURL, feed.URL)

}

func TestFeedsGetNext(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds/next")
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusOK, rsp.StatusCode)
	assertEqual(t, mimeJSON, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))
	//assertEqual(t, 98, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, rsp.ContentLength, n)
	feeds, err := deserializeFeeds(&buf)
	assertNoError(t, err)
	assertEqual(t, 1, len(feeds))
	feed := feeds[0]
	assertEqual(t, feedURL, feed.URL)

}

func TestFeedGetByURL(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds?url=http%3A%2F%2Flocalhost%3A5555%2Ffeed.xml")
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusOK, rsp.StatusCode)
	assertEqual(t, mimeJSON, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))
	//assertEqual(t, 106, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	t.Logf("feedByURL: %s\n", string(buf.Bytes()))
	assertNoError(t, err)
	assertEqual(t, rsp.ContentLength, n)
	feed, err := deserializeFeed(buf.Bytes())
	assertNoError(t, err)
	assertNotNil(t, feed)
	assertEqual(t, "http://localhost:5555/feed.xml", feed.URL)

}

func TestFeedGetByURL404(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds?url=http%3A%2F%2Flocalhost%3A5555%2Ffeed.XML")
	rsp, err := c.Do(req)
	assert404NotFoundAPI(t, rsp, err)

	expectedText := "Not Found\n"
	assertEqual(t, len(expectedText), int(rsp.ContentLength))
	bodyText, err := getBodyAsString(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, expectedText, bodyText)

}

func TestFeedGetByID(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds/"+feedID)
	rsp, err := c.Do(req)
	assertNoError(t, err)
	assertNotNil(t, rsp)
	assertEqual(t, http.StatusOK, rsp.StatusCode)
	assertEqual(t, mimeJSON, rsp.Header.Get(hContentType))
	assertEqual(t, "", rsp.Header.Get(hContentEncoding))
	//assertEqual(t, 106, int(rsp.ContentLength))

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(rsp.Body)
	assertNoError(t, err)
	assertEqual(t, rsp.ContentLength, n)
	feed, err := deserializeFeed(buf.Bytes())
	assertNoError(t, err)
	assertNotNil(t, feed)
	assertEqual(t, "http://localhost:5555/feed.xml", feed.URL)

}

func TestFeedGetByID404(t *testing.T) {

	assertNotNil(t, ws)

	c := newHTTPClient()

	req := newRequest(mGet, "/api/feeds/helloworld")
	rsp, err := c.Do(req)
	assert404NotFoundAPI(t, rsp, err)

}
