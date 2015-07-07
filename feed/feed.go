package feed

import (
	"time"
)

// Feed feed
type Feed struct {
	Author    *Person
	Entries   []*Entry
	Flavor    string
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
	Author     *Person
	Categories []string
	Content    string
	Created    *time.Time
	ID         string
	Links      map[string]string
	Summary    string
	Title      string
	Updated    *time.Time
}

// Person person
type Person struct {
	EMail string
	Name  string
	URI   string
}
