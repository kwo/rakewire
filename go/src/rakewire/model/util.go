package model

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ParseFeedsFromFile returns the uncommented, non-blank lines from a file
func ParseFeedsFromFile(filename string) ([]*Feed, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ParseFeedsFromReader(f), nil

}

// ParseFeedsFromReader parse url list to feeds
func ParseFeedsFromReader(r io.Reader) []*Feed {
	feeds := []*Feed{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			feeds = append(feeds, NewFeed(url))
		}
	}
	return feeds
}
