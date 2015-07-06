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
	ID     string     `xml:"id"`
	Date   time.Time  `xml:"updated"`
	Author atomAuthor `xml:"author"`
	Title  string     `xml:"title"`
	Links  []atomLink `xml:"link"`
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

	entry := Entry{}

	for i := 0; i < len(a.Entries); i++ {
		atomEntry := a.Entries[i]

		for j := 0; j < len(atomEntry.Links); j++ {
			atomLink := atomEntry.Links[j]
			link := Link{atomLink.Rel, atomLink.Href}
			entry.Links = append(entry.Links, &link)
		}

		entry.ID = atomEntry.ID
		entry.Title = atomEntry.Title
		entry.Date = &atomEntry.Date
		entry.Author = &Author{atomEntry.Author.Name, atomEntry.Author.EMail, atomEntry.Author.URI}
		f.Entries = append(f.Entries, &entry)
	}

	return &f
}
