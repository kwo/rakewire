package feedparser

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAndrew(t *testing.T) {

	t.SkipNow()

	f := testURL(t, "http://andrewhammel.typepad.com/how_many_which_ones_the_g/atom.xml")

	// then
	expectedEntries := 10
	if len(f.Entries) != expectedEntries {
		t.Errorf("Expected %d, actual %d", expectedEntries, len(f.Entries))
	}

	for _, entry := range f.Entries {
		t.Logf("Entry: %s", entry.Title)
	}

}

func TestStratechery(t *testing.T) {

	t.SkipNow()

	f := testURL(t, "https://stratechery.com/feed/")

	// then

	for _, entry := range f.Entries {
		t.Logf("Entry: %s", entry.Title)
	}

}

func TestReddit(t *testing.T) {

	t.SkipNow()

	f := testURL(t, "https://www.reddit.com/r/NichtDerPostillon.rss")

	// then

	for _, entry := range f.Entries {
		t.Logf("Entry: %s %s", entry.Title, entry.LinkAlternate)
	}

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

	htmlLink := feed.LinkAlternate
	if htmlLink != "https://www.tbray.org/ongoing/" {
		t.Errorf("Links without rel not being intrepreted as alternate link: expected: %s, actual: %s", "https://www.tbray.org/ongoing/", htmlLink)
	}

	selfLink := feed.LinkSelf
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

	f := testFeed(t, getAtomFeed(), "")

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

	f := testFeed(t, getRSSFeed(), "")

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

}

