package feedparser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAtom(t *testing.T) {
	f := testFile(t, "../../../test/feed/atomtest1.xml")

	assert.Equal(t, "atom", f.Flavor)
	assert.Equal(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml", f.ID)
	assert.Equal(t, time.Date(2005, time.November, 9, 11, 56, 34, 0, time.UTC), f.Updated)
	assert.Equal(t, "Sample Feed", f.Title.Text)
	assert.Equal(t, "text", f.Title.Type)

	assert.Equal(t, "For documentation <em>only</em>", f.Subtitle.Text)
	assert.Equal(t, "html", f.Subtitle.Type)

	assert.Equal(t, "http://example.org/icon.jpg", f.Icon)

	assert.Equal(t, "<p>Copyright 2005, Mark Pilgrim</p>", f.Rights.Text)
	assert.Equal(t, "html", f.Rights.Type)

	assert.Equal(t, "Sample Toolkit 4.0 (http://example.org/generator/)", f.Generator)

	assert.Equal(t, 2, len(f.Links))
	assert.Equal(t, "http://example.org/", f.Links["alternate"])
	assert.Equal(t, "http://www.example.org/atom10.xml", f.Links["self"])

	assert.Equal(t, 1, len(f.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", f.Authors[0])

	// entries

	assert.Equal(t, 1, len(f.Entries))
	e := f.Entries[0]
	assert.Equal(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml:3", e.ID)

	assert.Equal(t, "First entry title", e.Title.Text)
	assert.Equal(t, "text", e.Title.Type)

	assert.Equal(t, time.Date(2005, time.November, 9, 11, 56, 34, 0, time.UTC), e.Updated)
	assert.Equal(t, time.Date(2005, time.November, 9, 0, 23, 47, 0, time.UTC), e.Created)

	assert.Equal(t, 4, len(e.Links))
	assert.Equal(t, "http://example.org/entry/3", e.Links["alternate"])
	assert.Equal(t, "http://search.example.com/", e.Links["related"])
	assert.Equal(t, "http://toby.example.com/examples/atom10", e.Links["via"])
	assert.Equal(t, "http://www.example.com/movie.mp4", e.Links["enclosure"])

	assert.Equal(t, 2, len(e.Categories))
	assert.Equal(t, "Football", e.Categories[0])
	assert.Equal(t, "Basketball", e.Categories[1])

	assert.Equal(t, 1, len(e.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	assert.Equal(t, 2, len(e.Contributors))
	assert.Equal(t, "Joe <joe@example.org> (http://example.org/joe/)", e.Contributors[0])
	assert.Equal(t, "Sam <sam@example.org> (http://example.org/sam/)", e.Contributors[1])

	assert.Equal(t, "Watch out for nasty tricks", e.Summary.Text)
	assert.Equal(t, "text", e.Summary.Type)

	t.Log(e.Content.Text)
	assert.Equal(t, e.Content.Text, "Watch out for<span style=\"background-image: url(javascript:window.location=’http://example.org/’)\">nasty tricks</span>")
	assert.Equal(t, "html", e.Content.Type)

}

func TestRSS(t *testing.T) {
	//t.SkipNow()
	f := testFile(t, "../../../test/feed/wordpress.xml")

	assert.Equal(t, "rss2.0", f.Flavor)
	assert.Equal(t, "https://en.blog.wordpress.com/feed/", f.ID)
	assert.True(t, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(f.Updated))
	assert.Equal(t, "WordPress.com News", f.Title.Text)
	assert.Equal(t, "text", f.Title.Type)
	assert.Equal(t, "http://wordpress.com/", f.Generator)

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

func testFeed(t *testing.T, reader io.Reader) *Feed {
	feed, err := Parse(reader)
	require.Nil(t, err)
	require.NotNil(t, feed)
	return feed
}

func testFile(t *testing.T, filename string) *Feed {
	f, err := os.Open(filename)
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()
	return testFeed(t, f)
}

func testURL(t *testing.T, url string) *Feed {
	rsp, err := http.Get(url)
	require.Nil(t, err)
	require.NotNil(t, rsp)
	defer rsp.Body.Close()
	return testFeed(t, rsp.Body)
}
