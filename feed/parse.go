package feed

import (
	"encoding/xml"
	"log"
	"net/http"
)

// Parse feed
func Parse(feedURL string) (*Feed, error) {

	rsp, err := http.Get(feedURL)
	if err != nil {
		panic(err)
	}

	reader := rsp.Body
	defer reader.Close()
	decoder := xml.NewDecoder(reader)

	var feed *Feed

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
				feed = a.toFeed()
			} else {
				log.Printf("Unknown feed type: %s\n", element.Name.Local)
			}
		} // switch

	} // for loop

	return feed, nil

}
