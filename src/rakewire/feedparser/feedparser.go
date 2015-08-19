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

// Namespace constants in lowercase
const (
	NSAtom = "http://www.w3.org/2005/atom" // "http://www.w3.org/2005/Atom"
	NSNone = ""
	NSRss  = ""
	NSXml  = "http://www.w3.org/xml/1998/namespace" // "http://www.w3.org/XML/1998/namespace"
)

// Parse feed
func Parse(reader io.Reader) (*Feed, error) {

	// #TODO:0 attach used charset to feed object

	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReader
	parser.Strict = false

	var n xml.Name
	var attr map[xml.Name]string
	var baseURL string
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
			n = makeLowerName(t.Name)
			attr = mapAttr(t.Attr)
			switch {

			case entry == nil && n.Space == NSAtom && n.Local == "feed":
				f = &Feed{}
				f.Flavor = "atom"
				f.Links = make(map[string]string)
				baseURL = getAttr(attr, NSXml, "base")
			case entry == nil && n.Space == NSRss && n.Local == "rss":
				f = &Feed{}
				f.Flavor = "rss" + getAttr(attr, NSRss, "version")
				f.Links = make(map[string]string)
				baseURL = getAttr(attr, NSXml, "base")

			case entry == nil && n.Space == NSAtom && n.Local == "entry":
				entry = &Entry{}
				entry.Links = make(map[string]string)
			case entry == nil && n.Space == NSRss && n.Local == "item":
				entry = &Entry{}
				entry.Links = make(map[string]string)

			case n.Space == NSAtom && n.Local == "link":
				if entry == nil {
					f.Links[getAttr(attr, NSNone, "rel")] = makeURL(baseURL, getAttr(attr, NSNone, "href"))
				} else {
					entry.Links[getAttr(attr, NSNone, "rel")] = makeURL(baseURL, getAttr(attr, NSNone, "href"))
				}
			case n.Space == NSAtom && n.Local == "author":
				perso = &person{}

			}

		case xml.EndElement:
			n2 := makeLowerName(t.Name)
			switch {
			case n2.Space == NSAtom && n2.Local == "feed":
				// send feed
			case n2.Space == NSRss && n2.Local == "rss":
				// send feed
			case n2.Space == NSAtom && n2.Local == "entry":
				f.Entries = append(f.Entries, entry)
				entry = nil
			case n2.Space == NSRss && n2.Local == "item":
				f.Entries = append(f.Entries, entry)
				entry = nil
			case perso != nil && n2.Space == NSAtom && n2.Local == "author":
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

				case entry != nil && n.Space == NSAtom && n.Local == "content":
					entry.Content = makeText(text, attr)

				case entry == nil && n.Space == NSAtom && n.Local == "generator":
					f.Generator = makeGenerator(text, attr)

				case entry == nil && n.Space == NSAtom && n.Local == "icon":
					f.Icon = makeURL(baseURL, text)

				case n.Space == NSAtom && n.Local == "id":
					if entry == nil {
						f.ID = text
					} else {
						entry.ID = text
					}

				case n.Space == NSRss && n.Local == "pubdate":
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

				case entry == nil && n.Space == NSAtom && n.Local == "rights":
					f.Rights = makeText(text, attr)

				case entry == nil && n.Space == NSAtom && n.Local == "subtitle":
					f.Subtitle = makeText(text, attr)

				case entry != nil && n.Space == NSAtom && n.Local == "summary":
					entry.Summary = makeText(text, attr)

				case n.Space == NSAtom && n.Local == "title":
					if entry == nil {
						f.Title = makeText(text, attr)
					} else {
						entry.Title = makeText(text, attr)
					}

				case n.Space == NSAtom && n.Local == "updated":
					if entry == nil {
						f.Updated = parseTime(text)
					} else {
						entry.Updated = parseTime(text)
						if entry.Updated.After(f.Updated) {
							f.Updated = entry.Updated
						}
					}

				case perso != nil && n.Space == NSAtom && n.Local == "email":
					perso.email = text

				case perso != nil && n.Space == NSAtom && n.Local == "name":
					perso.name = text

				case perso != nil && n.Space == NSAtom && n.Local == "uri":
					perso.uri = text

				}

			} // end

		} // switch token

	} // loop

	return f, nil

}

func getAttr(attrs map[xml.Name]string, space string, local string) string {
	return attrs[xml.Name{Space: space, Local: local}]
}

func mapAttr(attrs []xml.Attr) map[xml.Name]string {
	result := make(map[xml.Name]string)
	for _, attr := range attrs {
		result[makeLowerName(attr.Name)] = attr.Value
	}
	return result
}

func makeGenerator(text string, attrs map[xml.Name]string) string {
	result := text
	if result != "" {
		if version := getAttr(attrs, NSNone, "version"); version != "" {
			result += " " + version
		}
		if uri := getAttr(attrs, NSNone, "uri"); uri != "" {
			result += " (" + uri + ")"
		}
	}
	return result
}

func makeLowerName(n xml.Name) xml.Name {
	return xml.Name{
		Local: strings.ToLower(n.Local),
		Space: strings.ToLower(n.Space),
	}
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

func makeText(text string, attr map[xml.Name]string) Text {
	result := Text{Text: text, Type: getAttr(attr, NSNone, "type")}
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
