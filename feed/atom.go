package feed

import (
	"rakewire.com/logging"
	"time"
)

var (
	logger = logging.New("feed")
)

// atom feed
type atomFeed struct {
	Author       atomAuthor     `xml:"author"`
	Categories   []atomCategory `xml:"category"`
	Contributors []atomAuthor   `xml:"contributor"`
	Entries      []atomEntry    `xml:"entry"`
	Generator    atomGenerator  `xml:"generator"`
	Icon         string         `xml:"icon"`
	ID           string         `xml:"id"`
	Links        []atomLink     `xml:"link"`
	Logo         string         `xml:"logo"`
	Rights       string         `xml:"rights"`
	Subtitle     string         `xml:"subtitle"`
	Title        string         `xml:"title"`
	Updated      time.Time      `xml:"updated"`
}

// atom entry
type atomEntry struct {
	Author       atomAuthor     `xml:"author"`
	Categories   []atomCategory `xml:"category"`
	Content      atomText       `xml:"content"`
	Contributors []atomAuthor   `xml:"contributor"`
	ID           string         `xml:"id"`
	Links        []atomLink     `xml:"link"`
	Published    time.Time      `xml:"published"`
	Rights       string         `xml:"rights"`
	Summary      atomText       `xml:"summary"`
	Title        string         `xml:"title"`
	Updated      time.Time      `xml:"updated"`
}

// atom author
type atomAuthor struct {
	EMail string `xml:"email"`
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
}

// atom link
type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// atom text
type atomText struct {
	URI  string `xml:"uri,attr"`
	Text string `xml:",chardata"`
	//Type string `xml:"type,attr"`
	XML string `xml:",innerxml"`
}

func (txt atomText) String() string {

	var result string

	if txt.Text != "" {
		result = txt.Text
	} else if txt.XML != "" {
		result = txt.XML
	} else if txt.URI != "" {
		result = txt.URI
	}

	return result

}

// atom generator
type atomGenerator struct {
	URI     string `xml:"uri,attr"`
	Name    string `xml:",chardata"`
	Version string `xml:"version,attr"`
}

func (g atomGenerator) String() string {

	var result string

	if g.Name != "" {
		result = g.Name
		if g.Version != "" {
			result += " " + g.Version
		}
		if g.URI != "" {
			result += " (" + g.URI + ")"
		}
	}

	return result

}

// atom category
type atomCategory struct {
	Label  string `xml:"label,attr"`
	Term   string `xml:",chardata"`
	Scheme string `xml:"scheme,attr"`
}

func (c atomCategory) String() string {
	return c.Term
}

func (a atomFeed) toFeed() (*Feed, error) {

	f := &Feed{}

	// #DOING:0 if no feed updated, populate with latest entry

	f.ID = a.ID
	f.Title = a.Title
	f.Subtitle = a.Subtitle
	f.Author = Person{a.Author.Name, a.Author.EMail, a.Author.URI}
	f.Icon = a.Icon
	f.Generator = a.Generator.String()
	f.Flavor = "atom"

	if !a.Updated.IsZero() {
		f.Updated = a.Updated
	}

	f.Links = make(map[string]string)
	for _, link := range a.Links {
		f.Links[link.Rel] = link.Href
	}

	for _, atomEntry := range a.Entries {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.ID = atomEntry.ID
		entry.Title = atomEntry.Title
		entry.Created = atomEntry.Published
		entry.Updated = atomEntry.Updated
		entry.Author = Person{atomEntry.Author.EMail, atomEntry.Author.Name, atomEntry.Author.URI}

		entry.Links = make(map[string]string)
		for _, link := range atomEntry.Links {
			entry.Links[link.Rel] = link.Href
		}

		for _, cat := range atomEntry.Categories {
			entry.Categories = append(entry.Categories, cat.String())
		}

		entry.Summary = atomEntry.Summary.String()
		entry.Content = atomEntry.Content.String()

	} // loop

	if f.Updated.IsZero() && len(f.Entries) > 0 {
		f.Updated = f.Entries[0].Updated
	}

	return f, nil

}
