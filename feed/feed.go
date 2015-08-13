package feed

import (
	"encoding/xml"
	"github.com/aaron-lebo/ocd/feeds/atom"
	"io"
	"log"
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
	Updated   *time.Time
}

// Entry entry
type Entry struct {
	Authors    []*Person
	Categories []string
	Content    string
	Created    *time.Time
	ID         string
	Links      map[string]string
	Rights     string
	Summary    string
	Title      string
	Updated    *time.Time
}

// Person person
type Person struct {
	Name  string
	EMail string
	URI   string
}

// Parse feed
func Parse(reader io.Reader) (*Feed, error) {

	decoder := xml.NewDecoder(reader)

	var feed *Feed
	var err error

	for {

		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		// Inspect the type of the token just read.
		switch element := t.(type) {
		case xml.StartElement:
			if element.Name.Local == "feed" {
				a := &atom.Feed{}
				decoder.DecodeElement(a, &element)
				feed, err = atomToFeed(a)
			} else if element.Name.Local == "rss" {
				r := &rssFeed{}
				decoder.DecodeElement(r, &element)
				decoder.DefaultSpace = "rss"
				feed, err = r.toFeed()
			} else {
				log.Printf("Unknown feed type: %s\n", element.Name.Local)
			}
		} // switch

	} // for loop

	return feed, err

}
