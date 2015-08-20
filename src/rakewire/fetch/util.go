package fetch

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	m "rakewire/model"
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

func parseDateHeader(value string) *time.Time {
	var result *time.Time
	m, err := http.ParseTime(value)
	if err == nil && !m.IsZero() {
		result = &m
	}
	return result
}

func usesGzip(header string) bool {
	return strings.Contains(header, "gzip")
}

func isFeedUpdated(newTime time.Time, lastTime *time.Time) bool {

	if !newTime.IsZero() && lastTime != nil {
		return newTime.After(*lastTime)
	} else if !newTime.IsZero() {
		return true
	}

	return false

}

func readBody(rsp *http.Response) (io.ReadCloser, error) {

	if rsp.Body == nil {
		return nil, nil
	}

	if usesGzip(rsp.Header.Get(hContentEncoding)) {
		var reader io.ReadCloser
		var err error
		reader, err = gzip.NewReader(rsp.Body)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}

	return rsp.Body, nil

}

func unzipReader(data io.Reader) ([]byte, error) {

	r, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}

	var uncompressedData bytes.Buffer
	if _, err = io.Copy(&uncompressedData, r); err != nil {
		return nil, err
	}

	return uncompressedData.Bytes(), nil

}
