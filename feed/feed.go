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
	Links     map[string]string
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
	Links      map[string]string
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
