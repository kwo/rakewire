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

// #DOING:10 implement stepping parser

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
	stack  *elements
	entry  *Entry
	feed   *Feed
	parser *xml.Decoder
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
	logger = logging.Null("feedparser")
)

// Parse feed
func (z *Parser) Parse(reader io.Reader) (*Feed, error) {

	// #TODO:0 attach used charset to feed object
	// #TODO:0 if parser exits in error, the response body won't be read fully which will render a http connection un-reusable

	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReader
	parser.Strict = false

	var f *Feed
	var entry *Entry
	stack := &elements{}

	for {

		token, err := parser.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := token.(type) {

		case xml.StartElement:
			e := stack.Push(t)
			logger.Printf("Start %t %s :: %s\n", f == nil, e.name.Local, stack.String())

			switch {
			case f == nil:

				if e.Match(nsAtom, "feed") || e.Match(nsRSS, "rss") {
					f = &Feed{}
					f.Links = make(map[string]string)
					switch {
					case e.Match(nsAtom, "feed"):
						f.Flavor = flavorAtom
					case e.Match(nsRSS, "rss"):
						f.Flavor = flavorRSS
					} // switch
				} // if

			case f != nil && entry == nil && stack.IsStackFeed(1):

				switch f.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := z.makePerson(e, &t, parser); value != "" {
							f.Authors = append(f.Authors, value)
						}
						stack.Pop()
					case e.Match(nsAtom, "entry"):
						entry = &Entry{}
						entry.Links = make(map[string]string)
					case e.Match(nsAtom, "generator"):
						f.Generator = z.makeGenerator(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "icon"):
						if text := z.makeText(e, &t, parser); text != "" {
							f.Icon = z.makeURL(stack.Attr(nsXML, "base"), text)
						}
						stack.Pop()
					case e.Match(nsAtom, "id"):
						f.ID = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						f.Links[key] = value
					case e.Match(nsAtom, "rights"):
						f.Rights = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "subtitle"):
						f.Subtitle = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "title"):
						f.Title = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "updated"):
						f.Updated = z.parseTime(z.makeText(e, &t, parser))
						stack.Pop()
					} // stack

				case flavorRSS:
					switch {
					case e.Match(nsRSS, "generator"):
						f.Generator = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsRSS, "guid"):
						f.ID = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsRSS, "pubdate"):
						if f.Updated.IsZero() {
							f.Updated = z.parseTime(z.makeText(e, &t, parser))
						}
						stack.Pop()
					case e.Match(nsRSS, "title"):
						f.Title = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsRSS, "item"):
						entry = &Entry{}
						entry.Links = make(map[string]string)
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						f.Links[key] = value
					} // stack

				} // flavor

			case entry != nil && stack.IsStackEntry(1):
				switch f.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := z.makePerson(e, &t, parser); value != "" {
							entry.Authors = append(entry.Authors, value)
						}
						stack.Pop()
					case e.Match(nsAtom, "category"):
						if value := z.makeCategory(e); value != "" {
							entry.Categories = append(entry.Categories, value)
						}
					case e.Match(nsAtom, "content"):
						entry.Content = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "contributor"):
						if value := z.makePerson(e, &t, parser); value != "" {
							entry.Contributors = append(entry.Contributors, value)
						}
						stack.Pop()
					case e.Match(nsAtom, "id"):
						entry.ID = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := z.makeURL(stack.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						entry.Links[key] = value
					case e.Match(nsAtom, "published"):
						entry.Created = z.parseTime(z.makeText(e, &t, parser))
						stack.Pop()
					case e.Match(nsAtom, "summary"):
						entry.Summary = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "title"):
						entry.Title = z.makeContent(e, &t, parser)
						stack.Pop()
					case e.Match(nsAtom, "updated"):
						if text := z.makeText(e, &t, parser); text != "" {
							entry.Updated = z.parseTime(text)
							if entry.Updated.After(f.Updated) {
								f.Updated = entry.Updated
							}
						}
						stack.Pop()
					} // stack
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "guid"):
						entry.ID = z.makeText(e, &t, parser)
						stack.Pop()
					case e.Match(nsRSS, "pubdate"):
						if text := z.makeText(e, &t, parser); text != "" {
							if entry.Updated.IsZero() {
								entry.Updated = z.parseTime(text)
								if entry.Updated.After(f.Updated) {
									f.Updated = entry.Updated
								}
							}
						}
						stack.Pop()
					} // stack

				} // flavor

			} // level

		case xml.EndElement:
			e, err := stack.PeekIf(t)
			if err != nil {
				return nil, err
			}
			logger.Printf("End   %t %s :: %s\n", stack.IsStackFeed(), e.name.Local, stack.String())

			switch {

			case f != nil && entry == nil && stack.IsStackFeed():
				switch f.Flavor {
				case flavorAtom:
					switch {
					case e.Match(nsAtom, "feed"):
						// finished: clean up atom feed here
					}
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "channel"):
						// #DOING:20 more possibilities for IDs
						if f.ID == "" {
							f.ID = f.Links["self"]
						}
						// finished: clean up rss feed here
						f.Flavor = flavorRSS + stack.Attr(nsRSS, "version")
					}
				}

			case entry != nil && stack.IsStackEntry():
				switch f.Flavor {
				case flavorAtom:
					switch {
					case e.Match(nsAtom, "entry"):
						f.Entries = append(f.Entries, entry)
						entry = nil
					}
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "item"):
						f.Entries = append(f.Entries, entry)
						entry = nil
					}
				}

			}
			stack.Pop() // at the very end of EndElement

		} // switch token

	} // loop

	return f, nil

}

func (z *Parser) makeCategory(e *element) string {
	term := strings.TrimSpace(e.Attr(nsNone, "term"))
	label := strings.TrimSpace(e.Attr(nsNone, "label"))
	if label != "" {
		return label
	}
	return term
}

func (z *Parser) makeContent(e *element, start *xml.StartElement, decoder *xml.Decoder) string {
	x := &content{}
	decoder.DecodeElement(x, start)
	return x.ToString()
}

func (z *Parser) makeGenerator(e *element, start *xml.StartElement, decoder *xml.Decoder) string {
	result := &generator{}
	decoder.DecodeElement(result, start)
	return result.ToString()
}

func (z *Parser) makePerson(e *element, start *xml.StartElement, decoder *xml.Decoder) string {
	result := &person{}
	decoder.DecodeElement(result, start)
	return result.ToString()
}

func (z *Parser) makeText(e *element, start *xml.StartElement, decoder *xml.Decoder) string {
	x := &text{}
	decoder.DecodeElement(x, start)
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

// taken from https://github.com/jteeuwen/go-pkg-rss/ timeparser.go
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
