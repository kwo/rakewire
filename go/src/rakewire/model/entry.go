package model

import (
	"time"
)

// Entry from a feed
type Entry struct {
	ID       string    `json:"id" kv:"key"`
	EntryID  string    `json:"-" kv:"EntryID:2"`
	FeedID   string    `json:"feedId" kv:"Date:1,EntryID:1"`
	Created  time.Time `json:"created" kv:"Date:2"`
	Updated  time.Time `json:"updated"`
	URL      string    `json:"url"`
	Author   string    `json:"author"`
	Title    string    `json:"title"`
	Contents string    `json:"contents"`
}
