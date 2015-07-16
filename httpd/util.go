package httpd

import (
	"bytes"
	"compress/gzip"
	"github.com/codegangsta/negroni"
	"io"
	"rakewire.com/logging"
)

func newInternalLogger() *negroni.Logger {
	return &negroni.Logger{logging.New("httpd")}
}

func zipBytes(data []byte) ([]byte, error) {

	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err = w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func unzipBytes(data []byte) ([]byte, error) {

	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var uncompressedData bytes.Buffer
	if _, err = io.Copy(&uncompressedData, r); err != nil {
		return nil, err
	}

	return uncompressedData.Bytes(), nil

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
