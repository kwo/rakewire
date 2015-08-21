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

const (
	nsAtom = "http://www.w3.org/2005/Atom"
	nsNone = ""
	nsRSS  = ""
	nsXML  = "http://www.w3.org/XML/1998/namespace"
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

			case elements.On(nsAtom, "feed"):
				f.Flavor = "atom"
			case elements.On(nsRSS, "rss"):
				f.Flavor = "rss" + elements.Peek().Attr(nsRSS, "version")

			case elements.On(nsAtom, "entry") && elements.In(nsAtom, "feed"):
				entry = &Entry{}
				entry.Links = make(map[string]string)
			case elements.On(nsRSS, "item") && elements.In(nsRSS, "channel"):
				entry = &Entry{}
				entry.Links = make(map[string]string)

			case elements.On(nsAtom, "author") || elements.On(nsAtom, "contributor"):
				perso = &person{}

			case elements.On(nsAtom, "category") && (elements.In(nsAtom, "entry") || elements.On(nsRSS, "item")):
				e := elements.Peek()
				value := makeCategory(e.Attr(nsNone, "term"), e.Attr(nsNone, "label"))
				if value != "" {
					entry.Categories = append(entry.Categories, value)
				}
			case elements.On(nsAtom, "content") && (elements.In(nsAtom, "entry") || elements.On(nsRSS, "item")):
				entry.Content = makeTextFromXML(parser, elements.Peek(), &t)
			case elements.On(nsAtom, "link"):
				e := elements.Peek()
				key := e.Attr(nsNone, "rel")
				value := makeURL(elements.Attr(nsXML, "base"), e.Attr(nsNone, "href"))
				if elements.In(nsAtom, "feed") || elements.In(nsRSS, "channel") {
					f.Links[key] = value
				} else if elements.In(nsAtom, "entry") || elements.On(nsRSS, "item") {
					entry.Links[key] = value
				}

			}

		case xml.EndElement:
			e, _ := elements.Pop(t)
			//fmt.Printf("pop:  %s\n", e.name.Local)
			switch {
			case e.Match(nsAtom, "feed"):
				// send feed
			case e.Match(nsRSS, "rss"):
				// #DOING:0 more possibilities for IDs
				if f.ID == "" {
					f.ID = f.Links["self"]
				}
				// send feed
			case e.Match(nsAtom, "entry") || e.Match(nsRSS, "item"):
				f.Entries = append(f.Entries, entry)
				entry = nil
			case e.Match(nsAtom, "author"):
				if elements.On(nsAtom, "feed") || elements.On(nsRSS, "channel") {
					f.Authors = append(f.Authors, makePerson(perso))
				} else if elements.On(nsAtom, "entry") || elements.On(nsRSS, "item") {
					entry.Authors = append(entry.Authors, makePerson(perso))
				}
				perso = nil
			case e.Match(nsAtom, "contributor"):
				if elements.On(nsAtom, "entry") || elements.On(nsRSS, "item") {
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

			case elements.In(nsAtom, "feed") || elements.In(nsRSS, "channel"):
				switch {
				case elements.On(nsAtom, "generator"):
					f.Generator = makeGenerator(text, elements.Peek())
				case elements.On(nsRSS, "generator"):
					f.Generator = text
				case elements.On(nsRSS, "guid"):
					f.ID = text
				case elements.On(nsAtom, "icon"):
					f.Icon = makeURL(elements.Attr(nsXML, "base"), text)
				case elements.On(nsAtom, "id"):
					f.ID = text
				case elements.On(nsRSS, "pubdate"):
					if f.Updated.IsZero() {
						f.Updated = parseTime(text)
					}
				// #DOING:0 replace all Text type nodes with StartElement method
				case elements.On(nsAtom, "rights"):
					f.Rights = makeText(text, elements.Peek())
				case elements.On(nsAtom, "subtitle"):
					f.Subtitle = makeText(text, elements.Peek())
				case elements.On(nsAtom, "title"):
					f.Title = makeText(text, elements.Peek())
				case elements.On(nsRSS, "title"):
					f.Title = Text{Type: "text", Text: text}
				case elements.On(nsAtom, "updated"):
					f.Updated = parseTime(text)
				} // feed switch

			case elements.In(nsAtom, "entry") || elements.In(nsRSS, "item"):
				switch {
				case elements.On(nsAtom, "content"):
					// don't overwrite content if taken from xhtml element
					if entry.Content.Text == "" {
						entry.Content = makeText(text, elements.Peek())
					}
				case elements.On(nsAtom, "id"):
					entry.ID = text
				case elements.On(nsRSS, "pubdate"):
					if entry.Updated.IsZero() {
						entry.Updated = parseTime(text)
						if entry.Updated.After(f.Updated) {
							f.Updated = entry.Updated
						}
					}
				case elements.On(nsAtom, "published"):
					entry.Created = parseTime(text)
				case elements.On(nsAtom, "summary"):
					entry.Summary = makeText(text, elements.Peek())
				case elements.On(nsAtom, "title"):
					entry.Title = makeText(text, elements.Peek())
				case elements.On(nsAtom, "updated"):
					entry.Updated = parseTime(text)
					if entry.Updated.After(f.Updated) {
						f.Updated = entry.Updated
					}
				} // entry switch

			case elements.In(nsAtom, "author") || elements.In(nsAtom, "contributor"):
				switch {
				case elements.On(nsAtom, "email"):
					perso.email = text
				case elements.On(nsAtom, "name"):
					perso.name = text
				case elements.On(nsAtom, "uri"):
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
		if version := e.Attr(nsNone, "version"); version != "" {
			result += " " + version
		}
		if uri := e.Attr(nsNone, "uri"); uri != "" {
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
	result := Text{Text: text, Type: e.Attr(nsNone, "type")}
	if result.Type == "" {
		result.Type = "text"
	} else if result.Type == "xhtml" {

	}
	return result
}

func makeTextFromXML(decoder *xml.Decoder, e *Element, start *xml.StartElement) Text {
	result := Text{Type: e.Attr(nsNone, "type")}
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
