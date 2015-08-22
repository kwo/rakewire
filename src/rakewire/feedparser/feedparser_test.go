package feedparser

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestelementsAtom(t *testing.T) {

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "feed"}})

	require.Equal(t, 1, elements.Level())
	assert.True(t, elements.IsStackFeed())

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "entry"}})
	require.Equal(t, 2, elements.Level())
	assert.True(t, elements.IsStackEntry())
	assert.True(t, elements.IsStackFeed(1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "id"}})
	require.Equal(t, 3, elements.Level())
	assert.True(t, elements.IsStackEntry(1))

}

func TestelementsRSS(t *testing.T) {

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "rss"}})

	require.Equal(t, 1, elements.Level())
	assert.False(t, elements.IsStackFeed())

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "channel"}})
	require.Equal(t, 2, elements.Level())
	assert.True(t, elements.IsStackFeed())

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "item"}})
	require.Equal(t, 3, elements.Level())
	assert.True(t, elements.IsStackEntry())
	assert.True(t, elements.IsStackFeed(1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "guid"}})
	require.Equal(t, 4, elements.Level())
	assert.True(t, elements.IsStackEntry(1))

}

func TestAtom(t *testing.T) {
	//t.SkipNow()
	f := testFile(t, "../../../test/feed/atomtest1.xml")

	assert.Equal(t, "atom", f.Flavor)
	assert.Equal(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml", f.ID)
	assert.True(t, time.Date(2005, time.November, 11, 11, 56, 34, 0, time.UTC).Equal(f.Updated))
	assert.Equal(t, "Sample Feed", f.Title)

	assert.Equal(t, "For documentation <em>only</em>", f.Subtitle)

	assert.Equal(t, "http://example.org/icon.jpg", f.Icon)

	assert.Equal(t, "<p>Copyright 2005, Mark Pilgrim</p>", f.Rights)

	assert.Equal(t, "Sample Toolkit 4.0 (http://example.org/generator/)", f.Generator)

	require.Equal(t, 2, len(f.Links))
	assert.Equal(t, "http://example.org/", f.Links["alternate"])
	assert.Equal(t, "http://www.example.org/atom10.xml", f.Links["self"])

	require.Equal(t, 1, len(f.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", f.Authors[0])

	// entries

	require.Equal(t, 2, len(f.Entries))
	e := f.Entries[0]

	assert.Equal(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml:3", e.ID)
	assert.Equal(t, "First entry title", e.Title)
	assert.Equal(t, time.Date(2005, time.November, 9, 11, 56, 34, 0, time.UTC), e.Updated)
	assert.Equal(t, time.Date(2005, time.November, 9, 0, 23, 47, 0, time.UTC), e.Created)

	require.Equal(t, 4, len(e.Links))
	assert.Equal(t, "http://example.org/entry/3", e.Links["alternate"])
	assert.Equal(t, "http://search.example.com/", e.Links["related"])
	assert.Equal(t, "http://toby.example.com/examples/atom10", e.Links["via"])
	assert.Equal(t, "http://www.example.com/movie.mp4", e.Links["enclosure"])

	require.Equal(t, 2, len(e.Categories))
	assert.Equal(t, "Football", e.Categories[0])
	assert.Equal(t, "Basketball", e.Categories[1])

	require.Equal(t, 1, len(e.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	require.Equal(t, 2, len(e.Contributors))
	assert.Equal(t, "Joe <joe@example.org> (http://example.org/joe/)", e.Contributors[0])
	assert.Equal(t, "Sam <sam@example.org> (http://example.org/sam/)", e.Contributors[1])

	assert.Equal(t, "Watch out for nasty tricks", e.Summary)
	assert.Equal(t, "Watch out for<span style=\"background-image: url(javascript:window.location=’http://example.org/’)\">nasty tricks</span>", e.Content)

	e = f.Entries[1]
	assert.Equal(t, "tag:feedparser.org,2005-11-11:/docs/examples/atom11.xml:1", e.ID)
	assert.Equal(t, "Second entry title", e.Title)
	assert.True(t, time.Date(2005, time.November, 11, 11, 56, 34, 0, time.UTC).Equal(f.Updated))
	assert.True(t, e.Created.IsZero())

	require.Equal(t, 0, len(e.Links))

	require.Equal(t, 0, len(e.Categories))

	require.Equal(t, 1, len(e.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	require.Equal(t, 0, len(e.Contributors))

	assert.Equal(t, "Test content as text", e.Summary)
	assert.Equal(t, "Test content", e.Content)

}

func TestRSS(t *testing.T) {
	//t.SkipNow()
	f := testFile(t, "../../../test/feed/wordpress.xml")

	assert.Equal(t, "rss2.0", f.Flavor)
	assert.Equal(t, "https://en.blog.wordpress.com/feed/", f.ID)
	assert.True(t, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(f.Updated))
	assert.Equal(t, "WordPress.com News", f.Title)
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
	p := &Parser{}
	feed, err := p.Parse(reader)
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
