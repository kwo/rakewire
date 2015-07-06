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
	Entries []*Entry
}

// Entry entry
type Entry struct {
	ID     string
	Date   *time.Time
	Author *Author
	Title  string
}

// Author author
type Author struct {
	Name  string
	EMail string
	URI   string
}
