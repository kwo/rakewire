package feedparser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestAtom(t *testing.T) {
	//t.SkipNow()
	testFile(t, "../test/feed.xml")
}

func TestRSS(t *testing.T) {
	//t.SkipNow()
	testFile(t, "../test/wordpress.xml")
}

func TestAtomMalformed1(t *testing.T) {
	t.SkipNow()
	testURL(t, "http://www.quirksmode.org/blog/atom.xml")
}

func TestAtomMalformed2(t *testing.T) {
	t.SkipNow()
	testURL(t, "https://coreos.com/atom.xml")
}

func TestHtml1(t *testing.T) {
	t.SkipNow()
	testURL(t, "https://clusterhq.com/?cat=1&feed=rss2")
}

func TestRSSMalformed1(t *testing.T) {
	t.SkipNow()
	testURL(t, "http://feeds.feedburner.com/auth0")
}

func testFeed(t *testing.T, reader io.Reader) {

	feed, err := Parse(reader)
	assert.Nil(t, err)
	assert.Nil(t, feed)

}

func testFile(t *testing.T, filename string) {
	f, err := os.Open(filename)
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()
	testFeed(t, f)
}

func testURL(t *testing.T, url string) {
	rsp, err := http.Get(url)
	require.Nil(t, err)
	require.NotNil(t, rsp)
	defer rsp.Body.Close()
	testFeed(t, rsp.Body)
}
