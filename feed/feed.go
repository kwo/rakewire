package feed

import (
	"time"
)

// Feed feed
type Feed struct {
	ID      string
	Title   string
	Date    *time.Time
	Author  *Author
	Links   []*Link
	Entries []*Entry
}

// Entry entry
type Entry struct {
	ID     string
	Date   *time.Time
	Author *Author
	Title  string
	Links  []*Link
}

// Author author
type Author struct {
	Name  string
	EMail string
	URI   string
}

// Link link
type Link struct {
	Rel  string
	Href string
}
