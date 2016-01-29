package feedparser

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Feed feed
type Feed struct {
	Authors       []string
	Entries       []*Entry
	Flavor        string
	Generator     string
	Icon          string
	ID            string
	Links         map[string]string
	LinkAlternate string
	LinkSelf      string
	Rights        string
	Subtitle      string
	Title         string
	Updated       time.Time
}

// Entry entry
type Entry struct {
	Authors       []string
	Categories    []string
	Content       string
	Contributors  []string
	Created       time.Time
	ID            string
	Links         map[string]string
	LinkAlternate string
	LinkSelf      string
	Summary       string
	Title         string
	Updated       time.Time
}

func (z *Entry) hash() string {

	// Needs to be reproducible so that on future parsings of the same entry, the same ID will ge generated.
	// Also needs to have enough variables to be unique within the feed.

	hash := md5.New()
	hash.Write([]byte(z.Content))
	hash.Write([]byte(z.Summary))
	hash.Write([]byte(z.Title))
	hash.Write([]byte(z.Links[linkAlternate]))
	return hex.EncodeToString(hash.Sum(nil))

}

// Parser can parse feeds
type Parser struct {
	decoder *xml.Decoder
	entry   *Entry
	feed    *Feed
	stack   *elements
}

const (
	nsAtom       = "http://www.w3.org/2005/Atom"
	nsContent    = "http://purl.org/rss/1.0/modules/content/"
	nsDublinCore = "http://purl.org/dc/elements/1.1/"
	nsNone       = empty
	nsRSS        = empty
	nsXML        = "http://www.w3.org/XML/1998/namespace"
)

const (
	flavorAtom = "atom"
	flavorRSS  = "rss"
)

const (
	linkSelf      = "self"
	linkAlternate = "alternate"
)

const (
	logName  = "[parsr]"
	logTrace = "[TRACE]"
)

const (
	empty = ""
)

var rssPerson = regexp.MustCompile(`^(.+)\s+\((.+)\)$`)

// NewParser returns a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse feed
func (z *Parser) Parse(reader io.Reader, contentType string) (*Feed, error) {

	z.decoder = xml.NewDecoder(reader)
	z.decoder.CharsetReader = NewFilterCharsetReader
	z.decoder.Strict = false

	z.stack = &elements{}
	z.feed = nil
	z.entry = nil

	var exitError error

Loop:
	for {

		token, err := z.decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			exitError = err
			break
		}

		switch t := token.(type) {

		case xml.StartElement:
			e := z.stack.Push(t)
			//log.Printf("%-7s %-7s Start %s :: %s", logTrace, logName, e.name.Local, z.stack.String())

			switch {
			case z.feed == nil:
				if err := z.doStartFeedNil(e, &t); err != nil {
					exitError = err
					break Loop
				}

			case z.feed != nil && z.entry == nil && z.stack.IsStackFeed(z.feed.Flavor, 1):
				switch z.feed.Flavor {
				case flavorAtom:
					z.doStartFeedAtom(e, &t)
				case flavorRSS:
					z.doStartFeedRSS(e, &t)
				} // flavor

			case z.entry != nil && z.stack.IsStackEntry(z.feed.Flavor, 1):
				switch z.feed.Flavor {
				case flavorAtom:
					z.doStartEntryAtom(e, &t)
				case flavorRSS:
					z.doStartEntryRSS(e, &t)
				} // flavor

			} // level

		case xml.EndElement:
			e, err := z.stack.PeekIf(t)
			if err != nil {
				exitError = err
				break Loop
			}
			//log.Printf("%-7s %-7s End   %s :: %s", logTrace, logName, e.name.Local, z.stack.String())

			switch {

			case z.feed != nil && z.entry == nil && z.stack.IsStackFeed(z.feed.Flavor, 0):
				switch z.feed.Flavor {
				case flavorAtom:
					z.doEndFeedAtom(e)
				case flavorRSS:
					z.doEndFeedRSS(e)
				}

			case z.entry != nil && z.stack.IsStackEntry(z.feed.Flavor, 0):
				switch z.feed.Flavor {
				case flavorAtom:
					z.doEndEntryAtom(e)
				case flavorRSS:
					z.doEndEntryRSS(e)
				}

			}
			z.stack.Pop() // at the very end of EndElement

		} // switch token

	} // loop

	if exitError != nil {
		ioutil.ReadAll(reader)
	}

	if closer, ok := reader.(io.Closer); ok {
		closer.Close()
	}

	if exitError == nil && z.feed == nil {
		exitError = fmt.Errorf("Cannot parse feed")
	}

	return z.feed, exitError

}

