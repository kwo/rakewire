package feedparser

import (
	"bytes"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data" // required by go-charset
	"io"
)

// NewFilterCharsetReader constructs a new filtered reader
// The filter will operate on UTF8, which is the same as ASCII for the control characters.
func NewFilterCharsetReader(characterset string, r io.Reader) (io.Reader, error) {
	cr, err := charset.NewReader(characterset, r)
	if err != nil {
		return nil, err
	}
	return NewFilterReader(cr), nil
}

// NewFilterReader constructs a new filtered reader
func NewFilterReader(r io.Reader) io.Reader {
	return &FilterReader{source: r}
}

// FilterReader filters invalid characters from an XML string
type FilterReader struct {
	source io.Reader
}

// Close closes the reader
func (z *FilterReader) Close() error {
	if closer, ok := z.source.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Close closes the reader
func (z *FilterReader) Read(p []byte) (int, error) {

	buf := make([]byte, len(p))
	n, err := z.source.Read(buf)
	if err != nil {
		return 0, err
	}

	clean := removeControlBytes(buf[:n])
	for i, b := range clean {
		p[i] = b
	}

	return len(clean), nil

}

func removeControlBytes(data []byte) []byte {
	buf := bytes.Buffer{}
	for _, b := range data {
		if b >= 32 && b != 127 {
			buf.WriteByte(b)
		}
	}
	return buf.Bytes()
}
