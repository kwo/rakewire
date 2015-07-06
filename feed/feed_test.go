package feed

import (
	//"github.com/stretchr/testify/assert"
	"fmt"
	"testing"
)

func TestFeed(t *testing.T) {
	var feed, _ = Parse("https://ostendorf.com/feed.xml")
	fmt.Printf("%s, %s, %s, %s\n", feed.Title, feed.ID, feed.Date, feed.Author.Name)
	for _, entry := range feed.Entries {
		fmt.Printf("%s, %s, %s, %s\n", entry.Title, entry.ID, entry.Date, entry.Author.Name)
	}
}
