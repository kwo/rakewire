package fetch

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
)

func parseDateHeader(value string) time.Time {
	var result time.Time
	m, err := http.ParseTime(value)
	if err == nil && !m.IsZero() {
		result = m
	}
	return result
}

func usesGzip(header string) bool {
	return strings.Contains(header, "gzip")
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
