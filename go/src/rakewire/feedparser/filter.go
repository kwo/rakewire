package feedparser

import (
	"bytes"
	"io"
)

// NewFilterReader constructs a new filtered reader
func NewFilterReader(r io.ReadCloser) io.ReadCloser {
	return &FilterReader{source: r}
}

// FilterReader filters invalid characters from an XML string
type FilterReader struct {
	source io.ReadCloser
}

// Close closes the reader
func (z *FilterReader) Close() error {
	return z.source.Close()
}

// Close closes the reader
func (z *FilterReader) Read(p []byte) (int, error) {

	buf := make([]byte, len(p))

	n, err := z.source.Read(buf)
	if err != nil {
		return 0, err
	}

	cleanBuf := removeControlBytes(buf[:n])
	for i, b := range cleanBuf {
		p[i] = b
	}

	return len(cleanBuf), nil

}

func removeControlBytes(data []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, data)
}