func (z *Parser) doStartFeedNil(e *element, start *xml.StartElement) error {
	if e.Match(nsAtom, "feed") || e.Match(nsRSS, "rss") {
		z.feed = &Feed{}
		z.feed.Links = make(map[string]string)
		switch {
		case e.Match(nsAtom, "feed"):
			z.feed.Flavor = flavorAtom
		case e.Match(nsRSS, "rss"):
			z.feed.Flavor = flavorRSS
		} // switch
	} else {
		return fmt.Errorf("Cannot parse %s:%s", e.name.Space, e.name.Local)
	}
	return nil
}

func (z *Parser) doStartFeedAtom(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsAtom, "author"):
		if value := z.makePersonAtom(e, start); value != empty {
			z.feed.Authors = append(z.feed.Authors, value)
		}
	case e.Match(nsAtom, "entry"):
		z.entry = &Entry{}
		z.entry.Links = make(map[string]string)
	case e.Match(nsAtom, "generator"):
		z.feed.Generator = z.makeGenerator(e, start)
	case e.Match(nsAtom, "icon"):
		if text := z.makeText(e, start); text != empty {
			z.feed.Icon = makeURL(z.stack.Attr(nsXML, "base"), text)
		}
	case e.Match(nsAtom, "id"):
		z.feed.ID = z.makeText(e, start)
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.feed.Links[key] = value
	case e.Match(nsAtom, "rights"):
		z.feed.Rights = z.makeContent(e, start)
	case e.Match(nsAtom, "subtitle"):
		z.feed.Subtitle = z.makeContent(e, start)
	case e.Match(nsAtom, "title"):
		z.feed.Title = z.makeContent(e, start)
		// ignore feed updated, and calculate from entries (see doEndEntryAtom)
		// case e.Match(nsAtom, "updated"):
		// 	z.feed.Updated = z.parseTime(z.makeText(e, start))
	} // z.stack
}

func (z *Parser) doStartFeedRSS(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.feed.Links[key] = value
	case e.Match(nsRSS, "copyright"):
		z.feed.Rights = z.makeText(e, start)
	case e.Match(nsRSS, "description"):
		z.feed.Subtitle = z.makeText(e, start)
	case e.Match(nsRSS, "generator"):
		z.feed.Generator = z.makeText(e, start)
	case e.Match(nsRSS, "image"):
		if image := z.makeRSSImage(e, start); image != nil {
			z.feed.Icon = image.URL
		}
		// ignore feed updated, and calculate from entries (see doEndEntryRSS)
	// case e.Match(nsRSS, "pubdate"):
	// 	z.feed.Updated = z.parseTime(z.makeText(e, start))
	case e.Match(nsRSS, "title"):
		z.feed.Title = z.makeText(e, start)
	case e.Match(nsRSS, "item"):
		z.entry = &Entry{}
		z.entry.Links = make(map[string]string)
	case e.Match(nsRSS, "link"):
		z.feed.Links[linkAlternate] = makeURL(z.stack.Attr(nsXML, "base"), z.makeText(e, start))
	}
}

