package feed

// Feed feed
type Feed struct {
	ID      string
	Title   string
	Date    string
	Author  Author
	Entries []Entry
}

// Entry entry
type Entry struct {
	ID     string
	Date   string
	Author Author
	Title  string
}

// Author author
type Author struct {
	Name  string
	EMail string
	URI   string
}
