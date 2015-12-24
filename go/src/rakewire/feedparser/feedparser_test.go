package feedparser

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	log.SetOutput(ioutil.Discard)

	m.Run()

}

func TestLinkWithoutRel(t *testing.T) {

	// reference https://www.tbray.org/ongoing/ongoing.atom

	t.Parallel()

	atom := `
	<?xml version='1.0' encoding='UTF-8'?>
	<feed xmlns='http://www.w3.org/2005/Atom'
	      xmlns:thr='http://purl.org/syndication/thread/1.0'
	      xml:lang='en-us'>
	 <title>ongoing by Tim Bray</title>
	 <link rel='hub' href='http://pubsubhubbub.appspot.com/' />
	 <id>https://www.tbray.org/ongoing/</id>
	 <link href='https://www.tbray.org/ongoing/' />
	 <link rel='self' href='https://www.tbray.org/ongoing/ongoing.atom' />
	 <link rel='replies'       thr:count='101'       href='https://www.tbray.org/ongoing/comments.atom' />
	 <logo>rsslogo.jpg</logo>
	 <icon>/favicon.ico</icon>
	 <updated>2015-12-24T01:04:28-08:00</updated>
	 <author><name>Tim Bray</name></author>
	 <subtitle>ongoing fragmented essay by Tim Bray</subtitle>
	</feed>
	`

	p := NewParser()
	feed, err := p.Parse(ioutil.NopCloser(strings.NewReader(atom)), "")
	if err != nil {
		t.Fatalf("Cannot parse Aton feed: %s", err.Error())
	}
	if feed == nil {
		t.Fatal("Nil feed object")
	}

	if count := len(feed.Entries); count != 0 {
		t.Fatalf("bad entry count, expected: %d, actual: %d", 0, count)
	}

	htmlLink := feed.GetLinkAlternate()
	if htmlLink != "https://www.tbray.org/ongoing/" {
		t.Errorf("Links without rel not being intrepreted as alternate link: expected: %s, actual: %s", "https://www.tbray.org/ongoing/", htmlLink)
	}

	selfLink := feed.GetLinkSelf()
	if selfLink != "https://www.tbray.org/ongoing/ongoing.atom" {
		t.Errorf("Links without rel not being intrepreted as alternate link: expected: %s, actual: %s", "https://www.tbray.org/ongoing/ongoing.atom", htmlLink)
	}

}

func TestRSSPerson(t *testing.T) {

	t.Parallel()

	rss := `
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
		<channel>
			<title>Blog</title>
			<link>http://localhost/</link>
			<description>blog desc</description>
			<item>
				<author>a.ut-hor@author.com (Arthur A. Author)</author>
				<author>other guy</author>
				<title>The Title</title>
				<link>http://localhost/?ts=a89ea735</link>
				<guid>http://localhost/?ts=a89ea735</guid>
			</item>
		</channel>
	</rss>
	`

	p := NewParser()
	feed, err := p.Parse(ioutil.NopCloser(strings.NewReader(rss)), "")
	if err != nil {
		t.Fatalf("Cannot parse RSS feed: %s", err.Error())
	}
	if feed == nil {
		t.Fatal("Nil feed object")
	}

	if count := len(feed.Entries); count != 1 {
		t.Fatalf("bad entry count, expected: %d, actual: %d", 1, count)
	}

	expectedAuthorCount := 2
	if count := len(feed.Entries[0].Authors); count != expectedAuthorCount {
		t.Fatalf("bad author count, expected: %d, actual: %d", expectedAuthorCount, count)
	}

	expectedAuthor := "Arthur A. Author <a.ut-hor@author.com>"
	if author := feed.Entries[0].Authors[0]; author != expectedAuthor {
		t.Errorf("author malformed, expected: %s, actual: %s", expectedAuthor, author)
	}

	expectedAuthor = "other guy"
	if author := feed.Entries[0].Authors[1]; author != expectedAuthor {
		t.Errorf("author malformed, expected: %s, actual: %s", expectedAuthor, author)
	}

}

func TestElementsAtom(t *testing.T) {

	t.Parallel()

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "feed"}})

	assertEqual(t, 1, elements.Level())
	assertEqual(t, true, elements.IsStackFeed(flavorAtom, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "entry"}})
	assertEqual(t, 2, elements.Level())
	assertEqual(t, true, elements.IsStackEntry(flavorAtom, 0))
	assertEqual(t, true, elements.IsStackFeed(flavorAtom, 1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsAtom, Local: "id"}})
	assertEqual(t, 3, elements.Level())
	assertEqual(t, true, elements.IsStackEntry(flavorAtom, 1))

}