func (z *Parser) doStartEntryAtom(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsAtom, "author"):
		if value := z.makePersonAtom(e, start); value != empty {
			z.entry.Authors = append(z.entry.Authors, value)
		}
	case e.Match(nsAtom, "category"):
		if value := makeCategory(e, start); value != empty {
			z.entry.Categories = append(z.entry.Categories, value)
		}
	case e.Match(nsAtom, "content"):
		z.entry.Content = z.makeContent(e, start)
	case e.Match(nsAtom, "contributor"):
		if value := z.makePersonAtom(e, start); value != empty {
			z.entry.Contributors = append(z.entry.Contributors, value)
		}
	case e.Match(nsAtom, "id"):
		z.entry.ID = z.makeText(e, start)
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.entry.Links[key] = value
	case e.Match(nsAtom, "published"):
		z.entry.Created = parseTime(z.makeText(e, start))
	case e.Match(nsAtom, "summary"):
		z.entry.Summary = z.makeContent(e, start)
	case e.Match(nsAtom, "title"):
		z.entry.Title = z.makeContent(e, start)
	case e.Match(nsAtom, "updated"):
		if text := z.makeText(e, start); text != empty {
			z.entry.Updated = parseTime(text)
		}
	}
}

func (z *Parser) doStartEntryRSS(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.entry.Links[key] = value
	case e.Match(nsContent, "encoded"):
		z.entry.Content = z.makeText(e, start)
	case e.Match(nsDublinCore, "creator"):
		if creator := z.makeText(e, start); creator != empty {
			z.entry.Authors = append(z.entry.Authors, creator)
		}
	case e.Match(nsDublinCore, "date"):
		if z.entry.Updated.IsZero() {
			if text := z.makeText(e, start); text != empty {
				z.entry.Updated = parseTime(text)
			}
		}
	case e.Match(nsRSS, "author"):
		if value := z.makePersonRSS(e, start); value != empty {
			z.entry.Authors = append(z.entry.Authors, value)
		}
	case e.Match(nsRSS, "category"):
		if value := z.makeText(e, start); value != empty {
			z.entry.Categories = append(z.entry.Categories, value)
		}
	case e.Match(nsRSS, "description"):
		z.entry.Summary = z.makeText(e, start)
	case e.Match(nsRSS, "guid"):
		z.entry.ID = z.makeText(e, start)
	case e.Match(nsRSS, "link"):
		z.entry.Links[linkAlternate] = makeURL(z.stack.Attr(nsXML, "base"), z.makeText(e, start))
	case e.Match(nsRSS, "pubdate"):
		if text := z.makeText(e, start); text != empty {
			z.entry.Updated = parseTime(text)
		}
	case e.Match(nsRSS, "title"):
		z.entry.Title = z.makeText(e, start)
	}
}

func (z *Parser) doEndFeedAtom(e *element) {
	switch {
	case e.Match(nsAtom, "feed"):
		// finished: clean up atom feed here
		z.feed.LinkSelf = z.feed.Links[linkSelf]
		z.feed.LinkAlternate = z.feed.Links[linkAlternate]
		if z.feed.LinkAlternate == "" {
			z.feed.LinkAlternate = z.feed.Links[""]
		}
	}
}

func (z *Parser) doEndFeedRSS(e *element) {
	switch {
	case e.Match(nsRSS, "channel"):
		if z.feed.ID == empty {
			z.feed.ID = z.feed.Links[linkSelf]
		}
		if z.feed.ID == empty {
			z.feed.ID = z.feed.Links[linkAlternate]
		}
		// finished: clean up rss feed here
		z.feed.Flavor = flavorRSS + z.stack.Attr(nsRSS, "version")
		z.feed.LinkSelf = z.feed.Links[linkSelf]
		z.feed.LinkAlternate = z.feed.Links[linkAlternate]
		if z.feed.LinkAlternate == "" {
			z.feed.LinkAlternate = z.feed.Links[""]
		}
	}
}

func (z *Parser) doEndEntryAtom(e *element) {
	switch {
	case e.Match(nsAtom, "entry"):
		if z.entry.Created.IsZero() {
			z.entry.Created = z.entry.Updated
		}
		if z.feed.Updated.Before(z.entry.Updated) {
			z.feed.Updated = z.entry.Updated
		}
		// guarentee entry always has an id
		if z.entry.ID == empty {
			z.entry.ID = z.entry.hash()
		}

		z.entry.LinkSelf = z.entry.Links[linkSelf]
		z.entry.LinkAlternate = z.entry.Links[linkAlternate]
		if z.entry.LinkAlternate == "" {
			z.entry.LinkAlternate = z.entry.Links[""]
		}

		z.feed.Entries = append(z.feed.Entries, z.entry)
		z.entry = nil

	}
}

