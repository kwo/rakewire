package feedparser

import (
	"code.google.com/p/go-charset/charset"
	// required by go-charset
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"io"
	"net/url"
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
	Rights    Text
	Subtitle  Text
	Title     Text
	Updated   time.Time
}

// Entry entry
type Entry struct {
	Authors    []string
	Categories []string
	Content    Text
	Created    time.Time
	ID         string
	Links      map[string]string
	Summary    Text
	Title      Text
	Updated    time.Time
}

// Text text
type Text struct {
	Type string
	Text string
}

type person struct {
	name  string
	email string
	uri   string
}

// Namespace constants
const (
	NSAtom = "http://www.w3.org/2005/Atom"
	NSNone = ""
	NSRss  = ""
	NSXml  = "http://www.w3.org/XML/1998/namespace"
)

// Parse feed
func Parse(reader io.Reader) (*Feed, error) {

	// #TODO:0 attach used charset to feed object

	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReader
	parser.Strict = false

	elements := &Elements{}
	var f *Feed
	var entry *Entry
	var perso *person

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
			elements.Push(t)
			switch {

			case elements.On(NSAtom, "feed"):
				f = &Feed{}
				f.Flavor = "atom"
				f.Links = make(map[string]string)
			case elements.On(NSRss, "rss"):
				f = &Feed{}
				f.Flavor = "rss" + elements.Peek().Attr(NSRss, "version")
				f.Links = make(map[string]string)

			case elements.In(NSAtom, "feed") && elements.On(NSAtom, "entry"):
				entry = &Entry{}
				entry.Links = make(map[string]string)
			case elements.In(NSRss, "rss") && elements.On(NSRss, "item"):
				entry = &Entry{}
				entry.Links = make(map[string]string)

			case elements.On(NSAtom, "link"):
				e := elements.Peek()
				key := e.Attr(NSNone, "rel")
				value := makeURL(elements.Attr(NSXml, "base"), e.Attr(NSNone, "href"))
				if entry == nil {
					f.Links[key] = value
				} else {
					entry.Links[key] = value
				}
			case elements.On(NSAtom, "author"):
				perso = &person{}

			}

		case xml.EndElement:
			n2, _ := elements.Pop(t)
			switch {
			case n2.Match(NSAtom, "feed"):
				// send feed
			case n2.Match(NSRss, "rss"):
				// send feed
			case n2.Match(NSAtom, "entry"):
				f.Entries = append(f.Entries, entry)
				entry = nil
			case n2.Match(NSRss, "item"):
				f.Entries = append(f.Entries, entry)
				entry = nil
			case perso != nil && n2.Match(NSAtom, "author"):
				if entry == nil {
					f.Authors = append(f.Authors, makePerson(perso))
				} else {
					entry.Authors = append(entry.Authors, makePerson(perso))
				}
				perso = nil
			}

		case xml.CharData:
			text := strings.TrimSpace(string([]byte(t)))
			if text == "" {
				continue
			} else {

				// feed header
				switch {

				case elements.In(NSAtom, "entry") && elements.On(NSAtom, "content"):
					entry.Content = makeText(text, elements.Peek())

				case elements.In(NSAtom, "feed") && elements.On(NSAtom, "generator"):
					f.Generator = makeGenerator(text, elements.Peek())

				case elements.In(NSAtom, "feed") && elements.On(NSAtom, "icon"):
					f.Icon = makeURL(elements.Attr(NSXml, "base"), text)

				case elements.On(NSAtom, "id"):
					if elements.In(NSAtom, "feed") {
						f.ID = text
					} else if elements.In(NSAtom, "entry") {
						entry.ID = text
					}

				case elements.On(NSRss, "pubdate"):
					if entry == nil {
						if f.Updated.IsZero() {
							f.Updated = parseTime(text)
						}
					} else {
						if entry.Updated.IsZero() {
							entry.Updated = parseTime(text)
							if entry.Updated.After(f.Updated) {
								f.Updated = entry.Updated
							}
						}
					}

				case elements.In(NSAtom, "feed") && elements.On(NSAtom, "rights"):
					f.Rights = makeText(text, elements.Peek())

				case elements.In(NSAtom, "feed") && elements.On(NSAtom, "subtitle"):
					f.Subtitle = makeText(text, elements.Peek())

				case elements.In(NSAtom, "entry") && elements.On(NSAtom, "summary"):
					entry.Summary = makeText(text, elements.Peek())

				case elements.On(NSAtom, "title"):
					if elements.In(NSAtom, "feed") {
						f.Title = makeText(text, elements.Peek())
					} else if elements.In(NSAtom, "entry") {
						entry.Title = makeText(text, elements.Peek())
					}

				case elements.On(NSAtom, "updated"):
					if entry == nil {
						f.Updated = parseTime(text)
					} else {
						entry.Updated = parseTime(text)
						if entry.Updated.After(f.Updated) {
							f.Updated = entry.Updated
						}
					}

				case elements.In(NSAtom, "author") && elements.On(NSAtom, "email"):
					perso.email = text

				case elements.In(NSAtom, "author") && elements.On(NSAtom, "name"):
					perso.name = text

				case elements.In(NSAtom, "author") && elements.On(NSAtom, "uri"):
					perso.uri = text

				}

			} // end

		} // switch token

	} // loop

	return f, nil

}

func makeGenerator(text string, e *Element) string {
	result := text
	if result != "" {
		if version := e.Attr(NSNone, "version"); version != "" {
			result += " " + version
		}
		if uri := e.Attr(NSNone, "uri"); uri != "" {
			result += " (" + uri + ")"
		}
	}
	return result
}

func makePerson(p *person) string {
	var result string
	if p != nil {
		result = p.name
		if p.email != "" {
			result += " <" + p.email + ">"
		}
		if p.uri != "" {
			result += " (" + p.uri + ")"
		}
	}
	return result
}

func makeText(text string, e *Element) Text {
	result := Text{Text: text, Type: e.Attr(NSNone, "type")}
	if result.Type == "" {
		result.Type = "text"
	}
	return result
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
