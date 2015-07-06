package feed

import (
	"time"
)

// Feed feed
type Feed struct {
	Author    *Author
	Entries   []*Entry
	Generator string
	Icon      string
	ID        string
	Links     []*Link
	Title     string
	Subtitle  string
	Updated   *time.Time
}

// Entry entry
type Entry struct {
	Author     *Author
	Categories []string
	Content    string
	Created    *time.Time
	ID         string
	Links      []*Link
	Summary    string
	Title      string
	Updated    *time.Time
}

// Author author
type Author struct {
	EMail string
	Name  string
	URI   string
}

// Link link
type Link struct {
	Href string
	Rel  string
}
