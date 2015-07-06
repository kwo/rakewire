package feed

import (
	"time"
)

// atom feed
type atomFeed struct {
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Date    time.Time   `xml:"updated"`
	Author  atomAuthor  `xml:"author"`
	Rights  string      `xml:"rights"`
	Links   []atomLink  `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

// atom entry
type atomEntry struct {
	ID      string     `xml:"id"`
	Date    time.Time  `xml:"updated"`
	Author  atomAuthor `xml:"author"`
	Title   string     `xml:"title"`
	Links   []atomLink `xml:"link"`
	Summary atomText   `xml:"summary"`
	Content atomText   `xml:"content"`
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

func (a atomFeed) toFeed() *Feed {

	var f Feed

	f.ID = a.ID
	f.Title = a.Title
	f.Date = &a.Date
	f.Author = &Author{a.Author.Name, a.Author.EMail, a.Author.URI}

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
		entry.Date = &atomEntry.Date
		entry.Author = &Author{atomEntry.Author.Name, atomEntry.Author.EMail, atomEntry.Author.URI}

		for j := 0; j < len(atomEntry.Links); j++ {
			atomLink := atomEntry.Links[j]
			link := Link{atomLink.Rel, atomLink.Href}
			entry.Links = append(entry.Links, &link)
		}

		entry.Summary = atomEntry.Summary.String()
		entry.Content = atomEntry.Content.String()

	}

	return &f

}
