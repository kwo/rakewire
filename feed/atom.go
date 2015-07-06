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
	Entries []atomEntry `xml:"entry"`
}

// atom entry
type atomEntry struct {
	ID     string     `xml:"id"`
	Date   time.Time  `xml:"updated"`
	Author atomAuthor `xml:"author"`
	Title  string     `xml:"title"`
}

// atom author
type atomAuthor struct {
	Name  string `xml:"name"`
	EMail string `xml:"email"`
	URI   string `xml:"uri"`
}

func (a atomFeed) toFeed() *Feed {

	var f Feed

	f.ID = a.ID
	f.Title = a.Title
	f.Date = &a.Date
	f.Author = &Author{a.Author.Name, a.Author.EMail, a.Author.URI}

	for i := 0; i < len(a.Entries); i++ {
		b := a.Entries[i]
		e := Entry{}
		e.ID = b.ID
		e.Title = b.Title
		e.Date = &b.Date
		e.Author = &Author{b.Author.Name, b.Author.EMail, b.Author.URI}
		f.Entries = append(f.Entries, &e)
	}

	return &f
}
