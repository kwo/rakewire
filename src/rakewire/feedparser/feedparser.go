package feedparser

import (
	"encoding/xml"
	"github.com/rogpeppe/go-charset/charset"
	// required by go-charset
	_ "github.com/rogpeppe/go-charset/data"
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
	Authors      []string
	Categories   []string
	Content      Text
	Contributors []string
	Created      time.Time
	ID           string
	Links        map[string]string
	Summary      Text
	Title        Text
	Updated      time.Time
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

type divelement struct {
	Div xmltext `xml:"div"`
}
type xmltext struct {
	Text string `xml:",innerxml"`
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

	var entry *Entry
	var perso *person
	elements := &Elements{}
	f := &Feed{}
	f.Links = make(map[string]string)

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
			//fmt.Printf("push: %s\n", e.name.Local)
			switch {

			case elements.On(NSAtom, "feed"):
				f.Flavor = "atom"
			case elements.On(NSRss, "rss"):
				f.Flavor = "rss" + elements.Peek().Attr(NSRss, "version")

			case elements.On(NSAtom, "entry") && elements.In(NSAtom, "feed"):
				entry = &Entry{}
				entry.Links = make(map[string]string)
			case elements.On(NSRss, "item") && elements.In(NSRss, "channel"):
				entry = &Entry{}
				entry.Links = make(map[string]string)

			case elements.On(NSAtom, "author") || elements.On(NSAtom, "contributor"):
				perso = &person{}

			case elements.On(NSAtom, "category") && (elements.In(NSAtom, "entry") || elements.On(NSRss, "item")):
				e := elements.Peek()
				value := makeCategory(e.Attr(NSNone, "term"), e.Attr(NSNone, "label"))
				if value != "" {
					entry.Categories = append(entry.Categories, value)
				}
			case elements.On(NSAtom, "content") && (elements.In(NSAtom, "entry") || elements.On(NSRss, "item")):
				entry.Content = makeTextFromXML(parser, elements.Peek(), &t)
			case elements.On(NSAtom, "link"):
				e := elements.Peek()
				key := e.Attr(NSNone, "rel")
				value := makeURL(elements.Attr(NSXml, "base"), e.Attr(NSNone, "href"))
				if elements.In(NSAtom, "feed") || elements.In(NSRss, "channel") {
					f.Links[key] = value
				} else if elements.In(NSAtom, "entry") || elements.On(NSRss, "item") {
					entry.Links[key] = value
				}

			}

		case xml.EndElement:
			e, _ := elements.Pop(t)
			//fmt.Printf("pop:  %s\n", e.name.Local)
			switch {
			case e.Match(NSAtom, "feed"):
				// send feed
			case e.Match(NSRss, "rss"):
				// #DOING:0 more possibilities for IDs
				if f.ID == "" {
					f.ID = f.Links["self"]
				}
				// send feed
			case e.Match(NSAtom, "entry") || e.Match(NSRss, "item"):
				f.Entries = append(f.Entries, entry)
				entry = nil
			case e.Match(NSAtom, "author"):
				if elements.On(NSAtom, "feed") || elements.On(NSRss, "channel") {
					f.Authors = append(f.Authors, makePerson(perso))
				} else if elements.On(NSAtom, "entry") || elements.On(NSRss, "item") {
					entry.Authors = append(entry.Authors, makePerson(perso))
				}
				perso = nil
			case e.Match(NSAtom, "contributor"):
				if elements.On(NSAtom, "entry") || elements.On(NSRss, "item") {
					entry.Contributors = append(entry.Contributors, makePerson(perso))
				}
				perso = nil
			}

		case xml.CharData:
			//fmt.Printf("char: %t %s\n", entry == nil, elements.Peek())
			text := strings.TrimSpace(string([]byte(t)))

			switch {

			case text == "":
				continue

			case elements.In(NSAtom, "feed") || elements.In(NSRss, "channel"):
				switch {
				case elements.On(NSAtom, "generator"):
					f.Generator = makeGenerator(text, elements.Peek())
				case elements.On(NSRss, "generator"):
					f.Generator = text
				case elements.On(NSRss, "guid"):
					f.ID = text
				case elements.On(NSAtom, "icon"):
					f.Icon = makeURL(elements.Attr(NSXml, "base"), text)
				case elements.On(NSAtom, "id"):
					f.ID = text
				case elements.On(NSRss, "pubdate"):
					if f.Updated.IsZero() {
						f.Updated = parseTime(text)
					}
				// #DOING:0 replace all Text type nodes with StartElement method
				case elements.On(NSAtom, "rights"):
					f.Rights = makeText(text, elements.Peek())
				case elements.On(NSAtom, "subtitle"):
					f.Subtitle = makeText(text, elements.Peek())
				case elements.On(NSAtom, "title"):
					f.Title = makeText(text, elements.Peek())
				case elements.On(NSRss, "title"):
					f.Title = Text{Type: "text", Text: text}
				case elements.On(NSAtom, "updated"):
					f.Updated = parseTime(text)
				} // feed switch

			case elements.In(NSAtom, "entry") || elements.In(NSRss, "item"):
				switch {
				case elements.On(NSAtom, "content"):
					// don't overwrite content if taken from xhtml element
					if entry.Content.Text == "" {
						entry.Content = makeText(text, elements.Peek())
					}
				case elements.On(NSAtom, "id"):
					entry.ID = text
				case elements.On(NSRss, "pubdate"):
					if entry.Updated.IsZero() {
						entry.Updated = parseTime(text)
						if entry.Updated.After(f.Updated) {
							f.Updated = entry.Updated
						}
					}
				case elements.On(NSAtom, "published"):
					entry.Created = parseTime(text)
				case elements.On(NSAtom, "summary"):
					entry.Summary = makeText(text, elements.Peek())
				case elements.On(NSAtom, "title"):
					entry.Title = makeText(text, elements.Peek())
				case elements.On(NSAtom, "updated"):
					entry.Updated = parseTime(text)
					if entry.Updated.After(f.Updated) {
						f.Updated = entry.Updated
					}
				} // entry switch

			case elements.In(NSAtom, "author") || elements.In(NSAtom, "contributor"):
				switch {
				case elements.On(NSAtom, "email"):
					perso.email = text
				case elements.On(NSAtom, "name"):
					perso.name = text
				case elements.On(NSAtom, "uri"):
					perso.uri = text
				} // author switch

			} // switch CharData

		} // switch token

	} // loop

	return f, nil

}

func makeCategory(term string, label string) string {
	term = strings.TrimSpace(term)
	label = strings.TrimSpace(label)
	if label != "" {
		return label
	}
	return term
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
	} else if result.Type == "xhtml" {

	}
	return result
}

func makeTextFromXML(decoder *xml.Decoder, e *Element, start *xml.StartElement) Text {
	result := Text{Type: e.Attr(NSNone, "type")}
	if result.Type == "xhtml" {
		x := &divelement{}
		err := decoder.DecodeElement(x, start)
		if err != nil {
			result.Type = ""
		} else {
			result.Type = "html"
			result.Text = strings.TrimSpace(x.Div.Text)
			// #TODO:0 use base to fix relative HREFs in XML
		}
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
