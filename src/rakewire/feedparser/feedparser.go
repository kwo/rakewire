package feedparser

import (
	"encoding/xml"
	"fmt"
	"github.com/rogpeppe/go-charset/charset"
	// required by go-charset
	_ "github.com/rogpeppe/go-charset/data"
	"io"
	"io/ioutil"
	"net/url"
	"rakewire/logging"
	"strings"
	"time"
)

// Feed feed
type Feed struct {
	Authors   []string
	Entries   []*Entry
	Flavor    string
	Generator string
	Icon      string
	ID        string
	Links     map[string]string
	Rights    string
	Subtitle  string
	Title     string
	Updated   time.Time
}

// Entry entry
type Entry struct {
	Authors      []string
	Categories   []string
	Content      string
	Contributors []string
	Created      time.Time
	ID           string
	Links        map[string]string
	Summary      string
	Title        string
	Updated      time.Time
}

// Parser can parse feeds
type Parser struct {
	decoder *xml.Decoder
	entry   *Entry
	feed    *Feed
	stack   *elements
}

const (
	nsAtom = "http://www.w3.org/2005/Atom"
	nsNone = ""
	nsRSS  = ""
	nsXML  = "http://www.w3.org/XML/1998/namespace"
)

const (
	flavorAtom = "atom"
	flavorRSS  = "rss"
)

var (
	logger = logging.Null("feedparser")
)

// NewParser returns a new parser
func NewParser() *Parser {
	p := &Parser{}
	return p
}

// Parse feed
func (z *Parser) Parse(reader io.ReadCloser) (*Feed, error) {

	// #DOING:50 attach used charset to feed object

	z.decoder = xml.NewDecoder(reader)
	z.decoder.CharsetReader = charset.NewReader
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
			logger.Printf("Start %s :: %s\n", e.name.Local, z.stack.String())

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
			logger.Printf("End   %s :: %s\n", e.name.Local, z.stack.String())

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

	logger.Printf("exitError: %s", exitError)

	if exitError != nil {
		ioutil.ReadAll(reader)
	}
	reader.Close()

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
		if value := z.makePerson(e, start); value != "" {
			z.feed.Authors = append(z.feed.Authors, value)
		}
	case e.Match(nsAtom, "entry"):
		z.entry = &Entry{}
		z.entry.Links = make(map[string]string)
	case e.Match(nsAtom, "generator"):
		z.feed.Generator = z.makeGenerator(e, start)
	case e.Match(nsAtom, "icon"):
		if text := z.makeText(e, start); text != "" {
			z.feed.Icon = z.makeURL(z.stack.Attr(nsXML, "base"), text)
		}
	case e.Match(nsAtom, "id"):
		z.feed.ID = z.makeText(e, start)
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.feed.Links[key] = value
	case e.Match(nsAtom, "rights"):
		z.feed.Rights = z.makeContent(e, start)
	case e.Match(nsAtom, "subtitle"):
		z.feed.Subtitle = z.makeContent(e, start)
	case e.Match(nsAtom, "title"):
		z.feed.Title = z.makeContent(e, start)
	case e.Match(nsAtom, "updated"):
		z.feed.Updated = z.parseTime(z.makeText(e, start))
	} // z.stack
}

func (z *Parser) doStartFeedRSS(e *element, start *xml.StartElement) {
	// #DOING:0 finish RSS parser
	switch {
	case e.Match(nsRSS, "description"):
		z.feed.Subtitle = z.makeText(e, start)
	case e.Match(nsRSS, "generator"):
		z.feed.Generator = z.makeText(e, start)
	case e.Match(nsRSS, "guid"):
		z.feed.ID = z.makeText(e, start)
	case e.Match(nsRSS, "pubdate"):
		z.feed.Updated = z.parseTime(z.makeText(e, start))
	case e.Match(nsRSS, "title"):
		z.feed.Title = z.makeText(e, start)
	case e.Match(nsRSS, "item"):
		z.entry = &Entry{}
		z.entry.Links = make(map[string]string)
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.feed.Links[key] = value
	case e.Match(nsRSS, "link"):
		z.feed.Links["alternate"] = z.makeText(e, start)
	} // z.stack
}