func (z *Parser) doEndEntryRSS(e *element) {
	switch {
	case e.Match(nsRSS, "item"):
		if z.entry.Summary != empty && z.entry.Content == empty {
			z.entry.Content = z.entry.Summary
			z.entry.Summary = empty
		}
		if z.entry.Created.IsZero() {
			z.entry.Created = z.entry.Updated
		}
		if z.feed.Updated.Before(z.entry.Updated) {
			z.feed.Updated = z.entry.Updated
		}
		// guarentee entry always has an id
		if z.entry.ID == empty {
			z.entry.ID = z.entry.hash()
		}

		z.entry.LinkSelf = z.entry.Links[linkSelf]
		z.entry.LinkAlternate = z.entry.Links[linkAlternate]
		if z.entry.LinkAlternate == "" {
			z.entry.LinkAlternate = z.entry.Links[""]
		}

		z.feed.Entries = append(z.feed.Entries, z.entry)
		z.entry = nil
	}
}

func makeCategory(e *element, start *xml.StartElement) string {
	term := strings.TrimSpace(e.Attr(nsNone, "term"))
	label := strings.TrimSpace(e.Attr(nsNone, "label"))
	if label != empty {
		return label
	}
	return term
}

func (z *Parser) makeContent(e *element, start *xml.StartElement) string {
	x := &content{}
	z.decoder.DecodeElement(x, start)
	z.stack.Pop()
	return x.ToString()
}

func (z *Parser) makeGenerator(e *element, start *xml.StartElement) string {
	result := &generator{}
	z.decoder.DecodeElement(result, start)
	z.stack.Pop()
	return result.ToString()
}

func (z *Parser) makePersonAtom(e *element, start *xml.StartElement) string {
	result := &person{}
	z.decoder.DecodeElement(result, start)
	z.stack.Pop()
	return result.ToString()
}

func (z *Parser) makePersonRSS(e *element, start *xml.StartElement) string {
	x := &text{}
	z.decoder.DecodeElement(x, start)
	z.stack.Pop()
	authorString := x.ToString()
	if matches := rssPerson.FindStringSubmatch(authorString); len(matches) == 3 {
		p := &person{
			Name:  matches[2],
			Email: matches[1],
		}
		return p.ToString()
	}
	return authorString
}

func (z *Parser) makeRSSImage(e *element, start *xml.StartElement) *rssImage {
	x := &rssImage{}
	err := z.decoder.DecodeElement(x, start)
	z.stack.Pop()
	if err == nil {
		return x
	}
	return nil
}

func (z *Parser) makeText(e *element, start *xml.StartElement) string {
	x := &text{}
	z.decoder.DecodeElement(x, start)
	z.stack.Pop()
	return x.ToString()
}

func makeURL(base string, urlstr string) string {
	u, err := url.Parse(urlstr)
	if err == nil {
		if base != empty && !u.IsAbs() {
			b, err := url.Parse(base)
			if err == nil {
				return b.ResolveReference(u).String()
			}
		}
	}
	return urlstr
}

// taken from https://github.com/jteeuwen/go-pkg-rss/ timedecoder.go
func parseTime(formatted string) (t time.Time) {
	var layouts = [...]string{
		"Mon, _2 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		"Mon, 2, Jan 2006 15:4",
		"02 Jan 2006 15:04:05 MST",
		"_2 Jan 2006 15:04:05 +0000", // found in the wild, differs slightly from RFC822Z
		"2006-01-02 15:04:05",        // found in the wild, apparently RFC3339 was too difficult
		"_2 Jan 2006",
		"2006-01-02",
	}
	formatted = strings.TrimSpace(formatted)
	for _, layout := range layouts {
		t, _ = time.Parse(layout, formatted)
		if !t.IsZero() {
			break
		}
	}
	return
}