func TestElementsRSS(t *testing.T) {

	t.Parallel()

	elements := &elements{}
	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "rss"}})

	assertEqual(t, 1, elements.Level())
	assertEqual(t, false, elements.IsStackFeed(flavorRSS, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "channel"}})
	assertEqual(t, 2, elements.Level())
	assertEqual(t, true, elements.IsStackFeed(flavorRSS, 0))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "item"}})
	assertEqual(t, 3, elements.Level())
	assertEqual(t, true, elements.IsStackEntry(flavorRSS, 0))
	assertEqual(t, true, elements.IsStackFeed(flavorRSS, 1))

	elements.stack = append(elements.stack, &element{name: xml.Name{Space: nsRSS, Local: "guid"}})
	assertEqual(t, 4, elements.Level())
	assertEqual(t, true, elements.IsStackEntry(flavorRSS, 1))

}

func TestAtom(t *testing.T) {

	t.Parallel()

	f := testFile(t, "../../../test/feed/atomtest1.xml")

	assertEqual(t, "atom", f.Flavor)
	assertEqual(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml", f.ID)
	assertEqual(t, true, time.Date(2005, time.November, 11, 11, 56, 34, 0, time.UTC).Equal(f.Updated))
	assertEqual(t, "Sample Feed", f.Title)

	assertEqual(t, "For documentation <em>only</em>", f.Subtitle)

	assertEqual(t, "http://example.org/icon.jpg", f.Icon)

	assertEqual(t, "<p>Copyright 2005, Mark Pilgrim</p>", f.Rights)

	assertEqual(t, "Sample Toolkit 4.0 (http://example.org/generator/)", f.Generator)

	assertEqual(t, 2, len(f.Links))
	assertEqual(t, "http://example.org/", f.Links["alternate"])
	assertEqual(t, "http://www.example.org/atom10.xml", f.Links["self"])

	assertEqual(t, 1, len(f.Authors))
	assertEqual(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", f.Authors[0])

	// entries

	assertEqual(t, 2, len(f.Entries))
	e := f.Entries[0]

	assertEqual(t, "tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml:3", e.ID)
	assertEqual(t, "First entry title", e.Title)
	assertEqual(t, time.Date(2005, time.November, 9, 11, 56, 34, 0, time.UTC), e.Updated)
	assertEqual(t, time.Date(2005, time.November, 9, 0, 23, 47, 0, time.UTC), e.Created)

	assertEqual(t, 4, len(e.Links))
	assertEqual(t, "http://example.org/entry/3", e.Links["alternate"])
	assertEqual(t, "http://search.example.com/", e.Links["related"])
	assertEqual(t, "http://toby.example.com/examples/atom10", e.Links["via"])
	assertEqual(t, "http://www.example.com/movie.mp4", e.Links["enclosure"])

	assertEqual(t, 2, len(e.Categories))
	assertEqual(t, "Football", e.Categories[0])
	assertEqual(t, "Basketball", e.Categories[1])

	assertEqual(t, 1, len(e.Authors))
	assertEqual(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	assertEqual(t, 2, len(e.Contributors))
	assertEqual(t, "Joe <joe@example.org> (http://example.org/joe/)", e.Contributors[0])
	assertEqual(t, "Sam <sam@example.org> (http://example.org/sam/)", e.Contributors[1])

	assertEqual(t, "Watch out for nasty tricks", e.Summary)
	assertEqual(t, "Watch out for<span style=\"background-image: url(javascript:window.location=’http://example.org/’)\">nasty tricks</span>", e.Content)

	e = f.Entries[1]
	assertEqual(t, "tag:feedparser.org,2005-11-11:/docs/examples/atom11.xml:1", e.ID)
	assertEqual(t, "Second entry title", e.Title)
	assertEqual(t, true, time.Date(2005, time.November, 11, 11, 56, 34, 0, time.UTC).Equal(f.Updated))
	assertEqual(t, true, e.Created.Equal(e.Updated))

	assertEqual(t, 0, len(e.Links))

	assertEqual(t, 0, len(e.Categories))

	assertEqual(t, 1, len(e.Authors))
	assertEqual(t, "Mark Pilgrim <mark@example.org> (http://diveintomark.org/)", e.Authors[0])

	assertEqual(t, 0, len(e.Contributors))

	assertEqual(t, "Test content as text", e.Summary)
	assertEqual(t, "Test content", e.Content)

}

func TestRSS(t *testing.T) {

	t.Parallel()

	f := testFile(t, "../../../test/feed/wordpress.xml")

	assertEqual(t, "rss2.0", f.Flavor)
	assertEqual(t, "https://en.blog.wordpress.com/feed/", f.ID)
	assertEqual(t, true, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(f.Updated))
	assertEqual(t, "WordPress.com News", f.Title)
	assertEqual(t, "The latest news on WordPress.com and the WordPress community.", f.Subtitle)
	assertEqual(t, "http://wordpress.com/", f.Generator)
	assertEqual(t, "https://secure.gravatar.com/blavatar/e6392390e3bcfadff3671c5a5653d95b?s=96&d=https%3A%2F%2Fs2.wp.com%2Fi%2Fbuttonw-com.png", f.Icon)

	assertEqual(t, 4, len(f.Links))
	assertEqual(t, "https://en.blog.wordpress.com", f.Links["alternate"])
	assertEqual(t, "https://en.blog.wordpress.com/feed/", f.Links["self"])
	assertEqual(t, "https://en.blog.wordpress.com/osd.xml", f.Links["search"])
	assertEqual(t, "https://en.blog.wordpress.com/?pushpress=hub", f.Links["hub"])

	assertEqual(t, 10, len(f.Entries))
	e := f.Entries[0]

	assertEqual(t, "http://en.blog.wordpress.com/?p=31505", e.ID)
	assertEqual(t, "New Theme: Libre", e.Title)
	assertEqual(t, true, time.Date(2015, time.July, 2, 17, 0, 4, 0, time.UTC).Equal(e.Updated))
	assertEqual(t, true, e.Created.Equal(e.Updated))

	assertEqual(t, 1, len(e.Categories))
	assertEqual(t, "Themes", e.Categories[0])

	assertEqual(t, 1, len(e.Links))
	assertEqual(t, "https://en.blog.wordpress.com/2015/07/02/libre/", e.Links["alternate"])

	assertEqual(t, "<em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site.<img alt=\"\" border=\"0\" src=\"https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1\" width=\"1\" height=\"1\" />", e.Summary)
	assertEqual(t, "<p>Happy Theme Thursday, all! Today I&#8217;m happy to introduce <em>Libre</em>, a new free theme designed by <a href=\"http://carolinethemes.com/\">yours truly</a>.</p>\n<h3><a href=\"http://wordpress.com/themes/libre/\">Libre</a></h3>\n<p><a href=\"http://wordpress.com/themes/libre/\"><img class=\"aligncenter size-full wp-image-31516\" src=\"https://en-blog.files.wordpress.com/2015/06/librelg.png?w=635&#038;h=476\" alt=\"Libre\" width=\"635\" height=\"476\" /></a></p>\n<p><em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site. The main navigation bar stays fixed to the top while your visitors read, keeping your most important content at hand. At the bottom of your site, three footer widget areas give your secondary content a comfortable home.</p>\n<p>Customize <em>Libre</em> with a logo or a header image to make it your own, or use one of two custom templates &#8212; including a full-width template\u00a0with no sidebar &#8212; to change up the look of your pages. <em>Libre</em> sports a clean, responsive design that works seamlessly on screens of any size.</p>\n<p><img class=\"aligncenter size-full wp-image-31517\" src=\"https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg?w=635&#038;h=252\" alt=\"Responsive design\" width=\"635\" height=\"252\" /></p>\n<p>Read more about <em>Libre</em> on the <a href=\"https://wordpress.com/themes/libre/\">Theme Showcase</a>, or activate it on your site from <em>Appearance → Themes</em>!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel=\"nofollow\" href=\"http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31505/\"><img alt=\"\" border=\"0\" src=\"http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31505/\" /></a> <img alt=\"\" border=\"0\" src=\"https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1\" width=\"1\" height=\"1\" />", e.Content)

}

func testFeed(t *testing.T, reader io.ReadCloser, contentType string) *Feed {
	p := NewParser()
	feed, err := p.Parse(reader, contentType)
	assertNoError(t, err)
	assertNotNil(t, feed)
	return feed
}

func testFile(t *testing.T, filename string) *Feed {
	f, err := os.Open(filename)
	assertNoError(t, err)
	assertNotNil(t, f)
	defer f.Close()
	return testFeed(t, f, "")
}

func testURL(t *testing.T, url string) *Feed {
	rsp, err := http.Get(url)
	contentType := rsp.Header.Get("Content-Type")
	assertNoError(t, err)
	assertNotNil(t, rsp)
	defer rsp.Body.Close()
	return testFeed(t, rsp.Body, contentType)
}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("Expected not nil value")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Not equal: expected %v, actual %v", a, b)
	}
}