func (z *Parser) doStartEntryAtom(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsAtom, "author"):
		if value := z.makePerson(e, start); value != "" {
			z.entry.Authors = append(z.entry.Authors, value)
		}
	case e.Match(nsAtom, "category"):
		if value := z.makeCategory(e, start); value != "" {
			z.entry.Categories = append(z.entry.Categories, value)
		}
	case e.Match(nsAtom, "content"):
		z.entry.Content = z.makeContent(e, start)
	case e.Match(nsAtom, "contributor"):
		if value := z.makePerson(e, start); value != "" {
			z.entry.Contributors = append(z.entry.Contributors, value)
		}
	case e.Match(nsAtom, "id"):
		z.entry.ID = z.makeText(e, start)
	case e.Match(nsAtom, "link"):
		key := e.Attr(nsNone, "rel")
		value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
		z.entry.Links[key] = value
	case e.Match(nsAtom, "published"):
		z.entry.Created = z.parseTime(z.makeText(e, start))
	case e.Match(nsAtom, "summary"):
		z.entry.Summary = z.makeContent(e, start)
	case e.Match(nsAtom, "title"):
		z.entry.Title = z.makeContent(e, start)
	case e.Match(nsAtom, "updated"):
		if text := z.makeText(e, start); text != "" {
			z.entry.Updated = z.parseTime(text)
			if z.entry.Updated.After(z.feed.Updated) {
				z.feed.Updated = z.entry.Updated
			}
		}
	} // z.stack
}

func (z *Parser) doStartEntryRSS(e *element, start *xml.StartElement) {
	switch {
	case e.Match(nsRSS, "guid"):
		z.entry.ID = z.makeText(e, start)
	case e.Match(nsRSS, "pubdate"):
		if text := z.makeText(e, start); text != "" {
			if z.entry.Updated.IsZero() {
				z.entry.Updated = z.parseTime(text)
				if z.entry.Updated.After(z.feed.Updated) {
					z.feed.Updated = z.entry.Updated
				}
			}
		}
	} // z.stack
}

func (z *Parser) doEndFeedAtom(e *element) {
	switch {
	case e.Match(nsAtom, "feed"):
		// finished: clean up atom feed here
	}
}

func (z *Parser) doEndFeedRSS(e *element) {
	switch {
	case e.Match(nsRSS, "channel"):
		// #DOING:10 more possibilities for IDs
		if z.feed.ID == "" {
			z.feed.ID = z.feed.Links["self"]
		}
		// finished: clean up rss feed here
		z.feed.Flavor = flavorRSS + z.stack.Attr(nsRSS, "version")
	}
}

func (z *Parser) doEndEntryAtom(e *element) {
	switch {
	case e.Match(nsAtom, "entry"):
		z.feed.Entries = append(z.feed.Entries, z.entry)
		z.entry = nil
	}
}

func (z *Parser) doEndEntryRSS(e *element) {
	switch {
	case e.Match(nsRSS, "item"):
		z.feed.Entries = append(z.feed.Entries, z.entry)
		z.entry = nil
	}
}

func (z *Parser) makeCategory(e *element, start *xml.StartElement) string {
	term := strings.TrimSpace(e.Attr(nsNone, "term"))
	label := strings.TrimSpace(e.Attr(nsNone, "label"))
	if label != "" {
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

func (z *Parser) makePerson(e *element, start *xml.StartElement) string {
	result := &person{}
	z.decoder.DecodeElement(result, start)
	z.stack.Pop()
	return result.ToString()
}

func (z *Parser) makeText(e *element, start *xml.StartElement) string {
	x := &text{}
	z.decoder.DecodeElement(x, start)
	z.stack.Pop()
	return x.ToString()
}

func (z *Parser) makeURL(base string, urlstr string) string {
	u, err := url.Parse(urlstr)
	if err == nil {
		if base != "" && !u.IsAbs() {
			b, err := url.Parse(base)
			if err == nil {
				return b.ResolveReference(u).String()
			}
		}
	}
	return urlstr
}

// taken from https://github.com/jteeuwen/go-pkg-rss/ timedecoder.go
func (z *Parser) parseTime(formatted string) (t time.Time) {
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
