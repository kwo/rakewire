package feed

import (
	"time"
)

// Feed feed
type Feed struct {
	ID        string
	Title     string
	Subtitle  string
	Updated   *time.Time
	Author    *Author
	Icon      string
	Rights    string
	Generator string
	Links     []*Link
	Entries   []*Entry
}

// Entry entry
type Entry struct {
	ID         string
	Created    *time.Time
	Updated    *time.Time
	Author     *Author
	Title      string
	Categories []string
	Links      []*Link
	Summary    string
	Content    string
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
