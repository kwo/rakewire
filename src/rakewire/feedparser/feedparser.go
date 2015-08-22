package feedparser

import (
	"encoding/xml"
	"github.com/rogpeppe/go-charset/charset"
	// required by go-charset
	_ "github.com/rogpeppe/go-charset/data"
	"io"
	"net/url"
	"rakewire/logging"
	"strings"
	"time"
)

// #DOING:30 implement stepping decoder

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

// Entry z.entry
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

// Status hold the status of a parse iteration
type Status struct {
	Feed    *Feed
	Entry   *Entry
	HasMore bool
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
	logger = logging.Null("feeddecoder")
)

// Parse feed
func (z *Parser) Parse(reader io.Reader) (*Feed, error) {

	// #TODO:0 attach used charset to feed object
	// #TODO:0 if decoder exits in error, the response body won't be read fully which will render a http connection un-reusable
	// #DOING:10 break Parse up into smaller functions

	z.decoder = xml.NewDecoder(reader)
	z.decoder.CharsetReader = charset.NewReader
	z.decoder.Strict = false

	z.stack = &elements{}

	for {

		token, err := z.decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := token.(type) {

		case xml.StartElement:
			e := z.stack.Push(t)
			logger.Printf("Start %s :: %s\n", e.name.Local, z.stack.String())

			switch {
			case z.feed == nil:

				if e.Match(nsAtom, "feed") || e.Match(nsRSS, "rss") {
					z.feed = &Feed{}
					z.feed.Links = make(map[string]string)
					switch {
					case e.Match(nsAtom, "feed"):
						z.feed.Flavor = flavorAtom
					case e.Match(nsRSS, "rss"):
						z.feed.Flavor = flavorRSS
					} // switch
				} // if

			case z.feed != nil && z.entry == nil && z.stack.IsStackFeed(z.feed.Flavor, 1):

				switch z.feed.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := z.makePerson(e, &t); value != "" {
							z.feed.Authors = append(z.feed.Authors, value)
						}
						z.stack.Pop()
					case e.Match(nsAtom, "entry"):
						z.entry = &Entry{}
						z.entry.Links = make(map[string]string)
					case e.Match(nsAtom, "generator"):
						z.feed.Generator = z.makeGenerator(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "icon"):
						if text := z.makeText(e, &t); text != "" {
							z.feed.Icon = z.makeURL(z.stack.Attr(nsXML, "base"), text)
						}
						z.stack.Pop()
					case e.Match(nsAtom, "id"):
						z.feed.ID = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						z.feed.Links[key] = value
					case e.Match(nsAtom, "rights"):
						z.feed.Rights = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "subtitle"):
						z.feed.Subtitle = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "title"):
						z.feed.Title = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "updated"):
						z.feed.Updated = z.parseTime(z.makeText(e, &t))
						z.stack.Pop()
					} // z.stack

				case flavorRSS:
					switch {
					case e.Match(nsRSS, "generator"):
						z.feed.Generator = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsRSS, "guid"):
						z.feed.ID = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsRSS, "pubdate"):
						if z.feed.Updated.IsZero() {
							z.feed.Updated = z.parseTime(z.makeText(e, &t))
						}
						z.stack.Pop()
					case e.Match(nsRSS, "title"):
						z.feed.Title = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsRSS, "item"):
						z.entry = &Entry{}
						z.entry.Links = make(map[string]string)
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						z.feed.Links[key] = value
					} // z.stack

				} // flavor

			case z.entry != nil && z.stack.IsStackEntry(z.feed.Flavor, 1):
				switch z.feed.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := z.makePerson(e, &t); value != "" {
							z.entry.Authors = append(z.entry.Authors, value)
						}
						z.stack.Pop()
					case e.Match(nsAtom, "category"):
						if value := z.makeCategory(e, &t); value != "" {
							z.entry.Categories = append(z.entry.Categories, value)
						}
					case e.Match(nsAtom, "content"):
						z.entry.Content = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "contributor"):
						if value := z.makePerson(e, &t); value != "" {
							z.entry.Contributors = append(z.entry.Contributors, value)
						}
						z.stack.Pop()
					case e.Match(nsAtom, "id"):
						z.entry.ID = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(z.stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						z.entry.Links[key] = value
					case e.Match(nsAtom, "published"):
						z.entry.Created = z.parseTime(z.makeText(e, &t))
						z.stack.Pop()
					case e.Match(nsAtom, "summary"):
						z.entry.Summary = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "title"):
						z.entry.Title = z.makeContent(e, &t)
						z.stack.Pop()
					case e.Match(nsAtom, "updated"):
						if text := z.makeText(e, &t); text != "" {
							z.entry.Updated = z.parseTime(text)
							if z.entry.Updated.After(z.feed.Updated) {
								z.feed.Updated = z.entry.Updated
							}
						}
						z.stack.Pop()
					} // z.stack
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "guid"):
						z.entry.ID = z.makeText(e, &t)
						z.stack.Pop()
					case e.Match(nsRSS, "pubdate"):
						if text := z.makeText(e, &t); text != "" {
							if z.entry.Updated.IsZero() {
								z.entry.Updated = z.parseTime(text)
								if z.entry.Updated.After(z.feed.Updated) {
									z.feed.Updated = z.entry.Updated
								}
							}
						}
						z.stack.Pop()
					} // z.stack

				} // flavor

			} // level

		case xml.EndElement:
			e, err := z.stack.PeekIf(t)
			if err != nil {
				return nil, err
			}
			logger.Printf("End   %s :: %s\n", e.name.Local, z.stack.String())

			switch {

			case z.feed != nil && z.entry == nil && z.stack.IsStackFeed(z.feed.Flavor, 0):
				switch z.feed.Flavor {
				case flavorAtom:
					switch {
					case e.Match(nsAtom, "feed"):
						// finished: clean up atom feed here
					}
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "channel"):
						// #DOING:20 more possibilities for IDs
						if z.feed.ID == "" {
							z.feed.ID = z.feed.Links["self"]
						}
						// finished: clean up rss feed here
						z.feed.Flavor = flavorRSS + z.stack.Attr(nsRSS, "version")
					}
				}

			case z.entry != nil && z.stack.IsStackEntry(z.feed.Flavor, 0):
				switch z.feed.Flavor {
				case flavorAtom:
					switch {
					case e.Match(nsAtom, "entry"):
						z.feed.Entries = append(z.feed.Entries, z.entry)
						z.entry = nil
					}
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "item"):
						z.feed.Entries = append(z.feed.Entries, z.entry)
						z.entry = nil
					}
				}

			}
			z.stack.Pop() // at the very end of EndElement

		} // switch token

	} // loop

	return z.feed, nil

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
	return x.ToString()
}

func (z *Parser) makeGenerator(e *element, start *xml.StartElement) string {
	result := &generator{}
	z.decoder.DecodeElement(result, start)
	return result.ToString()
}

func (z *Parser) makePerson(e *element, start *xml.StartElement) string {
	result := &person{}
	z.decoder.DecodeElement(result, start)
	return result.ToString()
}

func (z *Parser) makeText(e *element, start *xml.StartElement) string {
	x := &text{}
	z.decoder.DecodeElement(x, start)
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
