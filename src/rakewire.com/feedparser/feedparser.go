package feedparser

import (
	"code.google.com/p/go-charset/charset"
	// required by go-charset
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"fmt"
	"io"
	"rakewire.com/feed"
	"strings"
)

// Namespace constants
const (
	NSAtom = "http://www.w3.org/2005/Atom"
)

// Parse feed
func Parse(reader io.Reader) (feed *feed.Feed, err error) {

	// #TODO:0 attach used charset to feed object

	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReader
	parser.Strict = false

	for {

		token, err := parser.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := token.(type) {

		case xml.StartElement:
			ns := strings.ToLower(t.Name.Space)
			tag := strings.ToLower(t.Name.Local)
			fmt.Printf("Start %-20s %s\n", tag, ns)
			if ns == NSAtom && tag == "feed" {

			}

		case xml.EndElement:
			ns := strings.ToLower(t.Name.Space)
			tag := strings.ToLower(t.Name.Local)
			fmt.Printf("End   %-20s %s\n", tag, ns)

		case xml.CharData:
			text := string([]byte(t))
			if strings.TrimSpace(text) == "" {
				continue
			}
			//fmt.Println(text)

		} // switch token

	} // loop

	return

}
