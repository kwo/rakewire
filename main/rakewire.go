package main

import (
	"fmt"
	"os"
	"path"
	"rakewire.com/feed"
	"strings"
)

func main() {

	var feedURL string
	if len(os.Args[1:]) == 1 {
		feedURL = strings.TrimSpace(os.Args[1])
	}
	if len(feedURL) == 0 {
		fmt.Printf("Usage: %s <feed-url>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	feed.Parse(feedURL)

}
