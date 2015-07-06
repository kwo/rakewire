package feed

// atom feed
type atomfeed struct {
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Date    string      `xml:"updated"`
	Author  atomauthor  `xml:"author"`
	Entries []atomentry `xml:"entry"`
}

// atom entry
type atomentry struct {
	ID     string     `xml:"id"`
	Date   string     `xml:"updated"`
	Author atomauthor `xml:"author"`
	Title  string     `xml:"title"`
}

// atom author
type atomauthor struct {
	Name  string `xml:"name"`
	EMail string `xml:"email"`
	URI   string `xml:"uri"`
}
