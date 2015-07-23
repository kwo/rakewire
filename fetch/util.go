package fetch

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	m "rakewire.com/model"
	"strings"
	"time"
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

func checksum(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	d := hash.Sum(nil)
	return hex.EncodeToString(d)
}

func parseDateHeader(value string) *time.Time {
	var result *time.Time
	m, err := http.ParseTime(value)
	if err == nil && !m.IsZero() {
		result = &m
	}
	return result
}
