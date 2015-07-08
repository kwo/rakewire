package db

import (
	"bytes"
	"compress/gzip"
	"io"
)

func zip(data []byte) ([]byte, error) {

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

func unzip(data []byte) ([]byte, error) {

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
