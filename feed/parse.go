package feed

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// Parse feed
func Parse(feedURL string) {

	rsp, err := http.Get(feedURL)
	if err != nil {
		panic(err)
	}

	fmt.Printf("URL: %s\n", rsp.Request.URL)

	reader := rsp.Body
	defer reader.Close()
	decoder := xml.NewDecoder(reader)

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
				var feed atomfeed
				decoder.DecodeElement(&feed, &element)
				fmt.Printf("%s, %s, %s, %s\n", feed.Title, feed.ID, feed.Date, feed.Author.Name)
				for _, entry := range feed.Entries {
					fmt.Printf("%s, %s, %s, %s\n", entry.Title, entry.ID, entry.Date, entry.Author.Name)
				}
			} else {
				log.Printf("Unknown feed type: %s\n", element.Name.Local)
			}
		} // switch

	} // for loop

}
