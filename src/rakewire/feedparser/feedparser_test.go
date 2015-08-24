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

func TestElementsAtom(t *testing.T) {

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "feed"}})

	require.Equal(t, 1, elements.Level())
	assert.True(t, elements.IsStackFeed(flavorAtom, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "entry"}})
	require.Equal(t, 2, elements.Level())
	assert.True(t, elements.IsStackEntry(flavorAtom, 0))
	assert.True(t, elements.IsStackFeed(flavorAtom, 1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "id"}})
	require.Equal(t, 3, elements.Level())
	assert.True(t, elements.IsStackEntry(flavorAtom, 1))

}

func TestElementsRSS(t *testing.T) {

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "rss"}})

	require.Equal(t, 1, elements.Level())
	assert.False(t, elements.IsStackFeed(flavorRSS, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "channel"}})
	require.Equal(t, 2, elements.Level())
	assert.True(t, elements.IsStackFeed(flavorRSS, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "item"}})
	require.Equal(t, 3, elements.Level())
	assert.True(t, elements.IsStackEntry(flavorRSS, 0))
	assert.True(t, elements.IsStackFeed(flavorRSS, 1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "guid"}})
	require.Equal(t, 4, elements.Level())
	assert.True(t, elements.IsStackEntry(flavorRSS, 1))

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
	assert.True(t, e.Created.Equal(e.Updated))

	require.Equal(t, 0, len(e.Links))

	require.Equal(t, 0, len(e.Categories))

	require.Equal(t, 1, len(e.Authors))
	assert.Equal(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	require.Equal(t, 0, len(e.Contributors))

	assert.Equal(t, "Test content as text", e.Summary)
	assert.Equal(t, "Test content", e.Content)

}

func TestRSS(t *testing.T) {

	f := testFile(t, "../../../test/feed/wordpress.xml")

	assert.Equal(t, "rss2.0", f.Flavor)
	assert.Equal(t, "https://en.blog.wordpress.com/feed/", f.ID)
	assert.True(t, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(f.Updated))
	assert.Equal(t, "WordPress.com News", f.Title)
	assert.Equal(t, "The latest news on WordPress.com and the WordPress community.", f.Subtitle)
	assert.Equal(t, "http://wordpress.com/", f.Generator)
	assert.Equal(t, "https://secure.gravatar.com/blavatar/e6392390e3bcfadff3671c5a5653d95b?s=96&d=https%3A%2F%2Fs2.wp.com%2Fi%2Fbuttonw-com.png", f.Icon)

	require.Equal(t, 4, len(f.Links))
	assert.Equal(t, "https://en.blog.wordpress.com", f.Links["alternate"])
	assert.Equal(t, "https://en.blog.wordpress.com/feed/", f.Links["self"])
	assert.Equal(t, "https://en.blog.wordpress.com/osd.xml", f.Links["search"])
	assert.Equal(t, "https://en.blog.wordpress.com/?pushpress=hub", f.Links["hub"])

	require.Equal(t, 10, len(f.Entries))
	e := f.Entries[0]

	assert.Equal(t, "http://en.blog.wordpress.com/?p=31505", e.ID)
	assert.Equal(t, "New Theme: Libre", e.Title)
	assert.True(t, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(e.Updated))
	assert.True(t, e.Created.Equal(e.Updated))

	require.Equal(t, 1, len(e.Categories))
	require.Equal(t, "Themes", e.Categories[0])

	require.Equal(t, 1, len(e.Links))
	require.Equal(t, "https://en.blog.wordpress.com/2015/07/02/libre/", e.Links["alternate"])

	assert.Equal(t, "<em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site.<img alt=\"\" border=\"0\" src=\"https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1\" width=\"1\" height=\"1\" />", e.Summary)
	assert.Equal(t, "<p>Happy Theme Thursday, all! Today I&#8217;m happy to introduce <em>Libre</em>, a new free theme designed by <a href=\"http://carolinethemes.com/\">yours truly</a>.</p>\n<h3><a href=\"http://wordpress.com/themes/libre/\">Libre</a></h3>\n<p><a href=\"http://wordpress.com/themes/libre/\"><img class=\"aligncenter size-full wp-image-31516\" src=\"https://en-blog.files.wordpress.com/2015/06/librelg.png?w=635&#038;h=476\" alt=\"Libre\" width=\"635\" height=\"476\" /></a></p>\n<p><em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site. The main navigation bar stays fixed to the top while your visitors read, keeping your most important content at hand. At the bottom of your site, three footer widget areas give your secondary content a comfortable home.</p>\n<p>Customize <em>Libre</em> with a logo or a header image to make it your own, or use one of two custom templates &#8212; including a full-width template\u00a0with no sidebar &#8212; to change up the look of your pages. <em>Libre</em> sports a clean, responsive design that works seamlessly on screens of any size.</p>\n<p><img class=\"aligncenter size-full wp-image-31517\" src=\"https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg?w=635&#038;h=252\" alt=\"Responsive design\" width=\"635\" height=\"252\" /></p>\n<p>Read more about <em>Libre</em> on the <a href=\"https://wordpress.com/themes/libre/\">Theme Showcase</a>, or activate it on your site from <em>Appearance → Themes</em>!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel=\"nofollow\" href=\"http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31505/\"><img alt=\"\" border=\"0\" src=\"http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31505/\" /></a> <img alt=\"\" border=\"0\" src=\"https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1\" width=\"1\" height=\"1\" />", e.Content)

}

// func TestContentType(t *testing.T) {
// 	t.SkipNow()
// 	f, err := os.Open("../../../test/feed/intertwingly-net-blog.html")
// 	require.Nil(t, err)
// 	require.NotNil(t, f)
// 	defer f.Close()
//
// 	body, err := ioutil.ReadAll(f)
// 	require.Nil(t, err)
// 	require.NotNil(t, body)
//
// 	ct := http.DetectContentType(body)
// 	_, name, certain := charset.DetermineEncoding(body, "application/xml")
// 	//t.Logf("encoding: %s", encoding)
//
// 	assert.Equal(t, ct, "text/html; charset=utf-8")
// 	assert.Equal(t, "utf-8", name)
// 	assert.True(t, certain)
//
// }

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
	testURL(t, "http://intertwingly.net/blog/")
}

func TestRSSMalformed1(t *testing.T) {
	t.SkipNow()
	testURL(t, "http://feeds.feedburner.com/auth0")
}

func testFeed(t *testing.T, reader io.ReadCloser, contentType string) *Feed {
	p := NewParser()
	feed, err := p.Parse(reader, contentType)
	require.Nil(t, err)
	require.NotNil(t, feed)
	return feed
}

func testFile(t *testing.T, filename string) *Feed {
	f, err := os.Open(filename)
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()
	return testFeed(t, f, "")
}

func testURL(t *testing.T, url string) *Feed {
	rsp, err := http.Get(url)
	contentType := rsp.Header.Get("Content-Type")
	require.Nil(t, err)
	require.NotNil(t, rsp)
	defer rsp.Body.Close()
	return testFeed(t, rsp.Body, contentType)
}