func testFeed(t *testing.T, reader io.Reader, contentType string) *Feed {
	p := NewParser()
	feed, err := p.Parse(reader, contentType)
	if err != nil {
		t.Fatalf("Error parsing feed: %s", err.Error())
	} else if feed == nil {
		t.Fatal("Nil feed")
	}
	return feed
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

func getAtomFeed() io.Reader {

	data := `
	<?xml version="1.0" encoding="UTF-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom" xml:base="http://example.org/" xml:lang="en">

		<id>tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml</id>
	  <title type="text">Sample Feed</title>
		<subtitle type="html">For documentation &lt;em&gt;only&lt;/em&gt;</subtitle>
		<updated>2005-11-09T11:56:34Z</updated>
		<author>
			<name>Mark Pilgrim</name>
			<uri>http://diveintomark.org/</uri>
			<email>mark@example.org</email>
		</author>
	  <link rel="alternate" type="html" href="/"/>
	  <link rel="self" type="application/atom+xml" href="http://www.example.org/atom10.xml"/>
		<category term="sports" scheme="http://localhost/sports" label="Sports"/>
		<contributor>
			<name>Mary Pilgrim</name>
			<uri>http://diveintomary.org/</uri>
			<email>mary@example.org</email>
		</contributor>
	  <generator uri="http://example.org/generator/" version="4.0">Sample Toolkit</generator>
		<icon>/icon.jpg</icon>
		<logo>/logo.jpg</logo>
		<rights type="html">&lt;p&gt;Copyright 2005, Mark Pilgrim&lt;/p&gt;</rights>

	  <entry>
			<id>tag:feedparser.org,2005-11-09:/docs/examples/atom10.xml:3</id>
	    <title>First entry title</title>
			<updated>2005-11-09T11:56:34Z</updated>
	    <link rel="alternate" href="/entry/3"/>
	    <link rel="related"type="text/html" href="http://search.example.com/"/>
	    <link rel="via"type="text/html" href="http://toby.example.com/examples/atom10"/>
	    <link rel="enclosure" type="video/mpeg4" href="http://www.example.com/movie.mp4" length="42301"/>
			<category term="football" scheme="http://localhost/sports/football" label="Football"/>
			<category term="basketball" scheme="http://localhost/sports/basketball" label="Basketball"/>
	    <published>2005-11-09T00:23:47Z</published>
			<rights type="html">
			  &amp;copy; 2005 John Doe
			</rights>
	    <author>
	      <name>Mark Pilgrim</name>
	      <uri>http://diveintomark.org/</uri>
	      <email>mark@example.org</email>
	    </author>
	    <contributor>
	      <name>Joe</name>
	      <uri>http://example.org/joe/</uri>
	      <email>joe@example.org</email>
	    </contributor>
	    <contributor>
	      <name>Sam</name>
	      <uri>http://example.org/sam/</uri>
	      <email>sam@example.org</email>
	    </contributor>
	    <summary type="text">Watch out for nasty tricks</summary>
	    <content type="xhtml" xml:base="http://example.org/entry/3" xml:lang="en-US">
	      <div xmlns="http://www.w3.org/1999/xhtml">Watch out for<span style="background-image: url(javascript:window.location=’http://example.org/’)">nasty tricks</span>
	    	</div>
	  	</content>
			<source>
			  <id>http://example.org/</id>
			  <title>Fourty-Two</title>
			  <updated>2003-12-13T18:30:02Z</updated>
			  <rights>© 2005 Example, Inc.</rights>
			</source>
		</entry>

		<entry>
			<id>tag:feedparser.org,2005-11-11:/docs/examples/atom11.xml:1</id>
	    <title>Second entry title</title>
			<updated>2005-11-11T11:56:34Z</updated>
	    <author>
	      <name>Mark Pilgrim</name>
	      <uri>http://diveintomark.org/</uri>
	      <email>mark@example.org</email>
	    </author>
	    <summary type="text">Test content as text</summary>
	    <content type="text">
	      Test content
	  	</content>
		</entry>

	</feed>
	`

	return strings.NewReader(data)

}

func getRSSFeed() io.Reader {

	data := `
	<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"
		xmlns:content="http://purl.org/rss/1.0/modules/content/"
		xmlns:wfw="http://wellformedweb.org/CommentAPI/"
		xmlns:dc="http://purl.org/dc/elements/1.1/"
		xmlns:atom="http://www.w3.org/2005/Atom"
		xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
		xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
		xmlns:georss="http://www.georss.org/georss" xmlns:geo="http://www.w3.org/2003/01/geo/wgs84_pos#" xmlns:media="http://search.yahoo.com/mrss/"
		>

	<channel>
		<title>WordPress.com News</title>
		<atom:link href="https://en.blog.wordpress.com/feed/" rel="self" type="application/rss+xml" />
		<link>https://en.blog.wordpress.com</link>
		<description>The latest news on WordPress.com and the WordPress community.</description>
		<lastBuildDate>Mon, 06 Jul 2015 12:28:17 +0000</lastBuildDate>
		<language>en</language>
		<sy:updatePeriod>hourly</sy:updatePeriod>
		<sy:updateFrequency>1</sy:updateFrequency>
		<generator>http://wordpress.com/</generator>
	<cloud domain='en.blog.wordpress.com' port='80' path='/?rsscloud=notify' registerProcedure='' protocol='http-post' />
	<image>
			<url>https://secure.gravatar.com/blavatar/e6392390e3bcfadff3671c5a5653d95b?s=96&#038;d=https%3A%2F%2Fs2.wp.com%2Fi%2Fbuttonw-com.png</url>
			<title>WordPress.com News</title>
			<link>https://en.blog.wordpress.com</link>
		</image>
		<atom:link rel="search" type="application/opensearchdescription+xml" href="https://en.blog.wordpress.com/osd.xml" title="WordPress.com News" />
		<atom:link rel='hub' href='https://en.blog.wordpress.com/?pushpress=hub'/>
		<item>
			<title>New Theme: Libre</title>
			<link>https://en.blog.wordpress.com/2015/07/02/libre/</link>
			<comments>https://en.blog.wordpress.com/2015/07/02/libre/#comments</comments>
			<pubDate>Thu, 02 Jul 2015 17:00:04 +0000</pubDate>
			<dc:creator><![CDATA[Caroline Moore]]></dc:creator>
					<category><![CDATA[Themes]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31505</guid>
			<description><![CDATA[<em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>Happy Theme Thursday, all! Today I&#8217;m happy to introduce <em>Libre</em>, a new free theme designed by <a href="http://carolinethemes.com/">yours truly</a>.</p>
	<h3><a href="http://wordpress.com/themes/libre/">Libre</a></h3>
	<p><a href="http://wordpress.com/themes/libre/"><img class="aligncenter size-full wp-image-31516" src="https://en-blog.files.wordpress.com/2015/06/librelg.png?w=635&#038;h=476" alt="Libre" width="635" height="476" /></a></p>
	<p><em>Libre</em> brings a stylish, classic look to your personal blog or longform writing site. The main navigation bar stays fixed to the top while your visitors read, keeping your most important content at hand. At the bottom of your site, three footer widget areas give your secondary content a comfortable home.</p>
	<p>Customize <em>Libre</em> with a logo or a header image to make it your own, or use one of two custom templates &#8212; including a full-width template with no sidebar &#8212; to change up the look of your pages. <em>Libre</em> sports a clean, responsive design that works seamlessly on screens of any size.</p>
	<p><img class="aligncenter size-full wp-image-31517" src="https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg?w=635&#038;h=252" alt="Responsive design" width="635" height="252" /></p>
	<p>Read more about <em>Libre</em> on the <a href="https://wordpress.com/themes/libre/">Theme Showcase</a>, or activate it on your site from <em>Appearance → Themes</em>!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31505/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31505/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31505&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/07/02/libre/feed/</wfw:commentRss>
			<slash:comments>33</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg?w=150" medium="image">
				<media:title type="html">Responsive design</media:title>
			</media:content>

			<media:content url="https://2.gravatar.com/avatar/57509a6de7585125357530e2f3c3af1b?s=96&#38;d=retro" medium="image">
				<media:title type="html">sixhours</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/librelg.png" medium="image">
				<media:title type="html">Libre</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/libreresponsive.jpg" medium="image">
				<media:title type="html">Responsive design</media:title>
			</media:content>
		</item>
			<item>
			<title>Reinvented Video for WordPress</title>
			<link>https://en.blog.wordpress.com/2015/07/01/videopress-next/</link>
			<comments>https://en.blog.wordpress.com/2015/07/01/videopress-next/#comments</comments>
			<pubDate>Wed, 01 Jul 2015 22:52:57 +0000</pubDate>
			<dc:creator><![CDATA[Guillermo Rauch]]></dc:creator>
					<category><![CDATA[New Features]]></category>
			<category><![CDATA[Video]]></category>
			<category><![CDATA[VideoPress]]></category>
			<category><![CDATA[video post]]></category>
			<category><![CDATA[videos]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31384</guid>
			<description><![CDATA[The next generation VideoPress is here -- powerful, simple video hosting for your blog or website.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31384&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>Today, we’re announcing a complete overhaul of <a href="http://videopress.com/" target="_blank">VideoPress</a>, the service that has powered more than 3,000,000,000 video plays on WordPress.com and Jetpack-connected self-hosted WordPress sites around the world. We’ve made the next-generation VideoPress dynamic, responsive, and lightning fast to support the ever-evolving needs of content creators everywhere.</p>
	<p>Take a peek at what’s under the hood!</p>
	<h3>Embed anywhere, play anywhere</h3>
	<p>Out of the box, the new VideoPress is lightweight and responsive for beautiful playback on any screen, from smartphones to desktops. VideoPress works on all modern browsers and devices, and gives blog and site authors the power to engage their audiences no matter where they are.</p>
	<p>Not only do videos look amazing on WordPress sites, but you can also embed your videos anywhere on the web &#8212; other websites, social media, chat services &#8212; by using a permalink or a snippet of code.</p>
	<h3>Major speed enhancements</h3>
	<p>VideoPress now takes a fraction of the space it used to and is optimized for speed, so pages and posts with video content load faster. This is a huge plus for viewers who use slower connections or rely on bandwidth-strapped mobile devices.</p>
	<h3>It&#8217;s your content</h3>
	<p>At Automattic, we believe wholeheartedly that you own your content &#8212; video or otherwise.</p>
	<p>The new VideoPress puts your content front and center. The player is ad-free and unbranded to ensure your videos look and feel like an integral part of your website or blog, not like they belong to a third-party video platform. Unlike other video hosting services, VideoPress starts and ends on your video, keeping traffic on your site and giving you full control over the content to which your visitors are exposed.</p>
	<h3>New features, new look</h3>
	<p>A real-time “seek” feature lets you skim through any video and helps you start playing at the desired point in the video. Here&#8217;s a sneak peek that our developer community will particularly enjoy: the player&#8217;s skin and behavior is controlled entirely by JavaScript, HTML and CSS, opening up a malleable slate for customizations by themers in the future.</p>
	<p><iframe width='635' height='208' src='//videopress.com/embed/4tekemKK?hd=1&#038;loop=1&#038;autoplay=1' frameborder='0' allowfullscreen></iframe><script src='https://s0.wp.com/wp-content/plugins/video/assets/js/next/videopress-iframe.js?m=1435166243'></script></p>
	<h3>Everything you love, and more</h3>
	<p>VideoPress embraces the features that users have come to love, like privacy settings and rich stats, but also includes additional enhancements. Every video now features its own permalink, and the sharing pane has been redesigned to offer different options for embedding, like starting playback at specific times, looping, and autoplay.</p>
	<p><iframe width='635' height='356' src='//videopress.com/embed/iQnkO4EL?hd=1&#038;loop=1&#038;autoplay=1' frameborder='0' allowfullscreen></iframe><script src='https://s0.wp.com/wp-content/plugins/video/assets/js/next/videopress-iframe.js?m=1435166243'></script></p>
	<h3>OpenSource ❤︎</h3>
	<p>The rebuild of VideoPress revolves around a robust API for embedding videos easily in WordPress.com and Jetpack-enabled self-hosted websites. Though the API is private at the moment, many components of the new VideoPress libraries have been open sourced, including <a href="https://github.com/automattic/jpeg-stream" target="_blank">jpeg-stream</a>, <a href="https://github.com/automattic/pixel-stack" target="_blank">pixel-stack</a>, and <a href="https://github.com/automattic/video-thumb-grid" target="_blank">video-thumb-grid</a>.</p>
	<h3>Get the goods</h3>
	<p>Once a standalone upgrade, VideoPress is now available exclusively in paid WordPress.com Plans to streamline updates, payments, and security enhancements.</p>
	<p>To enable all the goodness of video support on your site, <strong><a href="https://wordpress.com/plans/" target="_blank">upgrade your WordPress.com Plan</a></strong> to Premium or Business. Then, upload and share videos as you please. If you already have a Premium or Business plan, sit back and enjoy. VideoPress will automatically handle the nuts and bolts to deliver an amazing viewing experience for your fans, followers, and friends!</p>
	<p><em>Correction: an earlier version of this post insinuated that the new VideoPress is accessible for Jetpack users who upgrade to WordPress.com Premium or Business. Though videos can certainly be embedded anywhere, the option for Jetpack users to upload videos directly to VideoPress via a WordPress.com Premium or Business plan is scheduled for a future Jetpack release.</em></p><br />Filed under: <a href='https://en.blog.wordpress.com/category/new-features/'>New Features</a>, <a href='https://en.blog.wordpress.com/category/video/'>Video</a>, <a href='https://en.blog.wordpress.com/category/videopress/'>VideoPress</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31384/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31384/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31384&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/07/01/videopress-next/feed/</wfw:commentRss>
			<slash:comments>8</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-25-at-9-45-41-pm.png?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-25-at-9-45-41-pm.png?w=150" medium="image">
				<media:title type="html">VideoPress</media:title>
			</media:content>

			<media:content url="https://0.gravatar.com/avatar/ff90b5a02769128a41b49da9b33ac547?s=96&#38;d=retro" medium="image">
				<media:title type="html">rauchg</media:title>
			</media:content>
		</item>
			<item>
			<title>#LoveWins! LGBTQ Bloggers Make Their Voices Heard</title>
			<link>https://en.blog.wordpress.com/2015/06/30/love-wins-lgbtq/</link>
			<comments>https://en.blog.wordpress.com/2015/06/30/love-wins-lgbtq/#comments</comments>
			<pubDate>Tue, 30 Jun 2015 17:30:00 +0000</pubDate>
			<dc:creator><![CDATA[Hugo Baeta]]></dc:creator>
					<category><![CDATA[Admin Bar]]></category>
			<category><![CDATA[Community]]></category>
			<category><![CDATA[lgbt]]></category>
			<category><![CDATA[lgbtq]]></category>
			<category><![CDATA[Love]]></category>
			<category><![CDATA[pride]]></category>
			<category><![CDATA[same-sex marriage]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31443</guid>
			<description><![CDATA[We're so proud to provide a platform for all the incredibly talented LGBTQ writers who are advocating for change and sharing their stories on WordPress.com.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31443&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>You might have noticed the rainbow banner across the top of WordPress.com over the weekend &#8212; our way of marking Pride month, celebrated by cities across the globe in June, as well as the US Supreme Court&#8217;s ruling legalizing same-sex marriage across all states. (The United States now joins <a href="http://www.pewforum.org/2015/06/26/gay-marriage-around-the-world-2013/">20 other countries</a>, including my own, Portugal, in fully recognizing same-sex marriage nationwide.)</p>
	<div class="embed-twitter">
	<blockquote class="twitter-tweet" width="550">
	<p lang="en" dir="ltr">Today is a big step in our march toward equality. Gay and lesbian couples now have the right to marry, just like anyone else. <a href="https://twitter.com/hashtag/LoveWins?src=hash">#LoveWins</a></p>
	<p>&mdash; President Obama (@POTUS) <a href="https://twitter.com/POTUS/status/614435467120001024">June 26, 2015</a></p></blockquote>
	<p><script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script></div>
	<p>Here at WordPress.com we strive to democratize publishing and empower freedom of speech. It&#8217;s amazing to see the thoughtful analyses of the Supreme Court&#8217;s decision already being published, like <a href="https://tropicsofmeta.wordpress.com/2015/06/29/what-ancient-history-can-teach-us-about-lovewins/" target="_blank">this excellent piece from <em>Tropics of Meta</em></a> putting the decision into long-term historical context. We&#8217;re also proud to provide a platform for all the incredibly talented LGBTQ writers who are advocating for change and sharing their stories, like&#8230;</p>
	<ul>
	<li><em><strong><a href="http://letsqueerthingsup.com/" target="_blank">Let&#8217;s Queer Things Up!</a></strong></em>, written by transgender activist Sam Dylan Finch. We appreciate how open Sam is about all aspects of his life, from his <a href="http://letsqueerthingsup.com/2015/05/31/a-psychiatrist-endangered-my-life-and-i-was-afraid-to-speak-out/" target="_blank">fraught relationship with psychiatry</a> to the <a href="http://letsqueerthingsup.com/2015/06/05/binding-while-broke-i-tried-all-these-cheapish-chest-binders-so-you-dont-have-to/" target="_blank">best affordable chest binders</a>.</li>
	<li><em><strong><a href="https://bullybloggers.wordpress.com/" target="_blank">Bully Bloggers</a></strong></em>, an excellent group blog on &#8220;everything queer&#8221; from five professors from NYU, Dartmouth College, USC, and the University of Arizona. If you need a break after their <a href="https://bullybloggers.wordpress.com/2015/06/20/the-republic-of-love/" target="_blank">analysis of the Irish same-sex marriage referendum</a>, plant your tongue in your cheek and visit their &#8220;<a href="https://bullybloggers.wordpress.com/freedom-society-page/" target="_blank">Freedom to Marry our Pets Society</a>&#8221; page.</li>
	<li><strong><a href="https://connerhabib.wordpress.com/" target="_blank">Conner Habib</a></strong><em>,</em> a writer, public speaker, and sex-workers&#8217; advocate who penned the incredibly popular (and thought-provoking) piece, &#8220;<a href="https://connerhabib.wordpress.com/2013/02/13/why-do-gay-porn-stars-kill-themselves/" target="_blank">Why do gay porn stars kill themselves?</a>&#8220;</li>
	</ul>
	<div data-shortcode="caption" id="attachment_31467" class="wp-caption aligncenter"><a href="https://en-blog.files.wordpress.com/2015/06/pride_2015-1.jpg"><img class="wp-image-31467 size-full" src="https://en-blog.files.wordpress.com/2015/06/pride_2015-1.jpg?w=635&#038;h=423" alt="Image by Pete Rosos, 2812 Photography." width="635" height="423" /></a><p class="wp-caption-text">Pride parade, San Francisco. Image by Pete Rosos, <a href="http://2812photography.com/2015/06/28/pride-2015/">2812 Photography</a>.</p></div>
	<p>WordPress.com bloggers were also out and about with their rainbow flags this weekend, participating in pride parades from Los Angeles to Lisbon and capturing the days in pictures:</p>
	<ul>
	<li><a href="https://ronmayhewphotography.wordpress.com/2015/06/28/key-west-pride/" target="_blank">Key West, Florida</a>, via Ron Mayhew Photography.</li>
	<li><a href="http://bozzettomedia.com/2015/06/29/toronto-pride-parade/" target="_blank">Toronto, Ontario</a>, via Michelle Bozzette Media.</li>
	<li><a href="http://transienteye.com/2015/06/29/barcelona-pride-2015/" target="_blank">Barcelona, Spain</a>, via <em>Transient Eye</em>.</li>
	<li><a href="http://seattletoshanghai.com/2015/06/28/pride-parade-in-seattle-is-a-rainbow-of-diversity/" target="_blank">Seattle, Washington</a>, via <em>Seattle to Shanghai</em>. (Bonus: the <a href="http://seattlepride.org/" target="_blank">Seattle Pride website</a> is powered by WordPress, as are <a href="http://atlantapride.org/" target="_blank">Atlanta&#8217;s</a> and <a href="http://vancouverpride.ca/" target="_blank">Vancouver&#8217;s</a>!)</li>
	<li><a href="http://qz.com/440054/photos-istanbuls-pride-parade-was-brutally-dispersed-with-water-cannons/" target="_blank">Istanbul, Turkey</a>, via <em>Quartz </em>(a less-than-joyous occasion that shows just how important it is to elevate LGBTQ voices).</li>
	<li><a href="http://parislondonstyle.com/2015/06/28/pride-london-2015/" target="_blank">London, England</a>, via <em>Paris London Style.</em></li>
	<li>And of course, <a href="http://2812photography.com/2015/06/28/pride-2015/" target="_blank">San Francisco, California</a>, via <em>2812 Photography</em>!</li>
	</ul>
	<div data-shortcode="caption" id="attachment_31468" class="wp-caption aligncenter"><a href="https://en-blog.files.wordpress.com/2015/06/key-west-pride-3.jpg"><img class="size-full wp-image-31468" src="https://en-blog.files.wordpress.com/2015/06/key-west-pride-3.jpg?w=635&#038;h=409" alt="Pride parade, Key West, Florida, USA. Photo by Ron Mayhew Photography." width="635" height="409" /></a><p class="wp-caption-text">Pride parade, Key West, Florida, USA.<br />Photo by <a href="https://ronmayhewphotography.wordpress.com/2015/06/28/key-west-pride/">Ron Mayhew Photography</a>.</p></div>
	<p>We also had the pleasure of launching a website for <a href="http://doubleduchess.com/">Double Duchess</a>, an energetic queer musical duo from the San Francisco Bay Area who just released their latest album. And c<span class="s1">heck out the new tees and tanks added to the <a href="http://hellomerch.com/collections/wordpress" target="_blank">WordPress swag store</a> in early June, in honor of Pride month &#8212; i</span>t&#8217;s the first time that we&#8217;ve made the same design available in a such wide, inclusive range of colors, fits, and body shapes! All the profits support the WordPress Foundation, which helps ensure that the free, downloadable version of WordPress is cutting edge, powerful, and ready for whoever wants to have a voice on the web.</p>
	<p><a href="http://hellomerch.com/collections/wordpress"><img class="aligncenter wp-image-31510 size-full" src="https://en-blog.files.wordpress.com/2015/06/wordpress-pride-swag-3.png?w=635" alt="WordPress Pride Swag"   /></a></p>
	<p>Automattic is committed to diversity as a company (<a href="http://automattic.com" target="_blank">and we&#8217;re hiring!</a>), and to providing tools that let anyone use the web to tell their truths and work for equality. Happy Pride!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/admin-bar/'>Admin Bar</a>, <a href='https://en.blog.wordpress.com/category/community/'>Community</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31443/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31443/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31443&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/30/love-wins-lgbtq/feed/</wfw:commentRss>
			<slash:comments>55</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/pride-featureimage-clear.png?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/pride-featureimage-clear.png?w=150" medium="image">
				<media:title type="html">pride-featureimage-clear</media:title>
			</media:content>

			<media:content url="https://0.gravatar.com/avatar/0218e99cde7283d859ea46a21e744319?s=96&#38;d=retro" medium="image">
				<media:title type="html">hugobaeta</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/pride_2015-1.jpg" medium="image">
				<media:title type="html">Image by Pete Rosos, 2812 Photography.</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/key-west-pride-3.jpg" medium="image">
				<media:title type="html">Pride parade, Key West, Florida, USA. Photo by Ron Mayhew Photography.</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/wordpress-pride-swag-3.png" medium="image">
				<media:title type="html">WordPress Pride Swag</media:title>
			</media:content>
		</item>
			<item>
			<title>July in Blogging U.: Blogging 101 and 201</title>
			<link>https://en.blog.wordpress.com/2015/06/29/july-in-blogging-u-blogging-101-and-201/</link>
			<comments>https://en.blog.wordpress.com/2015/06/29/july-in-blogging-u-blogging-101-and-201/#comments</comments>
			<pubDate>Mon, 29 Jun 2015 17:40:48 +0000</pubDate>
			<dc:creator><![CDATA[Michelle W.]]></dc:creator>
					<category><![CDATA[Better Blogging]]></category>
			<category><![CDATA[Community]]></category>
			<category><![CDATA[blogging u.]]></category>
			<category><![CDATA[blogging101]]></category>
			<category><![CDATA[Blogging201]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31454</guid>
			<description><![CDATA[Need help getting your blog off the ground or getting readers engaged?<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31454&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<blockquote class="alignleft"><p>Note: Blogging 101 and Blogging 201 cover the same content each time they are offered.</p></blockquote>
	<p>Have you just started blogging (welcome!), or are you looking to breathe new life into a blogging habit that&#8217;s fallen by the wayside? Blogging U. is a great way to get on track, with bite-size assignments, a supportive community, and staff to support you. We&#8217;re offering two courses in July &#8212; learn more:</p>
	<h3>Blogging 101: Zero to Hero — July 6 &#8211; 24</h3>
	<p><em>Blogging 101</em> is three weeks of bite-size blogging assignments that take you from “Blog?” to “Blog!” Every weekday, you’ll get a new assignment to help you publish a post, customize your blog, or engage with the community.</p>
	<p>You’ll walk away with a stronger focus for your blog, several published posts and a handful of drafts, a theme that reflects your personality, a small (but growing!) audience, a grasp of blogging etiquette — and a bunch of new friends.</p>
	<h3>Blogging 201: Branding and Growth &#8212; July 20 &#8211; 31</h3>
	<p><em>Blogging 201</em> is a two-week challenge that gives you the tools to define your brand, build your audience, use your stats to grow your traffic, and bring your older posts fresh attention. You don’t need to have completed <em>Blogging 101</em> to register, although it makes a great foundation.</p>
	<blockquote><p>Please note that we ask you not to register for Blogging 101 and 201 at the same time; Blogging 201 assumes that you have some readers and have already accomplished a lot of what we cover in the 101-level course. Both courses will be offered several more times throughout the year.</p></blockquote>
	<h3 class="p1">How do Blogging U. courses work?</h3>
	<p class="p1">Blogging U. courses exist for one reason: to help you meet your own blogging and writing goals.</p>
	<ul>
	<li>Courses are free, flexible, and open to all.</li>
	<li>You’ll get a new task to complete each day, along with our best advice and favorite resources. Do them on your own time, and interpret them however makes sense for your specific blog and personal goals — we’re not grading you, we’re not checking to make sure you complete every task, and there’s no “wrong” way to use the resources we give you.</li>
	<li>We’ll post new assignments here on <em>The Daily Post </em>each weekday at 12AM GMT. Each assignment will contain all the inspiration and instructions you need to complete it. Weekends are free.</li>
	<li>Each course will have a private community site, the <em>Commons,</em> for chatting, connecting, and seeking feedback and support. <em>Daily Post </em>staff and Happiness Engineers will be on hand to answer your questions and offer guidance and resources.</li>
	</ul>
	<h3>Ready to register?</h3>
	<p><strong><em>Registration for Blogging 101 is now closed, but you can still register for Blogging 201 with the form below.</em></strong></p>
	<p>Just fill out this short form! There&#8217;s no automated confirmation; you&#8217;ll receive a welcome email just prior to the start of your course. If you&#8217;re on a mobile device or reading this via email and don&#8217;t see the form, <a href="http://bloggingu.polldaddy.com/s/july-2015-blogging-201-extras" target="_blank">you can register with this link</a>.</p>
	<div class="pd-embed" data-settings="{&quot;type&quot;:&quot;iframe&quot;,&quot;auto&quot;:true,&quot;domain&quot;:&quot;bloggingu.polldaddy.com\/s\/&quot;,&quot;id&quot;:&quot;july-2015-blogging-201-extras&quot;}"></div>
	<script type="text/javascript"><!--//--><![CDATA[//><!--
	(function(d,c,j){if(!d.getElementById(j)){var pd=d.createElement(c),s;pd.id=j;pd.src=('https:'==d.location.protocol)?'https://polldaddy.com/survey.js':'http://i0.poll.fm/survey.js';s=d.getElementsByTagName(c)[0];s.parentNode.insertBefore(pd,s);}}(document,'script','pd-embed'));
	//--><!]]&gt;</script>
	<noscript><a href="http://polldaddy.com/s/554F14F52BB4DC6C">Take Our Survey</a></noscript><br />Filed under: <a href='https://en.blog.wordpress.com/category/better-blogging/'>Better Blogging</a>, <a href='https://en.blog.wordpress.com/category/community/'>Community</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31454/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31454/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31454&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/29/july-in-blogging-u-blogging-101-and-201/feed/</wfw:commentRss>
			<slash:comments>49</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2014/05/bu-featured-alt3.png?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2014/05/bu-featured-alt3.png?w=150" medium="image">
				<media:title type="html">blogging U</media:title>
			</media:content>

			<media:content url="https://2.gravatar.com/avatar/2367004060918e221dcb9799584e9279?s=96&#38;d=retro" medium="image">
				<media:title type="html">Michelle, Communications</media:title>
			</media:content>
		</item>
			<item>
			<title>Celebrating 10 Years of WordPress.com &amp; Automattic</title>
			<link>https://en.blog.wordpress.com/2015/06/26/celebrating-10-years-of-wordpress-com-automattic/</link>
			<comments>https://en.blog.wordpress.com/2015/06/26/celebrating-10-years-of-wordpress-com-automattic/#comments</comments>
			<pubDate>Fri, 26 Jun 2015 12:29:32 +0000</pubDate>
			<dc:creator><![CDATA[Mark Armstrong]]></dc:creator>
					<category><![CDATA[Automattic]]></category>
			<category><![CDATA[behind the scenes]]></category>
			<category><![CDATA[Community]]></category>
			<category><![CDATA[Milestone]]></category>
			<category><![CDATA[WordPress.com]]></category>
			<category><![CDATA[10 years]]></category>
			<category><![CDATA[birthday]]></category>
			<category><![CDATA[thank you]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31295</guid>
			<description><![CDATA[This year marks the 10th birthday of WordPress.com and our parent company, Automattic. We are proud to have served this community of millions: from writers, photographers, artists, and small and large publishers, to business owners and entrepreneurs. A quick bit of history: WordPress itself started as an open source project &#8230;<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31295&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p><span class='embed-youtube' style='text-align:center; display: block;'><iframe class='youtube-player' type='text/html' width='635' height='388' src='https://www.youtube.com/embed/EYs_FckMqow?version=3&#038;rel=1&#038;fs=1&#038;showsearch=0&#038;showinfo=1&#038;iv_load_policy=1&#038;wmode=transparent' frameborder='0' allowfullscreen='true'></iframe></span></p>
	<p>This year marks the 10th birthday of WordPress.com and our parent company, <a href="http://automattic.com/">Automattic</a>. We are proud to have served this community of millions: from writers, photographers, artists, and small and large publishers, to business owners and entrepreneurs.<span id="more-31295"></span></p>
	<p>A quick bit of history: WordPress itself started as an open source project <a href="https://wordpress.org/news/2013/05/ten-good-years/">in 2003</a>. WordPress co-founder Matt Mullenweg then built our company, and the free, hosted version &#8212; WordPress.com, what you see now &#8212; <a href="http://ma.tt/2015/06/ten-years-of-automattic/">opened for business in the summer of 2005</a>.</p>
	<p>Ten years, 2.5 billion posts, and 3 billion comments later, Automattic is stronger than ever &#8212; with WordPress.com and a host of other services aimed at helping independent publishers, bloggers, and business owners (the roster includes <a href="http://www.woothemes.com/woocommerce/">WooCommerce</a>, <a href="http://jetpack.me/">Jetpack</a>, <a href="https://akismet.com/">Akismet</a>, <a href="https://vaultpress.com/">VaultPress</a>, <a href="https://polldaddy.com">Polldaddy</a>, <a href="https://cloudup.com">Cloudup</a>, <a href="http://simplenote.com/">Simplenote</a>, <a href="http://longreads.com/">Longreads</a>, and more). All this, and we have <a href="https://en.blog.wordpress.com/2015/06/17/a-perfect-eff-score-were-proud-to-have-your-back/">a perfect record</a> from the Electronic Frontier Foundation for protecting our users’ rights.</p>
	<p>In addition to building a world-class publishing platform, our company has redefined what it means to be a global company working on the internet: we have nearly 400<a href="http://automattic.com/about/"> employees</a>, working from home (or their preferred coworking spaces), around the world. We use our own <a href="http://geto2.com/">p2 sites</a> to communicate with each other, and we all spend time doing support rotations with the people who matter most: Our users. (Also: <a href="http://automattic.com/work-with-us/">we’re hiring!</a>)</p>
	<p>And we&#8217;re just getting started: <a href="http://w3techs.com/technologies/details/cm-wordpress/all/all">24%</a> of websites on the internet now use WordPress to power their sites, and we believe there&#8217;s more for us to accomplish together to continue to make WordPress.com the engine that powers your creativity, your thoughts and ideas, and your business.</p>
	<p>Thank you to everyone who makes this community so special. Here’s to the next 10 years.</p>
	<h3>10 Years: Automattic &amp; WordPress.com by the Numbers</h3>
	<div class="ye-2014-numbers">
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-left-image-sites" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-post.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Total number of posts</div>
	<div class="ye-2014-big-number">2.5 billion</div>
	<div class="ye-2014-sub-text">Written with WordPress.com and Jetpack.</div>
	</blockquote>
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-right-image-posts" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-languages.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Number of languages</div>
	<div class="ye-2014-big-number">137</div>
	<div class="ye-2014-sub-text">Used across those 2.5 billion posts.</div>
	</blockquote>
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-left-image-sites" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-comment.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Number of comments</div>
	<div class="ye-2014-big-number">3 billion</div>
	<div class="ye-2014-sub-text">Conversations, encouraging words, passionate debates.</div>
	</blockquote>
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-right-image-posts" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-title.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Longest title</div>
	<div class="ye-2014-big-number">19,176 words</div>
	<div class="ye-2014-sub-text">On a WordPress.com post.</div>
	</blockquote>
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-left-image-sites" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-words1.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Longest post ever published</div>
	<div class="ye-2014-big-number">10+ million</div>
	<div class="ye-2014-sub-text">Words to share with others.</div>
	</blockquote>
	<blockquote class="ye-2014-numbers-item">
	<div class="ye-2014-right-image-posts" style="background:url('https://en-blog.files.wordpress.com/2015/06/10-support1.png') no-repeat center center;background-size:146px 152px;"></div>
	<div class="ye-2014-small-title">Total support messages</div>
	<div class="ye-2014-big-number">2.3 million</div>
	<div class="ye-2014-sub-text">Between customers and happiness engineers.</div>
	</blockquote>
	</div>
	<div class="ye-2014-automattic">
	<h3><a href="http://automattic.com/work-with-us/" target="_blank">Interested in being a part of our motley but merry crew? Work with us!</a></h3>
	<div class="ye-2014-banner-img" style="background-image:url('//en-blog.files.wordpress.com/2014/12/a8ccompanyphoto1.jpg');"></div>
	</div><br />Filed under: <a href='https://en.blog.wordpress.com/category/automattic/'>Automattic</a>, <a href='https://en.blog.wordpress.com/category/behind-the-scenes/'>behind the scenes</a>, <a href='https://en.blog.wordpress.com/category/community/'>Community</a>, <a href='https://en.blog.wordpress.com/category/milestone/'>Milestone</a>, <a href='https://en.blog.wordpress.com/category/wordpresscom/'>WordPress.com</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31295/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31295/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31295&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/26/celebrating-10-years-of-wordpress-com-automattic/feed/</wfw:commentRss>
			<slash:comments>117</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2014/01/a8c-lounge.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2014/01/a8c-lounge.jpg?w=150" medium="image">
				<media:title type="html">Automattic Lounge</media:title>
			</media:content>

			<media:content url="https://1.gravatar.com/avatar/19f62880113934b262c0bbb597644adb?s=96&#38;d=retro" medium="image">
				<media:title type="html">heymarkarms</media:title>
			</media:content>
		</item>
			<item>
			<title>Early Theme Adopters: Gazette</title>
			<link>https://en.blog.wordpress.com/2015/06/24/early-theme-adopters-gazette/</link>
			<comments>https://en.blog.wordpress.com/2015/06/24/early-theme-adopters-gazette/#comments</comments>
			<pubDate>Wed, 24 Jun 2015 16:00:26 +0000</pubDate>
			<dc:creator><![CDATA[Ben Huberman]]></dc:creator>
					<category><![CDATA[Customization]]></category>
			<category><![CDATA[Themes]]></category>
			<category><![CDATA[Customizing]]></category>
			<category><![CDATA[Early Theme Adopters]]></category>
			<category><![CDATA[Featured content]]></category>
			<category><![CDATA[Featured Image]]></category>
			<category><![CDATA[Gazette]]></category>
			<category><![CDATA[Site Logo]]></category>
			<category><![CDATA[theme]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31351</guid>
			<description><![CDATA[Whether you aim for a minimalist or visually rich site, <em>Gazette</em> is a magazine theme that makes your content shine.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31351&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p><em><a href="https://theme.wordpress.com/themes/gazette/">Gazette</a></em> is a theme that balances rich functionality with a pleasant, non-obtrusive look. Depending on your site&#8217;s needs, you can tweak it to look as stark and clean, or as warm and vibrant as you wish.</p>
	<p>Business owners, visual artists, and bloggers of all stripes can find something (or many things) to love about <em>Gazette.</em> The theme is particularly suited for those with image-heavy content: with striking <a href="http://en.support.wordpress.com/featured-images/">Featured Images</a>, several custom <a href="http://en.support.wordpress.com/posts/post-formats/">Post Formats</a>, and an optional featured content area on the homepage, it displays your posts with extra visual oomph.</p>
	<p>Not sure if it&#8217;s the theme for you? Here are four strikingly different sites that make the most of <em><a href="https://theme.wordpress.com/themes/gazette/">Gazette</a></em>&#8216;s features.</p>
	<h3><a href="https://dunc.wordpress.com/">simple tricks &amp; nonsense</a></h3>
	<p><a href="https://dunc.wordpress.com/"><img class="aligncenter size-large wp-image-31355" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-23-53-am.png?w=635&#038;h=290" alt="simple tricks &amp; nonsense" width="635" height="290" /></a><br />
	Covering the geekier side of film and pop culture, <i><a href="https://dunc.wordpress.com/">simple tricks &amp; nonsense</a></i> harnesses <em>Gazette</em>&#8216;s featured content area to highlight some of its greatest hits from the archives. The result is a colorful, engaging homepage that draws the reader in (the <a href="https://en.support.wordpress.com/custom-design/custom-fonts/">custom font</a> in the headings adds a nice touch, too).</p>
	<p>Scroll below the fold, and the site&#8217;s latest posts are there, each with its own featured image:</p>
	<p><a href="https://dunc.wordpress.com/"><img class="aligncenter size-large wp-image-31354" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-24-27-am.png?w=635&#038;h=407" alt="simple tricks &amp; nonsesne 2" width="635" height="407" /></a></p>
	<p>The colorful grid makes for a pleasing, smooth reading experience, where the reader&#8217;s eye is always engaged, but never distracted or disoriented by too much visual stimuli.</p>
	<h3><a href="https://depressioncomix.wordpress.com/">depression comix</a></h3>
	<p><a href="https://depressioncomix.wordpress.com/"><img class="aligncenter size-large wp-image-31356" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-22-34-am.png?w=635&#038;h=441" alt="depression comix" width="635" height="441" /></a></p>
	<p>On the more minimalist end of the spectrum we find <em><a href="https://depressioncomix.wordpress.com/">depression comix</a></em>, a weekly web comic depicting the daily struggles of people dealing with mental health issues.</p>
	<p>The site&#8217;s design complements its topic perfectly: using only a <a href="https://en.support.wordpress.com/site-logo/">Site Logo</a> in the header area, visitors are immediately plunged into a neat grid of posts that mimics the format of the comic itself. Here, <em>Gazette</em> recedes into the background, letting <a href="http://www.depressioncomix.com/contact/">artist Clay&#8217;s</a> work do the talking.</p>
	<h3>Calluna Studios</h3>
	<p><a href="https://callunastudios.wordpress.com/"><img class="aligncenter size-large wp-image-31358" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-20-49-am.png?w=635&#038;h=490" alt="calluna studios" width="635" height="490" /></a></p>
	<p>A professional photographer&#8217;s website, <em>Calluna Studios</em> highlights the work of Heather, the woman behind the camera, with a warm, clean look.</p>
	<p>The site&#8217;s <a href="https://en.support.wordpress.com/menus/">primary navigation menu</a> in the header is stripped to the bare essentials: links to the homepage and the contact page. Our eye immediately wanders to the lovely <a href="https://en.support.wordpress.com/themes/custom-header-image/">header image</a>, featuring the silhouette of a family enjoying a beautiful summer sunset outdoors.</p>
	<p><a href="https://callunastudios.wordpress.com"><img class="aligncenter size-large wp-image-31363" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-12-06-49-pm.png?w=635&#038;h=263" alt="calluna studios detail" width="635" height="263" /></a></p>
	<p>Heather makes great use of the different <a href="https://en.support.wordpress.com/posts/post-formats/">Post Formats</a> available in <em>Gazette</em> to highlight her projects in different ways and to avoid an overly monotonous homepage. The row of posts above, for example, features (from left to right) a standard post, an image post, and a gallery post.</p>
	<p>If you&#8217;d like to see the theme&#8217;s <a href="http://en.support.wordpress.com/featured-images/">Featured Images</a> in action in single-post view, make sure to check out one of <a href="https://callunastudios.wordpress.com/2015/06/18/659/">these</a> <a href="https://callunastudios.wordpress.com/2015/06/12/taryn-gil-toronto-maternity-session/">three</a> <a href="https://callunastudios.wordpress.com/2015/05/26/katie-sam/">posts</a>. (Hint: they look gorgeous!)</p>
	<h3><a href="http://back2spain.com/">Back to Spain</a></h3>
	<p><a href="http://back2spain.com/"><img class="aligncenter size-large wp-image-31357" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-21-48-am.png?w=635&#038;h=412" alt="back to spain" width="635" height="412" /></a></p>
	<p><em><a href="http://back2spain.com/">Back to Spain</a></em>, a food and lifestyle site by <a href="http://en.gravatar.com/caitlintherese">Caitlin</a>, a blogger passionate about cooking and travel, nails a perfect balance between minimalism and vibrancy.</p>
	<p>Caitlin&#8217;s homepage and her individual posts come alive with beautiful photography, and she uses <i>Gazette</i>&#8216;s full-width <a href="http://en.support.wordpress.com/featured-images/">Featured Images</a> to great effect throughout her site, like in this post for a <em>patatas arrugadas</em> recipe:</p>
	<p><a href="http://back2spain.com/2015/04/19/patatas-arrugadas-wrinkled-potatoes/"><img class="aligncenter size-large wp-image-31364" src="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-12-24-19-pm.png?w=635&#038;h=398" alt="patatas arrugadas" width="635" height="398" /></a><br />
	The rest of Caitlin&#8217;s site is just as smartly designed, with a well-organized <a href="https://en.support.wordpress.com/menus/">main menu</a> in the header and discreet social links in the footer&#8217;s widget area (one of two optional widget areas in <em>Gazette</em>).</p>
	<p><em>Are you using </em>Gazette <em>for your site? Which of the theme&#8217;s features is your favorite? Let us know in the comments!</em></p><br />Filed under: <a href='https://en.blog.wordpress.com/category/customization/'>Customization</a>, <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31351/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31351/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31351&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/24/early-theme-adopters-gazette/feed/</wfw:commentRss>
			<slash:comments>8</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-23-53-am.png?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-23-53-am.png?w=150" medium="image">
				<media:title type="html">simple tricks &#38; nonsense</media:title>
			</media:content>

			<media:content url="https://0.gravatar.com/avatar/663dcd498e8c5f255bfb230a7ba07678?s=96&#38;d=retro" medium="image">
				<media:title type="html">benhuberman</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-23-53-am.png?w=635" medium="image">
				<media:title type="html">simple tricks &#38; nonsense</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-24-27-am.png?w=635" medium="image">
				<media:title type="html">simple tricks &#38; nonsesne 2</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-22-34-am.png?w=635" medium="image">
				<media:title type="html">depression comix</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-20-49-am.png?w=635" medium="image">
				<media:title type="html">calluna studios</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-12-06-49-pm.png?w=635" medium="image">
				<media:title type="html">calluna studios detail</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-11-21-48-am.png?w=635" medium="image">
				<media:title type="html">back to spain</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/screen-shot-2015-06-23-at-12-24-19-pm.png?w=635" medium="image">
				<media:title type="html">patatas arrugadas</media:title>
			</media:content>
		</item>
			<item>
			<title>New Theme: Argent</title>
			<link>https://en.blog.wordpress.com/2015/06/18/argent/</link>
			<comments>https://en.blog.wordpress.com/2015/06/18/argent/#comments</comments>
			<pubDate>Thu, 18 Jun 2015 16:00:28 +0000</pubDate>
			<dc:creator><![CDATA[Aleksandra Laczek]]></dc:creator>
					<category><![CDATA[Themes]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31263</guid>
			<description><![CDATA[Showcase your best work with <em>Argent</em>.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31263&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>Today we’re happy to debut a new, free portfolio theme, <em>Argent</em>!</p>
	<h3><a href="http://wordpress.com/themes/argent">Argent</a></h3>
	<p><a href="http://wordpress.com/themes/argent/"><img class="alignnone size-full wp-image-31267" src="https://en-blog.files.wordpress.com/2015/06/argent-frontpage.jpg?w=635" alt="Argent WordPress Theme"   /></a></p>
	<p>Meet <em>Argent</em>, a new addition to our theme collection, designed by Automattic’s own <a href="https://profiles.wordpress.org/melchoyce">Mel Choyce</a>. <em>Argent&#8217;s</em> clean, modern portfolio theme is perfect for creative professionals like designers, artists, and photographers. Whether you&#8217;re showcasing a photo series or a design concept, <em>Argent&#8217;s</em> simple homepage template featuring portfolio projects will draw viewers to all of your wonderful work. Plus, the responsive layout allows for a seamless user experience and ensures that your portfolio looks stunning no matter the device or screen size.</p>
	<p>Read more about <em>Argent</em> on the <a href="https://wordpress.com/themes/argent/">Theme Showcase</a>, or activate it on your site from Appearance → Themes!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31263/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31263/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31263&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/18/argent/feed/</wfw:commentRss>
			<slash:comments>12</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/argent_responsive.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/argent_responsive.jpg?w=150" medium="image">
				<media:title type="html">Argent has a responsive design</media:title>
			</media:content>

			<media:content url="https://0.gravatar.com/avatar/698b667d7383ac83b109643d995193b5?s=96&#38;d=retro" medium="image">
				<media:title type="html">alex27pl</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/argent-frontpage.jpg" medium="image">
				<media:title type="html">Argent WordPress Theme</media:title>
			</media:content>
		</item>
			<item>
			<title>A perfect EFF score! We’re proud to have your back.</title>
			<link>https://en.blog.wordpress.com/2015/06/17/a-perfect-eff-score-were-proud-to-have-your-back/</link>
			<comments>https://en.blog.wordpress.com/2015/06/17/a-perfect-eff-score-were-proud-to-have-your-back/#comments</comments>
			<pubDate>Wed, 17 Jun 2015 19:02:58 +0000</pubDate>
			<dc:creator><![CDATA[Jenny Zhu]]></dc:creator>
					<category><![CDATA[Admin Bar]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31083</guid>
			<description><![CDATA[Concerns about online privacy and illicit government snooping are at the top of users’ minds, now more than ever. We appreciate that you trust us to safeguard your sensitive information on WordPress.com, and Automattic has a long-standing commitment to defending your rights and holding firm against legal bullying and over-reaching government requests. We &#8230;<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31083&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>Concerns about online privacy and illicit government snooping are at the top of users’ minds, now more than ever. We appreciate that you trust us to safeguard your sensitive information on WordPress.com, and Automattic has a long-standing commitment to defending your rights and holding firm against legal bullying and over-reaching government requests. We work to have the most stringent, user-friendly policies possible within the law, and to be as transparent as we can about information requests we receive and how we respond to them.</p>
	<p>Our friends at the <a href="https://www.eff.org/">Electronic Frontier Foundation</a> (EFF), an organization dedicated to defending your digital rights, recognized our efforts in their latest annual <a href="https://www.eff.org/who-has-your-back-government-data-requests-2015"><em>Who Has Your Back</em> report</a>, which evaluates the user privacy practices of prominent online service providers. We’re proud to receive a perfect score of five stars on the report, one of only nine (out of 24) companies to earn that honor. You can learn more about EFF&#8217;s evaluation criteria <a href="https://www.eff.org/who-has-your-back-government-data-requests-2015#evaluation-criteria">here</a>.</p>
	<p>We also received a perfect 5/5 score on the EFF’s most recent <a href="https://www.eff.org/pages/who-has-your-back-copyright-trademark-2014%20"><em>Intellectual Property Who Has Your Back</em> report</a>, making WordPress.com the only service to score a perfect 10/10 across both of the EFF’s surveys.</p>
	<p>In addition to our efforts on your behalf, we&#8217;re committed to giving back and making the internet as a whole a more secure place. Our <a href="http://transparency.automattic.com/legal-guidelines/">Legal Guidelines</a>, <a href="http://automattic.com/dmca-notice/">legal forms</a>, <a href="https://github.com/Automattic/legalmattic/tree/master/DMCA">DMCA templates</a>, <a href="https://wordpress.com/tos">Terms of Service</a>, and <a href="http://automattic.com/privacy/">Privacy Policies</a> are all licensed under Creative Commons licenses. Additionally, earlier this year, we posted a number of these documents to <a href="https://github.com/Automattic/legalmattic">GitHub</a> so other companies, startups, or small website owners can adopt and build on the policies we’ve adopted at Automattic.</p>
	<p>We’re proud to be recognized by the EFF for our work in this year’s <em>Who Has Your Back</em> report and are continually working to improve our practices to best serve the millions of website owners who put their trust in us.</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/admin-bar/'>Admin Bar</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31083/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31083/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31083&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/17/a-perfect-eff-score-were-proud-to-have-your-back/feed/</wfw:commentRss>
			<slash:comments>55</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/20150529-eff-have-your-back6.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/20150529-eff-have-your-back6.jpg?w=150" medium="image">
				<media:title type="html">20150529-EFF-Have-Your-Back</media:title>
			</media:content>

			<media:content url="https://0.gravatar.com/avatar/f7c00583275c2688796e7361faf544a6?s=96&#38;d=retro" medium="image">
				<media:title type="html">jennyzhu</media:title>
			</media:content>
		</item>
			<item>
			<title>Jenny Diski on Writing, Love, and Cancer</title>
			<link>https://en.blog.wordpress.com/2015/06/17/jenny-diski-on-writing-love-and-cancer/</link>
			<comments>https://en.blog.wordpress.com/2015/06/17/jenny-diski-on-writing-love-and-cancer/#comments</comments>
			<pubDate>Wed, 17 Jun 2015 14:07:31 +0000</pubDate>
			<dc:creator><![CDATA[Mark Armstrong]]></dc:creator>
					<category><![CDATA[Freshly Pressed]]></category>
			<category><![CDATA[WordPress.com]]></category>
			<category><![CDATA[Writing]]></category>
			<category><![CDATA[author]]></category>
			<category><![CDATA[books]]></category>
			<category><![CDATA[cancer]]></category>
			<category><![CDATA[essays]]></category>
			<category><![CDATA[jenny diski]]></category>
			<category><![CDATA[london review of books]]></category>
			<category><![CDATA[longreads]]></category>
			<category><![CDATA[memoir]]></category>
			<category><![CDATA[new york times magazine]]></category>
			<category><![CDATA[writer]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31271</guid>
			<description><![CDATA[The author and essayist reflects on a terminal cancer diagnosis, and all of the conflicted feelings that come with writing about it. <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31271&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<div data-shortcode="caption" id="attachment_31277" style="width: 645px" class="wp-caption alignnone"><a href="https://en-blog.files.wordpress.com/2015/06/diski1.png"><img class="size-large wp-image-31277" src="https://en-blog.files.wordpress.com/2015/06/diski1.png?w=635&#038;h=382" alt="Photo by Suki Dhanda" width="635" height="382" /></a><p class="wp-caption-text">Photo by <a href="https://sukidhanda.wordpress.com/">Suki Dhanda</a></p></div>
	<p>We’ve been following writer <a href="https://twitter.com/diski">Jenny Diski</a> for many years at the <em><a href="http://www.lrb.co.uk/contributors/jenny-diski">London Review of Books</a></em>, and more recently on her <a href="https://jennydiski.wordpress.com/">WordPress.com blog</a>. Just this past weekend Diski was featured in <a href="http://www.nytimes.com/2015/06/14/magazine/jenny-diskis-end-notes.html?_r=1">a profile by Giles Harvey</a> for the <em>New York Times Magazine</em>, about a subject she revealed in her own 2014 essays: she has inoperable lung cancer.<span id="more-31271"></span></p>
	<p>Her diagnosis began with a <em>Breaking Bad</em> joke (“So &#8212; we’d better get cooking the meth,” she said to her husband at the doctor’s office), and she rejects clichés, refusing to characterize her cancer as a &#8220;battle&#8221; or let others <a href="http://www.lrb.co.uk/v36/n17/jenny-diski/a-diagnosis">call her &#8220;brave&#8221;</a>:</p>
	<blockquote><p>One thing I state as soon as we’re out of the door: ‘Under no circumstances is anyone to say that I lost a battle with cancer. Or that I bore it bravely. I am not fighting, losing, winning or bearing.’ I will not personify the cancer cells inside me in any form. I reject all metaphors of attack or enmity in the midst, and will have nothing whatever to do with any notion of desert, punishment, fairness or unfairness, or any kind of moral causality. But I sense that I can’t avoid the cancer clichés simply by rejecting them. …</p>
	<p>I try but I can’t think of a single aspect of having cancer, start to finish, that isn’t an act in a pantomime in which my participation is guaranteed however I believe I choose to play each scene. I have been given this role. (There, see? Instant victim.) I have no choice but to perform and to be embarrassed to death.</p></blockquote>
	<p>She has made peace with writing about cancer (“I could either shut up, that’s the end, get on with dying. Or, get gripped, which is what happened”), and in the past year she’s also written about her time <a href="http://www.lrb.co.uk/v37/n01/jenny-diski/doris-and-me">living with Doris Lessing</a>, and on subjects such as love and beauty (<a href="https://jennydiski.wordpress.com/2014/12/12/depp-and-desire/">“Depp and Desire,”</a> translated from her column in the Swedish newspaper <em>Goteborgs-Posten</em>) and about meeting her husband, <a href="https://jennydiski.wordpress.com/2014/12/12/depp-and-desire/">whom she calls “The Poet”</a>:</p>
	<blockquote><p>When I was fifty I met The Poet, who is the same age as me. We had each left it until the last minute to find the relationship of our lives. Before that neither of us thought of ourselves as finally committed to a relationship, although we had had marriages and children. Our living happily ever after together, at such a late stage in our lives, is something we both smile at as improbable. It still surprises us, but it works. I don’t really know why. I came across something new, when we met, that both took in and transformed the youthful desire; we had the attraction but built a relationship on top of it that made the already but not quite diminished possibility at my age of looking at someone else in a room, wanting them, seeing it mirrored, and doing something about it, a voluntary surrender thereafter on my part.</p></blockquote>
	<p>For more on Diski, <a href="https://jennydiski.wordpress.com/">follow her blog</a>, and pick up her books: <a href="http://www.amazon.com/What-Dont-Know-About-Animals/dp/030018803X"><em>What I Don’t Know About Animals</em></a>, <a href="http://www.amazon.com/Skating-Antarctica-JENNY-DISKI/dp/1862070164"><em>Skating to Antarctica</em></a>, and <a href="http://www.amazon.com/Jenny-Diski/e/B001H6RZ6I/ref=dp_byline_cont_book_1">more</a>.</p>
	<p>Posts from Diski&#8217;s blog:</p>
	<h3><a href="https://jennydiski.wordpress.com/2015/02/23/sidebar-to-lrb-memoir-fish-there-are-fish/">&#8220;Fish, there Are Fish!&#8221;</a></h3>
	<p>A humorous sidebar to one of her LRB essays.</p>
	<h3><a href="https://jennydiski.wordpress.com/2015/02/19/a-sidebar-hows-it-going/">&#8220;How&#8217;s It Going?&#8221;</a></h3>
	<p>A window into Diski&#8217;s life, thoughts and emotions, in between the more formal essays.</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/freshly-pressed/'>Freshly Pressed</a>, <a href='https://en.blog.wordpress.com/category/wordpresscom/'>WordPress.com</a>, <a href='https://en.blog.wordpress.com/category/writing-2/'>Writing</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31271/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31271/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31271&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/17/jenny-diski-on-writing-love-and-cancer/feed/</wfw:commentRss>
			<slash:comments>4</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/7509061378_879c68927f_o.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/7509061378_879c68927f_o.jpg?w=150" medium="image">
				<media:title type="html">7509061378_879c68927f_o</media:title>
			</media:content>

			<media:content url="https://1.gravatar.com/avatar/19f62880113934b262c0bbb597644adb?s=96&#38;d=retro" medium="image">
				<media:title type="html">heymarkarms</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/diski1.png?w=635" medium="image">
				<media:title type="html">Photo by Suki Dhanda</media:title>
			</media:content>
		</item>
			<item>
			<title>New Theme: Cerauno</title>
			<link>https://en.blog.wordpress.com/2015/06/11/cerauno/</link>
			<comments>https://en.blog.wordpress.com/2015/06/11/cerauno/#comments</comments>
			<pubDate>Thu, 11 Jun 2015 16:00:55 +0000</pubDate>
			<dc:creator><![CDATA[Caroline Moore]]></dc:creator>
					<category><![CDATA[Themes]]></category>

			<guid isPermaLink="false">http://en.blog.wordpress.com/?p=31236</guid>
			<description><![CDATA[Bring your site into the spotlight with <em>Cerauno</em>, a free magazine theme.<img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31236&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></description>
					<content:encoded><![CDATA[<p>It&#8217;s a beautiful day for a new theme! Today, I&#8217;m happy to present a free magazine theme, <em>Cerauno</em>.</p>
	<h3><a href="http://wordpress.com/themes/cerauno">Cerauno</a></h3>
	<p><a href="http://wordpress.com/themes/cerauno"><img src="https://en-blog.files.wordpress.com/2015/06/ceraunolg.jpg?w=635&#038;h=705" alt="Cerauno" width="635" height="705" class="aligncenter size-full wp-image-31237" /></a></p>
	<p><em>Cerauno</em> is a polished, user-friendly theme designed by <a href="https://profiles.wordpress.org/melchoyce">Mel Choyce</a>.</p>
	<p>In Mel&#8217;s words:</p>
	<blockquote><p>I designed <em>Cerauno</em> with subject bloggers in mind, like food, fashion, or travel bloggers. I wanted to make a theme for someone who’s been blogging for a little while, but wants to boost their traffic and bring their site into the spotlight with a clean and authoritative design. I’m so excited to see it launch!</p></blockquote>
	<p>With plenty of options to tempt and tantalize, <em>Cerauno</em> is sure to please. Add secondary content in up to five widget areas, brand your content with a Site Logo or Custom Header, include links to your favorite social networks, and add Featured Images to grab your readers&#8217; attention.</p>
	<p><em>Cerauno</em> is also responsive, stretching or shrinking to accommodate any screen size.</p>
	<p><img src="https://en-blog.files.wordpress.com/2015/06/ceraunomobile.jpg?w=635&#038;h=252" alt="Cerauno&#039;s responsive design in action" width="635" height="252" class="aligncenter size-full wp-image-31239" /></p>
	<p>Check out <em>Cerauno</em> on the <a title="Cerauno Showcase" href="https://wordpress.com/themes/cerauno/">Theme Showcase</a>, or activate it on your site from <em>Appearance → Themes</em>!</p><br />Filed under: <a href='https://en.blog.wordpress.com/category/themes/'>Themes</a>  <a rel="nofollow" href="http://feeds.wordpress.com/1.0/gocomments/en.blog.wordpress.com/31236/"><img alt="" border="0" src="http://feeds.wordpress.com/1.0/comments/en.blog.wordpress.com/31236/" /></a> <img alt="" border="0" src="https://pixel.wp.com/b.gif?host=en.blog.wordpress.com&#038;blog=3584907&#038;post=31236&#038;subd=en.blog&#038;ref=&#038;feed=1" width="1" height="1" />]]></content:encoded>
				<wfw:commentRss>https://en.blog.wordpress.com/2015/06/11/cerauno/feed/</wfw:commentRss>
			<slash:comments>26</slash:comments>

			<media:thumbnail url="https://en-blog.files.wordpress.com/2015/06/ceraunomobile.jpg?w=150" />
			<media:content url="https://en-blog.files.wordpress.com/2015/06/ceraunomobile.jpg?w=150" medium="image">
				<media:title type="html">Cerauno&#039;s responsive design in action</media:title>
			</media:content>

			<media:content url="https://2.gravatar.com/avatar/57509a6de7585125357530e2f3c3af1b?s=96&#38;d=retro" medium="image">
				<media:title type="html">sixhours</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/ceraunolg.jpg" medium="image">
				<media:title type="html">Cerauno</media:title>
			</media:content>

			<media:content url="https://en-blog.files.wordpress.com/2015/06/ceraunomobile.jpg" medium="image">
				<media:title type="html">Cerauno&#039;s responsive design in action</media:title>
			</media:content>
		</item>
		</channel>
	</rss>
	`

	return strings.NewReader(data)

}
