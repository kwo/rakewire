package feed

import (
	"time"
)

// atom feed
type atomFeed struct {
	ID        string        `xml:"id"`
	Title     string        `xml:"title"`
	Subtitle  string        `xml:"subtitle"`
	Updated   time.Time     `xml:"updated"`
	Author    atomAuthor    `xml:"author"`
	Icon      string        `xml:"icon"`
	Generator atomGenerator `xml:"generator"`
	Rights    string        `xml:"rights"`
	Links     []atomLink    `xml:"link"`
	Entries   []atomEntry   `xml:"entry"`
}

// atom entry
type atomEntry struct {
	ID         string         `xml:"id"`
	Created    time.Time      `xml:"published"`
	Updated    time.Time      `xml:"updated"`
	Author     atomAuthor     `xml:"author"`
	Title      string         `xml:"title"`
	Categories []atomCategory `xml:"category"`
	Links      []atomLink     `xml:"link"`
	Summary    atomText       `xml:"summary"`
	Content    atomText       `xml:"content"`
}

// atom author
type atomAuthor struct {
	Name  string `xml:"name"`
	EMail string `xml:"email"`
	URI   string `xml:"uri"`
}

// atom link
type atomLink struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// atom text
type atomText struct {
	Text string `xml:",chardata"`
	XML  string `xml:",innerxml"`
	URI  string `xml:"uri,attr"`
	//Type string `xml:"type,attr"`
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
	Name    string `xml:",chardata"`
	URI     string `xml:"uri,attr"`
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
	Term   string `xml:",chardata"`
	Scheme string `xml:"scheme,attr"`
	Label  string `xml:"label,attr"`
}

func (c atomCategory) String() string {
	return c.Term
}

func (a atomFeed) toFeed() *Feed {

	var f Feed

	f.ID = a.ID
	f.Title = a.Title
	f.Subtitle = a.Subtitle
	f.Author = &Author{a.Author.Name, a.Author.EMail, a.Author.URI}
	f.Icon = a.Icon
	f.Rights = a.Rights
	f.Generator = a.Generator.String()

	if !a.Updated.IsZero() {
		f.Updated = &a.Updated
	}

	for j := 0; j < len(a.Links); j++ {
		atomLink := a.Links[j]
		link := Link{atomLink.Rel, atomLink.Href}
		f.Links = append(f.Links, &link)
	}

	for i := 0; i < len(a.Entries); i++ {

		entry := Entry{}
		atomEntry := a.Entries[i]
		f.Entries = append(f.Entries, &entry)

		entry.ID = atomEntry.ID
		entry.Title = atomEntry.Title
		if !atomEntry.Created.IsZero() {
			entry.Created = &atomEntry.Created
		}
		if !atomEntry.Updated.IsZero() {
			entry.Updated = &atomEntry.Updated
		}
		entry.Author = &Author{atomEntry.Author.Name, atomEntry.Author.EMail, atomEntry.Author.URI}

		for j := 0; j < len(atomEntry.Links); j++ {
			atomLink := atomEntry.Links[j]
			link := Link{atomLink.Rel, atomLink.Href}
			entry.Links = append(entry.Links, &link)
		}

		for j := 0; j < len(atomEntry.Categories); j++ {
			entry.Categories = append(entry.Categories, atomEntry.Categories[j].String())
		}

		entry.Summary = atomEntry.Summary.String()
		entry.Content = atomEntry.Content.String()

	}

	return &f

}
