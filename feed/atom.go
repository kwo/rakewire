package feed

import (
	"time"
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

	var f Feed

	f.ID = a.ID
	f.Title = a.Title
	f.Subtitle = a.Subtitle
	f.Author = &Person{a.Author.Name, a.Author.EMail, a.Author.URI}
	f.Icon = a.Icon
	f.Generator = a.Generator.String()

	if !a.Updated.IsZero() {
		f.Updated = &a.Updated
	}

	f.Links = make(map[string]string)
	for j := 0; j < len(a.Links); j++ {
		f.Links[a.Links[j].Rel] = a.Links[j].Href
	}

	for i := 0; i < len(a.Entries); i++ {

		entry := Entry{}
		atomEntry := a.Entries[i]
		f.Entries = append(f.Entries, &entry)

		entry.ID = atomEntry.ID
		entry.Title = atomEntry.Title
		if !atomEntry.Published.IsZero() {
			entry.Created = &atomEntry.Published
		}
		if !atomEntry.Updated.IsZero() {
			entry.Updated = &atomEntry.Updated
		}
		entry.Author = &Person{atomEntry.Author.EMail, atomEntry.Author.Name, atomEntry.Author.URI}

		entry.Links = make(map[string]string)
		for j := 0; j < len(atomEntry.Links); j++ {
			entry.Links[atomEntry.Links[j].Rel] = atomEntry.Links[j].Href
		}

		for j := 0; j < len(atomEntry.Categories); j++ {
			entry.Categories = append(entry.Categories, atomEntry.Categories[j].String())
		}

		entry.Summary = atomEntry.Summary.String()
		entry.Content = atomEntry.Content.String()

	}

	return &f, nil

}
