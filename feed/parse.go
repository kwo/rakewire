package feed

import (
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
)

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
				var a atomFeed
				decoder.DecodeElement(&a, &element)
				feed, err = a.toFeed()
			} else if element.Name.Local == "rss" {
				var r rssFeed
				decoder.DecodeElement(&r, &element)
				if r.Version == "2.0" {
					feed, err = r.Channel.toFeed()
				} else {
					err = errors.New("Not an RSS2 feed")
				}
			} else {
				log.Printf("Unknown feed type: %s\n", element.Name.Local)
			}
		} // switch

	} // for loop

	return feed, err

}

// ParseURL download url and parse feed
func ParseURL(feedURL string) (*Feed, error) {

	rsp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}

	reader := rsp.Body
	defer reader.Close()

	return Parse(reader)

}
