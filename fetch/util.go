package fetch

import (
	"bufio"
	"io"
	m "rakewire.com/model"
	"strings"
)

// URLListToFeeds parse url list to feeds
func URLListToFeeds(r io.Reader) *m.Feeds {

	result := m.NewFeeds()
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			req := m.NewFeed(url)
			result.Add(req)
		}
	}

	return result

}
