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

// #DOING: switch times back to pointer

// Feed feed
type Feed struct {
	Authors   []string
	Entries   []*Entry
	Flavor    string
	Generator string
	Icon      string
	ID        string
	Links     map[string]string
	Rights    *Text
	Subtitle  *Text
	Title     *Text
	Updated   time.Time
}

// Entry entry
type Entry struct {
	Authors      []string
	Categories   []string
	Content      *Text
	Contributors []string
	Created      time.Time
	ID           string
	Links        map[string]string
	Summary      *Text
	Title        *Text
	Updated      time.Time
}

// Text text
type Text struct {
	Type string
	Text string
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
func Parse(reader io.Reader) (*Feed, error) {

	// #TODO:0 attach used charset to feed object

	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReader
	parser.Strict = false

	var f *Feed
	var entry *Entry
	elements := &Elements{}

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
			e := elements.Push(t)
			logger.Printf("Start %t %s :: %s\n", f == nil, e.name.Local, elements.String())

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

			case f != nil && entry == nil && elements.IsStackFeed(1):

				switch f.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := makePerson(e, &t, parser).String(); value != "" {
							f.Authors = append(f.Authors, value)
						}
						elements.Pop()
					case e.Match(nsAtom, "entry"):
						entry = &Entry{}
						entry.Links = make(map[string]string)
					case e.Match(nsAtom, "generator"):
						if text := makeString(e, &t, parser); text != "" {
							f.Generator = makeGenerator(text, e)
						}
						elements.Pop()
					case e.Match(nsAtom, "icon"):
						if text := makeString(e, &t, parser); text != "" {
							f.Icon = makeURL(elements.Attr(nsXML, "base"), text)
						}
						elements.Pop()
					case e.Match(nsAtom, "id"):
						f.ID = makeString(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := makeURL(elements.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						f.Links[key] = value
					case e.Match(nsAtom, "rights"):
						f.Rights = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "subtitle"):
						f.Subtitle = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "title"):
						f.Title = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "updated"):
						if text := makeString(e, &t, parser); text != "" {
							f.Updated = parseTime(text)
						}
						elements.Pop()
					} // elements

				case flavorRSS:
					switch {
					case e.Match(nsRSS, "generator"):
						if text := makeString(e, &t, parser); text != "" {
							f.Generator = text
						}
						elements.Pop()
					case e.Match(nsRSS, "guid"):
						if text := makeString(e, &t, parser); text != "" {
							f.ID = text
						}
						elements.Pop()
					case e.Match(nsRSS, "pubdate"):
						if text := makeString(e, &t, parser); text != "" {
							if f.Updated.IsZero() {
								f.Updated = parseTime(text)
							}
						}
						elements.Pop()
					case e.Match(nsRSS, "title"):
						if text := makeString(e, &t, parser); text != "" {
							f.Title = &Text{Type: "text", Text: text}
						}
						elements.Pop()
					case e.Match(nsRSS, "item"):
						entry = &Entry{}
						entry.Links = make(map[string]string)
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := makeURL(elements.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						f.Links[key] = value
					} // elements

				} // flavor

			case entry != nil && elements.IsStackEntry(1):
				switch f.Flavor {

				case flavorAtom:
					switch {
					case e.Match(nsAtom, "author"):
						if value := makePerson(e, &t, parser).String(); value != "" {
							entry.Authors = append(entry.Authors, value)
						}
						elements.Pop()
					case e.Match(nsAtom, "category"):
						if value := makeCategory(e); value != "" {
							entry.Categories = append(entry.Categories, value)
						}
					case e.Match(nsAtom, "content"):
						entry.Content = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "contributor"):
						if value := makePerson(e, &t, parser).String(); value != "" {
							entry.Contributors = append(entry.Contributors, value)
						}
						elements.Pop()
					case e.Match(nsAtom, "id"):
						entry.ID = makeString(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "link"):
						key := e.Attr(nsNone, "rel")
						value := makeURL(elements.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
						entry.Links[key] = value
					case e.Match(nsAtom, "published"):
						if text := makeString(e, &t, parser); text != "" {
							entry.Created = parseTime(text)
						}
						elements.Pop()
					case e.Match(nsAtom, "summary"):
						entry.Summary = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "title"):
						entry.Title = makeText(e, &t, parser)
						elements.Pop()
					case e.Match(nsAtom, "updated"):
						if text := makeString(e, &t, parser); text != "" {
							entry.Updated = parseTime(text)
							if entry.Updated.After(f.Updated) {
								f.Updated = entry.Updated
							}
						}
						elements.Pop()
					} // elements
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "guid"):
						if text := makeString(e, &t, parser); text != "" {
							entry.ID = text
						}
						elements.Pop()
					case e.Match(nsRSS, "pubdate"):
						if text := makeString(e, &t, parser); text != "" {
							if entry.Updated.IsZero() {
								entry.Updated = parseTime(text)
								if entry.Updated.After(f.Updated) {
									f.Updated = entry.Updated
								}
							}
						}
						elements.Pop()
					} // elements

				} // flavor

			} // level

		case xml.EndElement:
			e, err := elements.PeekIf(t)
			if err != nil {
				return nil, err
			}
			logger.Printf("End   %t %s :: %s\n", elements.IsStackFeed(), e.name.Local, elements.String())

			switch {

			case f != nil && entry == nil && elements.IsStackFeed():
				switch f.Flavor {
				case flavorAtom:
					switch {
					case e.Match(nsAtom, "feed"):
						// finished: clean up atom feed here
					}
				case flavorRSS:
					switch {
					case e.Match(nsRSS, "channel"):
						// #DOING:0 more possibilities for IDs
						if f.ID == "" {
							f.ID = f.Links["self"]
						}
						// finished: clean up rss feed here
						f.Flavor = flavorRSS + elements.Attr(nsRSS, "version")
					}
				}

			case entry != nil && elements.IsStackEntry():
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
			elements.Pop() // at the very end of EndElement

		} // switch token

	} // loop

	return f, nil

}

func makeCategory(e *Element) string {
	term := strings.TrimSpace(e.Attr(nsNone, "term"))
	label := strings.TrimSpace(e.Attr(nsNone, "label"))
	if label != "" {
		return label
	}
	return term
}

func makeGenerator(text string, e *Element) string {
	result := text
	if result != "" {
		if version := e.Attr(nsNone, "version"); version != "" {
			result += " " + version
		}
		if uri := e.Attr(nsNone, "uri"); uri != "" {
			result += " (" + uri + ")"
		}
	}
	return result
}

func makePerson(e *Element, start *xml.StartElement, decoder *xml.Decoder) *person {
	result := &person{}
	decoder.DecodeElement(result, start)
	return result
}

func makeString(e *Element, start *xml.StartElement, decoder *xml.Decoder) string {
	x := &xmlString{}
	decoder.DecodeElement(x, start)
	return x.Text
}

func makeText(e *Element, start *xml.StartElement, decoder *xml.Decoder) *Text {
	x := &xmlText{}
	err := decoder.DecodeElement(x, start)
	if err == nil {
		return x.GetText()
	}
	return &Text{}
}

func makeURL(base string, urlstr string) string {
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
