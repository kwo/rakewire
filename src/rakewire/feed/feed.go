package feed

import (
	"bytes"
	"encoding/xml"
	"github.com/kwo/ocd/feeds/atom"
	"github.com/kwo/ocd/feeds/rss"
	"github.com/rogpeppe/go-charset/charset"
	// required by go-charset
	_ "github.com/rogpeppe/go-charset/data"
	"time"
)

// Feed feed
type Feed struct {
	Authors   []*Person
	Entries   []*Entry
	Flavor    string
	Generator string
	Icon      string
	ID        string
	Links     map[string]string
	Rights    string
	Subtitle  string
	Title     string
	Updated   time.Time
}

// Entry entry
type Entry struct {
	Authors    []*Person
	Categories []string
	Content    string
	Created    time.Time
	ID         string
	Links      map[string]string
	Summary    string
	Title      string
	Updated    time.Time
}

// Person person
type Person struct {
	Name  string
	EMail string
	URI   string
}

// Parse feed
func Parse(body []byte) (feed *Feed, err error) {

	// #TODO:0 attach used charset to feed object

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReader
	a := &atom.Feed{}
	err = decoder.Decode(a)
	if err == nil {
		return atomToFeed(a)
	}

	decoder = xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReader
	decoder.DefaultSpace = "rss"
	r := &rss.Rss{}
	err = decoder.Decode(r)
	if err == nil {
		return rssToFeed(r)
	}

	return

}
