package fetch

import (
	"io"
)

// ReadCounter is a ReadCloser that counts the total bytes read
type ReadCounter struct {
	io.ReadCloser
	Size int
}

func (z *ReadCounter) Read(p []byte) (int, error) {
	n, err := z.ReadCloser.Read(p)
	if err == nil {
		z.Size += n
	}
	return n, err
}
